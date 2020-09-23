package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestAggregateFunctions(t *testing.T) {
	type TT struct {
		description string
		f           Field
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	ur := USER_ROLES().As("ur")
	tests := []TT{
		{
			"Count",
			Count(),
			nil,
			"COUNT(*)",
			nil,
		},
		{
			"CountOver",
			CountOver(Window{}),
			nil,
			"COUNT(*) OVER ()",
			nil,
		},
		{
			"Sum",
			Sum(ur.USER_ID),
			nil,
			"SUM(ur.user_id)",
			nil,
		},
		{
			"SumOver",
			SumOver(ur.USER_ROLE_ID, PartitionBy(ur.USER_ID)),
			nil,
			"SUM(ur.user_role_id) OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"Avg",
			Avg(ur.USER_ID),
			nil,
			"AVG(ur.user_id)",
			nil,
		},
		{
			"AvgOver",
			AvgOver(ur.USER_ROLE_ID, PartitionBy(ur.USER_ID)),
			nil,
			"AVG(ur.user_role_id) OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"Min",
			Min(ur.USER_ROLE_ID),
			nil,
			"MIN(ur.user_role_id)",
			nil,
		},
		{
			"MinOver",
			MinOver(ur.USER_ROLE_ID, PartitionBy(ur.USER_ID)),
			nil,
			"MIN(ur.user_role_id) OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"Max",
			Max(ur.USER_ROLE_ID),
			nil,
			"MAX(ur.user_role_id)",
			nil,
		},
		{
			"MaxOver",
			MaxOver(ur.USER_ROLE_ID, PartitionBy(ur.USER_ID)),
			nil,
			"MAX(ur.user_role_id) OVER (PARTITION BY ur.user_id)",
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
			tt.f.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}
