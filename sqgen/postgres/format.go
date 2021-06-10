package postgres

import (
	"golang.org/x/tools/imports"
)

// uses goimports tool to format code
// goimports also runs gofmt, so no need to run both separately
func formatBytes(src []byte) ([]byte, error) {
	return imports.Process("", src, nil)
}
