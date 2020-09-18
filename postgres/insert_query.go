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
	DB          DB
	Mapper      func(*Row)
	Accumulator func()
	// Logging
	Log     Logger
	LogFlag LogFlag
	logSkip int
}

func (q InsertQuery) ToSQL() (string, []interface{}) {
	q.logSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

func (q InsertQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	var excludedTableQualifiers []string
	// WITH
	if !q.nested && q.SelectQuery != nil {
		AppendCTEs(buf, args, q.CTEs, q.SelectQuery.FromTable, q.SelectQuery.JoinTables)
	}
	// INSERT INTO
	buf.WriteString("INSERT INTO ")
	if q.IntoTable == nil {
		buf.WriteString("NULL")
	} else {
		q.IntoTable.AppendSQL(buf, args)
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
		q.InsertColumns.AppendSQLExclude(buf, args, excludedTableQualifiers)
		buf.WriteString(")")
	}
	// VALUES/SELECT
	switch {
	case len(q.RowValues) > 0:
		buf.WriteString(" VALUES ")
		q.RowValues.AppendSQL(buf, args)
	case q.SelectQuery != nil:
		buf.WriteString(" ")
		q.SelectQuery.Nested = true
		q.SelectQuery.AppendSQL(buf, args)
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
			q.ConflictFields.AppendSQLExclude(buf, args, excludedTableQualifiers)
			buf.WriteString(")")
			if len(q.ConflictPredicate.Predicates) > 0 {
				buf.WriteString(" WHERE ")
				q.ConflictPredicate.toplevel = true
				q.ConflictPredicate.AppendSQLExclude(buf, args, excludedTableQualifiers)
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
		q.Resolution.AppendSQLExclude(buf, args, excludedTableQualifiers)
		if len(q.ResolutionPredicate.Predicates) > 0 {
			buf.WriteString(" WHERE ")
			q.ResolutionPredicate.toplevel = true
			q.ResolutionPredicate.AppendSQLExclude(buf, args, nil)
		}
	default:
		buf.WriteString(" DO NOTHING")
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

func InsertInto(table BaseTable) InsertQuery {
	return InsertQuery{
		IntoTable: table,
	}
}

func (q InsertQuery) With(ctes ...CTE) InsertQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

func (q InsertQuery) InsertInto(table BaseTable) InsertQuery {
	q.IntoTable = table
	return q
}

func (q InsertQuery) Columns(fields ...Field) InsertQuery {
	q.InsertColumns = fields
	return q
}

func (q InsertQuery) Values(values ...interface{}) InsertQuery {
	q.RowValues = append(q.RowValues, values)
	return q
}

func (q InsertQuery) InsertRow(assignments ...FieldAssignment) InsertQuery {
	fields, values := make([]Field, len(assignments)), make([]interface{}, len(assignments))
	for i, assignment := range assignments {
		fields[i] = assignment.Field
		values[i] = assignment.Value
	}
	if len(q.InsertColumns) == 0 {
		q.InsertColumns = fields
	}
	q.RowValues = append(q.RowValues, values)
	return q
}

func (q InsertQuery) Select(selectQuery SelectQuery) InsertQuery {
	q.SelectQuery = &selectQuery
	return q
}

func (q InsertQuery) OnConflict(fields ...Field) InsertConflict {
	q.HandleConflict = true
	q.ConflictFields = fields
	return InsertConflict{insertQuery: &q}
}

func (q InsertQuery) OnConflictOnConstraint(name string) InsertConflict {
	q.HandleConflict = true
	q.ConflictConstraint = name
	return InsertConflict{insertQuery: &q}
}

type InsertConflict struct{ insertQuery *InsertQuery }

func (c InsertConflict) Where(predicates ...Predicate) InsertConflict {
	c.insertQuery.ConflictPredicate.Predicates = append(c.insertQuery.ConflictPredicate.Predicates, predicates...)
	return c
}

func (c InsertConflict) DoNothing() InsertQuery {
	if c.insertQuery == nil {
		return InsertQuery{}
	}
	return *c.insertQuery
}

func (c InsertConflict) DoUpdateSet(assignments ...Assignment) InsertQuery {
	if c.insertQuery == nil {
		return InsertQuery{}
	}
	c.insertQuery.Resolution = assignments
	return *c.insertQuery
}

func Excluded(field Field) CustomField {
	return CustomField{
		Format: "EXCLUDED." + field.GetName(),
	}
}

func (q InsertQuery) Where(predicates ...Predicate) InsertQuery {
	q.ResolutionPredicate.Predicates = append(q.ResolutionPredicate.Predicates, predicates...)
	return q
}

func (q InsertQuery) Returning(fields ...Field) InsertQuery {
	q.ReturningFields = append(q.ReturningFields, fields...)
	return q
}

func (q InsertQuery) ReturningOne() InsertQuery {
	q.ReturningFields = Fields{FieldLiteral("1")}
	return q
}

func (q InsertQuery) Returningx(mapper func(*Row), accumulator func()) InsertQuery {
	q.Mapper = mapper
	q.Accumulator = accumulator
	return q
}

func (q InsertQuery) ReturningRowx(mapper func(*Row)) InsertQuery {
	q.Mapper = mapper
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
	if q.Mapper == nil {
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
	q.Mapper(r)
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

func (q InsertQuery) NestThis() Query {
	q.nested = true
	return q
}
