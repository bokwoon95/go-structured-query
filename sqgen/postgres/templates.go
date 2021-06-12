package postgres

import (
	"fmt"
	"strings"
	"text/template"
)

func export(s string) string {
	str := strings.TrimPrefix(s, "_")
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ToUpper(str)
	return str
}

func quoteSpace(s string) string {
	if strings.Contains(s, " ") {
		return fmt.Sprintf(`"%s"`, s)
	}

	return s
}

// functions required to transforms strings within the template
// removes the need to declare a custom string alias with these methods attached 
// since all the template variables are referenced either with $ or . accessor in the template
// we won't have any naming collisions
var funcMap template.FuncMap = map[string]interface{}{
	"export":     export,
	"quoteSpace": quoteSpace,
}

func getTablesTemplate() (*template.Template, error) {
	return template.New("").Funcs(funcMap).Parse(tablesTemplate)
}

func getFunctionsTemplate() (*template.Template, error) {
	return template.New("").Funcs(funcMap).Parse(functionsTemplate)
}

type TablesTemplateData struct {
	PackageName string
	Imports     []string
	Tables      []Table
}

type FunctionsTemplateData struct {
	PackageName string
	Imports []string
	Functions []Function
}

var tablesTemplate = `// Code generated by 'sqgen-postgres tables'; DO NOT EDIT.
package {{$.PackageName}}

import (
	{{- range $_, $import := $.Imports}}
	{{$import}}
	{{- end}}
)
{{- range $_, $table := $.Tables}}
{{template "table_struct_definition" $table}}
{{template "table_constructor" $table}}
{{template "table_as" $table}}
{{- end}}

{{- define "table_struct_definition"}}
{{- with $table := .}}
{{- if eq $table.RawType "BASE TABLE"}}
// {{export $table.StructName}} references the {{$table.Schema}}.{{quoteSpace $table.Name}} table.
{{- else if eq $table.RawType "VIEW"}}
// {{export $table.StructName}} references the {{$table.Schema}}.{{quoteSpace $table.Name}} view.
{{- end}}
type {{export $table.StructName}} struct {
	*sq.TableInfo
	{{- range $_, $field := $table.Fields}}
	{{export $field.Name}} {{$field.Type}}
	{{- end}}
}
{{- end}}
{{- end}}

{{- define "table_constructor"}}
{{- with $table := .}}
{{- if eq $table.RawType "BASE TABLE"}}
// {{export $table.Constructor}} creates an instance of the {{$table.Schema}}.{{quoteSpace $table.Name}} table.
{{- else if eq $table.RawType "VIEW"}}
// {{export $table.Constructor}} creates an instance of the {{$table.Schema}}.{{quoteSpace $table.Name}} view.
{{- end}}
func {{export $table.Constructor}}() {{export $table.StructName}} {
	tbl := {{export $table.StructName}}{TableInfo: &sq.TableInfo{
		Schema: "{{$table.Schema}}",
		Name: "{{$table.Name}}",
	},}
	{{- range $_, $field := $table.Fields}}
	tbl.{{export $field.Name}} = {{$field.Constructor}}("{{$field.Name}}", tbl.TableInfo)
	{{- end}}
	return tbl
}
{{- end}}
{{- end}}

{{- define "table_as"}}
{{- with $table := .}}
{{- if eq $table.RawType "BASE TABLE"}}
// As modifies the alias of the underlying table.
{{- else if eq $table.RawType "VIEW"}}
// As modifies the alias of the underlying view.
{{- end}}
func (tbl {{export $table.StructName}}) As(alias string) {{export $table.StructName}} {
	tbl.TableInfo.Alias = alias
	return tbl
}
{{- end}}
{{- end}}`

var functionsTemplate = `// Code generated by 'sqgen-postgres functions'; DO NOT EDIT.
package {{$.PackageName}}
import (
	{{- range $_, $import := $.Imports}}
	{{$import}}
	{{- end}}
)
{{- range $_, $function := $.Functions}}
{{template "function_struct_definition" $function}}
{{template "function_constructor" $function}}
{{template "function_as" $function}}
{{- end}}
{{- define "function_struct_definition"}}
{{- with $function := .}}
// {{export $function.StructName}} references the {{$function.Schema}}.{{$function.Name}} function.
type {{export $function.StructName}} struct {
	*sq.FunctionInfo
	{{- range $_, $result := $function.Results}}
	{{export $result.Name}} {{$result.FieldType}}
	{{- end}}
}
{{- end}}
{{- end}}
{{- define "function_constructor"}}
{{- with $function := .}}
// {{export $function.Constructor}} creates an instance of the {{$function.Schema}}.{{$function.Name}} function.
func {{export $function.Constructor}}(
	{{- range $_, $arg := $function.Arguments}}
	{{$arg.Name}} {{$arg.GoType}},
	{{- end}}
	) {{export $function.StructName}} {
	return {{export $function.Constructor}}_({{range $i, $arg := $function.Arguments}}{{if not $i}}{{$arg.Name}}{{else}}, {{$arg.Name}}{{end}}{{end}})
}
// {{export $function.Constructor}}_ creates an instance of the {{$function.Schema}}.{{$function.Name}} function.
func {{export $function.Constructor}}_(
	{{- range $_, $arg := $function.Arguments}}
	{{$arg.Name}} interface{},
	{{- end}}
	) {{export $function.StructName}} {
	f := {{export $function.StructName}}{FunctionInfo: &sq.FunctionInfo{
		Schema: "{{$function.Schema}}",
		Name: "{{$function.Name}}",
		Arguments: []interface{}{{"{"}}{{range $i, $arg := $function.Arguments}}{{if not $i}}{{$arg.Name}}{{else}}, {{$arg.Name}}{{end}}{{end}}{{"}"}},
	},}
	{{- range $_, $result := $function.Results}}
	f.{{export $result.Name}} = {{$result.Constructor}}("{{$result.Name}}", f.FunctionInfo)
	{{- end}}
	return f
}
{{- end}}
{{- end}}
{{- define "function_as"}}
{{- with $function := .}}
// As modifies the alias of the underlying function.
func (f {{export $function.StructName}}) As(alias string) {{export $function.StructName}} {
	f.FunctionInfo.Alias = alias
	return f
}
{{- end}}
{{- end}}`
