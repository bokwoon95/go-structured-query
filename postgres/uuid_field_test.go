package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestUUIDField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           UUIDField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}

	tests := []TT{
		func() TT {
			desc := "literal uuid"
			f := UUID([16]byte{})
			wantQuery := "?"
			wantArgs := []interface{}{[16]byte{}}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewUUIDField("code", &TableInfo{
				Schema: "devlab",
				Name:   "users",
			})
			wantQuery := "users.code"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewUUIDField("code", &TableInfo{
				Schema: "devlab",
				Name:   "users",
				Alias:  "u",
			})
			wantQuery := "u.code"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewUUIDField("code", &TableInfo{
				Schema: "devlab",
				Name:   "users",
			})
			exclude := []string{"users"}
			wantQuery := "code"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewUUIDField("code", &TableInfo{
				Schema: "devlab",
				Name:   "users",
				Alias:  "u",
			})
			exclude := []string{"u"}
			wantQuery := "code"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewUUIDField("signup code", &TableInfo{
				Schema: "devlab",
				Name:   "app users",
			})
			wantQuery := `"app users"."signup code"`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "ASC"
			f := NewUUIDField("signup code", &TableInfo{
				Schema: "devlab",
				Name:   "app users",
			}).Asc()
			wantQuery := `"app users"."signup code" ASC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "DESC"
			f := NewUUIDField("signup code", &TableInfo{
				Schema: "devlab",
				Name:   "app users",
			}).Desc()
			wantQuery := `"app users"."signup code" DESC`
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

func TestUUIDField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}

	f := NewUUIDField("code", &TableInfo{
		Schema: "public",
		Name:   "users",
	})

	tests := []TT{
		{
			"set field",
			f.Set(f),
			nil,
			"users.code = users.code",
			nil,
		},
		{
			"set uuid",
			f.Set([16]byte{}),
			nil,
			"users.code = ?",
			[]interface{}{[16]byte{}},
		},
		{
			"setuuid uuid",
			f.SetUUID([16]byte{}),
			nil,
			"users.code = ?",
			[]interface{}{[16]byte{}},
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

func TestUUIDField_Predicates(t *testing.T) {
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
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.IsNull()
			wantQuery := "users.code IS NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.IsNotNull()
			wantQuery := "users.code IS NOT NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Eq"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.Eq(f)
			wantQuery := "users.code = users.code"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ne"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.Ne(f)
			wantQuery := "users.code <> users.code"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "EqUUID"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.EqUUID([16]byte{})
			wantQuery := "users.code = ?"
			wantArgs := []interface{}{[16]byte{}}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "NeUUID"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.NeUUID([16]byte{})
			wantQuery := "users.code <> ?"
			wantArgs := []interface{}{[16]byte{}}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In slice"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.In([]interface{}{[16]byte{}, [16]byte{}})
			wantQuery := "users.code IN (?, ?)"
			wantArgs := []interface{}{[16]byte{}, [16]byte{}}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In fields"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name:   "users",
			})
			p := f.In(Fields{f, f, f})
			wantQuery := "users.code IN (users.code, users.code, users.code)"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "In Row Values"
			f := NewUUIDField("code", &TableInfo{
				Schema: "public",
				Name: "users",
			})

			p := f.In(RowValue{1, 2, 3})
			wantQuery := "users.code IN (?, ?, ?)"
			wantArgs := []interface{}{1, 2, 3}
			return TT{desc, p, nil, wantQuery, wantArgs}
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
