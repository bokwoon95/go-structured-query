package sq

// FieldLiteral is a Field where its underlying string is literally plugged
// into the SQL query.
type FieldLiteral string

// AppendSQLExclude marshals the FieldLiteral into a buffer and an args slice.
func (f FieldLiteral) AppendSQLExclude(buf Buffer, _ *[]interface{}, _ []string) {
	buf.WriteString(string(f))
}

// GetAlias returns the alias of the FieldLiteral, which is always an empty
// string.
func (f FieldLiteral) GetAlias() string {
	return ""
}

// GetName returns the FieldLiteral's underlying string.
func (f FieldLiteral) GetName() string {
	return string(f)
}

// Fields represents the "field1, field2, etc..." SQL construct.
type Fields []Field

// AppendSQLExclude marshals PredicateCases into a buffer and an args slice. It
// propagates the excludedTableQualifiers down to its Fields.
func (fs Fields) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	for i := range fs {
		if i > 0 {
			buf.WriteString(", ")
		}
		if fs[i] == nil {
			buf.WriteString("NULL")
		} else {
			fs[i].AppendSQLExclude(buf, args, excludedTableQualifiers)
		}
	}
}

// AppendSQLExcludeWithAlias is exactly like AppendSQLExclude, but appends each
// field (i.e. field1 AS alias1, field2 AS alias2, ...) with its alias if it
// has one.
func (fs Fields) AppendSQLExcludeWithAlias(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	var alias string
	for i := range fs {
		if i > 0 {
			buf.WriteString(", ")
		}
		if fs[i] == nil {
			buf.WriteString("NULL")
		} else {
			fs[i].AppendSQLExclude(buf, args, excludedTableQualifiers)
			if alias = fs[i].GetAlias(); alias != "" {
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

// AppendSQLExclude marshals the FieldAssignment into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (set FieldAssignment) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	AppendSQLValue(buf, args, excludedTableQualifiers, set.Field)
	buf.WriteString(" = ")
	AppendSQLValue(buf, args, excludedTableQualifiers, set.Value)
}

// AssertAssignment implements the Assignment interface.
func (set FieldAssignment) AssertAssignment() {}

// Assignments is a list of Assignments.
type Assignments []Assignment

// AppendSQLExclude marshals the Assignments into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (assignments Assignments) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	for i := range assignments {
		if i > 0 {
			buf.WriteString(", ")
		}
		assignments[i].AppendSQLExclude(buf, args, excludedTableQualifiers)
	}
}
