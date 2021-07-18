package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func Test_appendSQLValue(t *testing.T) {
	type TT struct {
		description             string
		value                   interface{}
		excludedTableQualifiers []string
		wantQuery               string
		wantArgs                []interface{}
	}
	tests := []TT{
		{
			description: "normal slices get slice-expanded",
			value:       []int{1, 2, 3, 4},
			wantQuery:   "?, ?, ?, ?",
			wantArgs:    []interface{}{1, 2, 3, 4},
		},
		{
			description: "byte slices don't get slice-expanded",
			value:       []byte{1, 2, 3, 4},
			wantQuery:   "?",
			wantArgs:    []interface{}{[]byte{1, 2, 3, 4}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var gotArgs []interface{}
			appendSQLValue(buf, &gotArgs, tt.excludedTableQualifiers, tt.value)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}
