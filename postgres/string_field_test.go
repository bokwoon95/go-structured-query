package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestStringField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           StringField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "literal string"
			f := String("lorem ipsum")
			wantQuery := "?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			wantQuery := "users.email"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			wantQuery := "u.email"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			exclude := []string{"users"}
			wantQuery := "email"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			exclude := []string{"u"}
			wantQuery := "email"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewStringField("zip code", &TableInfo{Schema: "public", Name: "registered users"})
			wantQuery := `"registered users"."zip code"`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "ASC"
			f := NewStringField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Asc()
			wantQuery := `"registered users"."zip code" ASC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "DESC"
			f := NewStringField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Desc()
			wantQuery := `"registered users"."zip code" DESC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS FIRST"
			f := NewStringField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsFirst()
			wantQuery := `"registered users"."zip code" NULLS FIRST`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS LAST"
			f := NewStringField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsLast()
			wantQuery := `"registered users"."zip code" NULLS LAST`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			var _ Field = tt.f
			tt.f.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestStringField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
	tests := []TT{
		{
			"set field",
			f.Set(f),
			nil,
			"users.email = users.email",
			nil,
		},
		{
			"set string",
			f.Set("lorem ipsum"),
			nil,
			"users.email = ?",
			[]interface{}{"lorem ipsum"},
		},
		{
			"setstring string",
			f.SetString("lorem ipsum"),
			nil,
			"users.email = ?",
			[]interface{}{"lorem ipsum"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.a.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestStringField_Predicates(t *testing.T) {
	type TT struct {
		description string
		p           Predicate
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "IsNull"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.IsNull()
			wantQuery := "users.email IS NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.IsNotNull()
			wantQuery := "users.email IS NOT NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Eq"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.Eq(f)
			wantQuery := "users.email = users.email"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ne"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.Ne(f)
			wantQuery := "users.email <> users.email"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Gt"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.Gt(f)
			wantQuery := "users.email > users.email"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ge"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.Ge(f)
			wantQuery := "users.email >= users.email"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Lt"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.Lt(f)
			wantQuery := "users.email < users.email"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Le"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.Le(f)
			wantQuery := "users.email <= users.email"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "EqString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.EqString("lorem ipsum")
			wantQuery := "users.email = ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "NeString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.NeString("lorem ipsum")
			wantQuery := "users.email <> ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GtString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.GtString("lorem ipsum")
			wantQuery := "users.email > ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GeString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.GeString("lorem ipsum")
			wantQuery := "users.email >= ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LtString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.LtString("lorem ipsum")
			wantQuery := "users.email < ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LeString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.LeString("lorem ipsum")
			wantQuery := "users.email <= ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LikeString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.LikeString("lorem ipsum")
			wantQuery := "users.email LIKE ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "NotLikeString"
			f := NewStringField("email", &TableInfo{Schema: "public", Name: "users"})
			p := f.NotLikeString("lorem ipsum")
			wantQuery := "users.email NOT LIKE ?"
			wantArgs := []interface{}{"lorem ipsum"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In slice"
			f := Fieldf("users.email")
			p := f.In([]string{"a", "b", "c"})
			wantQuery := "users.email IN (?, ?, ?)"
			wantArgs := []interface{}{"a", "b", "c"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In Fields"
			f := Fieldf("users.email")
			p := f.In(Fields{f, f, f})
			wantQuery := "users.email IN (users.email, users.email, users.email)"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.p.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}
