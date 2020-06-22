package sq

import (
	"fmt"
	"strings"
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
	const val = "lorem ipsum"
	q := Queryf(val)
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
				Queries: []Query{nil, nil, nil},
			},
			"NULL UNION NULL UNION NULL",
			nil,
		},
		{
			"one query",
			VariadicQuery{
				Queries: []Query{q},
			},
			val,
			nil,
		},
		{
			"multiple queries",
			VariadicQuery{
				Queries: []Query{q, q, q},
			},
			fmt.Sprintf("%[1]s UNION %[1]s UNION %[1]s", val),
			nil,
		},
		{
			"multiple queries (explicit UNION)",
			VariadicQuery{
				Operator: QueryUnion,
				Queries:  []Query{q, q, q},
			},
			fmt.Sprintf("%[1]s UNION %[1]s UNION %[1]s", val),
			nil,
		},
		{
			"multiple queries (explicit UNION ALL)",
			VariadicQuery{
				Operator: QueryUnionAll,
				Queries:  []Query{q, q, q},
			},
			fmt.Sprintf("%[1]s UNION ALL %[1]s UNION ALL %[1]s", val),
			nil,
		},
		{
			"multiple queries (nested)",
			VariadicQuery{
				Operator: QueryUnionAll,
				Queries: []Query{
					VariadicQuery{
						Operator: QueryUnion,
						Queries:  []Query{q, q},
					},
					VariadicQuery{
						Operator: QueryUnion,
						Queries:  []Query{q, q},
					},
				},
			},
			fmt.Sprintf("(%[1]s UNION %[1]s) UNION ALL (%[1]s UNION %[1]s)", val),
			nil,
		},
		{
			"nested variadic query",
			VariadicQuery{
				Queries: []Query{q, q, q},
				Nested:  true,
			},
			fmt.Sprintf("(%[1]s UNION %[1]s UNION %[1]s)", val),
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			var _ Query = tt.q
			tt.q.AppendSQL(buf, &args)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestVariadicQueries_BasicTesting(t *testing.T) {
	is := is.New(t)
	q := Queryf("lorem ipsum")
	vq := Union(q, q, q).As("union_query")
	_ = vq.NestThis()
	_ = UnionAll(q, q, q)
	_ = Intersect(q, q, q)
	_ = IntersectAll(q, q, q)
	_ = Except(q, q, q)
	_ = ExceptAll(q, q, q)
	// ToSQL
	query, _ := vq.ToSQL()
	is.Equal("lorem ipsum UNION lorem ipsum UNION lorem ipsum", query)
	// GetAlias
	is.Equal("union_query", vq.GetAlias())
	// GetName
	is.Equal("", vq.GetName())
	// Get
	f := vq.Get("some_column")
	buf := &strings.Builder{}
	f.AppendSQLExclude(buf, nil, nil)
	is.Equal("union_query.some_column", buf.String())
}
