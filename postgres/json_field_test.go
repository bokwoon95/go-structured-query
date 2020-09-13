package sq

import (
	"database/sql/driver"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestJSONField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           JSONField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "literal value (struct)"
			val := struct {
				UserID      int    `json:"user_id"`
				Email       string `json:"email"`
				Displayname string `json:"displayname"`
			}{}
			f := MustJSON(val)
			wantQuery := "?"
			wantArgs := []interface{}{val}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewJSONField("data", &TableInfo{Schema: "public", Name: "users"})
			wantQuery := "users.data"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewJSONField("data", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			wantQuery := "u.data"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewJSONField("data", &TableInfo{Schema: "public", Name: "users"})
			exclude := []string{"users"}
			wantQuery := "data"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewJSONField("data", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			exclude := []string{"u"}
			wantQuery := "data"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewJSONField("zip code", &TableInfo{Schema: "public", Name: "registered users"})
			wantQuery := `"registered users"."zip code"`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "ASC"
			f := NewJSONField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Asc()
			wantQuery := `"registered users"."zip code" ASC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "DESC"
			f := NewJSONField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Desc()
			wantQuery := `"registered users"."zip code" DESC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS FIRST"
			f := NewJSONField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsFirst()
			wantQuery := `"registered users"."zip code" NULLS FIRST`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS LAST"
			f := NewJSONField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsLast()
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
			tt.f.AppendSQLExclude(buf, &args, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestJSONField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := NewJSONField("data", &TableInfo{Schema: "public", Name: "users"})
	val := struct {
		UserID      int    `json:"user_id"`
		Email       string `json:"email"`
		Displayname string `json:"displayname"`
	}{}
	tests := []TT{
		{
			"set field",
			f.Set(f),
			nil,
			"users.data = users.data",
			nil,
		},
		{
			"set json",
			f.Set(val),
			nil,
			"users.data = ?",
			[]interface{}{val},
		},
		{
			"setjson json",
			f.SetJSON(val),
			nil,
			"users.data = ?",
			[]interface{}{val},
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

func TestJSONField_Predicates(t *testing.T) {
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
			p := NewJSONField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).IsNull()
			wantQuery := `"registered users"."zip code" IS NULL`
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			p := NewJSONField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).IsNotNull()
			wantQuery := `"registered users"."zip code" IS NOT NULL`
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

type customValuer string

func (v customValuer) Value() (driver.Value, error) {
	return driver.Value(v), nil
}

func TestJSONField_Basic(t *testing.T) {
	is := is.New(t)

	var x customValuer = "lorem ipsum"
	f := JSONValue(x)
	is.Equal(x, f.value)
	_ = f.SetValue(x)
	is.Equal(`:"lorem ipsum":`, f.String())
	is.Equal("", f.GetName())
}
