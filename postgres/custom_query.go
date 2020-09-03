package sq

import "strings"

// CustomQuery is a Query that can render itself in an arbitrary way as defined
// by its Format string. Values are interpolated into the Format string as
// described in the (CustomQuery).CustomSprintf function.
//
// The difference between CustomTable and CustomQuery is that CustomTable is
// not meant for writing full queries, because it does not do any form of
// placeholder ?, ?, ? -> $1, $2, $3 etc rebinding.
type CustomQuery struct {
	Nested bool
	Alias  string
	Format string
	Values []interface{}
}

// ToSQL implements the Query interface.
func (q CustomQuery) ToSQL() (string, []interface{}) {
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

// AppendSQL marshals the CustomQuery into an SQL query.
func (q CustomQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	ExpandValues(buf, args, nil, q.Format, q.Values)
}

func Queryf(format string, values ...interface{}) CustomQuery {
	return CustomQuery{
		Format: format,
		Values: values,
	}
}

// As returns a new CustomQuery with the new alias i.e. 'query AS Alias'.
func (q CustomQuery) As(alias string) CustomQuery {
	q.Alias = alias
	return q
}

// GetAlias implements the Table interface. It returns the alias of the
// CustomQuery.
func (q CustomQuery) GetAlias() string {
	return q.Alias
}

// GetName implements the Table interface. It returns the name of the
// CustomQuery.
func (q CustomQuery) GetName() string {
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String()
}

// NestThis implements the Query interface.
func (q CustomQuery) NestThis() Query {
	q.Nested = true
	return q
}
