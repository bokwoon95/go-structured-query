package sqgen

import (
	"testing"
	"github.com/matryer/is"
)

func TestSliceToSQL(t *testing.T) {
	tt := []struct{
		name string
		args []string
		result string
	}{
		{
			name: "empty args returns empty string",
			args: nil,
			result: "",
		},
		{
			name: "one arg result",
			args: []string{""},
			result: "(?)",
		},
		{
			name: "two arg result",
			args: []string{"", ""},
			result: "(?, ?)",
		},
		{
			name: "three arg result",
			args: []string{"", "", ""},
			result: "(?, ?, ?)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(SliceToSQL(tc.args), tc.result)
		})
	}
}
