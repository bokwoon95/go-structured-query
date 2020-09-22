package sq

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestDeleteQuery_ToSQL(t *testing.T) {
	type TT struct {
		description string
		q           DeleteQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{"empty", DeleteQuery{}, "DELETE FROM NULL", nil},
		{"From", WithDefaultLog(Linterpolate).DeleteFrom(u), "DELETE FROM public.users AS u", nil},
		{
			"Joins",
			WithDefaultLog(Lverbose).
				DeleteFrom(u).
				Using(u).
				Join(u, u.USER_ID.Eq(u.USER_ID)).
				LeftJoin(u, u.USER_ID.Eq(u.USER_ID)).
				RightJoin(u, u.USER_ID.Eq(u.USER_ID)).
				FullJoin(u, u.USER_ID.Eq(u.USER_ID)).
				CustomJoin("CROSS JOIN", u).
				ReturningOne(),
			"DELETE FROM public.users AS u" +
				" USING public.users AS u" +
				" JOIN public.users AS u ON u.user_id = u.user_id" +
				" LEFT JOIN public.users AS u ON u.user_id = u.user_id" +
				" RIGHT JOIN public.users AS u ON u.user_id = u.user_id" +
				" FULL JOIN public.users AS u ON u.user_id = u.user_id" +
				" CROSS JOIN public.users AS u" +
				" RETURNING 1",
			nil,
		},
		func() TT {
			var tt TT
			tt.description = "assorted"
			cte1 := DeleteFrom(u).Where(Bool(true)).Returning(u.USER_ID).CTE("cte1")
			cte2 := DeleteFrom(u).Where(Bool(false)).ReturningOne().CTE("cte2")
			tt.q = WithDefaultLog(Lverbose).
				DeleteFrom(u).
				Using(u).
				CustomJoin("NATURAL JOIN", cte1).
				CustomJoin("NATURAL JOIN", cte2).
				Where(u.USER_ID.Eq(u.USER_ID)).
				Returning(u.USER_ID, u.DISPLAYNAME, u.EMAIL)
			tt.wantQuery = "WITH cte1 AS (DELETE FROM public.users AS u WHERE $1 RETURNING u.user_id)" +
				", cte2 AS (DELETE FROM public.users AS u WHERE $2 RETURNING 1)" +
				" DELETE FROM public.users AS u" +
				" USING public.users AS u" +
				" NATURAL JOIN cte1" +
				" NATURAL JOIN cte2" +
				" WHERE u.user_id = u.user_id" +
				" RETURNING u.user_id, u.displayname, u.email"
			tt.wantArgs = []interface{}{true, false}
			return tt
		}(),
		func() TT {
			desc := "aliasless table"
			u := USERS()
			q := WithDefaultLog(0).DeleteFrom(u)
			wantQuery := "DELETE FROM public.users"
			return TT{desc, q, wantQuery, nil}
		}(),
		func() TT {
			var tt TT
			tt.description = "subqueries"
			cte1 := DeleteFrom(u).Where(Bool(true)).Returning(u.USER_ID).CTE("cte1")
			cte2 := DeleteFrom(u).Where(Bool(false)).ReturningOne().CTE("cte2")
			subquery1 := DeleteFrom(u).Where(Bool(true)).Returning(u.USER_ID).Subquery("subquery1")
			subquery2 := DeleteFrom(u).Where(Bool(false)).ReturningOne().Subquery("subquery2")
			tt.q = WithDefaultLog(0).
				DeleteFrom(u).
				Using(subquery1).
				CustomJoin("NATURAL JOIN", cte1).
				CustomJoin("NATURAL JOIN", cte2).
				Join(subquery2, subquery1["user_id"].Eq(cte1["user_id"]))
			tt.wantQuery = "WITH cte1 AS (DELETE FROM public.users AS u WHERE $1 RETURNING u.user_id)" +
				", cte2 AS (DELETE FROM public.users AS u WHERE $2 RETURNING 1)" +
				" DELETE FROM public.users AS u" +
				" USING (DELETE FROM public.users AS u WHERE $3 RETURNING u.user_id) AS subquery1" +
				" NATURAL JOIN cte1" +
				" NATURAL JOIN cte2" +
				" JOIN (DELETE FROM public.users AS u WHERE $4 RETURNING 1) AS subquery2 ON subquery1.user_id = cte1.user_id"
			tt.wantArgs = []interface{}{true, false, true, false}
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

func TestDeleteQuery_Fetch(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "DeleteQuery_Fetch")
	is.NoErr(err)
	defer db.Close()
	s := SUBMISSIONS()

	// Missing DB
	err = DeleteFrom(s).
		ReturningRowx(func(row *Row) {}).
		Fetch(nil)
	is.True(err != nil)

	// SQL syntax error
	// use a tempDB so we don't foul up the current db transaction with the error
	tempDB, err := sql.Open("txdb", randomString(8))
	is.NoErr(err)
	var submissionID int
	err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		DeleteFrom(s).
		ReturningRowx(func(row *Row) {
			row.ScanInto(&submissionID, s.SUBMISSION_ID.Asc().NullsLast())
		}).
		Fetch(nil)
	is.True(err != nil)
	tempDB.Close()

	// No mapper
	err = WithDB(db).
		DeleteFrom(s).
		Fetch(nil)
	is.True(err != nil)

	// Empty mapper
	err = WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(0)).
		ReturningRowx(func(row *Row) {}).
		Fetch(nil)
	is.NoErr(err)

	// Wrong Scan type
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	err = WithDefaultLog(Lverbose).
		WithDB(tempDB).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {
			row.ScanInto(&submissionID, s.CREATED_AT)
		}).
		Fetch(nil)
	is.True(err != nil)
	tempDB.Close()

	// sql.ErrNoRows
	err = WithDefaultLog(Linterpolate).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(-99999)).
		ReturningRowx(func(row *Row) {
			row.Int(s.SUBMISSION_ID)
		}).
		Fetch(nil)
	is.True(errors.Is(err, sql.ErrNoRows))

	// simulate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	err = WithDefaultLog(Lverbose).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(0)).
		ReturningRowx(func(row *Row) {}).
		FetchContext(ctx, nil)
	is.True(errors.Is(err, context.DeadlineExceeded))

	// Mapper
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	err = WithDefaultLog(Lverbose).
		WithDB(tempDB).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(1)).
		ReturningRowx(func(row *Row) {
			submissionID = row.Int(s.SUBMISSION_ID)
		}).
		Fetch(nil)
	is.NoErr(err)
	is.Equal(1, submissionID)
	tempDB.Close()

	// Accumulator
	tempDB, err = sql.Open("txdb", randomString(8))
	is.NoErr(err)
	var submissionIDs []int
	err = WithDefaultLog(Lverbose).
		WithDB(tempDB).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.In([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})).
		Returningx(func(row *Row) {
			submissionID = row.Int(s.SUBMISSION_ID)
		}, func() {
			submissionIDs = append(submissionIDs, submissionID)
		}).
		Fetch(nil)
	is.NoErr(err)
	is.Equal(10, len(submissionIDs))
	tempDB.Close()

	// Panic with ExitPeacefully
	submissionIDs = submissionIDs[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.In([]int{3, 4})).
		Returningx(func(row *Row) {
			submissionID = row.Int(s.SUBMISSION_ID)
		}, func() {
			panic(ExitPeacefully)
		}).
		Fetch(nil)
	is.NoErr(err)
	is.Equal(0, len(submissionIDs))

	// Panic with any other ExitCode
	submissionIDs = submissionIDs[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.In([]int{5, 6})).
		Returningx(func(row *Row) {
			submissionID = row.Int(s.SUBMISSION_ID)
		}, func() {
			panic(ExitCode(1))
		}).
		Fetch(nil)
	is.True(errors.Is(err, ExitCode(1)))

	// Panic with error
	ErrTest := errors.New("this is a test error")
	submissionIDs = submissionIDs[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.In([]int{7, 8})).
		Returningx(func(row *Row) {
			submissionID = row.Int(s.SUBMISSION_ID)
		}, func() {
			panic(ErrTest)
		}).
		Fetch(nil)
	is.True(errors.Is(err, ErrTest))

	// Panic with 0
	submissionIDs = submissionIDs[:0]
	err = WithDefaultLog(Linterpolate).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.In([]int{9, 10})).
		Returningx(func(row *Row) {
			submissionID = row.Int(s.SUBMISSION_ID)
		}, func() {
			panic(0)
		}).
		Fetch(nil)
	is.Equal(fmt.Errorf("0").Error(), err.Error())
}

func TestDeleteQuery_Exec(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "DeleteQuery_Exec")
	is.NoErr(err)
	defer db.Close()
	s := SUBMISSIONS()

	// Missing DB
	_, err = DeleteFrom(s).
		ReturningRowx(func(row *Row) {}).
		Exec(nil, 0)
	is.True(err != nil)

	// SQL syntax error
	// use a tempDB so we don't foul up the current db transaction with the error
	tempDB, err := sql.Open("txdb", randomString(8))
	is.NoErr(err)
	_, err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		DeleteFrom(s).
		Returning(Fieldf("ERROR")).
		Exec(nil, ErowsAffected)
	is.True(err != nil)
	tempDB.Close()

	// simulate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	_, err = WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(0)).
		ExecContext(ctx, nil, ErowsAffected)
	is.True(errors.Is(err, context.DeadlineExceeded))

	// rowsAffected
	rowsAffected, err := WithDefaultLog(Lverbose).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(1)).
		Exec(nil, ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)

	rowsAffected, err = WithDefaultLog(Lverbose).
		WithDB(db).
		DeleteFrom(s).
		Where(s.SUBMISSION_ID.EqInt(2)).
		Exec(nil, ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)
}
