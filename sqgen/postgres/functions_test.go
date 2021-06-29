package postgres

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestBuildFunctionsQuery(t *testing.T) {
	type TT struct {
		name            string
		schemas         []string
		exclude         []string
		supportsProkind bool
		expectedQuery   string
		expectedArgs    []interface{}
	}

	var tests []TT

	tests = append(tests, TT{
		name:            "single schema, no exclude, supportsProkind",
		schemas:         []string{"public"},
		exclude:         nil,
		supportsProkind: true,
		expectedQuery:   "SELECT n.nspname, p.proname, pg_catalog.pg_get_function_result(p.oid) AS result, pg_catalog.pg_get_function_identity_arguments(p.oid) as arguments FROM pg_catalog.pg_proc AS p LEFT JOIN pg_catalog.pg_namespace AS n ON n.oid = p.pronamespace WHERE n.nspname IN ($1) AND p.prokind = 'f' ORDER BY n.nspname <> 'public', n.nspname, p.proname",
		expectedArgs:    []interface{}{"public"},
	})

	tests = append(tests, TT{
		name:            "single schema, no exclude, does not support prokind",
		schemas:         []string{"public"},
		exclude:         nil,
		supportsProkind: false,
		expectedQuery:   "SELECT n.nspname, p.proname, pg_catalog.pg_get_function_result(p.oid) AS result, pg_catalog.pg_get_function_identity_arguments(p.oid) as arguments FROM pg_catalog.pg_proc AS p LEFT JOIN pg_catalog.pg_namespace AS n ON n.oid = p.pronamespace WHERE n.nspname IN ($1) AND p.proisagg = false AND p.proiswindow = false AND p.prorettype <> 0 ORDER BY n.nspname <> 'public', n.nspname, p.proname",
		expectedArgs:    []interface{}{"public"},
	})

	tests = append(tests, TT{
		name:            "multiple schemas, no exclude, does not support prokind",
		schemas:         []string{"public", "geo"},
		exclude:         nil,
		supportsProkind: false,
		expectedQuery:   "SELECT n.nspname, p.proname, pg_catalog.pg_get_function_result(p.oid) AS result, pg_catalog.pg_get_function_identity_arguments(p.oid) as arguments FROM pg_catalog.pg_proc AS p LEFT JOIN pg_catalog.pg_namespace AS n ON n.oid = p.pronamespace WHERE n.nspname IN ($1, $2) AND p.proisagg = false AND p.proiswindow = false AND p.prorettype <> 0 ORDER BY n.nspname <> 'public', n.nspname, p.proname",
		expectedArgs:    []interface{}{"public", "geo"},
	})

	tests = append(tests, TT{
		name:            "multiple schemas, exclude functions, supports prokind",
		schemas:         []string{"public", "geo"},
		exclude:         []string{"create_user", "verify_user"},
		supportsProkind: true,
		expectedQuery:   "SELECT n.nspname, p.proname, pg_catalog.pg_get_function_result(p.oid) AS result, pg_catalog.pg_get_function_identity_arguments(p.oid) as arguments FROM pg_catalog.pg_proc AS p LEFT JOIN pg_catalog.pg_namespace AS n ON n.oid = p.pronamespace WHERE n.nspname IN ($1, $2) AND p.prokind = 'f' AND p.proname NOT IN ($3, $4) ORDER BY n.nspname <> 'public', n.nspname, p.proname",
		expectedArgs:    []interface{}{"public", "geo", "create_user", "verify_user"},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			query, args := buildFunctionsQuery(tt.schemas, tt.exclude, tt.supportsProkind)

			is.Equal(query, tt.expectedQuery)
			is.Equal(args, tt.expectedArgs)
		})
	}
}

func TestCheckProkindSupport(t *testing.T) {
	type TT struct {
		name    string
		version string
		support bool
		err     error
	}

	tests := []TT{
		{
			name:    "version above 11 has support",
			version: "11.0.5",
			support: true,
			err:     nil,
		},
		{
			name:    "version above 11 has support",
			version: "12.0.5",
			support: true,
			err:     nil,
		},
		{
			name:    "version below 11 does not support",
			version: "9.5",
			support: false,
			err:     nil,
		},
		{
			name:    "empty version returns error",
			version: "",
			support: false,
			err:     errors.New("could not find version number in string: ''"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			hasSupport, err := checkProkindSupport(tt.version)
			is.Equal(hasSupport, tt.support)
			is.Equal(err, tt.err)
		})
	}
}

