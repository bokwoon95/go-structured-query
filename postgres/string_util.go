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

// expandValues will expand each value one by one into successive question mark
// ? placeholders in the format string, writing the results into the buffer and
// args slice. It propagates the excludedTableQualifiers down to its child elements.
func expandValues(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string, format string, values []interface{}) {
	for i := strings.Index(format, "?"); i >= 0 && len(values) > 0; i = strings.Index(format, "?") {
		buf.WriteString(format[:i])
		// TODO: I don't know if ?? should be unescaped to ?
		// if len(format[i:]) > 1 && format[i:i+2] == "??" {
		// 	buf.WriteString("?")
		// 	format = format[i+2:]
		// 	continue
		// }
		appendSQLValue(buf, args, excludedTableQualifiers, values[0])
		format = format[i+1:]
		values = values[1:]
	}
	buf.WriteString(format)
}

// appendSQLValue will write the SQL representation of the interface{} value
// into the buffer and args slice. It propagates excludedTableQualifiers where
// relevant.
func appendSQLValue(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string, value interface{}) {
	switch v := value.(type) {
	case nil:
		buf.WriteString("NULL")
		return
	case interface {
		AppendSQLExclude(*strings.Builder, *[]interface{}, map[string]int, []string)
	}:
		v.AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
		return
	case interface {
		AppendSQL(*strings.Builder, *[]interface{}, map[string]int)
	}:
		v.AppendSQL(buf, args, nil)
		return
	}
	switch typ := reflect.TypeOf(value); typ.Kind() {
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			break
		}
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

// randomString is the RandStringBytesMaskImplSrcSB function taken from
// https://stackoverflow.com/a/31832326. It generates a random alphabetical
// string of length n.
func randomString(n int) string {
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

// TODO: write a version that takes in a buffer and writes into it instead
func interpolateSQLValue(arg interface{}) string {
	var str string // str is the SQL string representation of arg
	switch v := arg.(type) {
	case nil:
		str = "NULL"
	case bool:
		if v {
			str = "TRUE"
		} else {
			str = "FALSE"
		}
	case string:
		str = "'" + v + "'"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		str = fmt.Sprint(arg)
	case time.Time:
		str = "'" + v.Format(time.RFC3339Nano) + "'"
	case driver.Valuer:
		Interface, err := v.Value()
		if err != nil {
			str = ":" + err.Error() + ":"
		} else {
			switch Concrete := Interface.(type) {
			case string:
				str = "'" + Concrete + "'"
			case nil:
				str = "NULL"
			default:
				str = ":" + fmt.Sprintf("%#v", arg) + ":" // give up, don't know what it is, resort to fmt.Sprintf
			}
		}
	default:
		b, err := json.Marshal(arg)
		if err != nil {
			str = ":" + fmt.Sprintf("%#v", arg) + ":" // give up, don't know what it is, resort to fmt.Sprintf
		} else {
			str = "'" + string(b) + "'"
		}
	}
	return str
}

// AppendSQLDisplay
func appendSQLDisplay(arg interface{}) string {
	var str string // str is the SQL string representation of arg
	switch v := arg.(type) {
	case nil:
		str = "𝗡𝗨𝗟𝗟"
	case *sql.NullBool:
		if v.Valid {
			if v.Bool {
				str = "true"
			} else {
				str = "false"
			}
		} else {
			str = "𝗡𝗨𝗟𝗟"
		}
	case *sql.NullFloat64:
		if v.Valid {
			str = fmt.Sprintf("%f", v.Float64)
		} else {
			str = "𝗡𝗨𝗟𝗟"
		}
	case *sql.NullInt64:
		if v.Valid {
			str = strconv.FormatInt(v.Int64, 10)
		} else {
			str = "𝗡𝗨𝗟𝗟"
		}
	case *sql.NullString:
		if v.Valid {
			str = v.String
		} else {
			str = "𝗡𝗨𝗟𝗟"
		}
	case *sql.NullTime:
		if v.Valid {
			str = v.Time.String()
		} else {
			str = "𝗡𝗨𝗟𝗟"
		}
	default:
		str = fmt.Sprintf("%#v", arg)
	}
	return str
}

// questionToDollarPlaceholders will replace all MySQL style ? with Postgres
// style incrementing placeholders i.e. $1, $2, $3 etc. To escape a literal
// question mark ? , use two question marks ?? instead.
func questionToDollarPlaceholders(buf *strings.Builder, query string) {
	i := 0
	for {
		p := strings.Index(query, "?")
		if p < 0 {
			break
		}
		buf.WriteString(query[:p])
		if len(query[p:]) > 1 && query[p:p+2] == "??" {
			buf.WriteString("?")
			query = query[p+2:]
		} else {
			i++
			buf.WriteString("$" + strconv.Itoa(i))
			query = query[p+1:]
		}
	}
	buf.WriteString(query)
}

// questionInterpolate interpolates the question mark ? placeholders in a query
// string with the args in the args slice. It is vulnerable to SQL injection
// and should be used for display purposes only, not for actually running
// against a database.
func questionInterpolate(query string, args ...interface{}) string {
	buf := &strings.Builder{}
	// i is the position of the ? in the query
	for i := strings.Index(query, "?"); i >= 0 && len(args) > 0; i = strings.Index(query, "?") {
		buf.WriteString(query[:i])
		if len(query[i:]) > 1 && query[i:i+2] == "??" {
			buf.WriteString("?")
			query = query[i+2:]
			continue
		}
		buf.WriteString(interpolateSQLValue(args[0]))
		query = query[i+1:]
		args = args[1:]
	}
	buf.WriteString(query)
	return buf.String()
}

// dollarInterpolate interpolates the dollar $1 ($2, $3 etc) placeholders in a
// query string with the args in the args slice. It is vulnerable to SQL
// injection and should be used for display purposes only, not for actually
// running against a database.
func dollarInterpolate(query string, args ...interface{}) string {
	oldnewSets := make(map[int][]string)
	for i, arg := range args {
		str := interpolateSQLValue(arg)
		placeholder := "$" + strconv.Itoa(i+1)
		oldnewSets[len(placeholder)] = append(oldnewSets[len(placeholder)], placeholder, str)
	}
	result := query
	for i := len(oldnewSets) + 1; i >= 2; i-- {
		result = strings.NewReplacer(oldnewSets[i]...).Replace(result)
	}
	return result
}
