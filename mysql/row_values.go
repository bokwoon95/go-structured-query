package sq

// RowValues represents a list of RowValues i.e. (a, b, c...), (d, e, f...),
// (g, h, i...)
type RowValues []RowValue

// AppendSQL marshals the RowValues into a buffer and an args slice.
func (rs RowValues) AppendSQL(buf Buffer, args *[]interface{}) {
	for i := range rs {
		if i > 0 {
			buf.WriteString(", ")
		}
		rs[i].AppendSQL(buf, args)
	}
}

// RowValue represents an SQL Row Value Expression i.e. (a, b, c...)
type RowValue []interface{}

// AppendSQL marshals the RowValue into a buffer and an args slice.
func (r RowValue) AppendSQL(buf Buffer, args *[]interface{}) {
	r.AppendSQLExclude(buf, args, nil)
}

// AppendSQLExclude marshals the RowValue into a buffer and an args slice. It
// propagates the excludedTableQualifiers down to its child elements.
func (r RowValue) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	buf.WriteString("(")
	for i := range r {
		if i > 0 {
			buf.WriteString(", ")
		}
		AppendSQLValue(buf, args, excludedTableQualifiers, r[i])
	}
	buf.WriteString(")")
}

// In returns an 'X IN (Y)' Predicate.
func (r RowValue) In(v interface{}) CustomPredicate {
	var format string
	switch v.(type) {
	case RowValue:
		format = "? IN ?"
	default:
		format = "? IN (?)"
	}
	return CustomPredicate{
		Format: format,
		Values: []interface{}{r, v},
	}
}

// CustomAssignment is an Assignment that can render itself in an arbitrary way by calling
// ExpandValues on its Format and Values.
type CustomAssignment struct {
	Format string
	Values []interface{}
}

// AppendSQLExclude marshals the CustomAssignment into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (set CustomAssignment) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	ExpandValues(buf, args, excludedTableQualifiers, set.Format, set.Values)
}

// AssertAssignment implements the Assignment interface.
func (set CustomAssignment) AssertAssignment() {}

// Set returns an Assignment assigning v to the RowValue.
func (r RowValue) Set(v interface{}) CustomAssignment {
	var format string
	switch v.(type) {
	case RowValue:
		format = "? = ?"
	default:
		format = "? = (?)"
	}
	return CustomAssignment{
		Format: format,
		Values: []interface{}{r, v},
	}
}
