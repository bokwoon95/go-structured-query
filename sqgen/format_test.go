package sqgen

import (
	"testing"

	"github.com/matryer/is"
)

func TestFormatOutput(t *testing.T) {
	is := is.New(t)

	input := `package tables

func Something(column time.Time) bool {
	return true
}`

	expected := `package tables

import "time"

func Something(column time.Time) bool {
	return true
}
`

	out, err := FormatOutput([]byte(input))
	is.NoErr(err)
	is.Equal(string(out), expected)
}
