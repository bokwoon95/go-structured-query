package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func openAndPing(database string) (*sql.DB, error) {
	db, err := sql.Open("postgres", database)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("Could not ping the database, is the database reachable via %s? %w", database, err)
	}

	return db, nil
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

// used with "column IN (values)" queries
// expands []string{"val1", "val2"} to "(?, ?, ?)"
// should not accept nil slice
func sliceToSQL(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return "(?" + strings.Repeat(", ?", len(args)-1) + ")"
}
