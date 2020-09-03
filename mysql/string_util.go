package sq

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ExpandValues will expand each value one by one into successive question mark
// ? placeholders in the format string, writing the results into the buffer and
// args slice. It propagates the excludedTableQualifiers down to its child elements.
func ExpandValues(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string, format string, values []interface{}) {
	for i := strings.Index(format, "?"); i >= 0 && len(values) > 0; i = strings.Index(format, "?") {
		buf.WriteString(format[:i])
		// TODO: I don't know if ?? should be unescaped to ?
		// if len(format[i:]) > 1 && format[i:i+2] == "??" {
		// 	buf.WriteString("?")
		// 	format = format[i+2:]
		// 	continue
		// }
		AppendSQLValue(buf, args, excludedTableQualifiers, values[0])
		format = format[i+1:]
		values = values[1:]
	}
	buf.WriteString(format)
}

// AppendSQLValue will write the SQL representation of the interface{} value
// into the buffer and args slice. It propagates excludedTableQualifiers where
// relevant.
func AppendSQLValue(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string, value interface{}) {
	switch v := value.(type) {
	case nil:
		buf.WriteString("NULL")
		return
	case interface {
		AppendSQLExclude(*strings.Builder, *[]interface{}, []string)
	}:
		v.AppendSQLExclude(buf, args, excludedTableQualifiers)
		return
	case interface {
		AppendSQL(*strings.Builder, *[]interface{})
	}:
		v.AppendSQL(buf, args)
		return
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(value)
		if l := s.Len(); l == 0 {
			buf.WriteString("NULL")
		} else {
			buf.WriteString("?")
			buf.WriteString(strings.Repeat(", ?", l-1))
			for i := 0; i < l; i++ {
				*args = append(*args, s.Index(i).Interface())
			}
		}
		return
	}
	buf.WriteString("?")
	*args = append(*args, value)
}

// RandomString is the RandStringBytesMaskImprSrcSB function taken from
// https://stackoverflow.com/a/31832326. It generates a random alphabetical
// string of length n.
func RandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

// InterpolateSQLValue interpolates an interface value as its SQL
// representation into a buffer. This makes it vulnerable to SQL injection and
// should be used for display purposes ONLY, not for actually running against a
// database.
func InterpolateSQLValue(buf *strings.Builder, value interface{}) {
	switch v := value.(type) {
	case nil:
		buf.WriteString("NULL")
	case bool:
		if v {
			buf.WriteString("TRUE")
		} else {
			buf.WriteString("FALSE")
		}
	case string:
		buf.WriteString("'")
		buf.WriteString(v)
		buf.WriteString("'")
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		buf.WriteString(fmt.Sprint(value))
	case time.Time:
		buf.WriteString("'")
		buf.WriteString(v.Format(time.RFC3339Nano))
		buf.WriteString("'")
	case driver.Valuer:
		Interface, err := v.Value()
		if err != nil {
			buf.WriteString("(")
			buf.WriteString(err.Error())
			buf.WriteString(")")
		} else {
			switch Concrete := Interface.(type) {
			case string:
				buf.WriteString("'")
				buf.WriteString(Concrete)
				buf.WriteString("'")
			case nil:
				buf.WriteString("NULL")
			default:
				buf.WriteString("(")
				buf.WriteString(fmt.Sprintf("%#v", value)) // give up, don't know what it is, resort to fmt.Sprintf
				buf.WriteString(")")
			}
		}
	default:
		b, err := json.Marshal(value)
		if err != nil {
			buf.WriteString("(")
			buf.WriteString(fmt.Sprintf("%#v", value)) // give up, don't know what it is, resort to fmt.Sprintf
			buf.WriteString(")")
		} else {
			buf.WriteString("'")
			buf.Write(b)
			buf.WriteString("'")
		}
	}
}

// AppendSQLDisplay marshals an interface value into a buffer.
func AppendSQLDisplay(buf *strings.Builder, value interface{}) {
	switch v := value.(type) {
	case nil:
		buf.WriteString("𝗡𝗨𝗟𝗟")
	case *sql.NullBool:
		if v.Valid {
			if v.Bool {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}
		} else {
			buf.WriteString("𝗡𝗨𝗟𝗟")
		}
	case *sql.NullFloat64:
		if v.Valid {
			buf.WriteString(fmt.Sprintf("%f", v.Float64))
		} else {
			buf.WriteString("𝗡𝗨𝗟𝗟")
		}
	case *sql.NullInt64:
		if v.Valid {
			buf.WriteString(strconv.FormatInt(v.Int64, 10))
		} else {
			buf.WriteString("𝗡𝗨𝗟𝗟")
		}
	case *sql.NullString:
		if v.Valid {
			buf.WriteString(v.String)
		} else {
			buf.WriteString("𝗡𝗨𝗟𝗟")
		}
	case *sql.NullTime:
		if v.Valid {
			buf.WriteString(v.Time.String())
		} else {
			buf.WriteString("𝗡𝗨𝗟𝗟")
		}
	default:
		buf.WriteString(fmt.Sprintf("%#v", value))
	}
}

// QuestionInterpolate interpolates the question mark ? placeholders in a query
// string with the args in the args slice. It is vulnerable to SQL injection
// and should be used for display purposes only, not for actually running
// against a database.
func QuestionInterpolate(query string, args ...interface{}) string {
	buf := &strings.Builder{}
	// i is the position of the ? in the query
	for i := strings.Index(query, "?"); i >= 0 && len(args) > 0; i = strings.Index(query, "?") {
		buf.WriteString(query[:i])
		if len(query[i:]) > 1 && query[i:i+2] == "??" {
			buf.WriteString("?")
			query = query[i+2:]
			continue
		}
		InterpolateSQLValue(buf, args[0])
		query = query[i+1:]
		args = args[1:]
	}
	buf.WriteString(query)
	return buf.String()
}
