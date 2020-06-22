package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestBooleanField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           BooleanField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "literal value"
			f := Bool(true)
			wantQuery := "?"
			wantArgs := []interface{}{true}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewBooleanField("is_active", &TableInfo{Schema: "devlab", Name: "users"})
			wantQuery := "users.is_active"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewBooleanField("is_active", &TableInfo{Schema: "devlab", Name: "users", Alias: "u"})
			wantQuery := "u.is_active"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewBooleanField("is_active", &TableInfo{Schema: "devlab", Name: "users"})
			exclude := []string{"users"}
			wantQuery := "is_active"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewBooleanField("is_active", &TableInfo{Schema: "devlab", Name: "users", Alias: "u"})
			exclude := []string{"u"}
			wantQuery := "is_active"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"})
			wantQuery := "`registered users`.`zip code`"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "ASC"
			f := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"}).Asc()
			wantQuery := "`registered users`.`zip code` ASC"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "DESC"
			f := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"}).Desc()
			wantQuery := "`registered users`.`zip code` DESC"
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
			var _ Predicate = tt.f
			var _ = tt.f.String()
			tt.f.AppendSQLExclude(buf, &args, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestBooleanField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := NewBooleanField("is_active", &TableInfo{Schema: "devlab", Name: "users"})
	tests := []TT{
		{
			"set field",
			f.Set(f),
			nil,
			"users.is_active = users.is_active",
			nil,
		},
		{
			"set bool",
			f.Set(true),
			nil,
			"users.is_active = ?",
			[]interface{}{true},
		},
		{
			"setbool bool",
			f.SetBool(true),
			nil,
			"users.is_active = ?",
			[]interface{}{true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.a.AppendSQLExclude(buf, &args, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestBooleanField_Predicates(t *testing.T) {
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
			p := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"}).IsNull()
			wantQuery := "`registered users`.`zip code` IS NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			p := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"}).IsNotNull()
			wantQuery := "`registered users`.`zip code` IS NOT NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Eq"
			f := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"})
			p := f.Eq(f)
			wantQuery := "`registered users`.`zip code` = `registered users`.`zip code`"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ne"
			f := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"})
			p := f.Ne(f)
			wantQuery := "`registered users`.`zip code` <> `registered users`.`zip code`"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Not"
			f := NewBooleanField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"})
			p := f.Not()
			wantQuery := "NOT `registered users`.`zip code`"
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
