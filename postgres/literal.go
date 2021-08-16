package sq

import (
	"strings"
)

// Literal allows for the underlying string to be literally plugged into the SQL query
type Literal string

var _ SQLExcludeAppender = (*Literal)(nil)

func (l Literal) AppendSQLExclude(buf *strings.Builder, _ *[]interface{}, _ map[string]int, _ []string) {
	buf.WriteString(string(l))
}
