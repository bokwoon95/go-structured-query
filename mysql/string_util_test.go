package sq

import (
	"database/sql/driver"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

type valuer struct {
	a int
	b int
}

func (v valuer) Value() (driver.Value, error) {
	if v.a == 0 && v.b == 0 {
		return nil, nil
	}
	sum := v.a + v.b
	return strconv.Itoa(sum), nil
}

type marshaler struct {
	Lorem string `json:"lorem"`
	Dolor string `json:"dolor"`
	Amet  int    `json:"amet"`
}

func TestInterpolateSQLValue(t *testing.T) {
	type TT struct {
		description string
		value       interface{}
		want        string
	}
	tests := []TT{
		{"nil", nil, "NULL"},
		{"true", true, "TRUE"},
		{"false", false, "FALSE"},
		{"string", "lorem ipsum", "'lorem ipsum'"},
		{"int", 33, "33"},
		func() TT {
			desc := "time"
			now := time.Now()
			return TT{desc, now, "'" + now.Format(time.RFC3339Nano) + "'"}
		}(),
		{"driver.Valuer value", valuer{a: 3, b: 4}, "'7'"},
		{"driver.Valuer nil", valuer{a: 0, b: 0}, "NULL"},
		{
			"jsonable",
			marshaler{Lorem: "ipsum", Dolor: "sit", Amet: 5},
			`'{"lorem":"ipsum","dolor":"sit","amet":5}'`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			InterpolateSQLValue(buf, tt.value)
			is.Equal(tt.want, buf.String())
		})
	}
}
