package sq

import (
	"fmt"
	"log"
	"strings"
)

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
	TopLevel bool
	Operator VariadicQueryOperator
	Queries  []Query
	// DB
	DB          DB
	Mapper      func(*Row)
	Accumulator func()
	// Logging
	Log     Logger
	LogFlag LogFlag
	LogSkip int
}

func (vq VariadicQuery) ToSQL() (string, []interface{}) {
	vq.LogSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	vq.AppendSQL(buf, &args)
	return buf.String(), args
}

func (vq VariadicQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	if vq.Operator == "" {
		vq.Operator = QueryUnion
	}
	switch len(vq.Queries) {
	case 0:
		break
	case 1:
		switch q := vq.Queries[0].(type) {
		case nil:
			buf.WriteString("NULL")
		case VariadicQuery:
			q.TopLevel = true
			q.NestThis().AppendSQL(buf, args)
		default:
			q.NestThis().AppendSQL(buf, args)
		}
	default:
		if !vq.TopLevel {
			buf.WriteString("(")
		}
		for i, q := range vq.Queries {
			if i > 0 {
				buf.WriteString(" ")
				buf.WriteString(string(vq.Operator))
				buf.WriteString(" ")
			}
			switch q := q.(type) {
			case nil:
				buf.WriteString("NULL")
			case VariadicQuery:
				q.TopLevel = false
				q.NestThis().AppendSQL(buf, args)
			default:
				q.NestThis().AppendSQL(buf, args)
			}
		}
		if !vq.TopLevel {
			buf.WriteString(")")
		}
	}
	if !vq.Nested {
		if vq.Log != nil {
			query := buf.String()
			var logOutput string
			switch {
			case Lstats&vq.LogFlag != 0:
				logOutput = "\n----[ Executing query ]----\n" + buf.String() + " " + fmt.Sprint(*args) +
					"\n----[ with bind values ]----\n" + QuestionInterpolate(query, *args...)
			case Linterpolate&vq.LogFlag != 0:
				logOutput = QuestionInterpolate(query, *args...)
			default:
				logOutput = buf.String() + " " + fmt.Sprint(*args)
			}
			switch vq.Log.(type) {
			case *log.Logger:
				vq.Log.Output(vq.LogSkip+2, logOutput)
			default:
				vq.Log.Output(vq.LogSkip+1, logOutput)
			}
		}
	}
}

func (vq VariadicQuery) NestThis() Query {
	vq.Nested = true
	return vq
}

func Union(queries ...Query) VariadicQuery {
	return VariadicQuery{
		TopLevel: true,
		Operator: QueryUnion,
		Queries:  queries,
	}
}

func UnionAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		TopLevel: true,
		Operator: QueryUnionAll,
		Queries:  queries,
	}
}

func Intersect(queries ...Query) VariadicQuery {
	return VariadicQuery{
		TopLevel: true,
		Operator: QueryIntersect,
		Queries:  queries,
	}
}

func IntersectAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		TopLevel: true,
		Operator: QueryIntersectAll,
		Queries:  queries,
	}
}

func Except(queries ...Query) VariadicQuery {
	return VariadicQuery{
		TopLevel: true,
		Operator: QueryExcept,
		Queries:  queries,
	}
}

func ExceptAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		TopLevel: true,
		Operator: QueryExceptAll,
		Queries:  queries,
	}
}
