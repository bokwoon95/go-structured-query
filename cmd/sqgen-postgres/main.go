package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bokwoon95/go-structured-query/sqgen/postgres"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func main() {
	if err := sqgenCmd.Execute(); err != nil {
		dump(os.Stderr, err)
		os.Exit(1)
	}
}

// sqgenCmd is the root command for sqgen-postgres. It is referenced by the
// functionsCmd in functions.go and tablesCmd in tables.go.
var sqgenCmd = &cobra.Command{
	Use:           "sqgen-postgres",
	Short:         "Code generation for the sq package",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Generate tables from the database",
	RunE:  tablesRun,
}

var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Generate functions from the database",
	RunE:  functionsRun,
}

// currdir is the current directory of where the command was run from.
var currdir string = func() string {
	log.SetFlags(log.Lshortfile)
	currdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currdir
}()

var (
	tablesDatabase  *string
	tablesDirectory *string
	tablesDryrun    *bool
	tablesFile      *string
	tablesOverwrite *bool
	tablesPkg       *string
	tablesSchemas   *[]string
	tablesExclude   *[]string

	functionsDatabase  *string
	functionsDirectory *string
	functionsDryrun    *bool
	functionsFile      *string
	functionsOverwrite *bool
	functionsPkg       *string
	functionsSchemas   *[]string
	functionsExclude   *[]string
)

func init() {
	sqgenCmd.AddCommand(tablesCmd, functionsCmd)

	// initialize tables flags

	tablesDatabase = tablesCmd.Flags().String("database", "", "(required) Database URL")
	tablesDirectory = tablesCmd.Flags().
		String("directory", filepath.Join(currdir, "tables"), "(optional) Directory to place the generated file. Can be absolute or relative filepath")
	tablesDryrun = tablesCmd.Flags().
		Bool("dryrun", false, "(optional) Print the list of tables to be generated without generating the file")
	tablesFile = tablesCmd.Flags().
		String("file", "tables.go", "(optional) Name of the file to be generated. If file already exists, -overwrite flag must be specified to overwrite the file")
	tablesOverwrite = tablesCmd.Flags().
		Bool("overwrite", false, "(optional) Overwrite any files that already exist")
	tablesPkg = tablesCmd.Flags().
		String("pkg", "tables", "(optional) Package name of the file to be generated")
	tablesSchemas = tablesCmd.Flags().
		StringSlice("schemas", []string{"public"}, "(optional) A comma separated list of database schemas that you want to generate tables for. Please don't include any spaces")
	tablesExclude = tablesCmd.Flags().
		StringSlice("exclude", nil, "(optional) A comma separated list of case-insensitive table names that you wish to exclude from table generation. Please don't include any spaces")
	// required flag
	err := cobra.MarkFlagRequired(tablesCmd.LocalFlags(), "database")

	if err != nil {
		panic(err)
	}

	// initialize functions flags

	functionsDatabase = functionsCmd.Flags().String("database", "", "(required) Database URL")
	functionsDirectory = functionsCmd.Flags().
		String("directory", filepath.Join(currdir, "tables"), "(optional) Directory to place the generated file. Can be absolute or relative filepath")
	functionsDryrun = functionsCmd.Flags().
		Bool("dryrun", false, "(optional) Print the list of functions to be generated without generating the file")
	functionsFile = functionsCmd.Flags().
		String("file", "functions.go", "(optional) Name of the file to be generated. If file already exists, -overwrite flag must be specified to overwrite the file")
	functionsOverwrite = functionsCmd.Flags().
		Bool("overwrite", false, "(optional) Overwrite any files that already exist")
	functionsPkg = functionsCmd.Flags().
		String("pkg", "tables", "(optional) Package name of the file to be generated")
	functionsSchemas = functionsCmd.Flags().
		StringSlice("schemas", []string{"public"}, "(optional) A comma separated list of database schemas that you want to generate functions for. Please don't include any spaces")
	functionsExclude = functionsCmd.Flags().
		StringSlice("exclude", nil, "(optional) A comma separated list of case-insensitive function names that you wish to exclude from table generation. Please don't include any spaces")
	// required flag
	err = cobra.MarkFlagRequired(functionsCmd.LocalFlags(), "database")

	if err != nil {
		panic(err)
	}
}

// tablesRun is the main function to be run with `sqgen-postgres tables`
func tablesRun(cmd *cobra.Command, args []string) error {
	db, err := openAndPing(*tablesDatabase)

	if err != nil {
		return err
	}

	// dereference to get flag values
	config := postgres.Config{
		DB:      db,
		Package: *tablesPkg,
		Schemas: *tablesSchemas,
		Exclude: *tablesExclude,
		Logger:  log.New(os.Stderr, "", log.Ltime),
	}

	writer, err := getWriter(*tablesDryrun, *tablesOverwrite, *tablesDirectory, *tablesFile)

	if err != nil {
		return err
	}

	defer writer.Close()
	numTables, err := postgres.BuildTables(config, writer)

	if err != nil {
		return err
	}

	if !*tablesDryrun {
		fmt.Printf("[RESULT] %d tables written into %s\n", numTables, writer.Name())
	}
	return nil
}

// functionsRun is the main function to be run with `sqgen-postgres functions`
func functionsRun(cmd *cobra.Command, args []string) error {
	db, err := openAndPing(*functionsDatabase)

	if err != nil {
		return err
	}

	// dereference to get flag values
	config := postgres.Config{
		DB:      db,
		Package: *functionsPkg,
		Schemas: *functionsSchemas,
		Exclude: *functionsExclude,
		Logger:  log.New(os.Stderr, "", log.Ltime),
	}

	writer, err := getWriter(
		*functionsDryrun,
		*functionsOverwrite,
		*functionsDirectory,
		*functionsFile,
	)

	if err != nil {
		return err
	}

	defer writer.Close()

	numFunctions, err := postgres.BuildFunctions(config, writer)

	if err != nil {
		return err
	}

	if !*functionsDryrun {
		fmt.Printf("[RESULT] %d functions written into %s\n", numFunctions, writer.Name())
	}

	return nil
}

func openAndPing(database string) (*sql.DB, error) {
	db, err := sql.Open("postgres", database)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf(
			"Could not ping the database, is the database reachable via %s? %w",
			database,
			err,
		)
	}

	return db, nil
}

func getWriter(dryrun, overwrite bool, directory, file string) (*os.File, error) {
	if dryrun {
		return os.Stdout, nil
	}

	if !strings.HasSuffix(file, ".go") {
		file = file + ".go"
	}

	asboluteFilePath := filepath.Join(directory, file)
	if _, err := os.Stat(asboluteFilePath); err == nil && !overwrite {
		return nil, fmt.Errorf(
			"%s already exists. If you wish to overwrite it, provide the --overwrite flag",
			asboluteFilePath,
		)
	}

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return nil, fmt.Errorf("Could not create directory %s: %w", directory, err)
	}

	filename := filepath.Join(directory, file)
	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
}

/* Misc Utilities */

const recSep rune = 30 // ASCII Record Separator

// dump will dump the formatted error string (with each error in its own line)
// into w io.Writer.
func dump(w io.Writer, err error) {
	fmtedErr := strings.ReplaceAll(err.Error(), " "+string(recSep)+" ", "\n")
	fmt.Fprintln(w, fmtedErr)
}
