package sq

import (
	"testing"

	"github.com/matryer/is"
)

func TestCTEs_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		q           Query
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			var tt TT
			tt.description = "Select CTE"
			u := USERS().As("u")
			cte := Select(u.USER_ID, u.DISPLAYNAME, u.EMAIL).From(u).Where(u.USER_ID.LtInt(5)).CTE("cte")
			tt.q = Select(cte["user_id"], cte["displayname"]).From(cte).Where(cte["displayname"].Eq(cte["email"]))
			tt.wantQuery = "WITH cte AS" +
				" (SELECT u.user_id, u.displayname, u.email FROM devlab.users AS u WHERE u.user_id < ?)" +
				" SELECT cte.user_id, cte.displayname FROM cte WHERE cte.displayname = cte.email"
			tt.wantArgs = []interface{}{5}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "Select CTE aliased"
			u := USERS().As("u")
			apple := Select(u.USER_ID, u.DISPLAYNAME, u.EMAIL).From(u).Where(u.USER_ID.LtInt(5)).CTE("apple")
			banana := apple.As("banana")
			tt.q = Select(banana["user_id"], banana["displayname"], apple["email"]).
				From(banana).
				Join(apple, Int(1).EqInt(1)).
				Where(apple["displayname"].Eq(banana["email"]))
			tt.wantQuery = "WITH apple AS" +
				" (SELECT u.user_id, u.displayname, u.email FROM devlab.users AS u WHERE u.user_id < ?)" +
				" SELECT banana.user_id, banana.displayname, apple.email FROM apple AS banana JOIN apple ON ? = ? WHERE apple.displayname = banana.email"
			tt.wantArgs = []interface{}{5, 1, 1}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "Recursive CTE (explicit columns)"
			tens := RecursiveCTE("tens", "n")
			tens = tens.
				Initial(Select(Int(10))).
				UnionAll(
					Select(Fieldf("? + 10", tens["n"])).From(tens).Where(Predicatef("? + 10 <= 100", tens["n"])),
				)
			tt.q = Select(tens["n"]).From(tens)
			tt.wantQuery = "WITH RECURSIVE tens (n) AS" +
				" (SELECT ?" +
				" UNION ALL" +
				" SELECT tens.n + 10 FROM tens WHERE tens.n + 10 <= 100)" +
				" SELECT tens.n FROM tens"
			tt.wantArgs = []interface{}{10}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "Recursive CTE (implicit columns)"
			tens := RecursiveCTE("tens")
			tens = tens.
				Initial(Select(Int(10).As("n"))).
				UnionAll(
					Select(Fieldf("? + 10", tens["n"])).From(tens).Where(Predicatef("? + 10 <= 100", tens["n"])),
				)
			tt.q = Select(tens["n"]).From(tens)
			tt.wantQuery = "WITH RECURSIVE tens AS" +
				" (SELECT ? AS n" +
				" UNION ALL" +
				" SELECT tens.n + 10 FROM tens WHERE tens.n + 10 <= 100)" +
				" SELECT tens.n FROM tens"
			tt.wantArgs = []interface{}{10}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "UNIONing a non recursive CTE should have no effect"
			u := USERS().As("u")
			q1 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(1))
			q2 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(2))
			cte := Select(u.USER_ID, u.DISPLAYNAME, u.EMAIL).From(u).Where(u.USER_ID.LtInt(5)).CTE("cte").Initial(q1).Union(q2)
			tt.q = Select(cte["user_id"], cte["displayname"]).From(cte).Where(cte["displayname"].Eq(cte["email"]))
			tt.wantQuery = "WITH cte AS" +
				" (SELECT u.user_id, u.displayname, u.email FROM devlab.users AS u WHERE u.user_id < ?)" +
				" SELECT cte.user_id, cte.displayname FROM cte WHERE cte.displayname = cte.email"
			tt.wantArgs = []interface{}{5}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "VariadicQuery CTE (explicit columns)"
			u := USERS().As("u")
			q1 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(1))
			q2 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(2))
			q3 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(3))
			q := Union(q1, q2, q3).CTE("cte", "user_id", "email")
			tt.q = Select(q["user_id"], q["email"]).From(q)
			tt.wantQuery = "WITH cte (user_id, email) AS" +
				" (SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?" +
				" UNION" +
				" SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?" +
				" UNION" +
				" SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?)" +
				" SELECT cte.user_id, cte.email FROM cte"
			tt.wantArgs = []interface{}{1, 2, 3}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "VariadicQuery CTE (implicit columns from SELECT)"
			u := USERS().As("u")
			q1 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(1))
			q2 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(2))
			q3 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(3))
			q := UnionAll(q1, q2, q3).CTE("cte")
			tt.q = Select(q["user_id"], q["email"]).From(q)
			tt.wantQuery = "WITH cte AS" +
				" (SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM devlab.users AS u WHERE u.user_id = ?)" +
				" SELECT cte.user_id, cte.email FROM cte"
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
