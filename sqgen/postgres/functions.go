// contains the logic for the sqgen-postgres functions command
package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"regexp"
	"strings"
	"strconv"

	"github.com/bokwoon95/go-structured-query/sqgen"
)

// Function contains metadata for a plpgsql function.
type Function struct {
	Schema       string
	Name         string
	RawResults   string
	RawArguments string
	StructName   string
	Constructor  string
	Results      []FunctionField
	Arguments    []FunctionField
}

// FunctionField represents a Function that is also a Field.
type FunctionField struct {
	RawField    string
	Name        string
	FieldType   string
	GoType      string
	Constructor string
}

func BuildFunctions(config Config, writer io.Writer) error {
	db, err := openAndPing(config.Database)

	if err != nil {
		return sqgen.Wrap(err)
	}

	functions, err := executeFunctions(config, db)

	if err != nil {
		return sqgen.Wrap(err)
	}

	templateData := FunctionsTemplateData{
		PackageName: config.Package,
		Imports: []string{
			`sq "github.com/bokwoon95/go-structured-query/postgres"`,
		},
		Functions: functions,
	}

	t, err := getFunctionsTemplate()

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
		return err
	}

	_, err = writer.Write(src)

	return err
}

func executeFunctions(config Config, db *sql.DB) ([]Function, error) {
	pgVersion, err := queryPgVersion(db)

	if err != nil {
		return nil, sqgen.Wrap(err)
	}

	supportsProkind := checkProkindSupport(pgVersion)

	query, args := buildFunctionsQuery(config.Schemas, config.Exclude, supportsProkind)

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, sqgen.Wrap(err)
	}

	defer rows.Close()

	// "full" function name refers to schema + function name
	// map from full function name to slice of functions
	// each function in the slice refers to a function overload
	// if only one item, it's a non-overloaded function
	functionMap := make(map[string][]Function)

	// keeps track of how many times a function name appears, EXCLUDING OVERLOADS
	// i.e. only incremented the first time a function is encountered in a given schema
	functionNameCount := make(map[string]int)

	// keeps track of the order of functions as they appear in the sorted query
	// functionMap can't keep track of this order
	var orderedFunctions []string

	for rows.Next() {
		// scan into functionMap, 
		var schema, name, rawResults, rawArguments string

		err := rows.Scan(&schema, &name, &rawResults, &rawArguments)

		if err != nil {
			return nil, err
		}

		schema = strings.ReplaceAll(schema, " ", "_")
		name = strings.ReplaceAll(name, " ", "_")

		function := Function{
			Schema: schema,
			Name: name,
			RawResults: rawResults,
			RawArguments: rawArguments,
		}

		qualifiedName := fmt.Sprintf("%s.%s", schema, name)
		functionArr := functionMap[qualifiedName]

		// first time a function with this name was encountered in a given schema
		if len(functionArr) == 0 {
			functionNameCount[name]++
			
			// only need one item in this slice per set of function overloads
			// prevents generating the same function more than once
			orderedFunctions = append(orderedFunctions, qualifiedName)
		}

		functionMap[qualifiedName] = append(functionArr, function)
	}

	var functions []Function

	for _, fullFunctionName := range orderedFunctions {
		funcSlice := functionMap[fullFunctionName]

		if funcSlice == nil {
			continue
		}
		
		for i, function := range funcSlice {
			var overloadCount int
			isDuplicate := functionNameCount[function.Name] > 1

			if len(funcSlice) > 1 {
				overloadCount = i + 1
			}

			f, err := function.Populate(isDuplicate, overloadCount)

			if err != nil {
				config.Logger.Println(err)
				continue
			}

			functions = append(functions, *f)
		}

	}

	return functions, nil
}

