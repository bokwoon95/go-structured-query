package sq

import (
	"context"
	"database/sql"
	"strings"
)

// Table is an interface representing anything that you can SELECT FROM or
// JOIN.
type Table interface {
	AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int)
	GetAlias() string
	GetName() string
}

func getAliasOrName(val interface {
	GetAlias() string
	GetName() string
}) string {
	s := val.GetAlias()
	if s == "" {
		s = val.GetName()
	}
	return s
}

// Query is an interface that specialises the Table interface. It covers only
// queries like SELECT/INSERT/UPDATE/DELETE.
type Query interface {
	AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int)
	// When NestThis is called on a query, it signals to the query that it is
	// being nested as part of a larger query. The nested query should:
	// - hold off logging anything because the parent query will do it
	NestThis() Query
	ToSQL() (string, []interface{})
}

// BaseTable is an interface that specialises the Table interface. It covers
// only tables/views that exist in the database.
type BaseTable interface {
	Table
	AssertBaseTable()
}

// Field is an interface that represents either a Table column or an SQL value.
type Field interface {
	// Fields should respect the excludedTableQualifiers argument in ToSQL().
	// E.g. if the field 'name' belongs to a table called 'users' and the
	// excludedTableQualifiers contains 'users', the field should present itself
	// as 'name' and not 'users.name'. i.e. any table qualifiers in the list
	// must be excluded.
	//
	// This is to play nice with certain clauses in the INSERT and UPDATE
	// queries that expressly forbid table qualified columns.
	AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string)
	GetAlias() string
	GetName() string
}

// Predicate is an interface that evaluates to true or false in SQL.
type Predicate interface {
	Field
	Not() Predicate
}

// Assignment is an interface representing an SQL Assignment 'Field = Value'.
type Assignment interface {
	AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string)
	AssertAssignment()
}

// DB is an interface providing database querying abilities.
type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Logger is an interface that provides logging.
type Logger interface {
	Output(calldepth int, s string) error
}
