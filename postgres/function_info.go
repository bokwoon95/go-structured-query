package sq

import "strings"

type FunctionInfo struct {
	Schema    string
	Name      string
	Alias     string
	Arguments []interface{}
}

// AppendSQL adds the fully qualified function call into the buffer.
func (f *FunctionInfo) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	f.AppendSQLExclude(buf, args, nil)
}

// AppendSQLExclude adds the fully qualified function call into the buffer.
func (f *FunctionInfo) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	if f == nil {
		return
	}
	var format string
	if f.Schema != "" {
		if strings.ContainsAny(f.Schema, " \t") {
			format = `"` + f.Schema + `".`
		} else {
			format = f.Schema + "."
		}
	}
	switch len(f.Arguments) {
	case 0:
		format = format + f.Name + "()"
	default:
		format = format + f.Name + "(?" + strings.Repeat(", ?", len(f.Arguments)-1) + ")"
	}
	expandValues(buf, args, excludedTableQualifiers, format, f.Arguments)
}

func Functionf(name string, args ...interface{}) *FunctionInfo {
	return &FunctionInfo{
		Name:      name,
		Arguments: args,
	}
}

// GetAlias implements the Table interface. It returns the alias of the
// FunctionInfo.
func (f *FunctionInfo) GetAlias() string {
	return f.Alias
}

// GetName implements the Table interface. It returns the name of the
// FunctionInfo.
func (f *FunctionInfo) GetName() string {
	return f.Name
}