func buildFunctionsQuery(schemas, exclude []string, supportsProkind bool) (string, []interface{}) {
	query := "SELECT n.nspname, p.proname" +
		", pg_catalog.pg_get_function_result(p.oid) AS result" +
		", pg_catalog.pg_get_function_identity_arguments(p.oid) as arguments" +
		" FROM pg_catalog.pg_proc AS p" +
		" LEFT JOIN pg_catalog.pg_namespace AS n ON n.oid = p.pronamespace" +
		" WHERE n.nspname IN " + sqgen.SliceToSQL(schemas)

	// following block filters for only functions, not window/aggregate/procedures
	// support for prokind column in pg_proc changed in postgres 11
	// supportsProkind param indicates if that column on pg_proc is supported
	if supportsProkind {
		query += "  AND p.prokind = 'f'"
	} else {
		// p.prokind also rules out any procedures... use the p.prorettype <> 0 check to remove procedures from result set
		// see: https://git.postgresql.org/gitweb/?p=postgresql.git;a=commitdiff;h=fd1a421fe66173fb9b85d3fe150afde8e812cbe4
		query += "  AND p.proisagg = false AND p.proiswindow = false AND p.prorettype <> 0"
	}

	if len(exclude) > 0 {
		query += " AND p.proname NOT IN " + sqgen.SliceToSQL(exclude)
	}

	// sql custom ordering: https://stackoverflow.com/q/4088532
	query += "ORDER BY n.nspname <> 'public', n.nspname, p.proname"

	q := replacePlaceholders(query)

	args := make([]interface{}, len(schemas)+len(exclude))

	for i, schema := range schemas {
		args[i] = schema
	}

	for i, ex := range exclude {
		args[i+len(schemas)] = ex
	}

	return q, args
}

func queryPgVersion(db *sql.DB) (string, error) {
	var version string

	query := "SHOW server_version;"

	row := db.QueryRow(query)

	if err := row.Err(); err != nil {
		return "", err
	}

	if err := row.Scan(&version); err != nil {
		return "", err
	}

	return version, nil
}

func checkProkindSupport(version string) bool {
	re := regexp.MustCompile(`(\d{1,2})\..*`)
	matches := re.FindStringSubmatch(version)

	if len(matches) < 2 {
		panic("version matches not long enough")
	}

	// 0th match is the whole regexp
	// 1st match is the capturing group
	majorVersionStr := matches[1]

	majorVersion, err := strconv.Atoi(majorVersionStr)

	if err != nil {
		panic(err)
	}

	return majorVersion >= 11
}

// isDuplicate refers to if the function name is duplicated in an other schema
// if we have multiple function overloads within the same schema
// suffix the generated function with an index (starting at 1) that increments with each function overload
// the current overload index that we're on is the overloadCount
// if it is 0, means that the function is not overloaded, and we can skip adding the suffix
func (function Function) Populate(isDuplicate bool, overloadCount int) (*Function, error) {
	function.StructName = "FUNCTION_"

	if isDuplicate {
		schemaPrefix := strings.ToUpper(function.Schema) + "__"
		function.StructName += schemaPrefix
		function.Constructor += schemaPrefix
	}

	function.StructName += strings.ToUpper(function.Name)
	function.Constructor += strings.ToUpper(function.Name)

	if overloadCount > 0 {
		function.StructName += strconv.Itoa(overloadCount)
		function.Constructor += strconv.Itoa(overloadCount)
	}

	// Function Arguments

	if function.RawArguments != "" {
		rawFields := strings.Split(function.RawArguments, ",")

		for i := range rawFields {
			field := extractNameAndType(rawFields[i])

			rawField := strings.ToUpper(field.RawField)

			if strings.HasPrefix(rawField, "VARIADIC ") {
				err := fmt.Errorf("Skipping %s.%s because VARIADIC arguments are not supported '%s'", function.Schema, function.Name, field.RawField)
				return nil, err
			}

			if strings.HasPrefix(rawField, "IN ") || strings.HasPrefix(rawField, "OUT ") {
				err := fmt.Errorf("Skipping %s.%s because INOUT arguments are not supported '%s'", function.Schema, function.Name, function.RawArguments)
				return nil, err
			}

			if field.FieldType == "" {
				err := fmt.Errorf("Skipping %s.%s because user-defined type '%s' is not supported", function.Schema, function.Name, field.RawField)
				return nil, err
			}

			if field.Name == "" {
				field.Name = fmt.Sprintf("_arg%d", i+1)
			}

			function.Arguments = append(function.Arguments, field)
		}
	}

	// Function Return Types

	isTable := strings.HasPrefix(function.RawResults, "TABLE(") && strings.HasSuffix(function.RawResults, ")")

	if function.RawResults == "void" {
		// no return type
	} else if function.RawResults == "trigger" {
		err := fmt.Errorf("Skipping %s.%s because it is a trigger function", function.Schema, function.Name)
		return nil, err
	} else if isTable {
		rawResults := function.RawResults[6 : len(function.RawResults) - 1] // remove 'TABLE (' prefix and ')' suffix
		rawFields := strings.Split(rawResults, ",")

		for i := range rawFields {
			field := extractNameAndType(rawFields[i])
			
			if field.FieldType == "" {
				err := fmt.Errorf("Skipping %s.%s because return type '%s' is not supported", function.Schema, function.Name, field.RawField)
				return nil, err
			}

			if strings.ToLower(field.Name) != field.Name {
				err := fmt.Errorf("Skipping %s.%s because return type '%s' is case-sensitive", function.Schema, function.Name, field.Name)
				return nil, err
			}

			if field.Name == "" {
				field.Name = fmt.Sprintf("Result%d", i+1)
			}

			function.Results = append(function.Results, field)
		}
	} else {
		rawResults := strings.TrimPrefix(function.RawResults, "SETOF ")
		rawResults = strings.TrimSpace(rawResults)

		field := extractNameAndType(rawResults)

		if field.FieldType == "" {
			err := fmt.Errorf("Skipping %s.%s because return type '%s' is not supported", function.Schema, function.Name, function.RawResults)
			return nil, err
		}

		if strings.ToLower(field.Name) != field.Name {
			err := fmt.Errorf("Skipping %s.%s because return type '%s' is case-sensitive", function.Schema, function.Name, field.Name)
			return nil, err
		}

		field.Name = "Result"
		function.Results = []FunctionField{field}
	}

	return &function, nil
}

