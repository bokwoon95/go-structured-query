package sq

import "strings"

// FieldLiteral is a Field where its underlying string is literally plugged
// into the SQL query.
type FieldLiteral string

// ToSQL returns the underlying string of the FieldLiteral.
func (f FieldLiteral) AppendSQLExclude(buf *strings.Builder, _ *[]interface{}, _ []string) {
	buf.WriteString(string(f))
}

// GetAlias implements the Field interface. It always returns an empty string
// because FieldLiterals do not have aliases.
func (f FieldLiteral) GetAlias() string {
	return ""
}

// GetName implements the Field interface. It returns the FieldLiteral's
// underlying string as the name.
func (f FieldLiteral) GetName() string {
	return string(f)
}

// Fields represents the "field1, field2, etc..." SQL construct.
type Fields []Field

// AppendSQLExclude will write the a slice of Fields into the buffer and args as
// described in the Fields description. The list of table qualifiers to be
// excluded is propagated down to the individual Fields.
func (fs Fields) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	for i, field := range fs {
		if i > 0 {
			buf.WriteString(", ")
		}
		if field == nil {
			buf.WriteString("NULL")
		} else {
			field.AppendSQLExclude(buf, args, excludedTableQualifiers)
		}
	}
}

// AppendSQLExcludeWithAlias is exactly like AppendSQLExclude, but appends each
// field (i.e.  field1 AS alias1, field2 AS alias2, ...) with its alias if it
// has one.
func (fs Fields) AppendSQLExcludeWithAlias(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	var alias string
	for i, field := range fs {
		if i > 0 {
			buf.WriteString(", ")
		}
		if field == nil {
			buf.WriteString("NULL")
		} else {
			field.AppendSQLExclude(buf, args, excludedTableQualifiers)
			if alias = field.GetAlias(); alias != "" {
				buf.WriteString(" AS ")
				buf.WriteString(alias)
			}
		}
	}
}

// FieldAssignment represents a Field and Value set. Its usage appears in both
// the UPDATE and INSERT queries whenever values are assigned to columns e.g.
// 'field = value'.
type FieldAssignment struct {
	Field Field
	Value interface{}
}

// AppendSQLExclude will write the FieldAssignment into the buffer and args as
// described in the Assignments description.
func (set FieldAssignment) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	appendSQLValue(buf, args, excludedTableQualifiers, set.Field)
	buf.WriteString(" = ")
	switch v := set.Value.(type) {
	case Query:
		buf.WriteString("(")
		appendSQLValue(buf, args, excludedTableQualifiers, v.NestThis())
		buf.WriteString(")")
	default:
		appendSQLValue(buf, args, excludedTableQualifiers, set.Value)
	}
}

func (set FieldAssignment) AssertAssignment() {}

// Assignments is a list of Assignments, when translated to SQL it looks
// something like "SET field1 = value1, field2 = value2, etc...".
type Assignments []Assignment

// AppendSQLExclude will write the Assignments into the buffer and args as
// described in the Assignments description.
func (assignments Assignments) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	for i, assignment := range assignments {
		if i > 0 {
			buf.WriteString(", ")
		}
		assignment.AppendSQLExclude(buf, args, excludedTableQualifiers)
	}
}
