package sq

import "strings"

// RowValues represents a list of RowValues i.e. (a, b, c...), (d, e, f...),
// (g, h, i...)
type RowValues []RowValue

// AppendSQL marshals the RowValues into a buffer and an args slice.
func (rs RowValues) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	for i, rowvalue := range rs {
		if i > 0 {
			buf.WriteString(", ")
		}
		rowvalue.AppendSQL(buf, args)
	}
}

// RowValue represents an SQL Row Value Expression i.e. (a, b, c...)
type RowValue []interface{}

// AppendSQL marshals the RowValue into a buffer and an args slice.
func (r RowValue) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	r.AppendSQLExclude(buf, args, nil)
}

// AppendSQLExclude marshals the RowValue into a buffer and an args slice. It
// propagates the excludedTableQualifiers down to its child elements.
func (r RowValue) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	buf.WriteString("(")
	for i, value := range r {
		if i > 0 {
			buf.WriteString(", ")
		}
		appendSQLValue(buf, args, excludedTableQualifiers, value)
	}
	buf.WriteString(")")
}

// In returns an 'X IN (Y)' Predicate.
func (r RowValue) In(v interface{}) CustomPredicate {
	switch v.(type) {
	case RowValue:
		return CustomPredicate{
			Format: "? IN ?",
			Values: []interface{}{r, v},
		}
	default:
		return CustomPredicate{
			Format: "? IN (?)",
			Values: []interface{}{r, v},
		}
	}
}

// CustomAssignment is an Assignment that can render itself in an arbitrary way by calling
// expandValues on its Format and Values.
type CustomAssignment struct {
	Format string
	Values []interface{}
}

// AppendSQLExclude marshals the CustomAssignment into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (set CustomAssignment) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	expandValues(buf, args, excludedTableQualifiers, set.Format, set.Values)
}

// AssertAssignment implements the Assignment interface.
func (set CustomAssignment) AssertAssignment() {}

// Set returns an Assignment assigning v to the RowValue.
func (r RowValue) Set(v interface{}) CustomAssignment {
	switch v.(type) {
	case RowValue:
		return CustomAssignment{
			Format: "? = ?",
			Values: []interface{}{r, v},
		}
	default:
		return CustomAssignment{
			Format: "? = (?)",
			Values: []interface{}{r, v},
		}
	}
}
