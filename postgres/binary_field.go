package sq

import "strings"

// BinaryField either represents a BYTEA column or a literal []byte value.
type BinaryField struct {
	// BinaryField will be one of the following:

	// 1) Literal []byte value
	value *[]byte

	// 2) BYTEA column
	alias string
	table Table
	name  string
}

// AppendSQLExclude marshals the BinaryField into a buffer and an args slice. It
// will not table qualify itself if its table qualifer appears in the
// excludedTableQualifiers list.
func (f BinaryField) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	switch {
	case f.value != nil:
		// 1) Literal []byte value
		buf.WriteString("?")
		*args = append(*args, *f.value)
	default:
		// 2) BYTEA column
		tableQualifier := f.table.GetAlias()
		if tableQualifier == "" {
			tableQualifier = f.table.GetName()
		}
		for i := range excludedTableQualifiers {
			if tableQualifier == excludedTableQualifiers[i] {
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
}

// NewBinaryField returns a new BinaryField representing a BYTEA column.
func NewBinaryField(name string, table Table) BinaryField {
	return BinaryField{
		name:  name,
		table: table,
	}
}

// Bytes returns a new BinaryField representing a literal []byte value.
func Bytes(b []byte) BinaryField {
	return BinaryField{
		value: &b,
	}
}

// Set returns a FieldAssignment associating the BinaryField to the value i.e.
// 'field = value'.
func (f BinaryField) Set(v interface{}) FieldAssignment {
	switch v := v.(type) {
	case []byte:
		return FieldAssignment{
			Field: f,
			Value: Bytes(v),
		}
	default:
		return FieldAssignment{
			Field: f,
			Value: v,
		}
	}
}

// SetBytes returns a FieldAssignment associating the BinaryField to the int
// value i.e. 'field = value'.
func (f BinaryField) SetBytes(b []byte) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: Bytes(b),
	}
}

// IsNull returns an 'X IS NULL' Predicate.
func (f BinaryField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f BinaryField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// GetAlias implements the Field interface. It returns the Alias of the
// BinaryField.
func (f BinaryField) GetAlias() string {
	return f.alias
}

// GetName implements the Field interface. It returns the Name of the
// BinaryField.
func (f BinaryField) GetName() string {
	return f.name
}
