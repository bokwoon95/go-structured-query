package sq

import "strings"

// VariadicQueryOperator is an operator that can join a variadic number of
// queries together.
type VariadicQueryOperator string

// VariadicQueryOperators
const (
	QueryUnion        VariadicQueryOperator = "UNION"
	QueryUnionAll     VariadicQueryOperator = "UNION ALL"
	QueryIntersect    VariadicQueryOperator = "INTERSECT"
	QueryIntersectAll VariadicQueryOperator = "INTERSECT ALL"
	QueryExcept       VariadicQueryOperator = "EXCEPT"
	QueryExceptAll    VariadicQueryOperator = "EXCEPT ALL"
)

// VariadicQuery represents a variadic number of queries joined together by an
// VariadicQueryOperator.
type VariadicQuery struct {
	Nested   bool
	Alias    string
	Operator VariadicQueryOperator
	Queries  []Query
}

// ToSQL marshals the VariadicQuery into a query string and args slice.
func (q VariadicQuery) ToSQL() (string, []interface{}) {
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

// AppendSQL marshals the VariadicQuery into a buffer and args slice.
func (q VariadicQuery) AppendSQL(buf Buffer, args *[]interface{}) {
	if q.Operator == "" {
		q.Operator = QueryUnion
	}
	switch len(q.Queries) {
	case 0:
		break
	case 1:
		q.Queries[0].AppendSQL(buf, args)
	default:
		if q.Nested {
			buf.WriteString("(")
		}
		for i := range q.Queries {
			if i > 0 {
				buf.WriteString(" ")
				buf.WriteString(string(q.Operator))
				buf.WriteString(" ")
			}
			switch v := q.Queries[i].(type) {
			case nil:
				buf.WriteString("NULL")
			case VariadicQuery:
				v.Nested = true
				v.AppendSQL(buf, args)
			default:
				v.AppendSQL(buf, args)
			}
		}
		if q.Nested {
			buf.WriteString(")")
		}
	}
}

// As aliases the VariadicQuery i.e. 'query AS alias'.
func (q VariadicQuery) As(alias string) VariadicQuery {
	q.Alias = alias
	return q
}

// Get returns a Field from the VariadicQuery, identified by fieldName.
func (q VariadicQuery) Get(fieldName string) CustomField {
	return CustomField{
		Format: q.Alias + "." + fieldName,
	}
}

// GetAlias returns the alias of the VariadicQuery.
func (q VariadicQuery) GetAlias() string {
	return q.Alias
}

// GetName returns the name of the VariadicQuery, which is always an empty
// string.
func (q VariadicQuery) GetName() string {
	return ""
}

// NestThis indicates to the VariadicQuery that it is nested.
func (q VariadicQuery) NestThis() Query {
	q.Nested = true
	return q
}

// Union joins the list of queries together by the UNION operator.
func Union(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryUnion,
		Queries:  queries,
	}
}

// UnionAll joins the list of queries together by the UNION ALL operator.
func UnionAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryUnionAll,
		Queries:  queries,
	}
}

// Intersect joins the list of queries together by the INTERSECT operator.
func Intersect(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryIntersect,
		Queries:  queries,
	}
}

// IntersectAll joins the list of queries together by the INTERSECT ALL operator.
func IntersectAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryIntersectAll,
		Queries:  queries,
	}
}

// Except joins the list of queries together by the EXCEPT operator.
func Except(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryExcept,
		Queries:  queries,
	}
}

// ExceptAll joins the list of queries together by the EXCEPT ALL operator.
func ExceptAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryExceptAll,
		Queries:  queries,
	}
}
