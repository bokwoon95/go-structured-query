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

type UpdateQuery struct {
	Nested bool
	Alias  string
	// WITH
	CTEs CTEs
	// UPDATE
	UpdateTable BaseTable
	// SET
	Assignments Assignments
	// FROM
	FromTable  Table
	JoinTables JoinTables
	// WHERE
	WherePredicate VariadicPredicate
	// RETURNING
	ReturningFields Fields
	// DB
	DB          DB
	Mapper      func(*Row)
	Accumulator func()
	// Logging
	Log     Logger
	LogFlag LogFlag
	LogSkip int
}

func (q UpdateQuery) ToSQL() (string, []interface{}) {
	q.LogSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

func (q UpdateQuery) AppendSQL(buf Buffer, args *[]interface{}) {
	var excludedTableQualifiers []string
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
		q.UpdateTable.AppendSQL(buf, args)
		name := q.UpdateTable.GetName()
		alias := q.UpdateTable.GetAlias()
		if alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(alias)
			excludedTableQualifiers = append(excludedTableQualifiers, alias)
		} else {
			excludedTableQualifiers = append(excludedTableQualifiers, name)
		}
	}
	// SET
	if len(q.Assignments) > 0 {
		buf.WriteString(" SET ")
		q.Assignments.AppendSQLExclude(buf, args, excludedTableQualifiers)
	}
	// FROM
	if q.FromTable != nil {
		buf.WriteString(" FROM ")
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
	// RETURNING
	if len(q.ReturningFields) > 0 {
		buf.WriteString(" RETURNING ")
		q.ReturningFields.AppendSQLExcludeWithAlias(buf, args, nil)
	}
	if !q.Nested {
		query := buf.String()
		buf.Reset()
		QuestionToDollarPlaceholders(buf, query)
		if q.Log != nil {
			var logOutput string
			switch {
			case Lstats&q.LogFlag != 0:
				logOutput = "\n----[ Executing query ]----\n" + buf.String() + " " + fmt.Sprint(*args) +
					"\n----[ with bind values ]----\n" + QuestionInterpolate(query, *args...)
			case Linterpolate&q.LogFlag != 0:
				logOutput = QuestionInterpolate(query, *args...)
			default:
				logOutput = buf.String() + " " + fmt.Sprint(*args)
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

func (q UpdateQuery) As(alias string) UpdateQuery {
	q.Alias = alias
	return q
}

func Update(table BaseTable) UpdateQuery {
	return UpdateQuery{
		UpdateTable: table,
		Alias:       RandomString(8),
	}
}

func (q UpdateQuery) With(ctes ...CTE) UpdateQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

func (q UpdateQuery) Update(table BaseTable) UpdateQuery {
	q.UpdateTable = table
	return q
}

func (q UpdateQuery) Set(assignments ...Assignment) UpdateQuery {
	q.Assignments = append(q.Assignments, assignments...)
	return q
}

func (q UpdateQuery) From(table Table) UpdateQuery {
	q.FromTable = table
	return q
}

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

func (q UpdateQuery) Where(predicates ...Predicate) UpdateQuery {
	q.WherePredicate.Predicates = append(q.WherePredicate.Predicates, predicates...)
	return q
}

func (q UpdateQuery) Returning(fields ...Field) UpdateQuery {
	q.ReturningFields = append(q.ReturningFields, fields...)
	return q
}

func (q UpdateQuery) ReturningOne() UpdateQuery {
	q.ReturningFields = Fields{FieldLiteral("1")}
	return q
}

func (q UpdateQuery) Returningx(mapper func(*Row), accumulator func()) UpdateQuery {
	q.Mapper = mapper
	q.Accumulator = accumulator
	return q
}

func (q UpdateQuery) ReturningRowx(mapper func(*Row)) UpdateQuery {
	q.Mapper = mapper
	return q
}

func (q UpdateQuery) Fetch(db DB) (err error) {
	q.LogSkip += 1
	return q.FetchContext(nil, db)
}

func (q UpdateQuery) FetchContext(ctx context.Context, db DB) (err error) {
	if db == nil {
		if q.DB == nil {
			return errors.New("DB cannot be nil")
		}
		db = q.DB
	}
	if q.Mapper == nil {
		return fmt.Errorf("Cannot call Fetch without a mapper")
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
				q.Log.Output(q.LogSkip+2, logBuf.String())
			default:
				q.Log.Output(q.LogSkip+1, logBuf.String())
			}
		}
	}()
	r := &Row{}
	q.Mapper(r)
	q.ReturningFields = r.fields
	tmpbuf := &strings.Builder{}
	var tmpargs []interface{}
	q.LogSkip += 1
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
					DollarInterpolate(tmpbuf.String(), tmpargs...) + " => " +
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
				logBuf.WriteString(DollarInterpolate(tmpbuf.String(), tmpargs...))
				logBuf.WriteString(": ")
				logBuf.WriteString(AppendSQLDisplay(r.dest[i]))
			}
		}
		r.index = 0
		q.Mapper(r)
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

func (q UpdateQuery) Exec(db DB, flag ExecFlag) (rowsAffected int64, err error) {
	q.LogSkip += 1
	return q.ExecContext(nil, db, flag)
}

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

func (q UpdateQuery) GetAlias() string {
	return q.Alias
}

func (q UpdateQuery) GetName() string {
	return ""
}

func (q UpdateQuery) NestThis() Query {
	q.Nested = true
	return q
}
