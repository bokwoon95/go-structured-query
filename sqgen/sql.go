package sqgen

import "strings"

// used with "column IN (values)" queries
// expands []string{"val1", "val2"} to "(?, ?, ?)"
func SliceToSQL(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return "(?" + strings.Repeat(", ?", len(args)-1) + ")"
}
