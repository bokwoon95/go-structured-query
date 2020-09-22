package sq

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DeleteQuery represents a DELETE query.
type DeleteQuery struct {
	nested bool
	// WITH
	CTEs []CTE
	// DELETE FROM
	FromTable BaseTable
	// USING
	UsingTable Table
	JoinTables JoinTables
	// WHERE
	WherePredicate VariadicPredicate
	// RETURNING
	ReturningFields Fields
	// DB
	DB          DB
	RowMapper   func(*Row)
	Accumulator func()
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
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

// AppendSQL marshals the DeleteQuery into a buffer and args slice.
func (q DeleteQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	// WITH
	if !q.nested {
		appendCTEs(buf, args, q.CTEs, nil, q.JoinTables)
	}
	// DELETE FROM
	buf.WriteString("DELETE FROM ")
	if q.FromTable == nil {
		buf.WriteString("NULL")
	} else {
		switch v := q.FromTable.(type) {
		case Query:
			buf.WriteString("(")
			v.NestThis().AppendSQL(buf, args)
			buf.WriteString(")")
		default:
			q.FromTable.AppendSQL(buf, args)
		}
		alias := q.FromTable.GetAlias()
		if alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(alias)
		}
	}
	// USING
	if q.UsingTable != nil {
		buf.WriteString(" USING ")
		switch v := q.UsingTable.(type) {
		case Query:
			buf.WriteString("(")
			v.NestThis().AppendSQL(buf, args)
			buf.WriteString(")")
		default:
			q.FromTable.AppendSQL(buf, args)
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
		q.JoinTables.AppendSQL(buf, args)
	}
	// WHERE
	if len(q.WherePredicate.Predicates) > 0 {
		buf.WriteString(" WHERE ")
		q.WherePredicate.toplevel = true
		q.WherePredicate.AppendSQLExclude(buf, args, nil)
	}
	// RETURNING
	if len(q.ReturningFields) > 0 {
		buf.WriteString(" RETURNING ")
		q.ReturningFields.AppendSQLExcludeWithAlias(buf, args, nil)
	}
	if !q.nested {
		query := buf.String()
		buf.Reset()
		questionToDollarPlaceholders(buf, query)
		if q.Log != nil {
			var logOutput string
			switch {
			case Lstats&q.LogFlag != 0:
				logOutput = "\n----[ Executing query ]----\n" + buf.String() + " " + fmt.Sprint(*args) +
					"\n----[ with bind values ]----\n" + questionInterpolate(query, *args...)
			case Linterpolate&q.LogFlag != 0:
				logOutput = questionInterpolate(query, *args...)
			default:
				logOutput = buf.String() + " " + fmt.Sprint(*args)
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
func DeleteFrom(table BaseTable) DeleteQuery {
	return DeleteQuery{
		FromTable: table,
	}
}

// With appends the CTEs into the DeleteQuery.
func (q DeleteQuery) With(ctes ...CTE) DeleteQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

// DeleteFrom sets the table to be deleted from in the DeleteQuery.
func (q DeleteQuery) DeleteFrom(table BaseTable) DeleteQuery {
	q.FromTable = table
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

// Returning appends the fields to the RETURNING clause of the DeleteQuery.
func (q DeleteQuery) Returning(fields ...Field) DeleteQuery {
	q.ReturningFields = append(q.ReturningFields, fields...)
	return q
}

// ReturningOne sets the RETURNING clause to RETURNING 1 in the DeleteQuery.
func (q DeleteQuery) ReturningOne() DeleteQuery {
	q.ReturningFields = Fields{FieldLiteral("1")}
	return q
}

// Returningx sets the rowmapper and accumulator function of the DeleteQuery.
func (q DeleteQuery) Returningx(mapper func(*Row), accumulator func()) DeleteQuery {
	q.RowMapper = mapper
	q.Accumulator = accumulator
	return q
}

// Returningx sets the rowmapper function of the DeleteQuery.
func (q DeleteQuery) ReturningRowx(mapper func(*Row)) DeleteQuery {
	q.RowMapper = mapper
	return q
}

func (q DeleteQuery) Fetch(db DB) (err error) {
	q.logSkip += 1
	return q.FetchContext(nil, db)
}

func (q DeleteQuery) FetchContext(ctx context.Context, db DB) (err error) {
	if db == nil {
		if q.DB == nil {
			return errors.New("DB cannot be nil")
		}
		db = q.DB
	}
	if q.RowMapper == nil {
		return fmt.Errorf("cannot call Fetch/FetchContext without a mapper")
	}
	logBuf := &strings.Builder{}
	start := time.Now()
	var rowcount int
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case ExitCode:
				if v != ExitPeacefully {
					err = v
				}
			case error:
				err = v
			default:
				err = fmt.Errorf("%#v", r)
			}
			return
		}
		if q.Log == nil {
			return
		}
		elapsed := time.Since(start)
		if Lresults&q.LogFlag != 0 && rowcount > 5 {
			logBuf.WriteString("\n...")
		}
		if Lstats&q.LogFlag != 0 {
			logBuf.WriteString("\n(Fetched ")
			logBuf.WriteString(strconv.Itoa(rowcount))
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
	r := &Row{}
	q.RowMapper(r)
	q.ReturningFields = r.fields
	tmpbuf := &strings.Builder{}
	var tmpargs []interface{}
	q.logSkip += 1
	q.AppendSQL(tmpbuf, &tmpargs)
	if ctx == nil {
		r.rows, err = db.Query(tmpbuf.String(), tmpargs...)
	} else {
		r.rows, err = db.QueryContext(ctx, tmpbuf.String(), tmpargs...)
	}
	if err != nil {
		return err
	}
	defer r.rows.Close()
	if len(r.dest) == 0 {
		return nil
	}
	for r.rows.Next() {
		rowcount++
		err = r.rows.Scan(r.dest...)
		if err != nil {
			errbuf := &strings.Builder{}
			for i := range r.dest {
				tmpbuf.Reset()
				tmpargs = tmpargs[:0]
				r.fields[i].AppendSQLExclude(tmpbuf, &tmpargs, nil)
				errbuf.WriteString("\n" +
					strconv.Itoa(i) + ") " +
					dollarInterpolate(tmpbuf.String(), tmpargs...) + " => " +
					reflect.TypeOf(r.dest[i]).String())
			}
			return fmt.Errorf("Please check if your mapper function is correct:%s\n%w", errbuf.String(), err)
		}
		if q.Log != nil && Lresults&q.LogFlag != 0 && rowcount <= 5 {
			logBuf.WriteString("\n----[ Row ")
			logBuf.WriteString(strconv.Itoa(rowcount))
			logBuf.WriteString(" ]----")
			for i := range r.dest {
				tmpbuf.Reset()
				tmpargs = tmpargs[:0]
				r.fields[i].AppendSQLExclude(tmpbuf, &tmpargs, nil)
				logBuf.WriteString("\n")
				logBuf.WriteString(dollarInterpolate(tmpbuf.String(), tmpargs...))
				logBuf.WriteString(": ")
				logBuf.WriteString(appendSQLDisplay(r.dest[i]))
			}
		}
		r.index = 0
		q.RowMapper(r)
		if q.Accumulator == nil {
			break
		}
		q.Accumulator()
	}
	if rowcount == 0 && q.Accumulator == nil {
		return sql.ErrNoRows
	}
	if e := r.rows.Close(); e != nil {
		return e
	}
	return r.rows.Err()
}

func (q DeleteQuery) Exec(db DB, flag ExecFlag) (rowsAffected int64, err error) {
	q.logSkip += 1
	return q.ExecContext(nil, db, flag)
}

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
