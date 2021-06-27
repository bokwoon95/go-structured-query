package postgres

import (
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
