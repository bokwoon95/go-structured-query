package mysql

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

// TableField represents a field in a database table
//
// In the MySQL information_schema.columns table, the data_type column (mapped into TableField.RawType) has the general type of the column (i.e. char, integer, etc.), but it does not have information on the exact size of the column
// The information_schema.column_type column (mapped into TableField.RawTypeEx) contains the full data type for the column, including size
//
// We need to have this extra field to distinguish between boolean columns and number columns, since MySQL stores boolean values as tinyint(1), and we can only get the (1) part from information_schema.column_type
type TableField struct {
	Name        string
	RawType     string
	RawTypeEx   string
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
			`sq "github.com/bokwoon95/go-structured-query/mysql"`,
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
	query, args := buildTablesQuery(config.Schemas, config.Exclude)

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, sqgen.Wrap(err)
	}

	defer rows.Close()

	// map of full table name (including schema) to table pointer
	tableMap := make(map[string]*Table)

	// keeps track of how many time a table name appears (irrespective of schema)
	// used later to deduplicate table definitions with schema name
	tableNameCount := make(map[string]int)

	//keeps track of the order of tables as they appear in the sorted query (by schema name + table name)
	// required, as tableMap is inherently unordered
	var orderedTables []string

	for rows.Next() {
		var tableType, tableSchema, tableName, columnName, columnType, columnTypeEx string

		if err := rows.Scan(&tableType, &tableSchema, &tableName, &columnName, &columnType, &columnTypeEx); err != nil {
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

		field := TableField{
			Name:      columnName,
			RawType:   columnType,
			RawTypeEx: columnTypeEx,
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
	query := "SELECT t.table_type, c.table_schema, c.table_name, c.column_name, c.data_type, c.column_type" +
		" FROM information_schema.tables AS t" +
		" JOIN information_schema.columns AS c USING (table_schema, table_name)" +
		" WHERE table_schema IN " + sqgen.SliceToSQL(
		schemas,
	)

	if len(exclude) > 0 {
		query += " AND table_name NOT IN " + sqgen.SliceToSQL(exclude)
	}

	query += " ORDER BY c.table_schema, t.table_type, c.table_name, c.column_name"

	args := make([]interface{}, len(schemas)+len(exclude))
	for i, schema := range schemas {
		args[i] = schema
	}

	for i, ex := range exclude {
		args[i+len(schemas)] = ex
	}

	return query, args
}

func (table Table) Populate(config *Config, isDuplicate bool) Table {
	table.StructName = "TABLE_"

	if table.RawType == "VIEW" {
		table.StructName = "VIEW_"
	}

	if isDuplicate {
		table.StructName += strings.ToUpper(table.Schema) + "__"
		table.Constructor += strings.ToUpper(table.Schema) + "__"
	}

	table.StructName += strings.ToUpper(table.Name)
	table.Constructor += strings.ToUpper(table.Name)

	var fields []TableField

	for _, field := range table.Fields {
		f := field.Populate()

		if f.Type == "" {
			if config != nil {
				config.Logger.Printf(
					"Skipping %s.%s because type '%s' is unknown\n",
					table.Name,
					field.Name,
					field.RawType,
				)
			}

			continue
		}

		fields = append(fields, field)
	}

	table.Fields = fields

	return table
}

func (field TableField) Populate() TableField {
	// Boolean
	if field.RawTypeEx == "tinyint(1)" {
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
	case "decimal", "numeric", "float", "double": // float
		fallthrough
	case "integer", "int", "smallint", "tinyint", "mediumint", "bigint": // integer
		field.Type = FieldTypeNumber
		field.Constructor = FieldConstructorNumber
		return field
	}

	// String
	switch field.RawType {
	case "tinytext", "text", "mediumtext", "longtext", "char", "varchar":
		field.Type = FieldTypeString
		field.Constructor = FieldConstructorString
		return field
	}

	// Time
	switch field.RawType {
	case "date", "time", "datetime", "timestamp":
		field.Type = FieldTypeTime
		field.Constructor = FieldConstructorTime
		return field
	}

	// Enum
	switch field.RawType {
	case "enum":
		field.Type = FieldTypeEnum
		field.Constructor = FieldConstructorEnum
		return field
	}

	// Blob
	switch field.RawType {
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		field.Type = FieldTypeBinary
		field.Constructor = FieldConstructorBinary
		return field
	}

	return field
}
