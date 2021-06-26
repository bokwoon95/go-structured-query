// contains the logic for the sqgen-postgres tables command
package postgres

import (
	"bytes"
	"database/sql"
	"io"
	"strings"

	"github.com/bokwoon95/go-structured-query/sqgen"
)

type Table struct {
	Schema      string
	Name        string
	StructName  string
	RawType     string
	Constructor string
	Fields      []TableField
}

type TableField struct {
	Name        string
	RawType     string
	Type        string
	Constructor string
}

func BuildTables(config Config, writer io.Writer) error {
	db, err := openAndPing(config.Database)

	if err != nil {
		return sqgen.Wrap(err)
	}

	tables, err := executeTables(config, db)

	if err != nil {
		return sqgen.Wrap(err)
	}

	templateData := TablesTemplateData{
		PackageName: config.Package,
		Imports: []string{
			`sq "github.com/bokwoon95/go-structured-query/postgres"`,
		},
		Tables: tables,
	}

	t, err := getTablesTemplate()

	if err != nil {
		return sqgen.Wrap(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, templateData)

	if err != nil {
		return sqgen.Wrap(err)
	}

	src, err := sqgen.FormatOutput(buf.Bytes())

	if err != nil {
		return sqgen.Wrap(err)
	}

	_, err = writer.Write(src)
	return err
}

func executeTables(config Config, db *sql.DB) ([]Table, error) {
	// Prepare the query and args
	query, args := buildTablesQuery(config.Schemas, config.Exclude)
	// Query the database and aggregate the results into a []Table slice
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, sqgen.Wrap(err)
	}

	defer rows.Close()

	// map of full table name (including schema) to table pointer
	tableMap := make(map[string]*Table)

	// keeps track of how many times a table name appears
	// used later to deduplicate using the schema name
	tableNameCount := make(map[string]int)

	// keeps track of the order of tables as they appear in the sorted query (by fullTableName)
	// tableMap can't keep track of this order
	var orderedTables []string

	for rows.Next() {
		var tableType, tableSchema, tableName, columnName, columnType string

		if err := rows.Scan(&tableType, &tableSchema, &tableName, &columnName, &columnType); err != nil {
			return nil, err
		}

		// used to index the tableMap
		fullTableName := tableSchema + "." + tableName

		// add table to map if not already exists
		if _, ok := tableMap[fullTableName]; !ok {
			table := &Table{
				Schema:  tableSchema,
				Name:    tableName,
				RawType: tableType,
			}
			tableNameCount[tableName]++

			tableMap[fullTableName] = table
			orderedTables = append(orderedTables, fullTableName)
		}

		// create the field corresponding to row in query
		field := TableField{
			Name:    columnName,
			RawType: columnType,
		}

		tableMap[fullTableName].Fields = append(tableMap[fullTableName].Fields, field)
	}

	var tables []Table

	for _, fullTableName := range orderedTables {
		table := tableMap[fullTableName]
		isDuplicate := tableNameCount[table.Name] > 1
		t := table.Populate(&config, isDuplicate)

		tables = append(tables, t)
	}

	return tables, nil
}

func buildTablesQuery(schemas, exclude []string) (string, []interface{}) {
	query := "SELECT t.table_type, c.table_schema, c.table_name, c.column_name, c.data_type" +
		" FROM information_schema.tables AS t" +
		" JOIN information_schema.columns AS c USING (table_schema, table_name)" +
		" WHERE table_schema IN " + sqgen.SliceToSQL(schemas)

	if len(exclude) > 0 {
		query += " AND table_name NOT IN " + sqgen.SliceToSQL(exclude)
	}

	// sql custom ordering: https://stackoverflow.com/q/4088532
	query += " ORDER BY c.table_schema <> 'public', c.table_schema, t.table_type, c.table_name, c.column_name"

	q := replacePlaceholders(query)

	args := make([]interface{}, len(schemas)+len(exclude))

	// if schemas is len 4
	// max index is 3
	for i, schema := range schemas {
		args[i] = schema
	}

	for i, ex := range exclude {
		args[i+len(schemas)] = ex
	}

	return q, args
}

// used in templates

// Adds constructor and struct names to table, populates Fields
// isDuplicate parameter indicates if there is a table in another schema with the same name
func (table Table) Populate(config *Config, isDuplicate bool) Table {
	// Add struct type prefix to struct name. For a list of possible
	// RawTypes that can appear, consult this link (look for table_type):
	// https://www.postgresql.org/docs/current/infoschema-tables.html
	table.StructName = "TABLE_"
	if table.RawType == "VIEW" {
		table.StructName = "VIEW_"
	}

	// Add schema prefix to struct name and constructor if more than one table share same name
	if isDuplicate {
		table.StructName += strings.ToUpper(table.Schema + "__")
		table.Constructor += strings.ToUpper(table.Schema + "__")
	}

	table.StructName += strings.ToUpper(table.Name)
	table.Constructor += strings.ToUpper(table.Name)

	var fields []TableField

	for _, field := range table.Fields {
		f := field.Populate()

		if f.Type == "" {
			if config != nil {
				config.Logger.Printf("Skipping %s.%s because type '%s' is unknown\n", table.Name, field.Name, field.RawType)
			}
			continue
		}

		if strings.ToLower(f.Name) != f.Name {
			if config != nil {
				config.Logger.Printf("Skipping %s.%s because column name is case sensitive\n", table.Name, field.Name)
			}
			continue
		}

		fields = append(fields, f)
	}

	table.Fields = fields

	return table
}

// populate will fill in the .Type and .Constructor for a field based on
// the field's .RawType. For list of possible RawTypes that can appear, consult
// this link (Table 8.1): https://www.postgresql.org/docs/current/datatype.html.
func (field TableField) Populate() TableField {
	// Boolean
	if field.RawType == "boolean" {
		field.Type = FieldTypeBoolean
		field.Constructor = FieldConstructorBoolean
		return field
	}

	// JSON
	if strings.HasPrefix(field.RawType, "json") {
		field.Type = FieldTypeJSON
		field.Constructor = FieldConstructorJSON
		return field
	}

	// Number
	switch field.RawType {
	case "oid": // https://www.postgresql.org/docs/current/datatype-oid.html
		fallthrough
	case "decimal", "numeric", "real", "double precision": // float
		fallthrough
	case "smallint", "integer", "bigint", "smallserial", "serial", "bigserial": // integer
		field.Type = FieldTypeNumber
		field.Constructor = FieldConstructorNumber
		return field
	}

	// String
	switch {
	case field.RawType == "name": // https://dba.stackexchange.com/questions/217533/what-is-the-data-type-name-in-postgresql
		fallthrough
	case field.RawType == "text", strings.HasPrefix(field.RawType, "char"), strings.HasPrefix(field.RawType, "varchar"):
		field.Type = FieldTypeString
		field.Constructor = FieldConstructorString
		return field
	}

	// Time
	if strings.HasPrefix(field.RawType, "time") || field.RawType == "date" {
		field.Type = FieldTypeTime
		field.Constructor = FieldConstructorTime
		return field
	}

	// Enum
	if field.RawType == "USER-DEFINED" {
		field.Type = FieldTypeEnum
		field.Constructor = FieldConstructorEnum
		return field
	}

	// Array
	if field.RawType == "ARRAY" {
		field.Type = FieldTypeArray
		field.Constructor = FieldConstructorArray
		return field
	}

	// Bytea
	if field.RawType == "bytea" {
		field.Type = FieldTypeBinary
		field.Constructor = FieldConstructorBinary
		return field
	}

	return field
}
