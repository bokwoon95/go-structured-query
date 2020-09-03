package sq

import "strings"

// EnumField is a type alias for StringField.
type EnumField = StringField

// NewEnumField returns an EnumField representing an enum column.
func NewEnumField(name string, table Table) EnumField {
	return NewStringField(name, table)
}

// StringField either represents a string column or a literal string value.
type StringField struct {
	// StringField will be one of the following:

	// 1) Literal string value
	// Examples of literal string values:
	// | query | args |
	// |-------|------|
	// | ?     | abcd |
	value *string

	// 2) String column
	// Examples of boolean columns:
	// | query       | args |
	// |-------------|------|
	// | users.name  |      |
	// | name        |      |
	// | users.email |      |
	alias      string
	table      Table
	name       string
	descending *bool
}

// AppendSQLExclude marshals the StringField into a buffer and an args slice. It
// will not table qualify itself if its table qualifer appears in the
// excludedTableQualifiers list.
func (f StringField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	switch {
	case f.value != nil:
		// 1) Literal string value
		buf.WriteString("?")
		*args = append(*args, *f.value)
	default:
		// 2) String column
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

// NewStringField returns a new StringField representing a boolean column.
func NewStringField(name string, table Table) StringField {
	return StringField{
		name:  name,
		table: table,
	}
}

// String returns a new StringField representing a literal string value.
func String(s string) StringField {
	return StringField{
		value: &s,
	}
}

// Set returns a FieldAssignment associating the StringField to the value i.e.
// 'field = value'.
func (f StringField) Set(value interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

// SetString returns a FieldAssignment associating the StringField to the string
// value i.e. 'field = value'.
func (f StringField) SetString(s string) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: s,
	}
}

// As returns a new StringField with the new field Alias i.e. 'field AS Alias'.
func (f StringField) As(alias string) StringField {
	f.alias = alias
	return f
}

// Asc returns a new StringField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f StringField) Asc() StringField {
	desc := false
	f.descending = &desc
	return f
}

// Desc returns a new StringField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f StringField) Desc() StringField {
	desc := true
	f.descending = &desc
	return f
}

// IsNull returns an 'X IS NULL' Predicate.
func (f StringField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f StringField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// Eq returns an 'X = Y' Predicate. It only accepts StringField.
func (f StringField) Eq(field StringField) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, field},
	}
}

// Ne returns an 'X <> Y' Predicate. It only accepts StringField.
func (f StringField) Ne(field StringField) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, field},
	}
}

// Gt returns an 'X > Y' Predicate. It only accepts StringField.
func (f StringField) Gt(field StringField) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, field},
	}
}

// Ge returns an 'X >= Y' Predicate. It only accepts StringField.
func (f StringField) Ge(field StringField) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, field},
	}
}

// Lt returns an 'X < Y' Predicate. It only accepts StringField.
func (f StringField) Lt(field StringField) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, field},
	}
}

// Le returns an 'X <= Y' Predicate. It only accepts StringField.
func (f StringField) Le(field StringField) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, field},
	}
}

// EqString returns an 'X = Y' Predicate. It only accepts string.
func (f StringField) EqString(s string) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, s},
	}
}

// NeString returns an 'X <> Y' Predicate. It only accepts string.
func (f StringField) NeString(s string) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, s},
	}
}

// GtString returns an 'X > Y' Predicate. It only accepts string.
func (f StringField) GtString(s string) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, s},
	}
}

// GeString returns an 'X >= Y' Predicate. It only accepts string.
func (f StringField) GeString(s string) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, s},
	}
}

// LtString returns an 'X < Y' Predicate. It only accepts string.
func (f StringField) LtString(s string) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, s},
	}
}

// LeString returns an 'X <= Y' Predicate. It only accepts string.
func (f StringField) LeString(s string) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, s},
	}
}

// LikeString returns an 'A LIKE B' Predicate. It only accepts string.
func (f StringField) LikeString(s string) Predicate {
	return CustomPredicate{
		Format: "? LIKE ?",
		Values: []interface{}{f, s},
	}
}

// NotLikeString returns an 'A NOT LIKE B' Predicate. It only accepts string.
func (f StringField) NotLikeString(s string) Predicate {
	return CustomPredicate{
		Format: "? NOT LIKE ?",
		Values: []interface{}{f, s},
	}
}

// In returns an 'X IN (Y)' Predicate.
func (f StringField) In(v interface{}) Predicate {
	var format string
	switch v.(type) {
	case RowValue:
		format = "? IN ?"
	default:
		format = "? IN (?)"
	}
	return CustomPredicate{
		Format: format,
		Values: []interface{}{f, v},
	}
}

// String returns the string representation of the StringField.
func (f StringField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil)
	return QuestionInterpolate(buf.String(), args...)
}

// GetAlias returns the alias of the StringField.
func (f StringField) GetAlias() string {
	return f.alias
}

// GetName returns the name of the StringField.
func (f StringField) GetName() string {
	return f.name
}
