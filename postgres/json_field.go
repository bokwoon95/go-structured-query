package sq

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// JSONField either represents a JSON column or a literal value that can be
// marshalled into a JSON string.
type JSONField struct {
	// JSONField will be one of the following:

	// 1) Literal JSONable value (almost all structs can be converted to JSON)
	value interface{}

	// 2) JSON column
	alias      string
	table      Table
	name       string
	descending *bool
	nullsfirst *bool
}

// AppendSQLExclude marshals the JSONField into an SQL query and args as
// described in the JSONField internal struct comments.
func (f JSONField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	switch {
	case f.value != nil:
		// 1) Literal JSONable value
		buf.WriteString("?")
		*args = append(*args, f.value)
	default:
		// 2) JSON column
		tableQualifier := f.table.GetAlias()
		if tableQualifier == "" {
			tableQualifier = f.table.GetName()
		}
		for _, excludedTableQualifier := range excludedTableQualifiers {
			if tableQualifier == excludedTableQualifier {
				tableQualifier = ""
				break
			}
		}
		if tableQualifier != "" {
			if strings.ContainsAny(tableQualifier, " \t") {
				buf.WriteString(`"`)
				buf.WriteString(tableQualifier)
				buf.WriteString(`".`)
			} else {
				buf.WriteString(tableQualifier)
				buf.WriteString(".")
			}
		}
		if strings.ContainsAny(f.name, " \t") {
			buf.WriteString(`"`)
			buf.WriteString(f.name)
			buf.WriteString(`"`)
		} else {
			buf.WriteString(f.name)
		}
	}
	if f.descending != nil {
		if *f.descending {
			buf.WriteString(" DESC")
		} else {
			buf.WriteString(" ASC")
		}
	}
	if f.nullsfirst != nil {
		if *f.nullsfirst {
			buf.WriteString(" NULLS FIRST")
		} else {
			buf.WriteString(" NULLS LAST")
		}
	}
}

// NewJSONField returns a new JSONField representing a JSON column.
func NewJSONField(name string, table Table) JSONField {
	return JSONField{
		name:  name,
		table: table,
	}
}

// JSON returns a new JSONField representing a literal JSONable value. It
// returns an error indicating if the value can be marshalled into JSON.
func JSON(val interface{}) (JSONField, error) {
	f := JSONField{
		value: val,
	}
	_, err := json.Marshal(val)
	if err != nil {
		return f, err
	}
	return f, nil
}

// MustJSON is like JSON but it panics on error.
func MustJSON(val interface{}) JSONField {
	f, err := JSON(val)
	if err != nil {
		panic(err)
	}
	return f
}

// JSONValue returns a new JSONField representing a driver.Valuer value.
func JSONValue(val driver.Valuer) JSONField {
	return JSONField{
		value: val,
	}
}

// Set returns a FieldAssignment associating the JSONField to the value i.e.
// 'field = value'.
func (f JSONField) Set(value interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

// SetJSON returns a FieldAssignment associating the JSONField to the JSONable
// value i.e. 'field = value'. Internally it uses MustJSON, which means it
// will panic if the value cannot be marshalled into JSON.
func (f JSONField) SetJSON(value interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: MustJSON(value).value,
	}
}

// SetValue returns a FieldAssignment associating the JSONField to the driver.Valuer
// value i.e. 'field = value'.
func (f JSONField) SetValue(value driver.Valuer) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

// As returns a new JSONField with the new field Alias i.e. 'field AS Alias'.
func (f JSONField) As(alias string) JSONField {
	f.alias = alias
	return f
}

// Asc returns a new JSONField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f JSONField) Asc() JSONField {
	desc := false
	f.descending = &desc
	return f
}

// Desc returns a new JSONField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f JSONField) Desc() JSONField {
	desc := true
	f.descending = &desc
	return f
}

// NullsFirst returns a new JSONField indicating that it should be ordered
// with nulls first i.e. 'ORDER BY field NULLS FIRST'.
func (f JSONField) NullsFirst() JSONField {
	nullsfirst := true
	f.nullsfirst = &nullsfirst
	return f
}

// NullsLast returns a new JSONField indicating that it should be ordered
// with nulls last i.e. 'ORDER BY field NULLS LAST'.
func (f JSONField) NullsLast() JSONField {
	nullsfirst := false
	f.nullsfirst = &nullsfirst
	return f
}

// IsNull returns an 'X IS NULL' Predicate.
func (f JSONField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f JSONField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// String implements the fmt.Stringer interface. It returns the string
// representation of a JSONField.
func (f JSONField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil, nil)
	return questionInterpolate(buf.String(), args...)
}

// GetAlias implements the Field interface. It returns the Alias of the
// JSONField.
func (f JSONField) GetAlias() string {
	return f.alias
}

// GetName implements the Field interface. It returns the Name of the
// JSONField.
func (f JSONField) GetName() string {
	return f.name
}
