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

// InsertQuery represents an INSERT query.
type InsertQuery struct {
	nested bool
	// WITH
	CTEs []CTE
	// INSERT INTO
	IntoTable     BaseTable
	InsertColumns Fields
	// VALUES
	RowValues RowValues
	// SELECT
	SelectQuery *SelectQuery
	// ON CONFLICT
	HandleConflict      bool
	ConflictFields      Fields
	ConflictPredicate   VariadicPredicate
	ConflictConstraint  string
	Resolution          Assignments
	ResolutionPredicate VariadicPredicate
	// RETURNING
	ReturningFields Fields
	// DB
	DB           DB
	ColumnMapper func(*Column)
	RowMapper    func(*Row)
	Accumulator  func()
	// Logging
	Log     Logger
	LogFlag LogFlag
	logSkip int
}

// ToSQL marshals the InsertQuery into a query string and args slice.
func (q InsertQuery) ToSQL() (string, []interface{}) {
	q.logSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args, nil)
	return buf.String(), args
}

// AppendSQL marshals the InsertQuery into a buffer and args slice.
func (q InsertQuery) AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int) {
	var excludedTableQualifiers []string
	if q.ColumnMapper != nil {
		col := &Column{mode: colmodeInsert}
		q.ColumnMapper(col)
		q.InsertColumns = col.insertColumns
		q.RowValues = col.rowValues
	}
	// WITH
	if !q.nested && q.SelectQuery != nil {
		appendCTEs(buf, args, q.CTEs, q.SelectQuery.FromTable, q.SelectQuery.JoinTables)
	}
	// INSERT INTO
	buf.WriteString("INSERT INTO ")
	if q.IntoTable == nil {
		buf.WriteString("NULL")
	} else {
		q.IntoTable.AppendSQL(buf, args, nil)
		name := q.IntoTable.GetName()
		alias := q.IntoTable.GetAlias()
		if alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(alias)
			excludedTableQualifiers = append(excludedTableQualifiers, alias)
		} else {
			excludedTableQualifiers = append(excludedTableQualifiers, name)
		}
	}
	if len(q.InsertColumns) > 0 {
		buf.WriteString(" (")
		q.InsertColumns.AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
		buf.WriteString(")")
	}
	// VALUES/SELECT
	switch {
	case len(q.RowValues) > 0:
		buf.WriteString(" VALUES ")
		q.RowValues.AppendSQL(buf, args, nil)
	case q.SelectQuery != nil:
		buf.WriteString(" ")
		q.SelectQuery.nested = true
		q.SelectQuery.AppendSQL(buf, args, nil)
	}
	// ON CONFLICT
	var noConflict bool
	switch {
	case q.HandleConflict:
		buf.WriteString(" ON CONFLICT")
		switch {
		case q.ConflictConstraint != "":
			buf.WriteString(" ON CONSTRAINT ")
			buf.WriteString(q.ConflictConstraint)
		case len(q.ConflictFields) > 0:
			buf.WriteString(" (")
			q.ConflictFields.AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
			buf.WriteString(")")
			if len(q.ConflictPredicate.Predicates) > 0 {
				buf.WriteString(" WHERE ")
				q.ConflictPredicate.toplevel = true
				q.ConflictPredicate.AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
			}
		}
	default:
		noConflict = true
	}
	switch {
	case noConflict:
		break
	case len(q.Resolution) > 0:
		buf.WriteString(" DO UPDATE SET ")
		q.Resolution.AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
		if len(q.ResolutionPredicate.Predicates) > 0 {
			buf.WriteString(" WHERE ")
			q.ResolutionPredicate.toplevel = true
			q.ResolutionPredicate.AppendSQLExclude(buf, args, nil, nil)
		}
	default:
		buf.WriteString(" DO NOTHING")
	}
	// RETURNING
	if len(q.ReturningFields) > 0 {
		buf.WriteString(" RETURNING ")
		q.ReturningFields.AppendSQLExcludeWithAlias(buf, args, nil, nil)
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

// InsertInto creates a new InsertQuery.
func InsertInto(table BaseTable) InsertQuery {
	return InsertQuery{
		IntoTable: table,
	}
}

// With appends a list of CTEs into the InsertQuery.
func (q InsertQuery) With(ctes ...CTE) InsertQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

// InsertInto sets the insert table for the InsertQuery.
func (q InsertQuery) InsertInto(table BaseTable) InsertQuery {
	q.IntoTable = table
	return q
}

// Columns sets the insert columns for the InsertQuery.
func (q InsertQuery) Columns(fields ...Field) InsertQuery {
	q.InsertColumns = fields
	return q
}

// Values appends a new RowValue to the InsertQuery.
func (q InsertQuery) Values(values ...interface{}) InsertQuery {
	q.RowValues = append(q.RowValues, values)
	return q
}

// Valuesx sets the column mapper for the InsertQuery.
func (q InsertQuery) Valuesx(mapper func(*Column)) InsertQuery {
	q.ColumnMapper = mapper
	return q
}

// Select adds a SelectQuery to the InsertQuery.
func (q InsertQuery) Select(selectQuery SelectQuery) InsertQuery {
	q.SelectQuery = &selectQuery
	return q
}

// OnConflict specifies which Fields may potentially experience a conflict.
func (q InsertQuery) OnConflict(fields ...Field) InsertConflict {
	q.HandleConflict = true
	q.ConflictFields = fields
	return InsertConflict{insertQuery: &q}
}

// OnConflict specifies which constraint may potentially experience a conflict.
func (q InsertQuery) OnConflictOnConstraint(name string) InsertConflict {
	q.HandleConflict = true
	q.ConflictConstraint = name
	return InsertConflict{insertQuery: &q}
}

// InsertConflict holds the intermediate state of an InsertQuery that may
// experience a conflict.
type InsertConflict struct{ insertQuery *InsertQuery }

// Where appends the predicates to the WHERE clause of the InsertQuery conflict.
func (c InsertConflict) Where(predicates ...Predicate) InsertConflict {
	c.insertQuery.ConflictPredicate.Predicates = append(c.insertQuery.ConflictPredicate.Predicates, predicates...)
	return c
}

// DoNothing indicates that nothing should be done in case of any conflicts.
func (c InsertConflict) DoNothing() InsertQuery {
	if c.insertQuery == nil {
		return InsertQuery{}
	}
	return *c.insertQuery
}

// DoUpdateSet specifies the assignments to be done in case of a conflict.
func (c InsertConflict) DoUpdateSet(assignments ...Assignment) InsertQuery {
	if c.insertQuery == nil {
		return InsertQuery{}
	}
	c.insertQuery.Resolution = assignments
	return *c.insertQuery
}

// Excluded wraps a field to simulate the EXCLUDED.field Postgres construct for the
// ON CONFLICT DO UPDATE SET clause.
func Excluded(field Field) CustomField {
	return CustomField{
		Format: "EXCLUDED." + field.GetName(),
	}
}

// Where appends the predicates to the WHERE clause of InsertQuery conflict resolution.
func (q InsertQuery) Where(predicates ...Predicate) InsertQuery {
	q.ResolutionPredicate.Predicates = append(q.ResolutionPredicate.Predicates, predicates...)
	return q
}

// Returning appends the fields to the RETURNING clause of the InsertQuery.
func (q InsertQuery) Returning(fields ...Field) InsertQuery {
	q.ReturningFields = append(q.ReturningFields, fields...)
	return q
}

// ReturningOne sets the RETURNING clause to RETURNING 1 in the InsertQuery.
func (q InsertQuery) ReturningOne() InsertQuery {
	q.ReturningFields = Fields{FieldLiteral("1")}
	return q
}

// Returningx sets the rowmapper and accumulator function of the InsertQuery.
func (q InsertQuery) Returningx(mapper func(*Row), accumulator func()) InsertQuery {
	q.RowMapper = mapper
	q.Accumulator = accumulator
	return q
}

func (q InsertQuery) ReturningRowx(mapper func(*Row)) InsertQuery {
	q.RowMapper = mapper
	return q
}

func (q InsertQuery) Fetch(db DB) (err error) {
	q.logSkip += 1
	return q.FetchContext(nil, db)
}

func (q InsertQuery) FetchContext(ctx context.Context, db DB) (err error) {
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
	q.AppendSQL(tmpbuf, &tmpargs, nil)
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
				r.fields[i].AppendSQLExclude(tmpbuf, &tmpargs, nil, nil)
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
				r.fields[i].AppendSQLExclude(tmpbuf, &tmpargs, nil, nil)
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

func (q InsertQuery) Exec(db DB, flag ExecFlag) (rowsAffected int64, err error) {
	q.logSkip += 1
	return q.ExecContext(nil, db, flag)
}

func (q InsertQuery) ExecContext(ctx context.Context, db DB, flag ExecFlag) (rowsAffected int64, err error) {
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
			logBuf.WriteString("\n(Inserted ")
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

// NestThis indicates to the InsertQuery that it is nested.
func (q InsertQuery) NestThis() Query {
	q.nested = true
	return q
}
