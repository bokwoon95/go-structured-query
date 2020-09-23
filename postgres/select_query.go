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

// SelectType represents the various SQL selects.
type SelectType string

// SelectTypes
const (
	SelectTypeDefault    SelectType = "SELECT"
	SelectTypeDistinct   SelectType = "SELECT DISTINCT"
	SelectTypeDistinctOn SelectType = "SELECT DISTINCT ON"
)

// SelectQuery represents a SELECT query.
type SelectQuery struct {
	nested bool
	// WITH
	CTEs []CTE
	// SELECT
	SelectType   SelectType
	SelectFields Fields
	DistinctOn   Fields
	// FROM
	FromTable  Table
	JoinTables JoinTables
	// WHERE
	WherePredicate VariadicPredicate
	// GROUP BY
	GroupByFields Fields
	// HAVING
	HavingPredicate VariadicPredicate
	// WINDOW
	Windows Windows
	// ORDER BY
	OrderByFields Fields
	// LIMIT
	LimitValue *int64
	// OFFSET
	OffsetValue *int64
	// DB
	DB          DB
	RowMapper   func(*Row)
	Accumulator func()
	// Logging
	Log     Logger
	LogFlag LogFlag
	logSkip int
}

// ToSQL marshals the SelectQuery into a query string and args slice.
func (q SelectQuery) ToSQL() (string, []interface{}) {
	q.logSkip += 1
	buf := &strings.Builder{}
	var args []interface{}
	q.AppendSQL(buf, &args, nil)
	return buf.String(), args
}

