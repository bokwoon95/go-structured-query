package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestJoinTable_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		j           JoinTable
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "empty"
			j := CustomJoin("", nil)
			wantQuery := "JOIN NULL"
			return TT{desc, j, wantQuery, nil}
		}(),
		func() TT {
			desc := "join table"
			u := USERS()
			j := CustomJoin(JoinTypeLeft, u, u.USER_ID.EqInt(1), u.DISPLAYNAME.EqString("John"))
			wantQuery := "LEFT JOIN public.users ON users.user_id = ? AND users.displayname = ?"
			wantArgs := []interface{}{1, "John"}
			return TT{desc, j, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "custom join table with alias"
			u := USERS().As("u")
			j := CustomJoin("LEFT JOIN LATERAL", u, u.USER_ID.EqInt(1), u.DISPLAYNAME.EqString("John"))
			wantQuery := "LEFT JOIN LATERAL public.users AS u ON u.user_id = ? AND u.displayname = ?"
			wantArgs := []interface{}{1, "John"}
			return TT{desc, j, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "join query"
			u := USERS().As("u")
			q := Select(u.USER_ID, u.DISPLAYNAME, u.EMAIL).From(u).Subquery("subquery")
			j := CustomJoin(JoinTypeInner, q, q["user_id"].Eq(1), q["displayname"].Eq("John"))
			wantQuery := "JOIN (" +
				"SELECT u.user_id, u.displayname, u.email FROM public.users AS u" +
				") AS subquery ON subquery.user_id = ? AND subquery.displayname = ?"
			wantArgs := []interface{}{1, "John"}
			return TT{desc, j, wantQuery, wantArgs}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.j.AppendSQL(buf, &args)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestJoinTables_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		j           JoinTables
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "empty"
			j := JoinTables{}
			wantQuery := ""
			return TT{desc, j, wantQuery, nil}
		}(),
		func() TT {
			desc := "basic"
			u := USERS().As("u")
			j := JoinTables{
				CustomJoin(JoinTypeLeft, u, u.USER_ID.EqInt(1), u.DISPLAYNAME.EqString("John")),
				CustomJoin(JoinTypeRight, u, u.DISPLAYNAME.EqString("Jane"), u.USER_ID.EqInt(2)),
				CustomJoin(JoinTypeFull, u),
			}
			wantQuery := "LEFT JOIN public.users AS u ON u.user_id = ? AND u.displayname = ?" +
				" RIGHT JOIN public.users AS u ON u.displayname = ? AND u.user_id = ?" +
				" FULL JOIN public.users AS u"
			wantArgs := []interface{}{1, "John", "Jane", 2}
			return TT{desc, j, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "more joins"
			u := USERS().As("u")
			j := JoinTables{
				Join(u, u.USER_ID.EqInt(1), u.DISPLAYNAME.EqString("John")),
				LeftJoin(u, u.DISPLAYNAME.EqString("Jane"), u.USER_ID.EqInt(2)),
				RightJoin(u, u.DISPLAYNAME.EqString("Street"), u.USER_ID.EqInt(3)),
				FullJoin(u),
			}
			wantQuery := "JOIN public.users AS u ON u.user_id = ? AND u.displayname = ?" +
				" LEFT JOIN public.users AS u ON u.displayname = ? AND u.user_id = ?" +
				" RIGHT JOIN public.users AS u ON u.displayname = ? AND u.user_id = ?" +
				" FULL JOIN public.users AS u"
			wantArgs := []interface{}{1, "John", "Jane", 2, "Street", 3}
			return TT{desc, j, wantQuery, wantArgs}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.j.AppendSQL(buf, &args)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestJoinTable_Basic(t *testing.T) {
}
