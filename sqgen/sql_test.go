package sqgen

import (
	"testing"

	"github.com/matryer/is"
)

func TestSliceToSQL(t *testing.T) {
	type TT struct {
		name   string
		args   []string
		result string
	}
	tests := []TT{
		{
			name:   "empty args returns empty string",
			args:   nil,
			result: "",
		},
		{
			name:   "one arg result",
			args:   []string{""},
			result: "(?)",
		},
		{
			name:   "two arg result",
			args:   []string{"", ""},
			result: "(?, ?)",
		},
		{
			name:   "three arg result",
			args:   []string{"", "", ""},
			result: "(?, ?, ?)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(SliceToSQL(tt.args), tt.result)
		})
	}
}
