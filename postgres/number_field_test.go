package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestNumberField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           NumberField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "literal int"
			f := Int(1)
			wantQuery := "?"
			wantArgs := []interface{}{1}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "literal int64"
			f := Int64(1)
			wantQuery := "?"
			wantArgs := []interface{}{int64(1)}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "literal float64"
			f := Float64(33.27)
			wantQuery := "?"
			wantArgs := []interface{}{33.27}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			wantQuery := "users.user_id"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			wantQuery := "u.user_id"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			exclude := []string{"users"}
			wantQuery := "user_id"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			exclude := []string{"u"}
			wantQuery := "user_id"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewNumberField("zip code", &TableInfo{Schema: "public", Name: "registered users"})
			wantQuery := `"registered users"."zip code"`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "ASC"
			f := NewNumberField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Asc()
			wantQuery := `"registered users"."zip code" ASC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "DESC"
			f := NewNumberField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Desc()
			wantQuery := `"registered users"."zip code" DESC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS FIRST"
			f := NewNumberField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsFirst()
			wantQuery := `"registered users"."zip code" NULLS FIRST`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS LAST"
			f := NewNumberField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsLast()
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

func TestNumberField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
	tests := []TT{
		{
			"set field",
			f.Set(f),
			nil,
			"users.user_id = users.user_id",
			nil,
		},
		{
			"set int",
			f.Set(1),
			nil,
			"users.user_id = ?",
			[]interface{}{1},
		},
		{
			"setint int",
			f.SetInt(1),
			nil,
			"users.user_id = ?",
			[]interface{}{1},
		},
		{
			"setint64 int64",
			f.SetInt64(1),
			nil,
			"users.user_id = ?",
			[]interface{}{int64(1)},
		},
		{
			"setfloat64 float64",
			f.SetFloat64(33.27),
			nil,
			"users.user_id = ?",
			[]interface{}{33.27},
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

func TestNumberField_Predicates(t *testing.T) {
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
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.IsNull()
			wantQuery := "users.user_id IS NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.IsNotNull()
			wantQuery := "users.user_id IS NOT NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Eq"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.Eq(f)
			wantQuery := "users.user_id = users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ne"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.Ne(f)
			wantQuery := "users.user_id <> users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Gt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.Gt(f)
			wantQuery := "users.user_id > users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ge"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.Ge(f)
			wantQuery := "users.user_id >= users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Lt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.Lt(f)
			wantQuery := "users.user_id < users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Le"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.Le(f)
			wantQuery := "users.user_id <= users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "EqInt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.EqInt(1)
			wantQuery := "users.user_id = ?"
			wantArgs := []interface{}{1}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "NeInt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.NeInt(1)
			wantQuery := "users.user_id <> ?"
			wantArgs := []interface{}{1}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GtInt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.GtInt(1)
			wantQuery := "users.user_id > ?"
			wantArgs := []interface{}{1}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GeInt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.GeInt(1)
			wantQuery := "users.user_id >= ?"
			wantArgs := []interface{}{1}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LtInt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.LtInt(1)
			wantQuery := "users.user_id < ?"
			wantArgs := []interface{}{1}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LeInt"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.LeInt(1)
			wantQuery := "users.user_id <= ?"
			wantArgs := []interface{}{1}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "EqFloat64"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.EqFloat64(33.27)
			wantQuery := "users.user_id = ?"
			wantArgs := []interface{}{33.27}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "NeFloat64"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.NeFloat64(33.27)
			wantQuery := "users.user_id <> ?"
			wantArgs := []interface{}{33.27}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GtFloat64"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.GtFloat64(33.27)
			wantQuery := "users.user_id > ?"
			wantArgs := []interface{}{33.27}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GeFloat64"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.GeFloat64(33.27)
			wantQuery := "users.user_id >= ?"
			wantArgs := []interface{}{33.27}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LtFloat64"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.LtFloat64(33.27)
			wantQuery := "users.user_id < ?"
			wantArgs := []interface{}{33.27}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LeFloat64"
			f := NewNumberField("user_id", &TableInfo{Schema: "public", Name: "users"})
			p := f.LeFloat64(33.27)
			wantQuery := "users.user_id <= ?"
			wantArgs := []interface{}{33.27}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In slice"
			f := Fieldf("users.user_id")
			p := f.In([]int{1, 2, 3})
			wantQuery := "users.user_id IN (?, ?, ?)"
			wantArgs := []interface{}{1, 2, 3}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In Fields"
			f := Fieldf("users.user_id")
			p := f.In(Fields{f, f, f})
			wantQuery := "users.user_id IN (users.user_id, users.user_id, users.user_id)"
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
