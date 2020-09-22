package sq

import (
	"log"
	"os"
)

// LogFlag is a flag that affects the verbosity of the Logger output.
type LogFlag int

// LogFlags
const (
	Linterpolate LogFlag = 1 << iota
	Lstats
	Lresults
	// Lparse
	Lverbose = Lstats | Lresults
)

// ExecFlag is a flag that affects the behavior of Exec.
type ExecFlag int

// ExecFlags
const (
	ErowsAffected ExecFlag = 1 << iota
)

var defaultLogger = log.New(os.Stdout, "[sq] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)

// BaseQuery is a common query builder that can transform into a SelectQuery,
// InsertQuery, UpdateQuery or DeleteQuery depending on the method that you
// call on it.
type BaseQuery struct {
	DB      DB
	Log     Logger
	LogFlag LogFlag
	CTEs    []CTE
}

// WithDefaultLog creates a new BaseQuery with the default logger and the LogFlag
func WithDefaultLog(flag LogFlag) BaseQuery {
	return BaseQuery{
		Log:     defaultLogger,
		LogFlag: flag,
	}
}

// WithDB creates a new BaseQuery with the DB.
func WithDB(db DB) BaseQuery {
	return BaseQuery{
		DB: db,
	}
}

// With creates a new BaseQuery with the CTEs.
func With(CTEs ...CTE) BaseQuery {
	return BaseQuery{
		CTEs: CTEs,
	}
}

// WithDefaultLog adds the default logger and the LogFlag to the BaseQuery.
func (q BaseQuery) WithDefaultLog(flag LogFlag) BaseQuery {
	q.Log = defaultLogger
	q.LogFlag = flag
	return q
}

// WithDB adds the DB to the BaseQuery.
func (q BaseQuery) WithDB(db DB) BaseQuery {
	q.DB = db
	return q
}

// With adds the CTEs to the BaseQuery
func (q BaseQuery) With(CTEs ...CTE) BaseQuery {
	q.CTEs = append(q.CTEs, CTEs...)
	return q
}

// From transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) From(table Table) SelectQuery {
	return SelectQuery{
		FromTable: table,
		CTEs:      q.CTEs,
		DB:        q.DB,
		Log:       q.Log,
		LogFlag:   q.LogFlag,
	}
}

// Select transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) Select(fields ...Field) SelectQuery {
	return SelectQuery{
		SelectFields: fields,
		CTEs:         q.CTEs,
		DB:           q.DB,
		Log:          q.Log,
		LogFlag:      q.LogFlag,
	}
}

// SelectOne transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) SelectOne() SelectQuery {
	return SelectQuery{
		SelectFields: Fields{FieldLiteral("1")},
		CTEs:         q.CTEs,
		DB:           q.DB,
		Log:          q.Log,
		LogFlag:      q.LogFlag,
	}
}

// SelectAll transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) SelectAll() SelectQuery {
	return SelectQuery{
		SelectFields: Fields{FieldLiteral("*")},
		CTEs:         q.CTEs,
		DB:           q.DB,
		Log:          q.Log,
		LogFlag:      q.LogFlag,
	}
}

// SelectCount transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) SelectCount() SelectQuery {
	return SelectQuery{
		SelectFields: Fields{FieldLiteral("COUNT(*)")},
		CTEs:         q.CTEs,
		DB:           q.DB,
		Log:          q.Log,
		LogFlag:      q.LogFlag,
	}
}

// SelectDistinct transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) SelectDistinct(fields ...Field) SelectQuery {
	return SelectQuery{
		SelectType:   SelectTypeDistinct,
		SelectFields: fields,
		CTEs:         q.CTEs,
		DB:           q.DB,
		Log:          q.Log,
		LogFlag:      q.LogFlag,
	}
}

// SelectDistinctOn transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) SelectDistinctOn(distinctFields ...Field) func(...Field) SelectQuery {
	return func(fields ...Field) SelectQuery {
		return SelectQuery{
			SelectType:   SelectTypeDistinctOn,
			SelectFields: fields,
			DistinctOn:   distinctFields,
			CTEs:         q.CTEs,
			DB:           q.DB,
			Log:          q.Log,
			LogFlag:      q.LogFlag,
		}
	}
}

// Selectx transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) Selectx(mapper func(*Row), accumulator func()) SelectQuery {
	return SelectQuery{
		RowMapper:   mapper,
		Accumulator: accumulator,
		CTEs:        q.CTEs,
		DB:          q.DB,
		Log:         q.Log,
		LogFlag:     q.LogFlag,
	}
}

// SelectRowx transforms the BaseQuery into a SelectQuery.
func (q BaseQuery) SelectRowx(mapper func(*Row)) SelectQuery {
	return SelectQuery{
		RowMapper: mapper,
		CTEs:      q.CTEs,
		DB:        q.DB,
		Log:       q.Log,
		LogFlag:   q.LogFlag,
	}
}

// InsertInto transforms the BaseQuery into an InsertQuery.
func (q BaseQuery) InsertInto(table BaseTable) InsertQuery {
	return InsertQuery{
		IntoTable: table,
		CTEs:      q.CTEs,
		DB:        q.DB,
		Log:       q.Log,
		LogFlag:   q.LogFlag,
	}
}

// Update transforms the BaseQuery into an UpdateQuery.
func (q BaseQuery) Update(table BaseTable) UpdateQuery {
	return UpdateQuery{
		UpdateTable: table,
		CTEs:        q.CTEs,
		DB:          q.DB,
		Log:         q.Log,
		LogFlag:     q.LogFlag,
	}
}

// DeleteFrom transforms the BaseQuery into a DeleteQuery.
func (q BaseQuery) DeleteFrom(table BaseTable) DeleteQuery {
	return DeleteQuery{
		FromTable: table,
		CTEs:      q.CTEs,
		DB:        q.DB,
		Log:       q.Log,
		LogFlag:   q.LogFlag,
	}
}

// Union transforms the BaseQuery into a VariadicQuery.
func (q BaseQuery) Union(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryUnion,
		Queries:  queries,
		DB:       q.DB,
		Log:      q.Log,
		LogFlag:  q.LogFlag,
	}
}

// UnionAll transforms the BaseQuery into a VariadicQuery.
func (q BaseQuery) UnionAll(queries ...Query) VariadicQuery {
	return VariadicQuery{
		topLevel: true,
		Operator: QueryUnionAll,
		Queries:  queries,
		DB:       q.DB,
		Log:      q.Log,
		LogFlag:  q.LogFlag,
	}
}
