package sq

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestTimeField_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           TimeField
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	now := time.Now()
	tests := []TT{
		func() TT {
			desc := "literal time.Time"
			f := Time(now)
			wantQuery := "?"
			wantArgs := []interface{}{now}
			return TT{desc, f, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "table qualified"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			wantQuery := "users.created_at"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "table alias qualified"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			wantQuery := "u.created_at"
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (name)"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			exclude := []string{"users"}
			wantQuery := "created_at"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "excludedTableQualifiers (alias)"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users", Alias: "u"})
			exclude := []string{"u"}
			wantQuery := "created_at"
			return TT{desc, f, exclude, wantQuery, nil}
		}(),
		func() TT {
			desc := "quoted whitespace"
			f := NewTimeField("zip code", &TableInfo{Schema: "public", Name: "registered users"})
			wantQuery := `"registered users"."zip code"`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "ASC"
			f := NewTimeField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Asc()
			wantQuery := `"registered users"."zip code" ASC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "DESC"
			f := NewTimeField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).Desc()
			wantQuery := `"registered users"."zip code" DESC`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS FIRST"
			f := NewTimeField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsFirst()
			wantQuery := `"registered users"."zip code" NULLS FIRST`
			return TT{desc, f, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NULLS LAST"
			f := NewTimeField("zip code", &TableInfo{Schema: "public", Name: "registered users"}).NullsLast()
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

func TestTimeField_FieldAssignment(t *testing.T) {
	type TT struct {
		description string
		a           FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
	now := time.Now()
	tests := []TT{
		{
			"set field",
			f.Set(f),
			nil,
			"users.created_at = users.created_at",
			nil,
		},
		{
			"set string",
			f.Set("lorem ipsum"),
			nil,
			"users.created_at = ?",
			[]interface{}{"lorem ipsum"},
		},
		{
			"settime time",
			f.SetTime(now),
			nil,
			"users.created_at = ?",
			[]interface{}{now},
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

func TestTimeField_Predicates(t *testing.T) {
	type TT struct {
		description string
		p           Predicate
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	now := time.Now()
	tests := []TT{
		func() TT {
			desc := "IsNull"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.IsNull()
			wantQuery := "users.created_at IS NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "IsNotNull"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.IsNotNull()
			wantQuery := "users.created_at IS NOT NULL"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Eq"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.Eq(f)
			wantQuery := "users.created_at = users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ne"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.Ne(f)
			wantQuery := "users.created_at <> users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Gt"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.Gt(f)
			wantQuery := "users.created_at > users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Ge"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.Ge(f)
			wantQuery := "users.created_at >= users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Lt"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.Lt(f)
			wantQuery := "users.created_at < users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Le"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.Le(f)
			wantQuery := "users.created_at <= users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "Between"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.Between(f, f)
			wantQuery := "users.created_at BETWEEN users.created_at AND users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "NotBetween"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.NotBetween(f, f)
			wantQuery := "users.created_at NOT BETWEEN users.created_at AND users.created_at"
			return TT{desc, p, nil, wantQuery, nil}
		}(),
		func() TT {
			desc := "EqTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.EqTime(now)
			wantQuery := "users.created_at = ?"
			wantArgs := []interface{}{now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "NeTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.NeTime(now)
			wantQuery := "users.created_at <> ?"
			wantArgs := []interface{}{now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GtTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.GtTime(now)
			wantQuery := "users.created_at > ?"
			wantArgs := []interface{}{now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "GeTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.GeTime(now)
			wantQuery := "users.created_at >= ?"
			wantArgs := []interface{}{now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LtTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.LtTime(now)
			wantQuery := "users.created_at < ?"
			wantArgs := []interface{}{now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "LeTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.LeTime(now)
			wantQuery := "users.created_at <= ?"
			wantArgs := []interface{}{now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "BetweenTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.BetweenTime(now, now)
			wantQuery := "users.created_at BETWEEN ? AND ?"
			wantArgs := []interface{}{now, now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "NotBetweenTime"
			f := NewTimeField("created_at", &TableInfo{Schema: "public", Name: "users"})
			p := f.NotBetweenTime(now, now)
			wantQuery := "users.created_at NOT BETWEEN ? AND ?"
			wantArgs := []interface{}{now, now}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In slice"
			f := Fieldf("users.created_at")
			p := f.In([]string{"a", "b", "c"})
			wantQuery := "users.created_at IN (?, ?, ?)"
			wantArgs := []interface{}{"a", "b", "c"}
			return TT{desc, p, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			desc := "In Fields"
			f := Fieldf("users.created_at")
			p := f.In(Fields{f, f, f})
			wantQuery := "users.created_at IN (users.created_at, users.created_at, users.created_at)"
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
