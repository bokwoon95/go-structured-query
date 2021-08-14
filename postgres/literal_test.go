package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestLiteral(t *testing.T) {
	is := is.New(t)
	buf := &strings.Builder{}
	field := `'ID'`
	l := Literal(field)
	l.AppendSQLExclude(buf, nil, nil, nil)

	is.Equal(buf.String(), field)
}
