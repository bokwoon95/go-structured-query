package postgres

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

/* Error Handling Utilities */

const recSep rune = 30 // ASCII Record Separator

func wrap(err error) error {
	if err == nil {
		return nil
	}
	_, filename, linenbr, _ := runtime.Caller(1)
	return fmt.Errorf(string(recSep)+" %s:%d %w", filename, linenbr, err)
}

// replacePlaceholders will replace question mark placeholders with dollar
// placeholders e.g. ?, ?, ? -> $1, $2, $3 etc
func replacePlaceholders(query string) string {
	buf := &strings.Builder{}
	var i int
	for pos := strings.Index(query, "?"); pos >= 0; pos = strings.Index(query, "?") {
		i++
		buf.WriteString(query[:pos] + "$" + strconv.Itoa(i))
		query = query[pos+1:]
	}
	buf.WriteString(query)
	return buf.String()
}
