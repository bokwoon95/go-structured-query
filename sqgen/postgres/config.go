package postgres

type Config struct {
	// (required) Database URL
	Database string
	// Package name of the file to be generated
	Package string
	// Slice of database schemas that you want to generate tables for
	Schemas []string
	// Slice of case-insensitive table names or functions to exclude from generation
	Exclude []string
	// Used to log any skipped/unsupported column types
	Logger Logger
}

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