func extractNameAndType(rawField string) FunctionField {
	var field FunctionField
	field.RawField = strings.TrimSpace(rawField)
	if matches := regexp.
		MustCompile(`boolean` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(field.RawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Boolean
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeBoolSlice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeBoolean
			field.GoType = GoTypeBool
			field.Constructor = FieldConstructorBoolean
		}

	} else if matches := regexp.
		MustCompile(`json` + `(?:b)?` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(field.RawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// JSON
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeInterface
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeJSON
			field.GoType = GoTypeInterface
			field.Constructor = FieldConstructorJSON
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `smallint` +
			`|` + `oid` +
			`|` + `integer` +
			`|` + `bigint` +
			`|` + `smallserial` +
			`|` + `serial` +
			`|` + `bigserial` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(field.RawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Integer
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeIntSlice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeNumber
			field.GoType = GoTypeInt
			field.Constructor = FieldConstructorNumber
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `decimal` +
			`|` + `numeric` +
			`|` + `real` +
			`|` + `double precision` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(field.RawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Float
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeFloat64Slice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeNumber
			field.GoType = GoTypeFloat64
			field.Constructor = FieldConstructorNumber
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `text` +
			`|` + `name` +
			`|` + `char` + `(?:\(\d+\))?` +
			`|` + `character` + `(?:\(\d+\))?` +
			`|` + `varchar` + `(?:\(\d+\))?` +
			`|` + `character varying` + `(?:\(\d+\))?` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(field.RawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// String
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeStringSlice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeString
			field.GoType = GoTypeString
			field.Constructor = FieldConstructorString
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `date` +
			`|` + `(?:time|timestamp)` +
			`(?: \(\d+\))?` +
			`(?: without time zone| with time zone)?` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(field.RawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Time
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			// Do nothing
		} else {
			field.FieldType = FieldTypeTime
			field.GoType = GoTypeTime
			field.Constructor = FieldConstructorTime
		}

	} else if matches := regexp.
		MustCompile(`bytea` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(field.RawField); len(matches) == 2 {
		if matches[1] == "[]" {
			// Do nothing
		} else {
			field.FieldType = FieldTypeBinary
			field.GoType = GoTypeByteSlice
			field.Constructor = FieldConstructorBinary
		}
	}
	return field
}
