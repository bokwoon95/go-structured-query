package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestFunctionInfo_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		f           *FunctionInfo
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS()
	tests := []TT{
		{"nil", nil, "", nil},
		{"empty", &FunctionInfo{}, "()", nil},
		{
			"zero arguments",
			&FunctionInfo{
				Schema:    "shitty schema with spaces",
				Name:      "do_something",
				Arguments: []interface{}{},
			},
			`"shitty schema with spaces".do_something()`,
			nil,
		},
		{
			"one or more arguments",
			&FunctionInfo{
				Schema:    "devlab",
				Name:      "do_something",
				Arguments: []interface{}{u.USER_ID, 1, 2, "red fish", "blue fish"},
			},
			`devlab.do_something(users.user_id, ?, ?, ?, ?)`,
			[]interface{}{1, 2, "red fish", "blue fish"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.f.AppendSQL(buf, &args)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestFunctionInfo_Basic(t *testing.T) {
	is := is.New(t)

	f := Functionf("SUM", 5)
	f.Alias = "alias"
	is.Equal("alias", f.GetAlias())
	is.Equal("SUM", f.GetName())
}
