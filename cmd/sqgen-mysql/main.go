package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bokwoon95/go-structured-query/sqgen/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

func main() {
	if err := sqgenCmd.Execute(); err != nil {
		dump(os.Stderr, err)
		os.Exit(1)
	}
}

// sqgenCmd is the root command for sqgen-mysql. It is referenced by the
// functionsCmd in functions.go and tablesCmd in tables.go.
var sqgenCmd = &cobra.Command{
	Use:           "sqgen-mysql",
	Short:         "Code generation for the sq package",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Generate tables from the database",
	RunE:  tablesRun,
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
)

func init() {
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
		StringSlice("schemas", nil, "(required) A comma separated list of schemas (databases) that you want to generate tables for. In MySQL this is usually the database name you are using. Please don't include any spaces")
	tablesExclude = tablesCmd.Flags().
		StringSlice("exclude", nil, "(optional) A comma separated list of case-insensitive table names that you wish to exclude from table generation. Please don't include any spaces")

	// required flags
	err := cobra.MarkFlagRequired(tablesCmd.LocalFlags(), "database")

	if err != nil {
		panic(err)
	}

	err = cobra.MarkFlagRequired(tablesCmd.LocalFlags(), "schemas")

	if err != nil {
		panic(err)
	}
}

func tablesRun(cmd *cobra.Command, args []string) error {
	// dereference to get flag values

	if len(*tablesSchemas) == 0 {
		return fmt.Errorf("'%v' is not a valid comma separated list of schemas", tablesSchemas)
	}

	db, err := openAndPing(*tablesDatabase)

	if err != nil {
		return err
	}

	config := mysql.Config{
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

	numTables, err := mysql.BuildTables(config, writer)

	if err != nil {
		return err
	}

	if !*tablesDryrun {
		fmt.Printf("[RESULT] %d tables written into %s\n", numTables, writer.Name())
	}

	return nil
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

func openAndPing(database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", database)

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

/* Error Handling Utilities */

const recSep rune = 30 // ASCII Record Separator

// dump will dump the formatted error string (with each error in its own line)
// into w io.Writer.
func dump(w io.Writer, err error) {
	fmtedErr := strings.ReplaceAll(err.Error(), " "+string(recSep)+" ", "\n")
	fmt.Fprintln(w, fmtedErr)
}
