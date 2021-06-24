package sqgen

import (
	"testing"

	"github.com/matryer/is"
)

func TestExport(t *testing.T) {
	tt := []struct{
		name string
		s string
		result string
	}{
		{
			name: "can remove prefix",
			s: "_VALUE_",
			result: "VALUE_",
		},
		{
			name: "removes all spaces",
			s: "VALUE TO EXPORT",
			result: "VALUE_TO_EXPORT",
		},
		{
			name: "uppercases",
			s: "value",
			result: "VALUE",
		},
		{
			name: "all together",
			s: "_value to export_",
			result: "VALUE_TO_EXPORT_",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(Export(tc.s), tc.result)
		})
	}
}

func TestQuoteSpace(t *testing.T) {
	tt := []struct{
		name string
		s string
		result string
	}{
		{
			name: "no spaces",
			s: "no_spaces_included",
			result: "no_spaces_included",
		},
		{
			name: "has spaces",
			s: "some spaces included",
			result: `"some spaces included"`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(QuoteSpace(tc.s), tc.result)
		})
	}
}
