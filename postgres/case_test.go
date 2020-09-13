package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestPredicateCases_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           PredicateCases
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{
			"empty",
			PredicateCases{},
			nil,
			"CASE END",
			nil,
		},
		{
			"nil",
			CaseWhen(nil, nil),
			nil,
			"CASE WHEN NULL THEN NULL END",
			nil,
		},
		{
			"basic",
			CaseWhen(u.USER_ID.EqInt(1), Int(1)).
				When(u.EMAIL.GtString("lorem ipsum"), String("lorem ipsum")).
				When(u.DISPLAYNAME.Eq(u.EMAIL), u.USER_ID).
				Else(Float64(99.99)),
			nil,
			"CASE" +
				" WHEN u.user_id = ? THEN ?" +
				" WHEN u.email > ? THEN ?" +
				" WHEN u.displayname = u.email THEN u.user_id" +
				" ELSE ?" +
				" END",
			[]interface{}{1, 1, "lorem ipsum", "lorem ipsum", 99.99},
		},
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

func TestPredicateCases_Basic(t *testing.T) {
	is := is.New(t)

	p := CaseWhen(nil, nil).As("test")
	is.Equal("test", p.GetAlias())
	is.Equal("", p.GetName())
}

func TestSimpleCases_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           SimpleCases
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{
			"empty",
			SimpleCases{},
			nil,
			"CASE NULL END",
			nil,
		},
		{
			"nil",
			Case(nil).When(nil, nil),
			nil,
			"CASE NULL WHEN NULL THEN NULL END",
			nil,
		},
		{
			"basic",
			Case(u.PASSWORD).When(u.USER_ID, Int(1)).
				When(u.EMAIL, String("lorem ipsum")).
				When(u.DISPLAYNAME, u.USER_ID).
				Else(Float64(99.99)),
			nil,
			"CASE u.password" +
				" WHEN u.user_id THEN ?" +
				" WHEN u.email THEN ?" +
				" WHEN u.displayname THEN u.user_id" +
				" ELSE ?" +
				" END",
			[]interface{}{1, "lorem ipsum", 99.99},
		},
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

func TestSimpleCases_Basic(t *testing.T) {
	is := is.New(t)

	p := Case(nil).When(nil, nil).As("test")
	is.Equal("test", p.GetAlias())
	is.Equal("", p.GetName())
}
