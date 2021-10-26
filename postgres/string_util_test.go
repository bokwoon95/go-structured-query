package sq

import (
	"database/sql/driver"
	"strings"
	"testing"

	"github.com/matryer/is"
)

type testEnum string
type testEnumArray []testEnum

// Value implements database/sql/driver:Valuer interface
func (src testEnumArray) Value() (driver.Value, error) {
	return src, nil
}

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
		{
			description: "slices implementing driver.Value don't get slice-expanded",
			value:       testEnumArray{testEnum("A"), testEnum("B"), testEnum("C")},
			wantQuery:   "?",
			wantArgs:    []interface{}{testEnumArray{testEnum("A"), testEnum("B"), testEnum("C")}},
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
