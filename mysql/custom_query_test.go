package sq

import (
	"testing"

	"github.com/matryer/is"
)

func TestCustomQuery_BasicTesting(t *testing.T) {
	is := is.New(t)
	q := Queryf("ABC, easy as ?, ?, ?", 1, 2, "2 ep 2").As("gaben")
	// ToSQL
	query, args := q.ToSQL()
	is.Equal("ABC, easy as ?, ?, ?", query)
	is.Equal([]interface{}{1, 2, "2 ep 2"}, args)
	// GetAlias
	is.Equal("gaben", q.GetAlias())
	// GetName
	is.Equal("ABC, easy as ?, ?, ?", q.GetName())
}
