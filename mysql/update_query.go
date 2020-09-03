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

// UpdateQuery represents an UPDATE query.
type UpdateQuery struct {
	Nested bool
	Alias  string
	// WITH
	CTEs CTEs
	// UPDATE
	UpdateTable BaseTable
	// SET
	Assignments Assignments
	// JOIN
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
	LogSkip int
}

// ToSQL marshals the UpdateQuery into a query string and args slice.
func (q UpdateQuery) ToSQL() (string, []interface{}) {
	q.LogSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

// AppendSQL marshals the UpdateQuery into a buffer and args slice.
func (q UpdateQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	// WITH
	if len(q.CTEs) > 0 {
		q.CTEs.AppendSQL(buf, args)
		buf.WriteString(" ")
	}
	// UPDATE
	buf.WriteString("UPDATE ")
	if q.UpdateTable == nil {
		buf.WriteString("NULL")
	} else {
		switch v := q.UpdateTable.(type) {
		case Query:
			buf.WriteString("(")
			v.NestThis().AppendSQL(buf, args)
			buf.WriteString(")")
		default:
			q.UpdateTable.AppendSQL(buf, args)
		}
		alias := q.UpdateTable.GetAlias()
		if alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(alias)
		}
	}
	// SET
	if len(q.Assignments) > 0 {
		buf.WriteString(" SET ")
		q.Assignments.AppendSQLExclude(buf, args, nil)
	}
	// JOIN
	if len(q.JoinTables) > 0 {
		buf.WriteString(" ")
		q.JoinTables.AppendSQL(buf, args)
	}
	// WHERE
	if len(q.WherePredicate.Predicates) > 0 {
		buf.WriteString(" WHERE ")
		q.WherePredicate.Toplevel = true
		q.WherePredicate.AppendSQLExclude(buf, args, nil)
	}
	// ORDER BY
	if len(q.OrderByFields) > 0 {
		buf.WriteString(" ORDER BY ")
		q.OrderByFields.AppendSQLExclude(buf, args, nil)
	}
	// LIMIT
	if q.LimitValue != nil {
		buf.WriteString(" LIMIT ?")
		if *q.LimitValue < 0 {
			*q.LimitValue = -*q.LimitValue
		}
		*args = append(*args, *q.LimitValue)
	}
	if !q.Nested {
		if q.Log != nil {
			query := buf.String()
			var logOutput string
			switch {
			case Lstats&q.LogFlag != 0:
				logOutput = "\n----[ Executing query ]----\n" + query + " " + fmt.Sprint(*args) +
					"\n----[ with bind values ]----\n" + QuestionInterpolate(query, *args...)
			case Linterpolate&q.LogFlag != 0:
				logOutput = "Executing query: " + QuestionInterpolate(query, *args...)
			default:
				logOutput = "Executing query: " + query + " " + fmt.Sprint(*args)
			}
			switch q.Log.(type) {
			case *log.Logger:
				q.Log.Output(q.LogSkip+2, logOutput)
			default:
				q.Log.Output(q.LogSkip+1, logOutput)
			}
		}
	}
}

// As aliases the UpdateQuery i.e. 'query AS alias'.
func (q UpdateQuery) As(alias string) UpdateQuery {
	q.Alias = alias
	return q
}

// Update creates a new UpdateQuery.
func Update(table BaseTable) UpdateQuery {
	return UpdateQuery{
		UpdateTable: table,
		Alias:       RandomString(8),
	}
}

// With appends a list of CTEs into the UpdateQuery.
func (q UpdateQuery) With(ctes ...CTE) UpdateQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

// Update sets the update table for the UpdateQuery.
func (q UpdateQuery) Update(table BaseTable) UpdateQuery {
	q.UpdateTable = table
	return q
}

// Set appends the assignments to SET clause of the UpdateQuery.
func (q UpdateQuery) Set(assignments ...Assignment) UpdateQuery {
	q.Assignments = append(q.Assignments, assignments...)
	return q
}

// Join joins a new table to the UpdateQuery based on the predicates.
func (q UpdateQuery) Join(table Table, predicate Predicate, predicates ...Predicate) UpdateQuery {
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

// LeftJoin left joins a new table to the UpdateQuery based on the predicates.
func (q UpdateQuery) LeftJoin(table Table, predicate Predicate, predicates ...Predicate) UpdateQuery {
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

// RightJoin right joins a new table to the UpdateQuery based on the predicates.
func (q UpdateQuery) RightJoin(table Table, predicate Predicate, predicates ...Predicate) UpdateQuery {
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

// FullJoin full joins a table to the UpdateQuery based on the predicates.
func (q UpdateQuery) FullJoin(table Table, predicate Predicate, predicates ...Predicate) UpdateQuery {
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

// CustomJoin custom joins a table to the UpdateQuery. The join type can be
// specified with a string, e.g. "CROSS JOIN".
func (q UpdateQuery) CustomJoin(joinType JoinType, table Table, predicates ...Predicate) UpdateQuery {
	q.JoinTables = append(q.JoinTables, JoinTable{
		JoinType: joinType,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	})
	return q
}

// Where appends the predicates to the WHERE clause in the UpdateQuery.
func (q UpdateQuery) Where(predicates ...Predicate) UpdateQuery {
	q.WherePredicate.Predicates = append(q.WherePredicate.Predicates, predicates...)
	return q
}

// OrderBy appends the fields to the ORDER BY clause in the UpdateQuery.
func (q UpdateQuery) OrderBy(fields ...Field) UpdateQuery {
	q.OrderByFields = append(q.OrderByFields, fields...)
	return q
}

// Limit sets the limit in the UpdateQuery.
func (q UpdateQuery) Limit(limit int) UpdateQuery {
	num := int64(limit)
	q.LimitValue = &num
	return q
}

// Exec will execute the UpdateQuery with the given DB. It will only compute
// the rowsAffected if the ErowsAffected Execflag is passed to it.
func (q UpdateQuery) Exec(db DB, flag ExecFlag) (rowsAffected int64, err error) {
	q.LogSkip += 1
	return q.ExecContext(nil, db, flag)
}

// ExecContext will execute the UpdateQuery with the given DB and context. It will
// only compute the rowsAffected if the ErowsAffected Execflag is passed to it.
func (q UpdateQuery) ExecContext(ctx context.Context, db DB, flag ExecFlag) (rowsAffected int64, err error) {
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
			logBuf.WriteString("\n(Updated ")
			logBuf.WriteString(strconv.FormatInt(rowsAffected, 10))
			logBuf.WriteString(" rows in ")
			logBuf.WriteString(elapsed.String())
			logBuf.WriteString(")")
		}
		if logBuf.Len() > 0 {
			switch q.Log.(type) {
			case *log.Logger:
				q.Log.Output(q.LogSkip+2, logBuf.String())
			default:
				q.Log.Output(q.LogSkip+1, logBuf.String())
			}
		}
	}()
	var res sql.Result
	tmpbuf := &strings.Builder{}
	var tmpargs []interface{}
	q.LogSkip += 1
	q.AppendSQL(tmpbuf, &tmpargs)
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

// GetAlias returns the alias of the UpdateQuery.
func (q UpdateQuery) GetAlias() string {
	return q.Alias
}

// GetName returns the name of the UpdateQuery, which is always an empty
// string.
func (q UpdateQuery) GetName() string {
	return ""
}

// NestThis indicates to the UpdateQuery that it is nested.
func (q UpdateQuery) NestThis() Query {
	q.Nested = true
	return q
}
