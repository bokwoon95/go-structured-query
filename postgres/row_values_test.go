package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestRowValues_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		vl          RowValues
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		{
			"0",
			RowValues{},
			"",
			nil,
		},
		{
			"1x0",
			RowValues{
				RowValue{},
			},
			"()",
			nil,
		},
		{
			"1x1",
			RowValues{
				RowValue{1},
			},
			"(?)",
			[]interface{}{1},
		},
		{
			"1x3",
			RowValues{
				RowValue{1, 2, 3},
			},
			"(?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
		{
			"3x3",
			RowValues{
				RowValue{1, 2, 3},
				RowValue{1, 2, 3},
				RowValue{1, 2, 3},
			},
			"(?, ?, ?), (?, ?, ?), (?, ?, ?)",
			[]interface{}{1, 2, 3, 1, 2, 3, 1, 2, 3},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.vl.AppendSQL(buf, &args, nil)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestRowValue_In(t *testing.T) {
	type TT struct {
		description string
		p           Predicate
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{
			"empty",
			RowValue{}.In(nil),
			nil,
			"() IN (NULL)",
			nil,
		},
		{
			"IN RowValue",
			RowValue{u.USER_ID}.In(RowValue{1, 2, 3}),
			nil,
			"(u.user_id) IN (?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
		{
			"IN RowValues",
			RowValue{u.USER_ID, u.EMAIL, u.DISPLAYNAME}.In(RowValues{
				RowValue{1, "a@email.com", "a"},
				RowValue{2, "b@email.com", "b"},
				RowValue{3, "c@email.com", "c"},
			}),
			nil,
			"(u.user_id, u.email, u.displayname) IN ((?, ?, ?), (?, ?, ?), (?, ?, ?))",
			[]interface{}{1, "a@email.com", "a", 2, "b@email.com", "b", 3, "c@email.com", "c"},
		},
		{
			"IN Query",
			RowValue{u.USER_ID, u.EMAIL, u.DISPLAYNAME}.In(Select(Int(1), Int(2), Int(3))),
			nil,
			"(u.user_id, u.email, u.displayname) IN (SELECT ?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
		{
			"excludedTableQualifiers",
			RowValue{u.USER_ID, u.EMAIL, u.DISPLAYNAME}.In(RowValues{
				RowValue{1, "a@email.com", "a"},
				RowValue{2, "b@email.com", "b"},
				RowValue{3, "c@email.com", "c"},
			}),
			[]string{u.GetAlias(), u.GetName()},
			"(user_id, email, displayname) IN ((?, ?, ?), (?, ?, ?), (?, ?, ?))",
			[]interface{}{1, "a@email.com", "a", 2, "b@email.com", "b", 3, "c@email.com", "c"},
		},
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

func TestCustomAssignment_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		a           CustomAssignment
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{
			"Set RowValue",
			RowValue{u.USER_ID, u.DISPLAYNAME, u.EMAIL}.Set(RowValue{1, "bob", "bob@email.com"}),
			nil,
			"(u.user_id, u.displayname, u.email) = (?, ?, ?)",
			[]interface{}{1, "bob", "bob@email.com"},
		},
		{
			"Set Query",
			RowValue{u.USER_ID, u.DISPLAYNAME, u.EMAIL}.Set(Select(Int(1), String("bob"), String("bob@email.com"))),
			nil,
			"(u.user_id, u.displayname, u.email) = (SELECT ?, ?, ?)",
			[]interface{}{1, "bob", "bob@email.com"},
		},
		{
			"Set value",
			RowValue{u.USER_ID, u.DISPLAYNAME, u.EMAIL}.Set(String("lorem ipsum dolor sit amet")),
			nil,
			"(u.user_id, u.displayname, u.email) = (?)",
			[]interface{}{"lorem ipsum dolor sit amet"},
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
