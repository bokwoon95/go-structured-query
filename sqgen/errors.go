package sqgen

import (
	"fmt"
	"runtime"
)

/* Error Handling Utilities */

const recSep rune = 30 // ASCII Record Separator

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	_, filename, linenbr, _ := runtime.Caller(1)
	return fmt.Errorf(string(recSep)+" %s:%d %w", filename, linenbr, err)
}

