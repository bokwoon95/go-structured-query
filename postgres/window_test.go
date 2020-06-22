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
			OrderBy(ur.ROLE, ur.COHORT.Desc().NullsFirst()),
			"(ORDER BY ur.role, ur.cohort DESC NULLS FIRST)",
			nil,
		},
		{
			"full",
			PartitionBy(ur.USER_ID).OrderBy(ur.ROLE, ur.COHORT.Desc().NullsFirst()).Frame("UNBOUNDED PRECEDING"),
			"(PARTITION BY ur.user_id ORDER BY ur.role, ur.cohort DESC NULLS FIRST UNBOUNDED PRECEDING)",
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
