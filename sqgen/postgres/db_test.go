package postgres

import (
	"testing"

	"github.com/matryer/is"
)

func TestReplacePlaceholders(t *testing.T) {
	tt := []struct {
		name   string
		query  string
		result string
	}{
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

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(replacePlaceholders(tc.query), tc.result)
		})
	}
}
