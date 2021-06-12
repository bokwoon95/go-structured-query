package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bokwoon95/go-structured-query/sqgen/postgres"

	"github.com/spf13/cobra"
)

var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Generate functions from the database",
	RunE:  functionsRun,
}

func init() {
	sqgenCmd.AddCommand(functionsCmd)
	// Initialise flags
	functionsCmd.Flags().String("database", "", "(required) Database URL")
	functionsCmd.Flags().String("directory", filepath.Join(currdir, "tables"), "(optional) Directory to place the generated file. Can be absolute or relative filepath")
	functionsCmd.Flags().Bool("dryrun", false, "(optional) Print the list of functions to be generated without generating the file")
	functionsCmd.Flags().String("file", "functions.go", "(optional) Name of the file to be generated. If file already exists, -overwrite flag must be specified to overwrite the file")
	functionsCmd.Flags().Bool("overwrite", false, "(optional) Overwrite any files that already exist")
	functionsCmd.Flags().String("pkg", "tables", "(optional) Package name of the file to be generated")
	functionsCmd.Flags().StringSlice("schemas", []string{"public"}, "(optional) A comma separated list of database schemas that you want to generate functions for. Please don't include any spaces")
	functionsCmd.Flags().StringSlice("exclude", []string{}, "(optional) A comma separated list of case-insensitive function names that you wish to exclude from table generation. Please don't include any spaces")
	// Mark required flags
	cobra.MarkFlagRequired(functionsCmd.LocalFlags(), "database")
}

func functionsRun(cmd *cobra.Command, args []string) error {
	database, _ := cmd.Flags().GetString("database")
	directory, _ := cmd.Flags().GetString("directory")
	dryrun, _ := cmd.Flags().GetBool("dryrun")
	file, _ := cmd.Flags().GetString("file")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	pkg, _ := cmd.Flags().GetString("pkg")
	schemas, _ := cmd.Flags().GetStringSlice("schemas")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")

	if !strings.HasSuffix(file, ".go") {
		file = file + ".go"
	}
	asboluteFilePath := filepath.Join(directory, file)
	if _, err := os.Stat(asboluteFilePath); err == nil && !overwrite {
		return fmt.Errorf("%s already exists. If you wish to overwrite it, provide the --overwrite flag", asboluteFilePath)
	}

	config := postgres.Config{
		Database: database,
		Package:  pkg,
		Schemas:  schemas,
		Exclude:  exclude,
		Logger:   log.New(os.Stderr, "", log.Ltime),
	}

	var writer io.Writer

	if dryrun {
		writer = os.Stdout
	} else {
		err := os.MkdirAll(directory, 0755)

		if err != nil {
			return fmt.Errorf("Could not create directory %s: %w", directory, err)
		}

		filename := filepath.Join(directory, file)

		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

		if err != nil {
			return err
		}

		writer = f

		defer f.Close()
	}

	return postgres.BuildFunctions(config, writer)
}
