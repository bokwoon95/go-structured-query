package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestBinaryField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           BinaryField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "literal value"
			f := Bytes([]byte("hello world!"))
			wantQuery := "?"
			wantArgs := []interface{}{[]byte("hello world!")}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewBinaryField("data", &TableInfo{Schema: "devlab", Name: "users"})
			wantQuery := "users.data"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewBinaryField("data", &TableInfo{Schema: "devlab", Name: "users", Alias: "u"})
			wantQuery := "u.data"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewBinaryField("data", &TableInfo{Schema: "devlab", Name: "users"})
			exclude := []string{"users"}
			wantQuery := "data"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewBinaryField("data", &TableInfo{Schema: "devlab", Name: "users", Alias: "u"})
			exclude := []string{"u"}
			wantQuery := "data"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewBinaryField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"})
			wantQuery := "`registered users`.`zip code`"
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

func TestBinaryField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := NewBinaryField("data", &TableInfo{Schema: "devlab", Name: "users"})
	tests := []TT{
		{
			"set field",
			f.Set(f),
			nil,
			"users.data = users.data",
			nil,
		},
		{
			"set bytes",
			f.Set([]byte("hello world!")),
			nil,
			"users.data = ?",
			[]interface{}{[]byte("hello world!")},
		},
		{
			"setbytes bytes",
			f.SetBytes([]byte("hello world!")),
			nil,
			"users.data = ?",
			[]interface{}{[]byte("hello world!")},
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

func TestBinaryField_Predicates(t *testing.T) {
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
			p := NewBinaryField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"}).IsNull()
			wantQuery := "`registered users`.`zip code` IS NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			p := NewBinaryField("zip code", &TableInfo{Schema: "devlab", Name: "registered users"}).IsNotNull()
			wantQuery := "`registered users`.`zip code` IS NOT NULL"
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
