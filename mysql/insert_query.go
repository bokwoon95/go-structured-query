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

// InsertQuery represents an INSERT query.
type InsertQuery struct {
	Nested bool
	Alias  string
	// INSERT INTO
	Ignore        bool
	IntoTable     BaseTable
	InsertColumns Fields
	// VALUES
	RowValues RowValues
	// SELECT
	SelectQuery *SelectQuery
	// ON DUPLICATE KEY
	Resolution Assignments
	// DB
	DB DB
	// Logging
	Log     Logger
	LogFlag LogFlag
	LogSkip int
}

// ToSQL marshals the InsertQuery into a query string and args slice.
func (q InsertQuery) ToSQL() (string, []interface{}) {
	q.LogSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args)
	return buf.String(), args
}

// AppendSQL marshals the InsertQuery into a buffer and args slice.
func (q InsertQuery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	var excludedTableQualifiers []string
	// INSERT INTO
	if q.Ignore {
		buf.WriteString("INSERT IGNORE INTO ")
	} else {
		buf.WriteString("INSERT INTO ")
	}
	if q.IntoTable == nil {
		buf.WriteString("NULL")
	} else {
		q.IntoTable.AppendSQL(buf, args)
		name := q.IntoTable.GetName()
		alias := q.IntoTable.GetAlias()
		if alias != "" {
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
	// ON DUPLICATE KEY UPDATE
	if len(q.Resolution) > 0 {
		buf.WriteString(" ON DUPLICATE KEY UPDATE ")
		q.Resolution.AppendSQLExclude(buf, args, excludedTableQualifiers)
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

// InsertInto creates a new InsertQuery.
func InsertInto(table BaseTable) InsertQuery {
	return InsertQuery{
		IntoTable: table,
		Alias:     RandomString(8),
	}
}

// InsertIgnoreInto creates a new InsertQuery.
func InsertIgnoreInto(table BaseTable) InsertQuery {
	return InsertQuery{
		Ignore:    true,
		IntoTable: table,
		Alias:     RandomString(8),
	}
}

// InsertInto sets the insert table for the InsertQuery.
func (q InsertQuery) InsertInto(table BaseTable) InsertQuery {
	q.IntoTable = table
	return q
}

// InsertIgnoreInto sets the insert table for the InsertQuery.
func (q InsertQuery) InsertIgnoreInto(table BaseTable) InsertQuery {
	q.Ignore = true
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

// InsertRow appends a new RowValue to the InsertQuery. It also sets the insert
// columns if not already set.
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

// Select adds a SelectQuery to the InsertQuery.
func (q InsertQuery) Select(selectQuery SelectQuery) InsertQuery {
	q.SelectQuery = &selectQuery
	return q
}

// OnDuplicateKeyUpdate sets the assignments done on duplicate key for the
// InsertQuery.
func (q InsertQuery) OnDuplicateKeyUpdate(assignments ...Assignment) InsertQuery {
	q.Resolution = assignments
	return q
}

// Values wraps a field to simulate the VALUES(field) MySQL construct for the
// ON DUPLICATE KEY UPDATE clause.
func Values(field Field) CustomField {
	return CustomField{
		Format: "VALUES(" + field.GetName() + ")",
	}
}

// Exec will execute the InsertQuery with the given DB. It will only compute
// the lastInsertID if the ElastInsertID ExecFlag is passed to it. It will only
// compute the rowsAffected if the ErowsAffected Execflag is passed to it. To
// compute both, bitwise or the flags together i.e.
// ElastInsertID|ErowsAffected.
func (q InsertQuery) Exec(db DB, flag ExecFlag) (lastInsertID, rowsAffected int64, err error) {
	q.LogSkip += 1
	return q.ExecContext(nil, db, flag)
}

// ExecContext will execute the InsertQuery with the given DB and context. It
// will only compute the lastInsertID if the ElastInsertID ExecFlag is passed
// to it. It will only compute the rowsAffected if the ErowsAffected Execflag
// is passed to it. To compute both, bitwise or the flags together i.e.
// ElastInsertID|ErowsAffected.
func (q InsertQuery) ExecContext(ctx context.Context, db DB, flag ExecFlag) (lastInsertID, rowsAffected int64, err error) {
	if db == nil {
		if q.DB == nil {
			return lastInsertID, rowsAffected, errors.New("DB cannot be nil")
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
		return lastInsertID, rowsAffected, err
	}
	if res != nil && ElastInsertID&flag != 0 {
		lastInsertID, err = res.LastInsertId()
		if err != nil {
			return lastInsertID, rowsAffected, err
		}
	}
	if res != nil && ErowsAffected&flag != 0 {
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return lastInsertID, rowsAffected, err
		}
	}
	return lastInsertID, rowsAffected, nil
}

// NestThis indicates to the InsertQuery that it is nested.
func (q InsertQuery) NestThis() Query {
	q.Nested = true
	return q
}
