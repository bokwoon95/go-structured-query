package sq

import (
	"context"
	"database/sql"
	"errors"
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
		{"From", WithDefaultLog(Linterpolate).DeleteFrom(u), "DELETE FROM u", nil},
		func() TT {
			desc := "Joins"
			u, ur := USERS().As("u"), USER_ROLES()
			q := WithDefaultLog(Lverbose).
				DeleteFrom(u, ur).
				Using(u).
				Join(ur, ur.USER_ID.Eq(u.USER_ID)).
				LeftJoin(u, u.USER_ID.Eq(u.USER_ID)).
				RightJoin(u, u.USER_ID.Eq(u.USER_ID)).
				FullJoin(u, u.USER_ID.Eq(u.USER_ID)).
				CustomJoin("CROSS JOIN", u).
				Where(u.USER_ID.EqInt(1)).
				OrderBy(u.DISPLAYNAME, u.EMAIL.Desc()).
				Limit(-10)
			wantQuery := "DELETE FROM u, devlab.user_roles" +
				" USING devlab.users AS u" +
				" JOIN devlab.user_roles ON user_roles.user_id = u.user_id" +
				" LEFT JOIN devlab.users AS u ON u.user_id = u.user_id" +
				" RIGHT JOIN devlab.users AS u ON u.user_id = u.user_id" +
				" FULL JOIN devlab.users AS u ON u.user_id = u.user_id" +
				" CROSS JOIN devlab.users AS u" +
				" WHERE u.user_id = ?" +
				" ORDER BY u.displayname, u.email DESC" +
				" LIMIT ?"
			wantArgs := []interface{}{1, int64(10)}
			return TT{desc, q, wantQuery, wantArgs}
		}(),
		func() TT {
			var tt TT
			tt.description = "assorted"
			cte1 := SelectOne().From(u).CTE("cte1")
			cte2 := SelectDistinct(u.EMAIL).From(u).CTE("cte2")
			tt.q = WithDefaultLog(Lverbose).
				DeleteFrom(u).
				Using(u).
				CustomJoin("NATURAL JOIN", cte1).
				CustomJoin("NATURAL JOIN", cte2).
				Where(u.USER_ID.Eq(u.USER_ID))
			tt.wantQuery = "WITH cte1 AS (SELECT 1 FROM devlab.users AS u)" +
				", cte2 AS (SELECT DISTINCT u.email FROM devlab.users AS u)" +
				" DELETE FROM u" +
				" USING devlab.users AS u" +
				" NATURAL JOIN cte1" +
				" NATURAL JOIN cte2" +
				" WHERE u.user_id = u.user_id"
			return tt
		}(),
		func() TT {
			desc := "aliasless table"
			u := USERS()
			q := WithDefaultLog(0).DeleteFrom(u)
			wantQuery := "DELETE FROM devlab.users"
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
		Exec(nil, 0)
	is.True(err != nil)

	// SQL syntax error
	// use a tempDB so we don't foul up the current db transaction with the error
	tempDB, err := sql.Open("txdb", randomString(8))
	is.NoErr(err)
	_, err = WithDefaultLog(Linterpolate).
		WithDB(tempDB).
		DeleteFrom(s).
		Where(Predicatef("ERROR")).
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
		Exec(nil, ElastInsertID|ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)
}
