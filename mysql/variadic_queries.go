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
	nested   bool
	topLevel bool
	Operator VariadicQueryOperator
	Queries  []Query
	// DB
	DB          DB
	Mapper      func(*Row)
	Accumulator func()
	// Logging
	Log     Logger
	LogFlag LogFlag
	logSkip int
}

func (vq VariadicQuery) ToSQL() (string, []interface{}) {
	vq.logSkip += 1
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
			q.topLevel = true
			q.NestThis().AppendSQL(buf, args)
		default:
			q.NestThis().AppendSQL(buf, args)
		}
	default:
		if !vq.topLevel {
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
				q.topLevel = false
				q.NestThis().AppendSQL(buf, args)
			default:
				q.NestThis().AppendSQL(buf, args)
			}
		}
		if !vq.topLevel {
			buf.WriteString(")")
		}
	}
	if !vq.nested {
		if vq.Log != nil {
			query := buf.String()
			var logOutput string
			switch {
			case Lstats&vq.LogFlag != 0:
				logOutput = "\n----[ Executing query ]----\n" + buf.String() + " " + fmt.Sprint(*args) +
					"\n----[ with bind values ]----\n" + questionInterpolate(query, *args...)
			case Linterpolate&vq.LogFlag != 0:
				logOutput = questionInterpolate(query, *args...)
			default:
				logOutput = buf.String() + " " + fmt.Sprint(*args)
			}
			switch vq.Log.(type) {
			case *log.Logger:
				_ = vq.Log.Output(vq.logSkip+2, logOutput)
			default:
				_ = vq.Log.Output(vq.logSkip+1, logOutput)
			}
		}
	}
}

func (vq VariadicQuery) NestThis() Query {
	vq.nested = true
	return vq
}

func Union(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryUnion,
		Queries:  queries,
	}
}

func UnionAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryUnionAll,
		Queries:  queries,
	}
}

func Intersect(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryIntersect,
		Queries:  queries,
	}
}

func IntersectAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryIntersectAll,
		Queries:  queries,
	}
}

func Except(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryExcept,
		Queries:  queries,
	}
}

func ExceptAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryExceptAll,
		Queries:  queries,
	}
}
