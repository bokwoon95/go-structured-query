package sq

import "strings"

// JoinType represents the various types of SQL joins.
type JoinType string

// JoinTypes
const (
	JoinTypeInner JoinType = "JOIN"
	JoinTypeLeft  JoinType = "LEFT JOIN"
	JoinTypeRight JoinType = "RIGHT JOIN"
	JoinTypeFull  JoinType = "FULL JOIN"
)

// JoinTable represents an SQL join.
type JoinTable struct {
	JoinType     JoinType
	Table        Table
	OnPredicates VariadicPredicate
}

func Join(table Table, predicates ...Predicate) JoinTable {
	return JoinTable{
		JoinType: JoinTypeInner,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	}
}

func LeftJoin(table Table, predicates ...Predicate) JoinTable {
	return JoinTable{
		JoinType: JoinTypeLeft,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	}
}

func RightJoin(table Table, predicates ...Predicate) JoinTable {
	return JoinTable{
		JoinType: JoinTypeRight,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	}
}

func FullJoin(table Table, predicates ...Predicate) JoinTable {
	return JoinTable{
		JoinType: JoinTypeFull,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	}
}

// CustomJoin constructs a new JoinTable. Meant to be used if you want to do a custom
// join like CROSS JOIN, NATURAL JOIN, LEFT JOIN LATERAL etc.
func CustomJoin(joinType JoinType, table Table, predicates ...Predicate) JoinTable {
	return JoinTable{
		JoinType: joinType,
		Table:    table,
		OnPredicates: VariadicPredicate{
			Predicates: predicates,
		},
	}
}

func (join JoinTable) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	if join.JoinType == "" {
		join.JoinType = JoinTypeInner
	}
	buf.WriteString(string(join.JoinType) + " ")
	switch v := join.Table.(type) {
	case nil:
		buf.WriteString("NULL")
	case Query:
		buf.WriteString("(")
		v.NestThis().AppendSQL(buf, args)
		buf.WriteString(")")
	default:
		join.Table.AppendSQL(buf, args)
	}
	if join.Table != nil {
		alias := join.Table.GetAlias()
		if alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(alias)
		}
	}
	if len(join.OnPredicates.Predicates) > 0 {
		buf.WriteString(" ON ")
		join.OnPredicates.Toplevel = true
		join.OnPredicates.AppendSQLExclude(buf, args, nil)
	}
}

// JoinTables is a list of JoinTables.
type JoinTables []JoinTable

// AppendSQL will write the JOIN clause into the buffer and args. If there are
// no JoinTables it simply writes nothing into the buffer. It returns a flag
// indicating whether anything was written into the buffer.
func (joins JoinTables) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	for i, join := range joins {
		if i > 0 {
			buf.WriteString(" ")
		}
		join.AppendSQL(buf, args)
	}
}
