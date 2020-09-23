package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestTableInfo_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		t           *TableInfo
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		{"empty", nil, "", nil},
		{"has schema", &TableInfo{Schema: "public", Name: "users"}, "public.users", nil},
		{"no schema", &TableInfo{Name: "users"}, "users", nil},
		{
			// https://stackoverflow.com/q/506826
			// only villians put whitespaces in their schema/table/column names >.>
			"quoted whitespace",
			&TableInfo{Schema: "student registration", Name: "table with whitespace"},
			"`student registration`.`table with whitespace`",
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
			tt.t.AppendSQL(buf, &args, nil)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}
