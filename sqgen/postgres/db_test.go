package postgres

import (
	"testing"

	"github.com/matryer/is"
)

func TestReplacePlaceholders(t *testing.T) {
	type TT struct {
		name   string
		query  string
		result string
	}
	tests := []TT{
		{
			name:   "string without placeholders is untouched",
			query:  "select * from public.some_table where status = 200",
			result: "select * from public.some_table where status = 200",
		},
		{
			name:   "can replace 1 placeholder",
			query:  "select * from public.some_table where status = ?",
			result: "select * from public.some_table where status = $1",
		},
		{
			name:   "can replace 2 placeholder",
			query:  "select * from public.some_table where status = ? and another_field ilike ?",
			result: "select * from public.some_table where status = $1 and another_field ilike $2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(replacePlaceholders(tt.query), tt.result)
		})
	}
}
