package sq

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// DeleteQuery represents a DELETE query.
type DeleteQuery struct {
	nested bool
	Alias  string
	// WITH
	CTEs []CTE
	// DELETE FROM
	FromTables []BaseTable
	// USING
	UsingTable Table
	JoinTables JoinTables
	// WHERE
	WherePredicate VariadicPredicate
	// ORDER BY
	OrderByFields Fields
	// LIMIT
	LimitValue *int64
	// DB
	DB DB
	// Logging
	Log     Logger
	LogFlag LogFlag
	logSkip int
}

// ToSQL marshals the DeleteQuery into a query string and args slice.
func (q DeleteQuery) ToSQL() (string, []interface{}) {
	q.logSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args, nil)
	return buf.String(), args
}

// AppendSQL marshals the DeleteQuery into a buffer and args slice.
func (q DeleteQuery) AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int) {
	// WITH
	if !q.nested {
		appendCTEs(buf, args, q.CTEs, nil, q.JoinTables)
	}
	// DELETE FROM
	buf.WriteString("DELETE FROM ")
	if len(q.FromTables) == 0 {
		buf.WriteString("NULL")
	} else {
		for i, table := range q.FromTables {
			if i > 0 {
				buf.WriteString(", ")
			}
			if table == nil {
				buf.WriteString("NULL")
				continue
			}
			alias := table.GetAlias()
			if alias != "" {
				buf.WriteString(alias)
			} else {
				table.AppendSQL(buf, args, nil)
			}
		}
	}
	// USING
	if q.UsingTable != nil {
		buf.WriteString(" USING ")
		switch v := q.UsingTable.(type) {
		case Query:
			buf.WriteString("(")
			v.NestThis().AppendSQL(buf, args, nil)
			buf.WriteString(")")
		default:
			q.UsingTable.AppendSQL(buf, args, nil)
		}
		alias := q.UsingTable.GetAlias()
		if alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(alias)
		}
	}
	// JOIN
	if len(q.JoinTables) > 0 {
		buf.WriteString(" ")
		q.JoinTables.AppendSQL(buf, args, nil)
	}
	// WHERE
	if len(q.WherePredicate.Predicates) > 0 {
		buf.WriteString(" WHERE ")
		q.WherePredicate.toplevel = true
		q.WherePredicate.AppendSQLExclude(buf, args, nil, nil)
	}
	// ORDER BY
	if len(q.OrderByFields) > 0 {
		buf.WriteString(" ORDER BY ")
		q.OrderByFields.AppendSQLExclude(buf, args, nil, nil)
	}
	// LIMIT
	if q.LimitValue != nil {
		buf.WriteString(" LIMIT ?")
		if *q.LimitValue < 0 {
			*q.LimitValue = -*q.LimitValue
		}
		*args = append(*args, *q.LimitValue)
	}
	if !q.nested {
		if q.Log != nil {
			query := buf.String()
			var logOutput string
			switch {
			case Lstats&q.LogFlag != 0:
				logOutput = "\n----[ Executing query ]----\n" + query + " " + fmt.Sprint(*args) +
					"\n----[ with bind values ]----\n" + questionInterpolate(query, *args...)
			case Linterpolate&q.LogFlag != 0:
				logOutput = "Executing query: " + questionInterpolate(query, *args...)
			default:
				logOutput = "Executing query: " + query + " " + fmt.Sprint(*args)
			}
			switch q.Log.(type) {
			case *log.Logger:
				_ = q.Log.Output(q.logSkip+2, logOutput)
			default:
				_ = q.Log.Output(q.logSkip+1, logOutput)
			}
		}
	}
}

// NestThis indicates to the DeleteQuery that it is nested.
func (q DeleteQuery) NestThis() Query {
	q.nested = true
	return q
}

// DeleteFrom creates a new DeleteQuery.
func DeleteFrom(tables ...BaseTable) DeleteQuery {
	return DeleteQuery{
		FromTables: tables,
		Alias:      randomString(8),
	}
}

// With appends the CTEs into the DeleteQuery.
func (q DeleteQuery) With(ctes ...CTE) DeleteQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

// DeleteFrom adds new tables to delete from to the DeleteQuery.
func (q DeleteQuery) DeleteFrom(tables ...BaseTable) DeleteQuery {
	q.FromTables = append(q.FromTables, tables...)
	return q
}

// Using adds a new table to the DeleteQuery.
func (q DeleteQuery) Using(table Table) DeleteQuery {
	q.UsingTable = table
	return q
}

