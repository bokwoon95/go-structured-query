package sq

import (
	"testing"

	"github.com/matryer/is"
)

func TestSubquery(t *testing.T) {
	type TT struct {
		description string
		q           Query
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			var tt TT
			tt.description = "Basic subquery"
			u := USERS().As("u")
			subq := Select(u.USER_ID, u.DISPLAYNAME, u.EMAIL).From(u).Where(u.USER_ID.LtInt(5)).Subquery("subq")
			tt.q = Select(subq["user_id"], subq["displayname"]).From(subq).Where(subq["displayname"].Eq(subq["email"]))
			tt.wantQuery = "SELECT subq.user_id, subq.displayname FROM" +
				" (SELECT u.user_id, u.displayname, u.email FROM devlab.users AS u WHERE u.user_id < ?) AS subq" +
				" WHERE subq.displayname = subq.email"
			tt.wantArgs = []interface{}{5}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "exists subquery"
			u, ur := USERS().As("u"), USER_ROLES().As("ur")
			subq := SelectOne().From(ur).Where(ur.USER_ID.Eq(u.USER_ID), ur.ROLE.EqString("student")).Subquery("subq")
			tt.q = Select(u.USER_ID, u.DISPLAYNAME).From(u).Where(Exists(subq))
			tt.wantQuery = "SELECT u.user_id, u.displayname FROM devlab.users AS u" +
				" WHERE EXISTS(SELECT 1 FROM devlab.user_roles AS ur WHERE ur.user_id = u.user_id AND ur.role = ?)"
			tt.wantArgs = []interface{}{"student"}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "VariadicQuery subquery (implicit columns from SELECT)"
			u := USERS().As("u")
			q1 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(1))
			q2 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(2))
			q3 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(3))
			q := Union(q1, q2, q3).Subquery("subquery")
			tt.q = Select(q["user_id"], q["email"]).From(q)
			tt.wantQuery = "SELECT subquery.user_id, subquery.email FROM" +
				" (SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?" +
				" UNION" +
				" SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?" +
				" UNION" +
				" SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?) AS subquery"
			tt.wantArgs = []interface{}{1, 2, 3}
			return tt
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}
