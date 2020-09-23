package sq

import "strings"

// RowValues represents a list of RowValues (a, b, c...), (d, e, f...), (g, h, i...)
type RowValues []RowValue

// AppendSQL will write the VALUES clause into the buffer and args as described
// in the RowValues description. If there are no values it will not write
// anything into the buffer. It returns a flag indicating whether anything was
// written into the buffer.
func (rs RowValues) AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int) {
	for i, rowvalue := range rs {
		if i > 0 {
			buf.WriteString(", ")
		}
		rowvalue.AppendSQL(buf, args, nil)
	}
}

// RowValue represents an SQL Row Value Expression i.e. (a, b, c...)
type RowValue []interface{}

// AppendSQL marshals the RowValue into a buffer and an args slice.
func (r RowValue) AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int) {
	r.AppendSQLExclude(buf, args, nil, nil)
}

// AppendSQLExclude marshals the RowValue into a buffer and an args slice. It
// propagates the excludedTableQualifiers down to its child elements.
func (r RowValue) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
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
	switch v := v.(type) {
	case RowValue:
		return CustomPredicate{
			Format: "? IN ?",
			Values: []interface{}{r, v},
		}
	case Query:
		return CustomPredicate{
			Format: "? IN (?)",
			Values: []interface{}{r, v.NestThis()},
		}
	default:
		return CustomPredicate{
			Format: "? IN (?)",
			Values: []interface{}{r, v},
		}
	}
}

// GetName implements the Field interface.
func (r RowValue) GetName() string {
	return ""
}

// GetAlias implements the Field interface.
func (r RowValue) GetAlias() string {
	return ""
}

// CustomAssignment is an Assignment that can render itself in an arbitrary way by calling
// expandValues on its Format and Values.
type CustomAssignment struct {
	Format string
	Values []interface{}
}

// AppendSQLExclude marshals the CustomAssignment into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (set CustomAssignment) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	expandValues(buf, args, excludedTableQualifiers, set.Format, set.Values)
}

// AssertAssignment implements the Assignment interface.
func (set CustomAssignment) AssertAssignment() {}

// Set returns an Assignment assigning v to the RowValue.
func (r RowValue) Set(v interface{}) CustomAssignment {
	switch v := v.(type) {
	case RowValue:
		return CustomAssignment{
			Format: "? = ?",
			Values: []interface{}{r, v},
		}
	case Query:
		return CustomAssignment{
			Format: "? = (?)",
			Values: []interface{}{r, v.NestThis()},
		}
	default:
		return CustomAssignment{
			Format: "? = (?)",
			Values: []interface{}{r, v},
		}
	}
}
