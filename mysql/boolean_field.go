package sq

import "strings"

// BooleanField either represents a boolean column or a literal bool value.
type BooleanField struct {
	// BooleanField will be one of the following:

	// 1) Literal bool value
	// Examples of literal bool values:
	// | query | args |
	// |-------|------|
	// | ?     | true |
	value *bool

	// 2) Boolean column
	// Examples of boolean columns:
	// | query            | args |
	// |------------------|------|
	// | users.is_created |      |
	// | is_created       |      |
	alias      string
	table      Table
	name       string
	descending *bool
	negative   bool
}

// AppendSQLExclude marshals the BooleanField into a buffer and an args slice. It
// will not table qualify itself if its table qualifer appears in the
// excludedTableQualifiers list.
func (f BooleanField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	if f.negative {
		buf.WriteString("NOT ")
	}
	switch {
	case f.value != nil:
		// 1) Literal bool value
		buf.WriteString("?")
		*args = append(*args, *f.value)
	default:
		// 2) Boolean column
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
				buf.WriteString("`")
				buf.WriteString(tableQualifier)
				buf.WriteString("`.")
			} else {
				buf.WriteString(tableQualifier)
				buf.WriteString(".")
			}
		}
		if strings.ContainsAny(f.name, " \t") {
			buf.WriteString("`")
			buf.WriteString(f.name)
			buf.WriteString("`")
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

// NewBooleanField returns a new BooleanField representing a boolean column.
func NewBooleanField(name string, table Table) BooleanField {
	return BooleanField{
		name:  name,
		table: table,
	}
}

// Bool returns a new Boolean Field representing a literal bool value.
func Bool(b bool) BooleanField {
	return BooleanField{
		value: &b,
	}
}

// Set returns a FieldAssignment associating the BooleanField to the value i.e.
// 'field = value'.
func (f BooleanField) Set(val interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: val,
	}
}

// SetBool returns a FieldAssignment associating the BooleanField to the bool
// value i.e. 'field = value'.
func (f BooleanField) SetBool(val bool) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: val,
	}
}

// As aliases the BooleanField i.e. 'field AS Alias'.
func (f BooleanField) As(alias string) BooleanField {
	f.alias = alias
	return f
}

// Asc returns a new BooleanField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f BooleanField) Asc() BooleanField {
	desc := false
	f.descending = &desc
	return f
}

// Desc returns a new BooleanField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f BooleanField) Desc() BooleanField {
	desc := true
	f.descending = &desc
	return f
}

// IsNull returns an 'X IS NULL' Predicate.
func (f BooleanField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f BooleanField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// Eq returns an 'X = Y' Predicate. It only accepts BooleanField.
func (f BooleanField) Eq(field BooleanField) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, field},
	}
}

// Ne returns an 'X <> Y' Predicate. It only accepts BooleanField.
func (f BooleanField) Ne(field BooleanField) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, field},
	}
}

// String returns the string representation of the BooleanField.
func (f BooleanField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil)
	return QuestionInterpolate(buf.String(), args...)
}

// GetAlias returns the alias of the BooleanField.
func (f BooleanField) GetAlias() string {
	return f.alias
}

// GetName returns the name of the BooleanField.
func (f BooleanField) GetName() string {
	return f.name
}

// Not inverts the BooleanField i.e. 'NOT BooleanField'.
func (f BooleanField) Not() Predicate {
	f.negative = !f.negative
	return f
}
