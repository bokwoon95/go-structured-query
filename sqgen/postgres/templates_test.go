package postgres

import (
	"testing"
	"strings"
	"github.com/matryer/is"
)

func TestTablesTemplate(t *testing.T) {
	is := is.New(t)

	template, err := getTablesTemplate()
	is.NoErr(err)

	var writer strings.Builder

	data := TablesTemplateData{
		PackageName: "tables",
		Imports: []string{
			`sq "github.com/bokwoon95/go-structured-query"`,
		},
		Tables: nil,
	}

	err = template.Execute(&writer, data)
	is.NoErr(err)
}
