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

	t.Run("multiple schemas, no excluded tables", func(t *testing.T) {
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

func TestTablePopulate(t *testing.T) {
	type TT struct {
		name        string
		table       Table
		isDuplicate bool
		result      Table
	}
	tests := []TT{
		{
			name: "normal table name, not duplicate",
			table: Table{
				Name:   "users",
				Schema: "public",
			},
			isDuplicate: false,
			result: Table{
				Name:        "users",
				Schema:      "public",
				StructName:  "TABLE_USERS",
				Constructor: "USERS",
			},
		},
		{
			name: "normal table name, is duplicate",
			table: Table{
				Name:   "users",
				Schema: "public",
			},
			isDuplicate: true,
			result: Table{
				Name:        "users",
				Schema:      "public",
				StructName:  "TABLE_PUBLIC__USERS",
				Constructor: "PUBLIC__USERS",
			},
		},
		{
			name: "normal table name with different schema, is duplicate",
			table: Table{
				Name:   "users",
				Schema: "geo",
			},
			isDuplicate: true,
			result: Table{
				Name:        "users",
				Schema:      "geo",
				StructName:  "TABLE_GEO__USERS",
				Constructor: "GEO__USERS",
			},
		},
		{
			name: "normal view, is not duplicate",
			table: Table{
				Name:    "verified_users",
				Schema:  "public",
				RawType: "VIEW",
			},
			isDuplicate: false,
			result: Table{
				Name:        "verified_users",
				Schema:      "public",
				RawType:     "VIEW",
				StructName:  "VIEW_VERIFIED_USERS",
				Constructor: "VERIFIED_USERS",
			},
		},
		{
			name: "normal view, is duplicate",
			table: Table{
				Name:    "verified_users",
				Schema:  "public",
				RawType: "VIEW",
			},
			isDuplicate: true,
			result: Table{
				Name:        "verified_users",
				Schema:      "public",
				RawType:     "VIEW",
				StructName:  "VIEW_PUBLIC__VERIFIED_USERS",
				Constructor: "PUBLIC__VERIFIED_USERS",
			},
		},
		{
			name: "normal table name, not duplicate, skips unknown fields",
			table: Table{
				Name:   "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name:    "id",
						RawType: "some_unknown_type",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name:        "users",
				Schema:      "public",
				StructName:  "TABLE_USERS",
				Constructor: "USERS",
			},
		},
		{
			name: "normal table name, not duplicate, skips case-sensitive field names",
			table: Table{
				Name:   "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name:    "ID",
						RawType: "boolean",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name:        "users",
				Schema:      "public",
				StructName:  "TABLE_USERS",
				Constructor: "USERS",
			},
		},
		{
			name: "normal table name, not duplicate, doesn't skip supported fields",
			table: Table{
				Name:   "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name:    "id",
						RawType: "boolean",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name:        "users",
				Schema:      "public",
				StructName:  "TABLE_USERS",
				Constructor: "USERS",
				Fields: []TableField{
					{
						Name:        "id",
						RawType:     "boolean",
						Type:        FieldTypeBoolean,
						Constructor: FieldConstructorBoolean,
					},
				},
			},
		},
		{
			name: "normal table name, not duplicate, can populate multiple fields",
			table: Table{
				Name:   "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name:    "id",
						RawType: "integer",
					},
					{
						Name:    "first_name",
						RawType: "text",
					},
					{
						Name:    "last_name",
						RawType: "varchar",
					},
					{
						Name:    "date_created",
						RawType: "timestamp",
					},
					{
						Name:    "is_verified",
						RawType: "boolean",
					},
					{
						Name:    "data",
						RawType: "jsonb",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name:        "users",
				Schema:      "public",
				StructName:  "TABLE_USERS",
				Constructor: "USERS",
				Fields: []TableField{
					{
						Name:        "id",
						RawType:     "integer",
						Type:        FieldTypeNumber,
						Constructor: FieldConstructorNumber,
					},
					{
						Name:        "first_name",
						RawType:     "text",
						Type:        FieldTypeString,
						Constructor: FieldConstructorString,
					},
					{
						Name:        "last_name",
						RawType:     "varchar",
						Type:        FieldTypeString,
						Constructor: FieldConstructorString,
					},
					{
						Name:        "date_created",
						RawType:     "timestamp",
						Type:        FieldTypeTime,
						Constructor: FieldConstructorTime,
					},
					{
						Name:        "is_verified",
						RawType:     "boolean",
						Type:        FieldTypeBoolean,
						Constructor: FieldConstructorBoolean,
					},
					{
						Name:        "data",
						RawType:     "jsonb",
						Type:        FieldTypeJSON,
						Constructor: FieldConstructorJSON,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(tt.table.Populate(nil, tt.isDuplicate), tt.result)
		})
	}
}

func TestTableFieldPopulate(t *testing.T) {
	type TT struct {
		name   string
		field  TableField
		result TableField
	}
	tests := []TT{
		{
			name: "unknown field",
			field: TableField{
				Name:    "flag",
				RawType: "some_unknown_type",
			},
			result: TableField{
				Name:    "flag",
				RawType: "some_unknown_type",
			},
		},
		{
			name: "boolean field",
			field: TableField{
				Name:    "flag",
				RawType: "boolean",
			},
			result: TableField{
				Name:        "flag",
				RawType:     "boolean",
				Type:        FieldTypeBoolean,
				Constructor: FieldConstructorBoolean,
			},
		},
		{
			name: "json field",
			field: TableField{
				Name:    "data",
				RawType: "json",
			},
			result: TableField{
				Name:        "data",
				RawType:     "json",
				Type:        FieldTypeJSON,
				Constructor: FieldConstructorJSON,
			},
		},
		{
			name: "jsonb field",
			field: TableField{
				Name:    "data",
				RawType: "jsonb",
			},
			result: TableField{
				Name:        "data",
				RawType:     "jsonb",
				Type:        FieldTypeJSON,
				Constructor: FieldConstructorJSON,
			},
		},
		{
			name: "user-defined field",
			field: TableField{
				Name:    "type",
				RawType: "USER-DEFINED",
			},
			result: TableField{
				Name:        "type",
				RawType:     "USER-DEFINED",
				Type:        FieldTypeEnum,
				Constructor: FieldConstructorEnum,
			},
		},
		{
			name: "array field",
			field: TableField{
				Name:    "arr",
				RawType: "ARRAY",
			},
			result: TableField{
				Name:        "arr",
				RawType:     "ARRAY",
				Type:        FieldTypeArray,
				Constructor: FieldConstructorArray,
			},
		},
		{
			name: "bytea field",
			field: TableField{
				Name:    "hash",
				RawType: "bytea",
			},
			result: TableField{
				Name:        "hash",
				RawType:     "bytea",
				Type:        FieldTypeBinary,
				Constructor: FieldConstructorBinary,
			},
		},
	}

	numberFields := []string{
		"oid",
		"decimal",
		"numeric",
		"real",
		"double precision",
		"smallint",
		"integer",
		"bigint",
		"smallserial",
		"serial",
		"bigserial",
	}

	for _, rawType := range numberFields {
		tests = append(tests, TT{
			name: rawType + " field",
			field: TableField{
				Name:    "number",
				RawType: rawType,
			},
			result: TableField{
				Name:        "number",
				RawType:     rawType,
				Type:        FieldTypeNumber,
				Constructor: FieldConstructorNumber,
			},
		})
	}

	stringFields := []string{
		"name",
		"text",
		"char",
		"char(64)",
		"varchar",
		"varchar(64)",
	}

	for _, rawType := range stringFields {
		tests = append(tests, TT{
			name: rawType + " field",
			field: TableField{
				Name:    "name",
				RawType: rawType,
			},
			result: TableField{
				Name:        "name",
				RawType:     rawType,
				Type:        FieldTypeString,
				Constructor: FieldConstructorString,
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(tt.field.Populate(), tt.result)
		})
	}
}
