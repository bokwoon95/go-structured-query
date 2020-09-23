package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestCTE(t *testing.T) {
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
				" (SELECT u.user_id, u.displayname, u.email FROM public.users AS u WHERE u.user_id < $1)" +
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
				" (SELECT u.user_id, u.displayname, u.email FROM public.users AS u WHERE u.user_id < $1)" +
				" SELECT banana.user_id, banana.displayname, apple.email FROM apple AS banana JOIN apple ON $2 = $3 WHERE apple.displayname = banana.email"
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
				" (SELECT $1" +
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
				" (SELECT $1 AS n" +
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
			cte := Select(u.USER_ID, u.DISPLAYNAME, u.EMAIL).From(u).Where(u.USER_ID.LtInt(5)).CTE("cte")
			cte = cte.Initial(q1).Union(q2)
			tt.q = Select(cte["user_id"], cte["displayname"]).From(cte).Where(cte["displayname"].Eq(cte["email"]))
			tt.wantQuery = "WITH cte AS" +
				" (SELECT u.user_id, u.displayname, u.email FROM public.users AS u WHERE u.user_id < $1)" +
				" SELECT cte.user_id, cte.displayname FROM cte WHERE cte.displayname = cte.email"
			tt.wantArgs = []interface{}{5}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "Insert CTE"
			u := USERS().As("u")
			cte := InsertInto(u).
				Columns(u.USER_ID, u.DISPLAYNAME, u.EMAIL).
				Values(1, "apple", "banana").
				Returning(u.USER_ID).
				CTE("cte")
			tt.q = Select(cte["user_id"]).From(cte)
			tt.wantQuery = "WITH cte AS" +
				" (INSERT INTO public.users AS u (user_id, displayname, email) VALUES ($1, $2, $3) RETURNING u.user_id)" +
				" SELECT cte.user_id FROM cte"
			tt.wantArgs = []interface{}{1, "apple", "banana"}
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
				" (SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $1" +
				" UNION" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $2" +
				" UNION" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $3)" +
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
			q := Union(q1, q2, q3).CTE("cte")
			tt.q = Select(q["user_id"], q["email"]).From(q)
			tt.wantQuery = "WITH cte AS" +
				" (SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $1" +
				" UNION" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $2" +
				" UNION" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $3)" +
				" SELECT cte.user_id, cte.email FROM cte"
			tt.wantArgs = []interface{}{1, 2, 3}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "VariadicQuery CTE (implicit columns from INSERT)"
			u := USERS().As("u")
			q1 := InsertInto(u).Columns(u.USER_ID, u.EMAIL).Values(1, "apple").Returning(u.USER_ID, u.EMAIL)
			q2 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(2))
			q3 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(3))
			q := UnionAll(q1, q2, q3).CTE("cte")
			tt.q = Select(q["user_id"], q["email"]).From(q)
			tt.wantQuery = "WITH cte AS" +
				" (INSERT INTO public.users AS u (user_id, email) VALUES ($1, $2) RETURNING u.user_id, u.email" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $3" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $4)" +
				" SELECT cte.user_id, cte.email FROM cte"
			tt.wantArgs = []interface{}{1, "apple", 2, 3}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "VariadicQuery CTE (implicit columns from UPDATE)"
			u := USERS().As("u")
			q1 := Update(u).Set(u.USER_ID.Set(1), u.EMAIL.Set("apple")).Returning(u.USER_ID, u.EMAIL)
			q2 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(2))
			q3 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(3))
			q := UnionAll(q1, q2, q3).CTE("cte")
			tt.q = Select(q["user_id"], q["email"]).From(q)
			tt.wantQuery = "WITH cte AS" +
				" (UPDATE public.users AS u SET user_id = $1, email = $2 RETURNING u.user_id, u.email" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $3" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $4)" +
				" SELECT cte.user_id, cte.email FROM cte"
			tt.wantArgs = []interface{}{1, "apple", 2, 3}
			return tt
		}(),
		func() TT {
			var tt TT
			tt.description = "VariadicQuery CTE (implicit columns from DELETE)"
			u := USERS().As("u")
			q1 := DeleteFrom(u).Where(u.USER_ID.EqInt(1), u.EMAIL.EqString("apple")).Returning(u.USER_ID, u.EMAIL)
			q2 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(2))
			q3 := Select(u.USER_ID, u.EMAIL).From(u).Where(u.USER_ID.EqInt(3))
			q := UnionAll(q1, q2, q3).CTE("cte")
			tt.q = Select(q["user_id"], q["email"]).From(q)
			tt.wantQuery = "WITH cte AS" +
				" (DELETE FROM public.users AS u WHERE u.user_id = $1 AND u.email = $2 RETURNING u.user_id, u.email" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $3" +
				" UNION ALL" +
				" SELECT u.user_id, u.email FROM public.users AS u WHERE u.user_id = $4)" +
				" SELECT cte.user_id, cte.email FROM cte"
			tt.wantArgs = []interface{}{1, "apple", 2, 3}
			return tt
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.q.AppendSQL(buf, &args, nil)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}
