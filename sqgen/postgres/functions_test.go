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
		name string
		version string
		support bool
		err error
	}

	tests := []TT{
		{
			name: "version above 11 has support",
			version: "11.0.5",
			support: true,
			err: nil,
		},
		{
			name: "version above 11 has support",
			version: "12.0.5",
			support: true,
			err: nil,
		},
		{
			name: "version below 11 does not support",
			version: "9.5",
			support: false,
			err: nil,
		},
		{
			name: "empty version returns error",
			version: "",
			support: false,
			err: errors.New("could not find version number in string: ''"),
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
		name string
		matches []string
		result bool
	}

	tests := []TT{
		{
			name: "empty array",
			matches: nil,
			result: false,
		},
		{
			name: "non-array second match",
			matches: []string{"text", "(1)"},
			result: false,
		},
		{
			name: "array second match",
			matches: []string{"text", "[]"},
			result: true,
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
		name string
		rawField string
		matches []string
		result string
	}

	tests := []TT{
		{
			name: "text field",
			rawField: "col text",
			matches: []string{"text", ""},
			result: "col",
		},
		{
			name: "text field with leading space returns empty",
			rawField: " text",
			matches: []string{"text", ""},
			result: "",
		},
		{
			name: "integer field with trailing whitespace",
			rawField: "col integer",
			matches: []string{"integer", ""},
			result: "col",
		},
		{
			name: "integer field with leading space returns empty",
			rawField: " integer",
			matches: []string{"integer", ""},
			result: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(getFieldName(tt.rawField, tt.matches), tt.result)
		})
	}
}
