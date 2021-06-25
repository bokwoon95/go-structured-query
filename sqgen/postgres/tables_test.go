package postgres

import (
	"testing"
	"github.com/matryer/is"
)

func TestBuildTablesQuery(t *testing.T) {
	t.Run("single schema, no excluded tables", func(t *testing.T) {
		is := is.New(t)

		schemas := []string{"public"}
		exclude := []string{}

		query, args := buildTablesQuery(schemas, exclude)

		expectedQuery := "SELECT t.table_type, c.table_schema, c.table_name, c.column_name, c.data_type FROM information_schema.tables AS t JOIN information_schema.columns AS c USING (table_schema, table_name) WHERE table_schema IN ($1) ORDER BY c.table_schema <> 'public', c.table_schema, t.table_type, c.table_name, c.column_name"
		expectedArgs := []interface{}{"public"}

		is.Equal(query, expectedQuery)
		is.Equal(args, expectedArgs)
	})

	t.Run("multiple schemas, no excluded tables", func (t *testing.T) {
		is := is.New(t)

		schemas := []string{"public", "geo"}
		exclude := []string{}

		query, args := buildTablesQuery(schemas, exclude)

		expectedQuery := "SELECT t.table_type, c.table_schema, c.table_name, c.column_name, c.data_type FROM information_schema.tables AS t JOIN information_schema.columns AS c USING (table_schema, table_name) WHERE table_schema IN ($1, $2) ORDER BY c.table_schema <> 'public', c.table_schema, t.table_type, c.table_name, c.column_name"
		expectedArgs := []interface{}{"public", "geo"}

		is.Equal(query, expectedQuery)
		is.Equal(args, expectedArgs)
	})

	t.Run("multiple schemas, excluded tables", func(t *testing.T) {
		is := is.New(t)
		schemas := []string{"public", "geo"}
		exclude := []string{"schema_migrations", "meta"}

		query, args := buildTablesQuery(schemas, exclude)

		expectedQuery := "SELECT t.table_type, c.table_schema, c.table_name, c.column_name, c.data_type FROM information_schema.tables AS t JOIN information_schema.columns AS c USING (table_schema, table_name) WHERE table_schema IN ($1, $2) AND table_name NOT IN ($3, $4) ORDER BY c.table_schema <> 'public', c.table_schema, t.table_type, c.table_name, c.column_name"

		expectedArgs := []interface{}{"public", "geo", "schema_migrations", "meta"}

		is.Equal(query, expectedQuery)
		is.Equal(args, expectedArgs)
	})
}
