package sq

import (
	"strings"
)

// UUIDField represents a UUID column or a literal UUID value.
type UUIDField struct {
	// 1) Literal UUID value
	// Examples of literal UUID values:
	// | query | args |
	// |-------|---------|
	// | ?     | [16]byte(123e4567-e89b-12d3-a456-426614174000) |
	value *[16]byte

	// 2) UUID column
	alias      string
	table      Table
	name       string
	descending *bool
}

// AppendSQLExclude marshals the UUIDField into a buffer and an args slice.
// It will not table qualify itself if its table qualifer appears in the
// excludedTableQualifiers list.
func (f UUIDField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	switch {
	case f.value != nil:
		// 1) Literal UUID value
		buf.WriteString("?")
		*args = append(*args, *f.value)
	default:
		// 2) UUID column
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
}

// NewUUIDField returns a new UUIDField representing a UUID column.
func NewUUIDField(name string, table Table) UUIDField {
	return UUIDField{
		name:  name,
		table: table,
	}
}

// UUID returns a new UUIDField representing a literal UUID value.
func UUID(u [16]byte) UUIDField {
	return UUIDField{
		value: &u,
	}
}

// Set returns a FieldAssignment associating the UUIDField to the value
// i.e. 'field = value'.
func (f UUIDField) Set(value interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

// SetUUID returns a fieldAssignment associating the UUIDField to the [16]byte value
// i.e. 'field = value'
func (f UUIDField) SetUUID(u [16]byte) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: u,
	}
}

// IsNull returns an 'X IS NULL' Predicate.
func (f UUIDField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f UUIDField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// Eq returns an 'X = Y' Predicate. It only accepts UUIDField.
func (f UUIDField) Eq(field UUIDField) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, field},
	}
}

// Eq returns an 'X <> Y' Predicate. It only accepts UUIDField.
func (f UUIDField) Ne(field UUIDField) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, field},
	}
}

// Eq returns an 'X = Y' Predicate. It only accepts [16]byte
func (f UUIDField) EqUUID(u [16]byte) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, u},
	}
}

// Eq returns an 'X <> Y' Predicate. It only accepts [16]byte
func (f UUIDField) NeUUID(u [16]byte) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, u},
	}
}

// In returns an 'X IN (Y) Predicate'.
func (f UUIDField) In(v interface{}) Predicate {
	var format string
	var values []interface{}

	switch v := v.(type) {
	case RowValue:
		format = "? IN ?"
		values = []interface{}{f, v}
	case Query:
		format = "? IN (?)"
		values = []interface{}{f, v.NestThis()}
	default:
		format = "? IN (?)"
		values = []interface{}{f, v}
	}
	return CustomPredicate{
		Format: format,
		Values: values,
	}
}

// Asc returns a new UUIDField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'
func (f UUIDField) Asc() UUIDField {
	desc := false
	f.descending = &desc
	return f
}

// Asc returns a new UUIDField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f UUIDField) Desc() UUIDField {
	desc := true
	f.descending = &desc
	return f
}

// String returns the string representation of the UUIDField
func (f UUIDField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil, nil)
	return questionInterpolate(buf.String(), args...)
}

// GetAlias retusn the alias of the UUIDField
func (f UUIDField) GetAlias() string {
	return f.alias
}

// GetName returns the name of the UUIDField
func (f UUIDField) GetName() string {
	return f.name
}