// Join joins a new table to the DeleteQuery based on the predicates.
func (q DeleteQuery) Join(table Table, predicate Predicate, predicates ...Predicate) DeleteQuery {
	predicates = append([]Predicate{predicate}, predicates...)
	q.JoinTables = append(q.JoinTables, JoinTable{
		JoinType: JoinTypeInner,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	})
	return q
}

// LeftJoin left joins a new table to the DeleteQuery based on the predicates.
func (q DeleteQuery) LeftJoin(table Table, predicate Predicate, predicates ...Predicate) DeleteQuery {
	predicates = append([]Predicate{predicate}, predicates...)
	q.JoinTables = append(q.JoinTables, JoinTable{
		JoinType: JoinTypeLeft,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	})
	return q
}

// RightJoin right joins a new table to the DeleteQuery based on the predicates.
func (q DeleteQuery) RightJoin(table Table, predicate Predicate, predicates ...Predicate) DeleteQuery {
	predicates = append([]Predicate{predicate}, predicates...)
	q.JoinTables = append(q.JoinTables, JoinTable{
		JoinType: JoinTypeRight,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	})
	return q
}

// FullJoin full joins a table to the DeleteQuery based on the predicates.
func (q DeleteQuery) FullJoin(table Table, predicate Predicate, predicates ...Predicate) DeleteQuery {
	predicates = append([]Predicate{predicate}, predicates...)
	q.JoinTables = append(q.JoinTables, JoinTable{
		JoinType: JoinTypeFull,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	})
	return q
}

// CustomJoin custom joins a table to the DeleteQuery. The join type can be
// specified with a string, e.g. "CROSS JOIN".
func (q DeleteQuery) CustomJoin(joinType JoinType, table Table, predicates ...Predicate) DeleteQuery {
	q.JoinTables = append(q.JoinTables, JoinTable{
		JoinType: joinType,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	})
	return q
}

// Where appends the predicates to the WHERE clause in the DeleteQuery.
func (q DeleteQuery) Where(predicates ...Predicate) DeleteQuery {
	q.WherePredicate.Predicates = append(q.WherePredicate.Predicates, predicates...)
	return q
}

// OrderBy appends the fields to the ORDER BY clause in the DeleteQuery.
func (q DeleteQuery) OrderBy(fields ...Field) DeleteQuery {
	q.OrderByFields = append(q.OrderByFields, fields...)
	return q
}

// Limit sets the limit in the DeleteQuery.
func (q DeleteQuery) Limit(limit int) DeleteQuery {
	num := int64(limit)
	q.LimitValue = &num
	return q
}

// Exec will execute the DeleteQuery with the given DB. It will only compute
// the rowsAffected if the ErowsAffected Execflag is passed to it.
func (q DeleteQuery) Exec(db DB, flag ExecFlag) (rowsAffected int64, err error) {
	q.logSkip += 1
	return q.ExecContext(nil, db, flag)
}

// ExecContext will execute the DeleteQuery with the given DB and context. It will
// only compute the rowsAffected if the ErowsAffected Execflag is passed to it.
func (q DeleteQuery) ExecContext(ctx context.Context, db DB, flag ExecFlag) (rowsAffected int64, err error) {
	if db == nil {
		if q.DB == nil {
			return rowsAffected, errors.New("DB cannot be nil")
		}
		db = q.DB
	}
	logBuf := &strings.Builder{}
	start := time.Now()
	defer func() {
		if q.Log == nil {
			return
		}
		elapsed := time.Since(start)
		if Lstats&q.LogFlag != 0 && ErowsAffected&flag != 0 {
			logBuf.WriteString("\n(Deleted ")
			logBuf.WriteString(strconv.FormatInt(rowsAffected, 10))
			logBuf.WriteString(" rows in ")
			logBuf.WriteString(elapsed.String())
			logBuf.WriteString(")")
		}
		if logBuf.Len() > 0 {
			switch q.Log.(type) {
			case *log.Logger:
				_ = q.Log.Output(q.logSkip+2, logBuf.String())
			default:
				_ = q.Log.Output(q.logSkip+1, logBuf.String())
			}
		}
	}()
	var res sql.Result
	tmpbuf := &strings.Builder{}
	var tmpargs []interface{}
	q.logSkip += 1
	q.AppendSQL(tmpbuf, &tmpargs, nil)
	if ctx == nil {
		res, err = db.Exec(tmpbuf.String(), tmpargs...)
	} else {
		res, err = db.ExecContext(ctx, tmpbuf.String(), tmpargs...)
	}
	if err != nil {
		return rowsAffected, err
	}
	if res != nil && ErowsAffected&flag != 0 {
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return rowsAffected, err
		}
	}
	return rowsAffected, nil
}
