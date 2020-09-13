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
				TopLevel: true,
				Queries:  []Query{nil, nil, nil},
			},
			"NULL UNION NULL UNION NULL",
			nil,
		},
		{
			"one query",
			WithLog(customLogger, 0).Union(q1),
			"SELECT $1",
			[]interface{}{1},
		},
		{
			"multiple queries (implicit UNION)",
			VariadicQuery{
				TopLevel: true,
				Queries:  []Query{q1, q2, q3},
			},
			"SELECT $1 UNION SELECT $2 UNION SELECT $3",
			[]interface{}{1, 2, 3},
		},
		{
			"multiple queries (explicit UNION)",
			WithDefaultLog(Lstats).Union(q1, q2, q3),
			"SELECT $1 UNION SELECT $2 UNION SELECT $3",
			[]interface{}{1, 2, 3},
		},
		{
			"multiple queries (explicit UNION ALL)",
			WithDefaultLog(Linterpolate).UnionAll(q1, q2, q3),
			"SELECT $1 UNION ALL SELECT $2 UNION ALL SELECT $3",
			[]interface{}{1, 2, 3},
		},
		{
			"variadic query containing multiple variadic queries (toplevel)",
			VariadicQuery{
				TopLevel: true,
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
			"(SELECT $1 UNION SELECT $2) UNION ALL (SELECT $3 UNION SELECT $4)",
			[]interface{}{1, 2, 2, 3},
		},
		{
			"variadic query containing one variadic query (toplevel)",
			VariadicQuery{
				TopLevel: true,
				Queries: []Query{
					VariadicQuery{
						Operator: QueryUnion,
						Queries:  []Query{q1, q2, q3},
					},
				},
			},
			"SELECT $1 UNION SELECT $2 UNION SELECT $3",
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
	is.Equal(true, vq.TopLevel)
	is.Equal(QueryIntersect, vq.Operator)

	vq = IntersectAll(q1, q2, q3)
	is.Equal(true, vq.TopLevel)
	is.Equal(QueryIntersectAll, vq.Operator)

	vq = Except(q1, q2, q3)
	is.Equal(true, vq.TopLevel)
	is.Equal(QueryExcept, vq.Operator)

	vq = ExceptAll(q1, q2, q3)
	is.Equal(true, vq.TopLevel)
	is.Equal(QueryExceptAll, vq.Operator)
}
