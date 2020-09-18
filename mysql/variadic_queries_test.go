package sq

import (
	"testing"

	"github.com/matryer/is"
)

func TestVariadicQueries_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		q           VariadicQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	q1 := Select(Int(1))
	q2 := Select(Int(2))
	q3 := Select(Int(3))
	tests := []TT{
		{
			"no queries",
			VariadicQuery{},
			"",
			nil,
		},
		{
			"nil queries",
			VariadicQuery{
				topLevel: true,
				Queries:  []Query{nil, nil, nil},
			},
			"NULL UNION NULL UNION NULL",
			nil,
		},
		{
			"one query",
			WithLog(customLogger, 0).Union(q1),
			"SELECT ?",
			[]interface{}{1},
		},
		{
			"multiple queries (implicit UNION)",
			VariadicQuery{
				topLevel: true,
				Queries:  []Query{q1, q2, q3},
			},
			"SELECT ? UNION SELECT ? UNION SELECT ?",
			[]interface{}{1, 2, 3},
		},
		{
			"multiple queries (explicit UNION)",
			WithDefaultLog(Lstats).Union(q1, q2, q3),
			"SELECT ? UNION SELECT ? UNION SELECT ?",
			[]interface{}{1, 2, 3},
		},
		{
			"multiple queries (explicit UNION ALL)",
			WithDefaultLog(Linterpolate).UnionAll(q1, q2, q3),
			"SELECT ? UNION ALL SELECT ? UNION ALL SELECT ?",
			[]interface{}{1, 2, 3},
		},
		{
			"variadic query containing multiple variadic queries (toplevel)",
			VariadicQuery{
				topLevel: true,
				Operator: QueryUnionAll,
				Queries: []Query{
					VariadicQuery{
						Operator: QueryUnion,
						Queries:  []Query{q1, q2},
					},
					VariadicQuery{
						Operator: QueryUnion,
						Queries:  []Query{q2, q3},
					},
				},
			},
			"(SELECT ? UNION SELECT ?) UNION ALL (SELECT ? UNION SELECT ?)",
			[]interface{}{1, 2, 2, 3},
		},
		{
			"variadic query containing one variadic query (toplevel)",
			VariadicQuery{
				topLevel: true,
				Queries: []Query{
					VariadicQuery{
						Operator: QueryUnion,
						Queries:  []Query{q1, q2, q3},
					},
				},
			},
			"SELECT ? UNION SELECT ? UNION SELECT ?",
			[]interface{}{1, 2, 3},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			var _ Query = tt.q
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestVariadicQueries_Basic(t *testing.T) {
	is := is.New(t)
	q1 := Select(Int(1))
	q2 := Select(Int(2))
	q3 := Select(Int(3))
	var vq VariadicQuery

	vq = Intersect(q1, q2, q3)
	is.Equal(true, vq.topLevel)
	is.Equal(QueryIntersect, vq.Operator)

	vq = IntersectAll(q1, q2, q3)
	is.Equal(true, vq.topLevel)
	is.Equal(QueryIntersectAll, vq.Operator)

	vq = Except(q1, q2, q3)
	is.Equal(true, vq.topLevel)
	is.Equal(QueryExcept, vq.Operator)

	vq = ExceptAll(q1, q2, q3)
	is.Equal(true, vq.topLevel)
	is.Equal(QueryExceptAll, vq.Operator)
}
