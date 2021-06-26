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

func TestTablePopulate(t *testing.T) {
	tt := []struct{
		name string
		table Table
		isDuplicate bool
		result Table
	}{
		{
			name: "normal table name, not duplicate",
			table: Table{
				Name: "users",
				Schema: "public",
			},
			isDuplicate: false,
			result: Table{
				Name: "users",
				Schema: "public",
				StructName: "TABLE_USERS",
				Constructor: "USERS",
			},
		},
		{
			name: "normal table name, is duplicate",
			table: Table{
				Name: "users",
				Schema: "public",
			},
			isDuplicate: true,
			result: Table{
				Name: "users",
				Schema: "public",
				StructName: "TABLE_PUBLIC__USERS",
				Constructor: "PUBLIC__USERS",
			},
		},
		{
			name: "normal table name with different schema, is duplicate",
			table: Table{
				Name: "users",
				Schema: "geo",
			},
			isDuplicate: true,
			result: Table{
				Name: "users",
				Schema: "geo",
				StructName: "TABLE_GEO__USERS",
				Constructor: "GEO__USERS",
			},
		},
		{
			name: "normal view, is not duplicate",
			table: Table{
				Name: "verified_users",
				Schema: "public",
				RawType: "VIEW",
			},
			isDuplicate: false,
			result: Table{
				Name: "verified_users",
				Schema: "public",
				RawType: "VIEW",
				StructName: "VIEW_VERIFIED_USERS",
				Constructor: "VERIFIED_USERS",
			},
		},
		{
			name: "normal view, is duplicate",
			table: Table{
				Name: "verified_users",
				Schema: "public",
				RawType: "VIEW",
			},
			isDuplicate: true,
			result: Table{
				Name: "verified_users",
				Schema: "public",
				RawType: "VIEW",
				StructName: "VIEW_PUBLIC__VERIFIED_USERS",
				Constructor: "PUBLIC__VERIFIED_USERS",
			},
		},
		{
			name: "normal table name, not duplicate, skips unknown fields",
			table: Table{
				Name: "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name: "id",
						RawType: "some_unknown_type",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name: "users",
				Schema: "public",
				StructName: "TABLE_USERS",
				Constructor: "USERS",
			},
		},
		{
			name: "normal table name, not duplicate, skips case-sensitive field names",
			table: Table{
				Name: "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name: "ID",
						RawType: "boolean",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name: "users",
				Schema: "public",
				StructName: "TABLE_USERS",
				Constructor: "USERS",
			},
		},
		{
			name: "normal table name, not duplicate, doesn't skip supported fields",
			table: Table{
				Name: "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name: "id",
						RawType: "boolean",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name: "users",
				Schema: "public",
				StructName: "TABLE_USERS",
				Constructor: "USERS",
				Fields: []TableField{
					{
						Name: "id",
						RawType: "boolean",
						Type: FieldTypeBoolean,
						Constructor: FieldConstructorBoolean,
					},
				},
			},
		},
		{
			name: "normal table name, not duplicate, can populate multiple fields",
			table: Table{
				Name: "users",
				Schema: "public",
				Fields: []TableField{
					{
						Name: "id",
						RawType: "integer",
					},
					{
						Name: "first_name",
						RawType: "text",
					},
					{
						Name: "last_name",
						RawType: "varchar",
					},
					{
						Name: "date_created",
						RawType: "timestamp",
					},
					{
						Name: "is_verified",
						RawType: "boolean",
					},
					{
						Name: "data",
						RawType: "jsonb",
					},
				},
			},
			isDuplicate: false,
			result: Table{
				Name: "users",
				Schema: "public",
				StructName: "TABLE_USERS",
				Constructor: "USERS",
				Fields: []TableField{
					{
						Name: "id",
						RawType: "integer",
						Type: FieldTypeNumber,
						Constructor: FieldConstructorNumber,
					},
					{
						Name: "first_name",
						RawType: "text",
						Type: FieldTypeString,
						Constructor: FieldConstructorString,
					},
					{
						Name: "last_name",
						RawType: "varchar",
						Type: FieldTypeString,
						Constructor: FieldConstructorString,
					},
					{
						Name: "date_created",
						RawType: "timestamp",
						Type: FieldTypeTime,
						Constructor: FieldConstructorTime,
					},
					{
						Name: "is_verified",
						RawType: "boolean",
						Type: FieldTypeBoolean,
						Constructor: FieldConstructorBoolean,
					},
					{
						Name: "data",
						RawType: "jsonb",
						Type: FieldTypeJSON,
						Constructor: FieldConstructorJSON,
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(tc.table.Populate(nil, tc.isDuplicate), tc.result)
		})
	}
}
