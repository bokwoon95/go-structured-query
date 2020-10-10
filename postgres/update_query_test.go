package sq

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestUpdateQuery_ToSQL(t *testing.T) {
	type TT struct {
		description string
		q           UpdateQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{"empty", UpdateQuery{}, "UPDATE NULL", nil},
		{"Update", WithDefaultLog(Linterpolate).Update(u).Set(u.USER_ID.SetInt(1)), "UPDATE public.users AS u SET user_id = $1", []interface{}{1}},
		{
			"Joins",
			WithDefaultLog(Lverbose).
				Update(u).
				From(SelectOne().From(u).Subquery("subquery")).
				Join(u, u.USER_ID.Eq(u.USER_ID)).
				LeftJoin(u, u.USER_ID.Eq(u.USER_ID)).
				RightJoin(u, u.USER_ID.Eq(u.USER_ID)).
				FullJoin(u, u.USER_ID.Eq(u.USER_ID)).
				CustomJoin("CROSS JOIN", u),
			"UPDATE public.users AS u FROM (SELECT 1 FROM public.users AS u) AS subquery" +
				" JOIN public.users AS u ON u.user_id = u.user_id" +
				" LEFT JOIN public.users AS u ON u.user_id = u.user_id" +
				" RIGHT JOIN public.users AS u ON u.user_id = u.user_id" +
				" FULL JOIN public.users AS u ON u.user_id = u.user_id" +
				" CROSS JOIN public.users AS u",
			nil,
		},
		func() TT {
			var tt TT
			tt.description = "assorted"
			cte1 := Update(u).Where(Bool(true)).ReturningOne().CTE("cte1")
			cte2 := Update(u).Where(Bool(true)).ReturningOne().CTE("cte2")
			tt.q = WithDefaultLog(Lverbose).
				Update(u).
				From(Update(u).Where(Bool(false)).ReturningOne().Subquery("subquery")).
				CustomJoin("NATURAL JOIN", cte1).
				CustomJoin("NATURAL JOIN", cte2).
				Where(u.USER_ID.Eq(u.USER_ID)).
				Returning(u.USER_ID, u.DISPLAYNAME, u.EMAIL)
			tt.wantQuery = "WITH cte1 AS (UPDATE public.users AS u WHERE $1 RETURNING 1)" +
				", cte2 AS (UPDATE public.users AS u WHERE $2 RETURNING 1)" +
				" UPDATE public.users AS u" +
				" FROM (UPDATE public.users AS u WHERE $3 RETURNING 1) AS subquery" +
				" NATURAL JOIN cte1" +
				" NATURAL JOIN cte2" +
				" WHERE u.user_id = u.user_id" +
				" RETURNING u.user_id, u.displayname, u.email"
			tt.wantArgs = []interface{}{true, true, false}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "Setx Basic"
			user := User{
				UserID:      1,
				Displayname: "Bob",
				Email:       "bob@email.com",
				Password:    "cant_hack_me",
			}
			u := USERS().As("u")
			tt.q = WithDefaultLog(Lverbose).
				Update(u).
				Setx(func(col *Column) {
					col.SetString(u.DISPLAYNAME, user.Displayname)
					col.SetString(u.EMAIL, user.Email)
					col.SetString(u.PASSWORD, user.Password)
				}).
				Where(u.USER_ID.EqInt(user.UserID)).
				Returning(u.USER_ID, u.DISPLAYNAME, u.EMAIL, u.PASSWORD)
			tt.wantQuery = "UPDATE public.users AS u" +
				" SET displayname = $1, email = $2, password = $3" +
				" WHERE u.user_id = $4" +
				" RETURNING u.user_id, u.displayname, u.email, u.password"
			tt.wantArgs = []interface{}{user.Displayname, user.Email, user.Password, user.UserID}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "Setx RowValue Basic"
			user := User{
				UserID:      1,
				Displayname: "Bob",
				Email:       "bob@email.com",
				Password:    "cant_hack_me",
			}
			u := USERS().As("u")
			tt.q = WithDefaultLog(Lverbose).
				Update(u).
				Setx(func(col *Column) {
					col.Set(
						RowValue{u.DISPLAYNAME, u.EMAIL, u.PASSWORD},
						RowValue{user.Displayname, user.Email, user.Password},
					)
				}).
				Where(u.USER_ID.EqInt(user.UserID)).
				Returning(u.USER_ID, u.DISPLAYNAME, u.EMAIL, u.PASSWORD)
			tt.wantQuery = "UPDATE public.users AS u" +
				" SET (displayname, email, password) = ($1, $2, $3)" +
				" WHERE u.user_id = $4" +
				" RETURNING u.user_id, u.displayname, u.email, u.password"
			tt.wantArgs = []interface{}{user.Displayname, user.Email, user.Password, user.UserID}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "Setx RowValue Select"
			user := User{
				UserID:      1,
				Displayname: "Bob",
				Email:       "bob@email.com",
				Password:    "cant_hack_me",
			}
			u1, u2 := USERS().As("u1"), USERS().As("u2")
			tt.q = WithDefaultLog(Lverbose).
				Update(u1).
				Setx(func(col *Column) {
					col.Set(
						RowValue{u1.DISPLAYNAME, u1.EMAIL, u1.PASSWORD},
						Select(u2.DISPLAYNAME, u2.EMAIL, u2.PASSWORD).
							From(u2).
							Where(u2.USER_ID.EqInt(99)),
					)
				}).
				Where(u1.USER_ID.EqInt(user.UserID)).
				Returning(u1.USER_ID, u1.DISPLAYNAME, u1.EMAIL, u1.PASSWORD)
			tt.wantQuery = "UPDATE public.users AS u1" +
				" SET (displayname, email, password) =" +
				" (SELECT u2.displayname, u2.email, u2.password FROM public.users AS u2 WHERE u2.user_id = $1)" +
				" WHERE u1.user_id = $2" +
				" RETURNING u1.user_id, u1.displayname, u1.email, u1.password"
			tt.wantArgs = []interface{}{99, user.UserID}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "ToSQL ColumnMapper panic translates to empty query and panicked value in args"
			user := User{}
			u := USERS().As("u")
			var errEmptyEmail = errors.New("email cannot be empty")
			tt.q = WithDefaultLog(Lverbose).
				Update(u).
				Setx(func(col *Column) {
					if user.Email == "" {
						panic(errEmptyEmail)
					}
					col.SetString(u.DISPLAYNAME, user.Displayname)
					col.SetString(u.EMAIL, user.Email)
					col.SetString(u.PASSWORD, user.Password)
				}).
				Where(u.USER_ID.EqInt(1)).
				Returning(u.USER_ID)
			tt.wantQuery = ""
			tt.wantArgs = []interface{}{errEmptyEmail}
			return tt
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			var _ Query = tt.q
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestUpdateQuery_Fetch(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "UpdateQuery_Fetch")
	is.NoErr(err)
	defer db.Close()
	u := USERS()

	// Missing DB
	err = Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {}).
		Fetch(nil)
	is.True(err != nil)

	// SQL syntax error
	// use a tempDB so we don't foul up the current db transaction with the error
	tempDB, err := sql.Open("txdb", randomString(8))
	is.NoErr(err)
	var userID int
	err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {
			row.ScanInto(&userID, u.USER_ID.Asc().NullsLast())
		}).
		Fetch(nil)
	is.True(err != nil)
	tempDB.Close()

	// No mapper
	err = WithDB(db).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		Fetch(nil)
	is.True(err != nil)

	// Empty mapper
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	err = WithDB(tempDB).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {}).
		Fetch(nil)
	is.NoErr(err)
	tempDB.Close()

	// Wrong Scan type
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	err = WithDefaultLog(Lverbose).
		WithDB(tempDB).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {
			row.ScanInto(&userID, u.DISPLAYNAME)
		}).
		Fetch(nil)
	is.True(err != nil)
	tempDB.Close()

	// sql.ErrNoRows
	err = WithDefaultLog(Linterpolate).
		WithDB(db).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(-99999)).
		ReturningRowx(func(row *Row) {
			row.Int(u.USER_ID)
		}).
		Fetch(nil)
	is.True(errors.Is(err, sql.ErrNoRows))

	// simulate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	err = WithDefaultLog(Lverbose).
		WithDB(db).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {}).
		FetchContext(ctx, nil)
	is.True(errors.Is(err, context.DeadlineExceeded))

	// RowMapper
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	var email string
	err = WithDefaultLog(Lverbose).
		WithDB(tempDB).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {
			userID = row.Int(u.USER_ID)
		}).
		Fetch(nil)
	is.NoErr(err)
	is.Equal(1, userID)
	tempDB.Close()

	// Accumulator
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	var emails []string
	err = WithDefaultLog(Lverbose).
		WithDB(tempDB).
		Update(u).
		Set(u.EMAIL.Set(Fieldf("?::TEXT", u.USER_ID))).
		Where(u.USER_ID.In([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})).
		Returningx(func(row *Row) {
			email = row.String(u.EMAIL)
		}, func() {
			emails = append(emails, email)
		}).
		Fetch(nil)
	is.NoErr(err)
	is.Equal(10, len(emails))
	tempDB.Close()

	// ColumnMapper
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	user := User{
		Displayname: "Bob",
		Email:       "bob@email.com",
		Password:    "cant_hack_me",
	}
	err = WithDefaultLog(Lverbose).
		WithDB(tempDB).
		Update(u).
		Setx(func(col *Column) {
			col.SetString(u.DISPLAYNAME, user.Displayname)
			col.SetString(u.EMAIL, user.Email)
			col.SetString(u.PASSWORD, user.Password)
		}).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) { userID = row.Int(u.USER_ID) }).
		Fetch(nil)
	is.NoErr(err)
	is.Equal(1, userID)
	tempDB.Close()

	// Panic with validation error in ColumnMapper
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	var errEmptyEmail = errors.New("email cannot be empty")
	user = User{
		Displayname: "Bob",
		Email:       "", // Empty email
		Password:    "cant_hack_me",
	}
	_, err = WithDefaultLog(Lverbose).
		Update(u).
		Setx(func(col *Column) {
			if user.Email == "" {
				panic(errEmptyEmail)
			}
			col.SetString(u.DISPLAYNAME, user.Displayname)
			col.SetString(u.EMAIL, user.Email)
			col.SetString(u.PASSWORD, user.Password)
		}).
		Where(u.USER_ID.EqInt(1)).
		Exec(tempDB, 0)
	is.Equal(err, errEmptyEmail)
	tempDB.Close()

	// Panic with ExitPeacefully
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	emails = emails[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		Update(u).
		Set(u.EMAIL.Set(Fieldf("?::TEXT", u.USER_ID))).
		Where(u.USER_ID.In([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})).
		Returningx(func(row *Row) {
			email = row.String(u.EMAIL)
		}, func() {
			panic(ExitPeacefully)
		}).
		Fetch(nil)
	is.NoErr(err)
	is.Equal(0, len(emails))
	tempDB.Close()

	// Panic with any other ExitCode
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	emails = emails[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		Update(u).
		Set(u.EMAIL.Set(Fieldf("?::TEXT", u.USER_ID))).
		Where(u.USER_ID.In([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})).
		Returningx(func(row *Row) {
			email = row.String(u.EMAIL)
		}, func() {
			panic(ExitCode(1))
		}).
		Fetch(nil)
	is.True(errors.Is(err, ExitCode(1)))
	tempDB.Close()

	// Panic with error
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	ErrTest := errors.New("this is a test error")
	emails = emails[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		Update(u).
		Set(u.EMAIL.Set(Fieldf("?::TEXT", u.USER_ID))).
		Where(u.USER_ID.In([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})).
		Returningx(func(row *Row) {
			email = row.String(u.EMAIL)
		}, func() {
			panic(ErrTest)
		}).
		Fetch(nil)
	is.True(errors.Is(err, ErrTest))
	is.Equal(0, len(emails))
	tempDB.Close()

	// Panic with 0
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	emails = emails[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		Update(u).
		Set(u.EMAIL.Set(Fieldf("?::TEXT", u.USER_ID))).
		Where(u.USER_ID.In([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})).
		Returningx(func(row *Row) {
			email = row.String(u.EMAIL)
		}, func() {
			panic(0)
		}).
		Fetch(nil)
	is.Equal(fmt.Errorf("0").Error(), err.Error())
	tempDB.Close()
}

func TestUpdateQuery_Exec(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "UpdateQuery_Exec")
	is.NoErr(err)
	defer db.Close()
	u := USERS()

	// Missing DB
	_, err = Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {}).
		Exec(nil, ErowsAffected)
	is.True(err != nil)

	// SQL syntax error
	// use a tempDB so we don't foul up the current db transaction with the error
	tempDB, err := sql.Open("txdb", randomString(8))
	is.NoErr(err)
	_, err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		Update(u).
		Set(u.USER_ID.Set(Fieldf("ERROR"))).
		Where(u.USER_ID.EqInt(1)).
		Exec(nil, ErowsAffected)
	is.True(err != nil)
	tempDB.Close()

	// simulate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	_, err = WithDefaultLog(Lverbose).
		WithDB(db).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		ExecContext(ctx, nil, ErowsAffected)
	is.True(errors.Is(err, context.DeadlineExceeded))

	// rowsAffected
	rowsAffected, err := WithDefaultLog(Lverbose).
		WithDB(db).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		Exec(nil, ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)

	// again
	rowsAffected, err = WithDefaultLog(Lverbose).
		WithDB(db).
		Update(u).
		Set(u.USER_ID.SetInt(1)).
		Where(u.USER_ID.EqInt(1)).
		Exec(nil, ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)
}
