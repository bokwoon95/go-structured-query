package sq

import "strings"

type VariadicQueryOperator string

const (
	QueryUnion        VariadicQueryOperator = "UNION"
	QueryUnionAll     VariadicQueryOperator = "UNION ALL"
	QueryIntersect    VariadicQueryOperator = "INTERSECT"
	QueryIntersectAll VariadicQueryOperator = "INTERSECT ALL"
	QueryExcept       VariadicQueryOperator = "EXCEPT"
	QueryExceptAll    VariadicQueryOperator = "EXCEPT ALL"
)

type VariadicQuery struct {
	Nested   bool
	Alias    string
	Operator VariadicQueryOperator
	Queries  []Query
}

func (q VariadicQuery) ToSQL() (string, []interface{}) {
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

func (q VariadicQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
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

func (q VariadicQuery) As(alias string) VariadicQuery {
	q.Alias = alias
	return q
}

func (q VariadicQuery) Get(fieldName string) CustomField {
	return CustomField{
		Format: q.Alias + "." + fieldName,
	}
}

func (q VariadicQuery) GetAlias() string {
	return q.Alias
}

func (q VariadicQuery) GetName() string {
	return ""
}

func (q VariadicQuery) NestThis() Query {
	q.Nested = true
	return q
}

func Union(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryUnion,
		Queries:  queries,
	}
}

func UnionAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryUnionAll,
		Queries:  queries,
	}
}

func Intersect(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryIntersect,
		Queries:  queries,
	}
}

func IntersectAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryIntersectAll,
		Queries:  queries,
	}
}

func Except(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryExcept,
		Queries:  queries,
	}
}

func ExceptAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		Operator: QueryExceptAll,
		Queries:  queries,
	}
}