// AppendSQL marshals the SelectQuery into a buffer and args slice.
func (q SelectQuery) AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int) {
	// WITH
	if !q.nested {
		appendCTEs(buf, args, q.CTEs, q.FromTable, q.JoinTables)
	}
	// SELECT
	if q.SelectType == "" {
		q.SelectType = SelectTypeDefault
	}
	buf.WriteString(string(q.SelectType))
	if q.SelectType == SelectTypeDistinctOn {
		buf.WriteString(" (")
		q.DistinctOn.AppendSQLExclude(buf, args, nil, nil)
		buf.WriteString(")")
	}
	if len(q.SelectFields) > 0 {
		buf.WriteString(" ")
		q.SelectFields.AppendSQLExcludeWithAlias(buf, args, nil, nil)
	}
	// FROM
	if q.FromTable != nil {
		buf.WriteString(" FROM ")
		switch v := q.FromTable.(type) {
		case Query:
			buf.WriteString("(")
			v.NestThis().AppendSQL(buf, args, nil)
			buf.WriteString(")")
		default:
			q.FromTable.AppendSQL(buf, args, nil)
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
		q.JoinTables.AppendSQL(buf, args, nil)
	}
	// WHERE
	if len(q.WherePredicate.Predicates) > 0 {
		buf.WriteString(" WHERE ")
		q.WherePredicate.toplevel = true
		q.WherePredicate.AppendSQLExclude(buf, args, nil, nil)
	}
	// GROUP BY
	if len(q.GroupByFields) > 0 {
		buf.WriteString(" GROUP BY ")
		q.GroupByFields.AppendSQLExclude(buf, args, nil, nil)
	}
	// HAVING
	if len(q.HavingPredicate.Predicates) > 0 {
		buf.WriteString(" HAVING ")
		q.HavingPredicate.toplevel = true
		q.HavingPredicate.AppendSQLExclude(buf, args, nil, nil)
	}
	// WINDOW
	if len(q.Windows) > 0 {
		buf.WriteString(" WINDOW ")
		q.Windows.AppendSQL(buf, args, nil)
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
	// OFFSET
	if q.OffsetValue != nil {
		buf.WriteString(" OFFSET ?")
		if *q.OffsetValue < 0 {
			*q.OffsetValue = -*q.OffsetValue
		}
		*args = append(*args, *q.OffsetValue)
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

// From creates a new SelectQuery.
func From(table Table) SelectQuery {
	return SelectQuery{
		FromTable: table,
	}
}

// Select creates a new SelectQuery.
func Select(fields ...Field) SelectQuery {
	return SelectQuery{
		SelectFields: fields,
	}
}

// SelectOne creates a new SelectQuery.
func SelectOne() SelectQuery {
	return SelectQuery{
		SelectFields: Fields{FieldLiteral("1")},
	}
}

// SelectDistinct creates a new SelectQuery.
func SelectDistinct(fields ...Field) SelectQuery {
	return SelectQuery{
		SelectType:   SelectTypeDistinct,
		SelectFields: fields,
	}
}

// SelectDistinctOn creates a new SelectQuery.
func SelectDistinctOn(distinctFields ...Field) func(...Field) SelectQuery {
	return func(fields ...Field) SelectQuery {
		return SelectQuery{
			SelectType:   SelectTypeDistinctOn,
			SelectFields: fields,
			DistinctOn:   distinctFields,
		}
	}
}

// Selectx creates a new SelectQuery.
func Selectx(mapper func(*Row), accumulator func()) SelectQuery {
	return SelectQuery{
		RowMapper:   mapper,
		Accumulator: accumulator,
	}
}

// SelectRowx creates a new SelectQuery.
func SelectRowx(mapper func(*Row)) SelectQuery {
	return SelectQuery{
		RowMapper: mapper,
	}
}

// With appends a list of CTEs into the SelectQuery.
func (q SelectQuery) With(ctes ...CTE) SelectQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

// Select adds the fields to the SelectFields in the SelectQuery.
func (q SelectQuery) Select(fields ...Field) SelectQuery {
	q.SelectFields = append(q.SelectFields, fields...)
	return q
}

// SelectOne sets the SELECT clause to SELECT 1.
func (q SelectQuery) SelectOne() SelectQuery {
	q.SelectFields = Fields{FieldLiteral("1")}
	return q
}

// SelectAll sets the SELECT clause to SELECT *.
func (q SelectQuery) SelectAll() SelectQuery {
	q.SelectFields = Fields{FieldLiteral("*")}
	return q
}

// SelectCount sets the SELECT clause to SELECT COUNT(*).
func (q SelectQuery) SelectCount() SelectQuery {
	q.SelectFields = Fields{FieldLiteral("COUNT(*)")}
	return q
}

// SelectDistinct adds the fields to the SelectFields in the SelectQuery.
func (q SelectQuery) SelectDistinct(fields ...Field) SelectQuery {
	q.SelectType = SelectTypeDistinct
	q.SelectFields = fields
	return q
}

// SelectDistinctOn adds the distinctFields to the DistinctOn fields and the
// fields to the SelectFields in the SelectQuery.
func (q SelectQuery) SelectDistinctOn(distinctFields ...Field) func(...Field) SelectQuery {
	return func(fields ...Field) SelectQuery {
		q.SelectType = SelectTypeDistinctOn
		q.SelectFields = fields
		q.DistinctOn = distinctFields
		return q
	}
}

// From sets the table in the SelectQuery.
func (q SelectQuery) From(table Table) SelectQuery {
	q.FromTable = table
	return q
}

// Join joins a new table to the SelectQuery based on the predicates.
func (q SelectQuery) Join(table Table, predicate Predicate, predicates ...Predicate) SelectQuery {
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

// LeftJoin left joins a new table to the SelectQuery based on the predicates.
func (q SelectQuery) LeftJoin(table Table, predicate Predicate, predicates ...Predicate) SelectQuery {
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

// RightJoin right joins a new table to the SelectQuery based on the predicates.
func (q SelectQuery) RightJoin(table Table, predicate Predicate, predicates ...Predicate) SelectQuery {
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

// FullJoin full joins a table to the SelectQuery based on the predicates.
func (q SelectQuery) FullJoin(table Table, predicate Predicate, predicates ...Predicate) SelectQuery {
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

// CustomJoin custom joins a table to the SelectQuery. The join type can be
// specified with a string, e.g. "CROSS JOIN".
func (q SelectQuery) CustomJoin(joinType JoinType, table Table, predicates ...Predicate) SelectQuery {
	q.JoinTables = append(q.JoinTables, JoinTable{
		JoinType: joinType,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	})
	return q
}

// Where appends the predicates to the WHERE clause in the SelectQuery.
func (q SelectQuery) Where(predicates ...Predicate) SelectQuery {
	q.WherePredicate.Predicates = append(q.WherePredicate.Predicates, predicates...)
	return q
}

// GroupBy appends the fields to the GROUP BY clause in the SelectQuery.
func (q SelectQuery) GroupBy(fields ...Field) SelectQuery {
	q.GroupByFields = append(q.GroupByFields, fields...)
	return q
}

// Having appends the predicates to the HAVING clause in the SelectQuery.
func (q SelectQuery) Having(predicates ...Predicate) SelectQuery {
	q.HavingPredicate.Predicates = append(q.HavingPredicate.Predicates, predicates...)
	return q
}

// Window appends the windows to the WINDOW clause in the SelectQuery.
func (q SelectQuery) Window(windows ...Window) SelectQuery {
	q.Windows = append(q.Windows, windows...)
	return q
}

// OrderBy appends the fields to the ORDER BY clause in the SelectQuery.
func (q SelectQuery) OrderBy(fields ...Field) SelectQuery {
	q.OrderByFields = append(q.OrderByFields, fields...)
	return q
}

// Limit sets the limit in the SelectQuery.
func (q SelectQuery) Limit(limit int) SelectQuery {
	num := int64(limit)
	q.LimitValue = &num
	return q
}

// Offset sets the offset in the SelectQuery.
func (q SelectQuery) Offset(offset int) SelectQuery {
	num := int64(offset)
	q.OffsetValue = &num
	return q
}

// Selectx sets the mapper function and accumulator function in the SelectQuery.
func (q SelectQuery) Selectx(mapper func(*Row), accumulator func()) SelectQuery {
	q.RowMapper = mapper
	q.Accumulator = accumulator
	return q
}

// SelectRowx sets the mapper function in the SelectQuery.
func (q SelectQuery) SelectRowx(mapper func(*Row)) SelectQuery {
	q.RowMapper = mapper
	return q
}

// Fetch will run SelectQuery with the given DB. It then maps the results based
// on the mapper function (and optionally runs the accumulator function).
func (q SelectQuery) Fetch(db DB) (err error) {
	q.logSkip += 1
	return q.FetchContext(nil, db)
}

// FetchContext will run SelectQuery with the given DB and context. It then
// maps the results based on the mapper function (and optionally runs the
// accumulator function).
func (q SelectQuery) FetchContext(ctx context.Context, db DB) (err error) {
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
	q.SelectFields = r.fields
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

// Exec will execute the SelectQuery with the given DB. It will only compute
// the rowsAffected if the ErowsAffected Execflag is passed to it.
func (q SelectQuery) Exec(db DB, flag ExecFlag) (rowsAffected int64, err error) {
	q.logSkip += 1
	return q.ExecContext(nil, db, flag)
}

// ExecContext will execute the SelectQuery with the given DB and context. It will
// only compute the rowsAffected if the ErowsAffected Execflag is passed to it.
func (q SelectQuery) ExecContext(ctx context.Context, db DB, flag ExecFlag) (rowsAffected int64, err error) {
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
			logBuf.WriteString("\n(Selected ")
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

// NestThis indicates to the SelectQuery that it is nested.
func (q SelectQuery) NestThis() Query {
	q.nested = true
	return q
}
