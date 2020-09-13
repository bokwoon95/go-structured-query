package sq

import (
	"context"
	"database/sql"
	"errors"
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
		{"Update", WithDefaultLog(Linterpolate).Update(u).Set(u.USER_ID.SetInt(1)), "UPDATE devlab.users AS u SET u.user_id = ?", []interface{}{1}},
		{
			"Joins",
			WithDefaultLog(Lverbose).
				Update(u).
				Join(SelectOne().From(u).Subquery("subquery"), Bool(true)).
				Join(u, u.USER_ID.Eq(u.USER_ID)).
				LeftJoin(u, u.USER_ID.Eq(u.USER_ID)).
				RightJoin(u, u.USER_ID.Eq(u.USER_ID)).
				FullJoin(u, u.USER_ID.Eq(u.USER_ID)).
				CustomJoin("CROSS JOIN", u),
			"UPDATE devlab.users AS u" +
				" JOIN (SELECT 1 FROM devlab.users AS u) AS subquery ON ?" +
				" JOIN devlab.users AS u ON u.user_id = u.user_id" +
				" LEFT JOIN devlab.users AS u ON u.user_id = u.user_id" +
				" RIGHT JOIN devlab.users AS u ON u.user_id = u.user_id" +
				" FULL JOIN devlab.users AS u ON u.user_id = u.user_id" +
				" CROSS JOIN devlab.users AS u",
			[]interface{}{true},
		},
		func() TT {
			var tt TT
			tt.description = "assorted"
			cte1 := SelectOne().From(u).CTE("cte1")
			cte2 := SelectDistinct(u.EMAIL).From(u).CTE("cte2")
			tt.q = WithLog(customLogger, Lverbose).
				Update(u).
				Join(u, Bool(true)).
				CustomJoin("NATURAL JOIN", cte1).
				CustomJoin("NATURAL JOIN", cte2).
				Where(u.USER_ID.Eq(u.USER_ID)).
				OrderBy(u.DISPLAYNAME, u.EMAIL.Desc()).
				Limit(-10)
			tt.wantQuery = "WITH cte1 AS (SELECT 1 FROM devlab.users AS u)" +
				", cte2 AS (SELECT DISTINCT u.email FROM devlab.users AS u)" +
				" UPDATE devlab.users AS u" +
				" JOIN devlab.users AS u ON ?" +
				" NATURAL JOIN cte1" +
				" NATURAL JOIN cte2" +
				" WHERE u.user_id = u.user_id" +
				" ORDER BY u.displayname, u.email DESC" +
				" LIMIT ?"
			tt.wantArgs = []interface{}{true, int64(10)}
			return tt
		}(),
		func() TT {
			desc := "aliasless table"
			u := USERS()
			q := WithDefaultLog(0).Update(u)
			wantQuery := "UPDATE devlab.users"
			return TT{desc, q, wantQuery, nil}
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
		Exec(nil, ErowsAffected)
	is.True(err != nil)

	// SQL syntax error
	// use a tempDB so we don't foul up the current db transaction with the error
	tempDB, err := sql.Open("txdb", RandomString(8))
	is.NoErr(err)
	_, err = WithLog(customLogger, Linterpolate).
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
	tempDB, err = sql.Open("txdb", RandomString(8))
	is.NoErr(err)
	rowsAffected, err := WithDefaultLog(Lverbose).
		WithDB(tempDB).
		Update(u).
		Set(
			u.DISPLAYNAME.SetString("aaa"),
			u.USER_ID.Set(Fieldf("last_insert_id(?)", u.USER_ID)),
		).
		Where(u.USER_ID.EqInt(3)).
		Exec(nil, ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)
	var lastInsertID int
	err = SelectRowx(func(row *Row) { row.ScanInto(&lastInsertID, Fieldf("last_insert_id()")) }).Fetch(tempDB)
	is.NoErr(err)
	is.Equal(3, lastInsertID)
	tempDB.Close()

	rowsAffected, err = WithLog(customLogger, Lverbose).
		WithDB(db).
		Update(u).
		Set(u.DISPLAYNAME.SetString("bbb")).
		Where(u.USER_ID.EqInt(4)).
		Exec(nil, ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)
}
