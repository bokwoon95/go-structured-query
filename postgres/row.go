package sq

import (
	"database/sql"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// ExitCode represents a reason for terminating the rows.Next() loop.
type ExitCode int

// ExitCodes
const (
	ExitPeacefully ExitCode = iota
)

// Error implements the error interface.
func (e ExitCode) Error() string {
	return "exit " + strconv.Itoa(int(e))
}

// Row represents the state of a row after a call to rows.Next().
type Row struct {
	rows    *sql.Rows
	index   int
	fields  []Field
	dest    []interface{}
	tmpdest []interface{}
}

/* custom */

// ScanInto scans the field into a dest, where dest is a pointer.
func (r *Row) ScanInto(dest interface{}, field Field) {
	if r.rows == nil {
		r.fields = append(r.fields, field)
		switch dest.(type) {
		case *bool, *sql.NullBool:
			r.dest = append(r.dest, &sql.NullBool{})
		case *float64, *sql.NullFloat64:
			r.dest = append(r.dest, &sql.NullFloat64{})
		case *int32, *sql.NullInt32:
			r.dest = append(r.dest, &sql.NullInt32{})
		case *int, *int64, *sql.NullInt64:
			r.dest = append(r.dest, &sql.NullInt64{})
		case *string, *sql.NullString:
			r.dest = append(r.dest, &sql.NullString{})
		case *time.Time, *sql.NullTime:
			r.dest = append(r.dest, &sql.NullTime{})
		default:
			r.dest = append(r.dest, dest)
		}
		return
	}
	switch ptr := dest.(type) {
	case *bool:
		nullbool := r.dest[r.index].(*sql.NullBool)
		*ptr = nullbool.Bool
	case *sql.NullBool:
		nullbool := r.dest[r.index].(*sql.NullBool)
		*ptr = *nullbool
	case *float64:
		nullfloat64 := r.dest[r.index].(*sql.NullFloat64)
		*ptr = nullfloat64.Float64
	case *sql.NullFloat64:
		nullfloat64 := r.dest[r.index].(*sql.NullFloat64)
		*ptr = *nullfloat64
	case *int:
		nullint64 := r.dest[r.index].(*sql.NullInt64)
		*ptr = int(nullint64.Int64)
	case *int32:
		nullint32 := r.dest[r.index].(*sql.NullInt32)
		*ptr = nullint32.Int32
	case *sql.NullInt32:
		nullint32 := r.dest[r.index].(*sql.NullInt32)
		*ptr = *nullint32
	case *int64:
		nullint64 := r.dest[r.index].(*sql.NullInt64)
		*ptr = nullint64.Int64
	case *sql.NullInt64:
		nullint64 := r.dest[r.index].(*sql.NullInt64)
		*ptr = *nullint64
	case *string:
		nullstring := r.dest[r.index].(*sql.NullString)
		*ptr = nullstring.String
	case *sql.NullString:
		nullstring := r.dest[r.index].(*sql.NullString)
		*ptr = *nullstring
	case *time.Time:
		nulltime := r.dest[r.index].(*sql.NullTime)
		*ptr = nulltime.Time
	case *sql.NullTime:
		nulltime := r.dest[r.index].(*sql.NullTime)
		*ptr = *nulltime
	default:
		var nothing interface{}
		if len(r.tmpdest) != len(r.dest) {
			r.tmpdest = make([]interface{}, len(r.dest))
			for i := range r.tmpdest {
				r.tmpdest[i] = &nothing
			}
		}
		r.tmpdest[r.index] = dest
		err := r.rows.Scan(r.tmpdest...)
		if err != nil {
			_, sourcefile, linenbr, _ := runtime.Caller(1)
			panic(fmt.Errorf("row.ScanInto failed on %s:%d: %w", sourcefile, linenbr, err))
		}
		r.tmpdest[r.index] = &nothing
	}
	r.index++
}

// ScanArray accepts a pointer to a slice and scans a postgres array into it.
// Only []bool, []float64, []int64 or []string slices are supported.
func (r *Row) ScanArray(slice interface{}, field Field) {
	var nothing interface{}
	if r.rows == nil {
		r.fields = append(r.fields, field)
		r.dest = append(r.dest, pq.Array(slice))
		return
	}
	if len(r.tmpdest) != len(r.dest) {
		r.tmpdest = make([]interface{}, len(r.dest))
		for i := range r.tmpdest {
			r.tmpdest[i] = &nothing
		}
	}
	r.tmpdest[r.index] = pq.Array(slice)
	err := r.rows.Scan(r.tmpdest...)
	if err != nil {
		_, sourcefile, linenbr, _ := runtime.Caller(1)
		panic(fmt.Errorf("row.ScanArray failed on %s:%d: %w", sourcefile, linenbr, err))
	}
	r.tmpdest[r.index] = &nothing
	r.index++
}

/* bool */

// Bool returns the bool value of the Predicate. BooleanFields are considered
// predicates, so you can use them here.
func (r *Row) Bool(predicate Predicate) bool {
	return r.NullBool(predicate).Bool
}

// BoolValid returns the bool value indicating if the Predicate is non-NULL.
// BooleanFields are considered Predicates, so you can use them here.
func (r *Row) BoolValid(predicate Predicate) bool {
	return r.NullBool(predicate).Valid
}

// NullBool returns the sql.NullBool value of the Predicate.
func (r *Row) NullBool(predicate Predicate) sql.NullBool {
	if r.rows == nil {
		buf := &strings.Builder{}
		var args []interface{}
		predicate.AppendSQLExclude(buf, &args, nil, nil)
		r.fields = append(r.fields, CustomPredicate{
			Format: buf.String(),
			Values: args,
		})
		r.dest = append(r.dest, &sql.NullBool{})
		return sql.NullBool{}
	}
	nullbool := r.dest[r.index].(*sql.NullBool)
	r.index++
	return *nullbool
}

/* float64 */

// Float64 returns the float64 value of the NumberField.
func (r *Row) Float64(field NumberField) float64 {
	return rowNullFloat64(r, field).Float64
}

// Float64Valid returns the bool value indicating if the NumberField is
// non-NULL.
func (r *Row) Float64Valid(field NumberField) bool {
	return rowNullFloat64(r, field).Valid
}

// NullFloat64 returns the sql.NullFloat64 value of the NumberField.
func (r *Row) NullFloat64(field NumberField) sql.NullFloat64 {
	return rowNullFloat64(r, field)
}

// rowFloat64 returns the float64 value of a Field
func rowFloat64(r *Row, field Field) float64 {
	return rowNullFloat64(r, field).Float64
}

// rowFloat64Valid returns the bool value indicating if the Field is non-NULL.
func rowFloat64Valid(r *Row, field Field) bool {
	return rowNullFloat64(r, field).Valid
}

// rowNullFloat64 returns the sql.NullFloat64 value of the Field.
func rowNullFloat64(r *Row, field Field) sql.NullFloat64 {
	if r.rows == nil {
		r.fields = append(r.fields, field)
		r.dest = append(r.dest, &sql.NullFloat64{})
		return sql.NullFloat64{}
	}
	nullfloat64 := r.dest[r.index].(*sql.NullFloat64)
	r.index++
	return *nullfloat64
}

/* int */

// Int returns the int value of the NumberField.
func (r *Row) Int(field NumberField) int {
	return int(rowNullInt64(r, field).Int64)
}

// IntValid returns the bool value indicating if the NumberField is non-NULL.
func (r *Row) IntValid(field NumberField) bool {
	return rowNullInt64(r, field).Valid
}

// rowInt returns the int value of the Field.
func rowInt(r *Row, field Field) int {
	return int(rowNullInt64(r, field).Int64)
}

// rowIntValid returns the bool value indicating if the Field is non-NULL.
func rowIntValid(r *Row, field Field) bool {
	return rowNullInt64(r, field).Valid
}

/* int64 */

// Int64 returns the int64 value of the NumberField.
func (r *Row) Int64(field NumberField) int64 {
	return rowNullInt64(r, field).Int64
}

// Int64Valid returns the bool value indicating if the NumberField is non-NULL.
func (r *Row) Int64Valid(field NumberField) bool {
	return r.NullInt64(field).Valid
}

// NullInt64 returns the sql.NullInt64 value of the NumberField.
func (r *Row) NullInt64(field NumberField) sql.NullInt64 {
	return rowNullInt64(r, field)
}

// rowInt64 returns the int64 value of the Field.
func rowInt64(r *Row, field NumberField) int64 {
	return rowNullInt64(r, field).Int64
}

// rowInt64Valid returns the bool value indicating if the Field is non-NULL.
func rowInt64Valid(r *Row, field NumberField) bool {
	return rowNullInt64(r, field).Valid
}

// rowNullInt64 returns the sql.NullInt64 value of the Field.
func rowNullInt64(r *Row, field Field) sql.NullInt64 {
	if r.rows == nil {
		r.fields = append(r.fields, field)
		r.dest = append(r.dest, &sql.NullInt64{})
		return sql.NullInt64{}
	}
	nullint64 := r.dest[r.index].(*sql.NullInt64)
	r.index++
	return *nullint64
}

/* string */

// String returns the string value of the StringField.
func (r *Row) String(field StringField) string {
	return rowNullString(r, field).String
}

// StringValid returns the bool value indicating if the StringField is
// non-NULL.
func (r *Row) StringValid(field StringField) bool {
	return rowNullString(r, field).Valid
}

// NullString returns the sql.NullString value of the StringField.
func (r *Row) NullString(field StringField) sql.NullString {
	return rowNullString(r, field)
}

// rowString returns the string value of the Field.
func rowString(r *Row, field Field) string {
	return rowNullString(r, field).String
}

// rowStringValid returns the bool value indicating if the Field is non-NULL.
func rowStringValid(r *Row, field Field) bool {
	return rowNullString(r, field).Valid
}

// rowNullString returns the sql.NullString value of the Field.
func rowNullString(r *Row, field Field) sql.NullString {
	if r.rows == nil {
		r.fields = append(r.fields, field)
		r.dest = append(r.dest, &sql.NullString{})
		return sql.NullString{}
	}
	nullstring := r.dest[r.index].(*sql.NullString)
	r.index++
	return *nullstring
}

/* time.Time */

// Time returns the time.Time value of the TimeField.
func (r *Row) Time(field TimeField) time.Time {
	return rowNullTime(r, field).Time
}

// TimeValid returns a bool value indicating if the TimeField is non-NULL.
func (r *Row) TimeValid(field TimeField) bool {
	return rowNullTime(r, field).Valid
}

// NullTime returns the sql.NullTime value of the TimeField.
func (r *Row) NullTime(field TimeField) sql.NullTime {
	return rowNullTime(r, field)
}

// rowTime returns the time.Time value of the Field.
func rowTime(r *Row, field Field) time.Time {
	return rowNullTime(r, field).Time
}

// rowTimeValid returns a bool value indicating if the Field is non-NULL.
func rowTimeValid(r *Row, field Field) bool {
	return rowNullTime(r, field).Valid
}

// rowNullTime returns the sql.NullTime value of the Field.
func rowNullTime(r *Row, field Field) sql.NullTime {
	if r.rows == nil {
		r.fields = append(r.fields, field)
		r.dest = append(r.dest, &sql.NullTime{})
		return sql.NullTime{}
	}
	nulltime := r.dest[r.index].(*sql.NullTime)
	r.index++
	return *nulltime
}

// UUID returns the [16]byte value of the UUIDField
func (r *Row) UUID(field UUIDField) [16]byte {
	if r.rows == nil {
		r.fields = append(r.fields, field)
		r.dest = append(r.dest, &[]byte{})
		return [16]byte{}
	}

	uuid := r.dest[r.index].(*[]byte)

	var dest [16]byte
	copy(dest[:], *uuid)

	r.index++
	return dest
}
