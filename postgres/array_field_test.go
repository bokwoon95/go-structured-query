package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestArrayField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           ArrayField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "empty []bool literal"
			field := Array([]bool{})
			wantQuery := "ARRAY[]::BOOLEAN[]"
			return TT{desc, field, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "empty []float64 literal"
			field := Array([]float64{})
			wantQuery := "ARRAY[]::FLOAT[]"
			return TT{desc, field, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "empty []int literal"
			field := Array([]int{})
			wantQuery := "ARRAY[]::INT[]"
			return TT{desc, field, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "empty []int64 literal"
			field := Array([]int64{})
			wantQuery := "ARRAY[]::BIGINT[]"
			return TT{desc, field, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "empty []string literal"
			field := Array([]string{})
			wantQuery := "ARRAY[]::TEXT[]"
			return TT{desc, field, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "[]bool literal"
			field := Array([]bool{true, true, false, true})
			wantQuery := "ARRAY[?, ?, ?, ?]"
			wantArgs := []interface{}{true, true, false, true}
			return TT{desc, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "[]float64 literal"
			field := Array([]float64{22.7, 3.15, 4.0})
			wantQuery := "ARRAY[?, ?, ?]"
			wantArgs := []interface{}{22.7, 3.15, 4.0}
			return TT{desc, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "[]int literal"
			field := Array([]int{1, 2, 3, 4})
			wantQuery := "ARRAY[?, ?, ?, ?]"
			wantArgs := []interface{}{1, 2, 3, 4}
			return TT{desc, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "[]int64 literal"
			field := Array([]int64{1, 2, 3, 4})
			wantQuery := "ARRAY[?, ?, ?, ?]"
			wantArgs := []interface{}{int64(1), int64(2), int64(3), int64(4)}
			return TT{desc, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "[]string literal"
			field := Array([]string{"apple", "banana", "cucumber"})
			wantQuery := "ARRAY[?, ?, ?]"
			wantArgs := []interface{}{"apple", "banana", "cucumber"}
			return TT{desc, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users"})
			wantQuery := "users.user_list"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			wantQuery := "u.user_list"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users"})
			exclude := []string{"users"}
			wantQuery := "user_list"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			exclude := []string{"u"}
			wantQuery := "user_list"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"})
			wantQuery := `"registered users"."zip code"`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "ASC"
			f := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Asc()
			wantQuery := `"registered users"."zip code" ASC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "DESC"
			f := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Desc()
			wantQuery := `"registered users"."zip code" DESC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS FIRST"
			f := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsFirst()
			wantQuery := `"registered users"."zip code" NULLS FIRST`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS LAST"
			f := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsLast()
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

func TestArrayField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users"})
	tests := []TT{
		{
			"set field",
			f.Set(Array([]string{"tom", "dick", "harry"})),
			nil,
			"users.user_list = ARRAY[?, ?, ?]",
			[]interface{}{"tom", "dick", "harry"},
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

func TestArrayField_Predicates(t *testing.T) {
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
			p := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).IsNull()
			wantQuery := `"registered users"."zip code" IS NULL`
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			p := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).IsNotNull()
			wantQuery := `"registered users"."zip code" IS NOT NULL`
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Eq"
			f := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"})
			p := f.Eq(f)
			wantQuery := `"registered users"."zip code" = "registered users"."zip code"`
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ne"
			f := NewArrayField("zip code", &TableInfo{Schema: "public", Name: "registered users"})
			p := f.Ne(f)
			wantQuery := `"registered users"."zip code" <> "registered users"."zip code"`
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Gt"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users"})
			p := f.Gt(f)
			wantQuery := "users.user_list > users.user_list"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ge"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users"})
			p := f.Ge(f)
			wantQuery := "users.user_list >= users.user_list"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Lt"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users"})
			p := f.Lt(f)
			wantQuery := "users.user_list < users.user_list"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Le"
			f := NewArrayField("user_list", &TableInfo{Schema: "public", Name: "users"})
			p := f.Le(f)
			wantQuery := "users.user_list <= users.user_list"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Contains"
			p := Array([]int{1, 2, 3}).Contains(Array([]int{2, 3}))
			wantQuery := "ARRAY[?, ?, ?] @> ARRAY[?, ?]"
			wantArgs := []interface{}{1, 2, 3, 2, 3}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "ContainedBy"
			p := Array([]int{1, 2, 3}).ContainedBy(Array([]int{2, 3}))
			wantQuery := "ARRAY[?, ?, ?] <@ ARRAY[?, ?]"
			wantArgs := []interface{}{1, 2, 3, 2, 3}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "Overlaps"
			p := Array([]int{1, 2, 3}).Overlaps(Array([]int{2, 3}))
			wantQuery := "ARRAY[?, ?, ?] && ARRAY[?, ?]"
			wantArgs := []interface{}{1, 2, 3, 2, 3}
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
