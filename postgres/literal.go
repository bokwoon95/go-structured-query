package sq

import (
	"strings"
)

type Literal string

var _ SQLExcludeAppender = (*Literal)(nil)

func (l Literal) AppendSQLExclude(buf *strings.Builder, _ *[]interface{}, _ map[string]int, _ []string) {
	buf.WriteString(string(l))
}
