package sq

import "strings"

// RowValues represents a list of RowValues (a, b, c...), (d, e, f...), (g, h, i...)
type RowValues []RowValue

// AppendSQL will write the VALUES clause into the buffer and args as described
// in the RowValues description. If there are no values it will not write
// anything into the buffer. It returns a flag indicating whether anything was
// written into the buffer.
func (rs RowValues) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	for i := range rs {
		if i > 0 {
			buf.WriteString(", ")
		}
		rs[i].AppendSQL(buf, args)
	}
}

type RowValue []interface{}

func (r RowValue) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	r.AppendSQLExclude(buf, args, nil)
}

func (r RowValue) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	buf.WriteString("(")
	for i := range r {
		if i > 0 {
			buf.WriteString(", ")
		}
		AppendSQLValue(buf, args, excludedTableQualifiers, r[i])
	}
	buf.WriteString(")")
}

func (r RowValue) In(v interface{}) CustomPredicate {
	var format string
	var values []interface{}
	switch v := v.(type) {
	case RowValue:
		format = "? IN ?"
		values = []interface{}{r, v}
	case Query:
		format = "? IN (?)"
		values = []interface{}{r, v.NestThis()}
	default:
		format = "? IN (?)"
		values = []interface{}{r, v}
	}
	return CustomPredicate{
		Format: format,
		Values: values,
	}
}

type CustomAssignment struct {
	Format string
	Values []interface{}
}

func (set CustomAssignment) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	ExpandValues(buf, args, excludedTableQualifiers, set.Format, set.Values)
}

func (set CustomAssignment) AssertAssignment() {}

func (r RowValue) Set(v interface{}) CustomAssignment {
	var format string
	var values []interface{}
	switch v := v.(type) {
	case RowValue:
		format = "? = ?"
		values = []interface{}{r, v}
	case Query:
		format = "? = (?)"
		values = []interface{}{r, v.NestThis()}
	default:
		format = "? = (?)"
		values = []interface{}{r, v}
	}
	return CustomAssignment{
		Format: format,
		Values: values,
	}
}
