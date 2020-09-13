package sq

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestBaseQuery(t *testing.T) {
	is := is.New(t)

	// (BaseQuery).With will append CTEs, not overwrite it
	q := With(CTE{}).With(CTE{}, CTE{}).With(CTE{})
	is.Equal(4, len(q.CTEs))

	var base BaseQuery
	var buf = &strings.Builder{}
	var args []interface{}
	var sel SelectQuery
	var ins InsertQuery
	var upd UpdateQuery
	var del DeleteQuery

	// WithDefaultLog
	base = WithDefaultLog(Lstats).WithDefaultLog(Lstats)
	is.Equal(defaultLogger, base.Log)
	is.Equal(Lstats, base.LogFlag)

	// With
	base = With(CTE{}, CTE{}, CTE{})
	is.Equal(3, len(base.CTEs))

	// WithLog
	l := log.New(os.Stdout, "", 0)
	base = WithLog(l, Lverbose).WithLog(l, Lverbose)
	is.Equal(l, base.Log)
	is.Equal(Lverbose, base.LogFlag)

	// SelectOne
	sel = BaseQuery{}.SelectOne()
	buf.Reset()
	sel.AppendSQL(buf, &args)
	is.Equal("SELECT 1", buf.String())

	// SelectAll
	sel = BaseQuery{}.SelectAll()
	buf.Reset()
	sel.AppendSQL(buf, &args)
	is.Equal("SELECT *", buf.String())

	// SelectCount
	sel = BaseQuery{}.SelectCount()
	buf.Reset()
	sel.AppendSQL(buf, &args)
	is.Equal("SELECT COUNT(*)", buf.String())

	// SelectDistinct
	sel = BaseQuery{}.SelectDistinct()
	buf.Reset()
	sel.AppendSQL(buf, &args)
	is.Equal("SELECT DISTINCT", buf.String())

	// Selectx
	mapper := func(_ *Row) {}
	accumulator := func() {}
	sel = BaseQuery{}.Selectx(mapper, accumulator)
	buf.Reset()
	sel.AppendSQL(buf, &args)
	is.Equal(mapper, sel.Mapper)
	is.Equal(accumulator, sel.Accumulator)

	// SelectRowx
	sel = BaseQuery{}.SelectRowx(mapper)
	buf.Reset()
	sel.AppendSQL(buf, &args)
	is.Equal(mapper, sel.Mapper)
	is.Equal(nil, sel.Accumulator)

	// InsertInto
	ins = BaseQuery{}.InsertInto(nil)
	buf.Reset()
	ins.AppendSQL(buf, &args)
	is.Equal("INSERT INTO NULL", buf.String())

	// Update
	upd = BaseQuery{}.Update(nil)
	buf.Reset()
	upd.AppendSQL(buf, &args)
	is.Equal("UPDATE NULL", buf.String())

	// DeleteFrom
	del = BaseQuery{}.DeleteFrom(nil)
	buf.Reset()
	del.AppendSQL(buf, &args)
	is.Equal("DELETE FROM NULL", buf.String())
}
