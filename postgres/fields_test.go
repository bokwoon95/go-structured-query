package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestFieldLiteral_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           FieldLiteral
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		{"count", FieldLiteral("COUNT(*)"), nil, "COUNT(*)", nil},
		{"one", FieldLiteral("1"), nil, "1", nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.f.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestFields_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		f           Fields
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS()
	tests := []TT{
		{
			"basic",
			Fields{u.EMAIL, u.DISPLAYNAME, u.PASSWORD},
			nil,
			"users.email, users.displayname, users.password",
			nil,
		},
		{
			"ignores aliases",
			Fields{u.EMAIL.As("e"), u.DISPLAYNAME.As("d"), u.PASSWORD.As("p")},
			nil,
			"users.email, users.displayname, users.password",
			nil,
		},
		{
			"nil fields",
			Fields{u.EMAIL, nil, nil},
			nil,
			"users.email, NULL, NULL",
			nil,
		},
		{
			"excludedTableQualifiers",
			Fields{u.EMAIL, u.DISPLAYNAME, u.PASSWORD},
			[]string{u.GetName(), u.GetAlias()},
			"email, displayname, password",
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.f.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestFields_AppendSQLExcludeWithAlias(t *testing.T) {
	type TT struct {
		description string
		f           Fields
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS()
	tests := []TT{
		{
			"basic",
			Fields{u.EMAIL.As("e"), u.DISPLAYNAME.As("d"), u.PASSWORD.As("p")},
			nil,
			"users.email AS e, users.displayname AS d, users.password AS p",
			nil,
		},
		{
			"nil fields",
			Fields{u.EMAIL.As("e"), nil, nil},
			nil,
			"users.email AS e, NULL, NULL",
			nil,
		},
		{
			"excludedTableQualifiers",
			Fields{u.EMAIL.As("e"), u.DISPLAYNAME.As("d"), u.PASSWORD.As("p")},
			[]string{u.GetName(), u.GetAlias()},
			"email AS e, displayname AS d, password AS p",
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.f.AppendSQLExcludeWithAlias(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestFieldsAssignment_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		set         FieldAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{
			"field assign field",
			u.USER_ID.Set(u.DISPLAYNAME),
			nil,
			"u.user_id = u.displayname",
			nil,
		},
		{
			"field assign value",
			u.USER_ID.Set(1),
			nil,
			"u.user_id = ?",
			[]interface{}{1},
		},
		{
			"nil assign field",
			FieldAssignment{nil, u.DISPLAYNAME},
			nil,
			"NULL = u.displayname",
			nil,
		},
		{
			"field assign nil",
			FieldAssignment{u.USER_ID, nil},
			nil,
			"u.user_id = NULL",
			nil,
		},
		{
			"excludedTableQualifiers",
			u.USER_ID.Set(u.DISPLAYNAME),
			[]string{u.GetAlias(), u.GetName()},
			"user_id = displayname",
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.set.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestFieldsAssignments_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		assignments Assignments
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{
			"empty",
			nil,
			nil,
			"",
			nil,
		},
		{
			"basic",
			Assignments{
				u.USER_ID.Set(u.DISPLAYNAME),
				u.USER_ID.Set(1),
				u.PASSWORD.Set(u.USER_ID),
			},
			nil,
			"u.user_id = u.displayname, u.user_id = ?, u.password = u.user_id",
			[]interface{}{1},
		},
		{
			"excludedTableQualifiers",
			Assignments{
				u.USER_ID.Set(u.DISPLAYNAME),
				u.USER_ID.Set(1),
				u.PASSWORD.Set(u.USER_ID),
			},
			[]string{u.GetAlias(), u.GetName()},
			"user_id = displayname, user_id = ?, password = user_id",
			[]interface{}{1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.assignments.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}
