package sq

import "strings"

// RowValues represents a list of RowValues (a, b, c...), (d, e, f...), (g, h, i...)
type RowValues []RowValue

// AppendSQL will write the VALUES clause into the buffer and args as described
// in the RowValues description. If there are no values it will not write
// anything into the buffer. It returns a flag indicating whether anything was
// written into the buffer.
func (rs RowValues) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	for i, rowvalue := range rs {
		if i > 0 {
			buf.WriteString(", ")
		}
		rowvalue.AppendSQL(buf, args)
	}
}

type RowValue []interface{}

func (r RowValue) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	r.AppendSQLExclude(buf, args, nil)
}

func (r RowValue) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	buf.WriteString("(")
	for i, value := range r {
		if i > 0 {
			buf.WriteString(", ")
		}
		AppendSQLValue(buf, args, excludedTableQualifiers, value)
	}
	buf.WriteString(")")
}

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

type CustomAssignment struct {
	Format string
	Values []interface{}
}

func (set CustomAssignment) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	ExpandValues(buf, args, excludedTableQualifiers, set.Format, set.Values)
}

func (set CustomAssignment) AssertAssignment() {}

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
