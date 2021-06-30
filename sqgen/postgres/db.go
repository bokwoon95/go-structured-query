package postgres

import (
	"strconv"
	"strings"
)

// replacePlaceholders will replace question mark placeholders with dollar
// placeholders e.g. ?, ?, ? -> $1, $2, $3 etc
func replacePlaceholders(query string) string {
	var buf strings.Builder
	var i int
	for pos := strings.Index(query, "?"); pos >= 0; pos = strings.Index(query, "?") {
		i++
		buf.WriteString(query[:pos] + "$" + strconv.Itoa(i))
		query = query[pos+1:]
	}
	buf.WriteString(query)
	return buf.String()
}
