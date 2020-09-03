package sq

import "strings"

// CustomQuery is a Query that can render itself in an arbitrary way by calling
// ExpandValues on its Format and Values.
type CustomQuery struct {
	Nested bool
	Alias  string
	Format string
	Values []interface{}
}

// ToSQL marshals the CustomQuery into a query string and args slice.
func (q CustomQuery) ToSQL() (string, []interface{}) {
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

// AppendSQL marshals CustomQuery into a buffer and an args slice.
func (q CustomQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	ExpandValues(buf, args, nil, q.Format, q.Values)
}

// Queryf creates a new CustomQuery.
func Queryf(format string, values ...interface{}) CustomQuery {
	return CustomQuery{
		Format: format,
		Values: values,
	}
}

// As aliases the CustomQuery i.e. 'query AS Alias'.
func (q CustomQuery) As(alias string) CustomQuery {
	q.Alias = alias
	return q
}

// GetAlias returns the alias of the CustomQuery.
func (q CustomQuery) GetAlias() string {
	return q.Alias
}

// GetName returns the name of the CustomQuery.
func (q CustomQuery) GetName() string {
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String()
}

// NestThis indicates to the Query that it is nested.
func (q CustomQuery) NestThis() Query {
	q.Nested = true
	return q
}
