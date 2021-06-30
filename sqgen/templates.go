package sqgen

import (
	"fmt"
	"strings"
	"text/template"
)

// functions required to transforms strings within the template
// removes the need to declare a custom string alias with these methods attached
// since all the template variables are referenced either with $ or . accessor in the template
// we won't have any naming collisions

func Export(s string) string {
	str := strings.TrimPrefix(s, "_")
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ToUpper(str)
	return str
}

func QuoteSpace(s string) string {
	if strings.Contains(s, " ") {
		return fmt.Sprintf(`"%s"`, s)
	}

	return s
}

var FuncMap template.FuncMap = map[string]interface{}{
	"export":     Export,
	"quoteSpace": QuoteSpace,
}
