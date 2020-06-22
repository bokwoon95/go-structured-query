package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestCustomField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           CustomField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "nested"
			f := CustomField{
				Format: "? = ?",
				Values: []interface{}{Fieldf("MAX(?, ?)", 67, Fieldf("ABS(?)", -88)), 5},
			}
			wantQuery := "MAX(?, ABS(?)) = ?"
			wantArgs := []interface{}{67, -88, 5}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "Asc"
			f := Fieldf("the quick brown fox").Asc()
			wantQuery := "the quick brown fox ASC"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Desc"
			f := Fieldf("the quick brown fox").Desc()
			wantQuery := "the quick brown fox DESC"
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
			tt.f.AppendSQLExclude(buf, &args, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestCustomField_Predicates(t *testing.T) {
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
			p := Fieldf("users.user_id").IsNull()
			wantQuery := "users.user_id IS NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			p := Fieldf("users.user_id").IsNotNull()
			wantQuery := "users.user_id IS NOT NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Eq"
			f := Fieldf("users.user_id")
			p := f.Eq(f)
			wantQuery := "users.user_id = users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ne"
			f := Fieldf("users.user_id")
			p := f.Ne(f)
			wantQuery := "users.user_id <> users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Gt"
			f := Fieldf("users.user_id")
			p := f.Gt(f)
			wantQuery := "users.user_id > users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ge"
			f := Fieldf("users.user_id")
			p := f.Ge(f)
			wantQuery := "users.user_id >= users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Lt"
			f := Fieldf("users.user_id")
			p := f.Lt(f)
			wantQuery := "users.user_id < users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Le"
			f := Fieldf("users.user_id")
			p := f.Le(f)
			wantQuery := "users.user_id <= users.user_id"
			return TT{desc, p, nil, wantQuery, nil}
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
			tt.p.AppendSQLExclude(buf, &args, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestCustomField_In(t *testing.T) {
	type TT struct {
		description string
		p           Predicate
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := Fieldf("id")
	tests := []TT{
		{
			"IN RowValue",
			f.In(RowValue{1, 2, 3}),
			nil,
			"id IN (?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
		{
			"IN Fields",
			f.In(RowValue{f, f, f}),
			nil,
			"id IN (id, id, id)",
			nil,
		},
		{
			"IN slice",
			f.In([]int{1, 2, 3}),
			nil,
			"id IN (?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
		{
			"IN subquery",
			f.In(Select(Int(1), Int(2), Int(3))),
			nil,
			"id IN (SELECT ?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.p.AppendSQLExclude(buf, &args, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestCustomField_BasicTesting(t *testing.T) {
	is := is.New(t)
	var f CustomField

	f = Fieldf("ABC, easy as ?, ?, ?", 1, 2, "2 ep 2").As("gaben")
	// GetName
	is.Equal("ABC, easy as ?, ?, ?", f.GetName())
	// GetAlias
	is.Equal("gaben", f.GetAlias())
	// String
	is.Equal("ABC, easy as 1, 2, '2 ep 2'", f.String())
}