func TestIsArrayType(t *testing.T) {
	type TT struct {
		name    string
		matches []string
		result  bool
	}

	tests := []TT{
		{
			name:    "empty array",
			matches: nil,
			result:  false,
		},
		{
			name:    "non-array second match",
			matches: []string{"text", "(1)"},
			result:  false,
		},
		{
			name:    "array second match",
			matches: []string{"text", "[]"},
			result:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(isArrayType(tt.matches), tt.result)
		})
	}
}

func TestGetFieldName(t *testing.T) {
	type TT struct {
		name     string
		rawField string
		matches  []string
		result   string
	}

	tests := []TT{
		{
			name:     "text field",
			rawField: "col text",
			matches:  []string{"text", ""},
			result:   "col",
		},
		{
			name:     "text field with leading space returns empty",
			rawField: " text",
			matches:  []string{"text", ""},
			result:   "",
		},
		{
			name:     "integer field with trailing whitespace",
			rawField: "col integer",
			matches:  []string{"integer", ""},
			result:   "col",
		},
		{
			name:     "integer field with leading space returns empty",
			rawField: " integer",
			matches:  []string{"integer", ""},
			result:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(getFieldName(tt.rawField, tt.matches), tt.result)
		})
	}
}

func TestExtractNameAndType(t *testing.T) {
	type TT struct {
		name     string
		rawField string
		result   FunctionField
	}

	tests := []TT{
		{
			name:     "anon boolean field",
			rawField: " boolean",
			result: FunctionField{
				Name:        "",
				RawField:    "boolean",
				FieldType:   FieldTypeBoolean,
				GoType:      GoTypeBool,
				Constructor: FieldConstructorBoolean,
			},
		},
		{
			name:     "named boolean field",
			rawField: "flag boolean",
			result: FunctionField{
				Name:        "flag",
				RawField:    "flag boolean",
				FieldType:   FieldTypeBoolean,
				GoType:      GoTypeBool,
				Constructor: FieldConstructorBoolean,
			},
		},
		{
			name:     "boolean slice field",
			rawField: "flags boolean[]",
			result: FunctionField{
				Name:        "flags",
				RawField:    "flags boolean[]",
				FieldType:   FieldTypeArray,
				GoType:      GoTypeBoolSlice,
				Constructor: FieldConstructorArray,
			},
		},
		{
			name:     "anon json field",
			rawField: " json",
			result: FunctionField{
				Name:        "",
				RawField:    "json",
				FieldType:   FieldTypeJSON,
				GoType:      GoTypeInterface,
				Constructor: FieldConstructorJSON,
			},
		},
		{
			name:     "named json field",
			rawField: "data json",
			result: FunctionField{
				Name:        "data",
				RawField:    "data json",
				FieldType:   FieldTypeJSON,
				GoType:      GoTypeInterface,
				Constructor: FieldConstructorJSON,
			},
		},
		{
			name:     "anon jsonb field",
			rawField: " jsonb",
			result: FunctionField{
				Name:        "",
				RawField:    "jsonb",
				FieldType:   FieldTypeJSON,
				GoType:      GoTypeInterface,
				Constructor: FieldConstructorJSON,
			},
		},
		{
			name:     "named jsonb field",
			rawField: "data jsonb",
			result: FunctionField{
				Name:        "data",
				RawField:    "data jsonb",
				FieldType:   FieldTypeJSON,
				GoType:      GoTypeInterface,
				Constructor: FieldConstructorJSON,
			},
		},
		{
			name:     "json array field",
			rawField: "data_arr json[]",
			result: FunctionField{
				Name:        "data_arr",
				RawField:    "data_arr json[]",
				FieldType:   FieldTypeArray,
				GoType:      GoTypeInterface,
				Constructor: FieldConstructorArray,
			},
		},
	}

	intFields := []string{
		"smallint",
		"oid",
		"integer",
		"bigint",
		"smallserial",
		"bigserial",
	}

	for _, rawField := range intFields {
		tests = append(tests, TT{
			name:     "anon " + rawField + " field",
			rawField: rawField,
			result: FunctionField{
				Name:        "",
				RawField:    rawField,
				FieldType:   FieldTypeNumber,
				GoType:      GoTypeInt,
				Constructor: FieldConstructorNumber,
			},
		})

		tests = append(tests, TT{
			name:     "named " + rawField + " field",
			rawField: "value " + rawField,
			result: FunctionField{
				Name:        "value",
				RawField:    "value " + rawField,
				FieldType:   FieldTypeNumber,
				GoType:      GoTypeInt,
				Constructor: FieldConstructorNumber,
			},
		})

		tests = append(tests, TT{
			name:     "anon " + rawField + " array field",
			rawField: rawField + "[]",
			result: FunctionField{
				Name:        "",
				RawField:    rawField + "[]",
				FieldType:   FieldTypeArray,
				GoType:      GoTypeIntSlice,
				Constructor: FieldConstructorArray,
			},
		})
	}

	floatFields := []string{
		"decimal",
		"numeric",
		"real",
		"double precision",
	}

	for _, rawField := range floatFields {
		tests = append(tests, TT{
			name:     "anon " + rawField + " field",
			rawField: rawField,
			result: FunctionField{
				Name:        "",
				RawField:    rawField,
				FieldType:   FieldTypeNumber,
				GoType:      GoTypeFloat64,
				Constructor: FieldConstructorNumber,
			},
		})

		tests = append(tests, TT{
			name:     "named " + rawField + " field",
			rawField: "value " + rawField,
			result: FunctionField{
				Name:        "value",
				RawField:    "value " + rawField,
				FieldType:   FieldTypeNumber,
				GoType:      GoTypeFloat64,
				Constructor: FieldConstructorNumber,
			},
		})

		tests = append(tests, TT{
			name:     "anon " + rawField + " array field",
			rawField: rawField + "[]",
			result: FunctionField{
				Name:        "",
				RawField:    rawField + "[]",
				FieldType:   FieldTypeArray,
				GoType:      GoTypeFloat64Slice,
				Constructor: FieldConstructorArray,
			},
		})
	}

	stringFields := []string{
		"text",
		"name",
		"char",
		"char(64)",
		"character",
		"character(8)",
		"varchar",
		"varchar(128)",
		"character varying",
		"character varying(256)",
	}

	for _, rawField := range stringFields {
		tests = append(tests, TT{
			name:     "anon " + rawField + " field",
			rawField: rawField,
			result: FunctionField{
				Name:        "",
				RawField:    rawField,
				FieldType:   FieldTypeString,
				GoType:      GoTypeString,
				Constructor: FieldConstructorString,
			},
		})

		tests = append(tests, TT{
			name:     "named " + rawField + " field",
			rawField: "first_name " + rawField,
			result: FunctionField{
				Name:        "first_name",
				RawField:    "first_name " + rawField,
				FieldType:   FieldTypeString,
				GoType:      GoTypeString,
				Constructor: FieldConstructorString,
			},
		})

		tests = append(tests, TT{
			name:     "anon " + rawField + " array field",
			rawField: rawField + "[]",
			result: FunctionField{
				Name:        "",
				RawField:    rawField + "[]",
				FieldType:   FieldTypeArray,
				GoType:      GoTypeStringSlice,
				Constructor: FieldConstructorArray,
			},
		})
	}

	timeFields := []string{
		"date",
		"time",
		"timestamp",
		"date with time zone",
		"time with time zone",
		"timestamp with time zone",
		"date without time zone",
		"time without time zone",
		"timestamp without time zone",
		"date (32)",
		"time (64)",
		"timestamp (128)",
		"date (32) with time zone",
		"time (64) with time zone",
		"timestamp (128) with time zone",
		"date (128) without time zone",
		"time (64) without time zone",
		"timestamp (32) without time zone",
	}

	for _, rawField := range timeFields {
		tests = append(tests, TT{
			name:     "anon " + rawField + " field",
			rawField: rawField,
			result: FunctionField{
				Name:        "",
				RawField:    rawField,
				FieldType:   FieldTypeTime,
				GoType:      GoTypeTime,
				Constructor: FieldConstructorTime,
			},
		})

		tests = append(tests, TT{
			name:     "named " + rawField + " field",
			rawField: "value " + rawField,
			result: FunctionField{
				Name:        "value",
				RawField:    "value " + rawField,
				FieldType:   FieldTypeTime,
				GoType:      GoTypeTime,
				Constructor: FieldConstructorTime,
			},
		})

		// no fields added if it's an array of times
		tests = append(tests, TT{
			name:     "anon " + rawField + " array field (skipped)",
			rawField: rawField + "[]",
			result: FunctionField{
				Name:     "",
				RawField: rawField + "[]",
			},
		})
	}

	tests = append(tests, TT{
		name:     "anon bytea field",
		rawField: "bytea",
		result: FunctionField{
			Name:        "",
			RawField:    "bytea",
			FieldType:   FieldTypeBinary,
			GoType:      GoTypeByteSlice,
			Constructor: FieldConstructorBinary,
		},
	})

	tests = append(tests, TT{
		name:     "named bytea field",
		rawField: "hash bytea",
		result: FunctionField{
			Name:        "hash",
			RawField:    "hash bytea",
			FieldType:   FieldTypeBinary,
			GoType:      GoTypeByteSlice,
			Constructor: FieldConstructorBinary,
		},
	})

	tests = append(tests, TT{
		name:     "anon bytea array field (skipped)",
		rawField: "bytea[]",
		result: FunctionField{
			Name:     "",
			RawField: "bytea[]",
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(extractNameAndType(tt.rawField), tt.result)
		})
	}
}

