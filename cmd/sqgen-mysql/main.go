package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

// currdir is the current directory of where the command was run from.
var currdir string = func() string {
	log.SetFlags(log.Lshortfile)
	currdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currdir
}()

// sqgenCmd is the root command for sqgen-mysql. It is referenced by the
// functionsCmd in functions.go and tablesCmd in tables.go.
var sqgenCmd = &cobra.Command{
	Use:           "sqgen-mysql",
	Short:         "Code generation for the sq package",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func main() {
	if err := sqgenCmd.Execute(); err != nil {
		dump(os.Stderr, err)
		os.Exit(1)
	}
}

/* Error Handling Utilities */

const recSep rune = 30 // ASCII Record Separator

// wrap will wrap an error and return a new error that is annotated with the
// file/linenumber of where wrap() was called.
func wrap(err error) error {
	if err == nil {
		return nil
	}
	_, filename, linenbr, _ := runtime.Caller(1)
	return fmt.Errorf(string(recSep)+" %s:%d %w", filename, linenbr, err)
}

// dump will dump the formatted error string (with each error in its own line)
// into w io.Writer.
func dump(w io.Writer, err error) {
	fmtedErr := strings.ReplaceAll(err.Error(), " "+string(recSep)+" ", "\n")
	fmt.Fprintln(w, fmtedErr)
}
