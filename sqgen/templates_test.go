package sqgen

import (
	"testing"

	"github.com/matryer/is"
)

func TestExport(t *testing.T) {
	type TT struct {
		name   string
		s      string
		result string
	}
	tests := []TT{
		{
			name:   "can remove prefix",
			s:      "_VALUE_",
			result: "VALUE_",
		},
		{
			name:   "removes all spaces",
			s:      "VALUE TO EXPORT",
			result: "VALUE_TO_EXPORT",
		},
		{
			name:   "uppercases",
			s:      "value",
			result: "VALUE",
		},
		{
			name:   "all together",
			s:      "_value to export_",
			result: "VALUE_TO_EXPORT_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(Export(tt.s), tt.result)
		})
	}
}

func TestQuoteSpace(t *testing.T) {
	type TT struct {
		name   string
		s      string
		result string
	}
	tests := []TT{
		{
			name:   "no spaces",
			s:      "no_spaces_included",
			result: "no_spaces_included",
		},
		{
			name:   "has spaces",
			s:      "some spaces included",
			result: `"some spaces included"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(QuoteSpace(tt.s), tt.result)
		})
	}
}