func TestFunctionPopulate(t *testing.T) {
	type TT struct {
		name           string
		function       Function
		isDuplicate    bool
		overloadCount  int
		functionResult *Function
		err            error
	}

	tests := []TT{
		{
			name: "regular function, not duplicate, no overload",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "first_name text, hash bytea, is_verified boolean, jsonb",
				RawResults:   " integer",
			},
			isDuplicate:   false,
			overloadCount: 0,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "first_name text, hash bytea, is_verified boolean, jsonb",
				RawResults:   " integer",
				StructName:   "FUNCTION_CREATE_USER",
				Constructor:  "CREATE_USER",
				Arguments: []FunctionField{
					{
						Name:        "first_name",
						RawField:    "first_name text",
						FieldType:   FieldTypeString,
						Constructor: FieldConstructorString,
						GoType:      GoTypeString,
					},
					{
						Name:        "hash",
						RawField:    "hash bytea",
						FieldType:   FieldTypeBinary,
						Constructor: FieldConstructorBinary,
						GoType:      GoTypeByteSlice,
					},
					{
						Name:        "is_verified",
						RawField:    "is_verified boolean",
						FieldType:   FieldTypeBoolean,
						Constructor: FieldConstructorBoolean,
						GoType:      GoTypeBool,
					},
					{
						Name:        "_arg4",
						RawField:    "jsonb",
						FieldType:   FieldTypeJSON,
						Constructor: FieldConstructorJSON,
						GoType:      GoTypeInterface,
					},
				},
				Results: []FunctionField{
					{
						Name:        "Result",
						RawField:    "integer",
						FieldType:   FieldTypeNumber,
						Constructor: FieldConstructorNumber,
						GoType:      GoTypeInt,
					},
				},
			},
			err: nil,
		},
		{
			name: "regular function, duplicate, no overload",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
			},
			isDuplicate:   true,
			overloadCount: 0,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
				StructName:   "FUNCTION_PUBLIC__CREATE_USER",
				Constructor:  "PUBLIC__CREATE_USER",
				Arguments:    nil,
				Results:      nil,
			},
			err: nil,
		},
		{
			name: "regular function, duplicate, overload 1",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
			},
			isDuplicate:   true,
			overloadCount: 1,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
				StructName:   "FUNCTION_PUBLIC__CREATE_USER1",
				Constructor:  "PUBLIC__CREATE_USER1",
				Arguments:    nil,
				Results:      nil,
			},
			err: nil,
		},
		{
			name: "regular function, duplicate, overload 2",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
			},
			isDuplicate:   true,
			overloadCount: 2,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
				StructName:   "FUNCTION_PUBLIC__CREATE_USER2",
				Constructor:  "PUBLIC__CREATE_USER2",
				Arguments:    nil,
				Results:      nil,
			},
			err: nil,
		},
		{
			name: "regular function, not duplicate, overload 1",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
			},
			isDuplicate:   false,
			overloadCount: 1,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "",
				RawResults:   "void",
				StructName:   "FUNCTION_PUBLIC__CREATE_USER1",
				Constructor:  "PUBLIC__CREATE_USER1",
				Arguments:    nil,
				Results:      nil,
			},
			err: nil,
		},
		{
			name: "function with variadic params is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "VARIADIC integer",
				RawResults:   "void",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because VARIADIC arguments are not supported 'VARIADIC integer'",
			),
		},
		{
			name: "function with IN param is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "IN integer",
				RawResults:   "void",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because INOUT arguments are not supported 'IN integer'",
			),
		},
		{
			name: "function with OUT param is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "OUT integer",
				RawResults:   "void",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because INOUT arguments are not supported 'OUT integer'",
			),
		},
		{
			name: "function with INOUT param is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "INOUT integer",
				RawResults:   "void",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because INOUT arguments are not supported 'INOUT integer'",
			),
		},
		{
			name: "function with unknown param type is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "some_unknown_type",
				RawResults:   "void",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because user-defined parameter type 'some_unknown_type' is not supported",
			),
		},
		{
			name: "function with void return type has nil Results",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "void",
			},
			isDuplicate:   false,
			overloadCount: 0,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "void",
				StructName:   "FUNCTION_CREATE_USER",
				Constructor:  "CREATE_USER",
				Arguments: []FunctionField{
					{
						Name:        "name",
						RawField:    "name text",
						FieldType:   FieldTypeString,
						Constructor: FieldConstructorString,
						GoType:      GoTypeString,
					},
				},
				Results: nil,
			},
			err: nil,
		},
		{
			name: "trigger function is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "trigger",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because it is a trigger function",
			),
		},
		{
			name: "function with table return type is supported",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "TABLE(first_name text, boolean)",
			},
			isDuplicate:   false,
			overloadCount: 0,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "TABLE(first_name text, boolean)",
				StructName:   "FUNCTION_CREATE_USER",
				Constructor:  "CREATE_USER",
				Arguments: []FunctionField{
					{
						Name:        "name",
						RawField:    "name text",
						FieldType:   FieldTypeString,
						Constructor: FieldConstructorString,
						GoType:      GoTypeString,
					},
				},
				Results: []FunctionField{
					{
						Name:        "first_name",
						RawField:    "first_name text",
						FieldType:   FieldTypeString,
						Constructor: FieldConstructorString,
						GoType:      GoTypeString,
					},
					{
						Name:        "Result2",
						RawField:    "boolean",
						FieldType:   FieldTypeBoolean,
						Constructor: FieldConstructorBoolean,
						GoType:      GoTypeBool,
					},
				},
			},
			err: nil,
		},
		{
			name: "function with table return type with unknown column type is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "TABLE(first_name text, some_unknown_type)",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because return type 'some_unknown_type' is not supported",
			),
		},
		{
			name: "function with SETOF return type is supported",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "SETOF text",
			},
			isDuplicate:   false,
			overloadCount: 0,
			functionResult: &Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "SETOF text",
				StructName:   "FUNCTION_CREATE_USER",
				Constructor:  "CREATE_USER",
				Arguments: []FunctionField{
					{
						Name:        "name",
						RawField:    "name text",
						FieldType:   FieldTypeString,
						Constructor: FieldConstructorString,
						GoType:      GoTypeString,
					},
				},
				Results: []FunctionField{
					{
						Name:        "Result",
						RawField:    "text",
						FieldType:   FieldTypeString,
						Constructor: FieldConstructorString,
						GoType:      GoTypeString,
					},
				},
			},
			err: nil,
		},
		{
			name: "function with SETOF return type with unknown type is skipped",
			function: Function{
				Name:         "create_user",
				Schema:       "public",
				RawArguments: "name text",
				RawResults:   "SETOF some_unknown_type",
			},
			isDuplicate:    false,
			overloadCount:  0,
			functionResult: nil,
			err: errors.New(
				"Skipping public.create_user because SETOF return type 'some_unknown_type' is not supported",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			function, err := tt.function.Populate(tt.isDuplicate, tt.overloadCount)

			is.Equal(err, tt.err)
			is.Equal(function, tt.functionResult)
		})
	}
}
