package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestWindow_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		w           Window
		wantQuery   string
		wantArgs    []interface{}
	}
	ur := USER_ROLES().As("ur")
	tests := []TT{
		{
			"empty",
			Window{},
			"()",
			nil,
		},
		{
			"PartitionBy",
			PartitionBy(ur.USER_ID),
			"(PARTITION BY ur.user_id)",
			nil,
		},
		{
			"OrderBy",
			OrderBy(ur.ROLE, ur.COHORT.Desc()),
			"(ORDER BY ur.role, ur.cohort DESC)",
			nil,
		},
		{
			"full",
			PartitionBy(ur.USER_ID).OrderBy(ur.ROLE, ur.COHORT.Desc()).Frame("UNBOUNDED PRECEDING"),
			"(PARTITION BY ur.user_id ORDER BY ur.role, ur.cohort DESC UNBOUNDED PRECEDING)",
			nil,
		},
		{
			"name",
			PartitionBy(ur.USER_ID).OrderBy(ur.ROLE).As("my_window").Name(),
			"my_window",
			nil,
		},
		func() TT {
			desc := "randomly generated name"
			w := Window{}.Name()
			return TT{desc, w, w.WindowName, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.w.AppendSQL(buf, &args)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestWindowFunctions(t *testing.T) {
	type TT struct {
		description string
		f           Field
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	ur := USER_ROLES().As("ur")
	w := PartitionBy(ur.USER_ID)
	tests := []TT{
		{
			"RowNumberOver",
			RowNumberOver(w),
			nil,
			"ROW_NUMBER() OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"RankOver",
			RankOver(w),
			nil,
			"RANK() OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"DenseRankOver",
			DenseRankOver(w),
			nil,
			"DENSE_RANK() OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"PercentRankOver",
			PercentRankOver(w),
			nil,
			"PERCENT_RANK() OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"CumeDistOver",
			CumeDistOver(w),
			nil,
			"CUME_DIST() OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"LeadOver",
			LeadOver(ur.COHORT, nil, nil, w),
			nil,
			"LEAD(ur.cohort, ?, NULL) OVER (PARTITION BY ur.user_id)",
			[]interface{}{1},
		},
		{
			"LagOver",
			LagOver(ur.COHORT, nil, nil, w),
			nil,
			"LAG(ur.cohort, ?, NULL) OVER (PARTITION BY ur.user_id)",
			[]interface{}{1},
		},
		{
			"NtileOver",
			NtileOver(3, w),
			nil,
			"NTILE(?) OVER (PARTITION BY ur.user_id)",
			[]interface{}{3},
		},
		{
			"FirstValueOver",
			FirstValueOver(ur.COHORT, w),
			nil,
			"FIRST_VALUE(ur.cohort) OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"LastValueOver",
			LastValueOver(ur.COHORT, w),
			nil,
			"LAST_VALUE(ur.cohort) OVER (PARTITION BY ur.user_id)",
			nil,
		},
		{
			"NthValueOver",
			NthValueOver(ur.COHORT, 3, w),
			nil,
			"NTH_VALUE(ur.cohort, ?) OVER (PARTITION BY ur.user_id)",
			[]interface{}{3},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.f.AppendSQLExclude(buf, &args, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}
