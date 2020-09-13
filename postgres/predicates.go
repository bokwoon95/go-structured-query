package sq

import "strings"

func Not(predicate Predicate) Predicate {
	return predicate.Not()
}

// CustomPredicate is a Query that can render itself in an arbitrary way as defined
// by its Format string. Values are interpolated into the Format string as
// described in the (CustomPredicate).CustomSprintf function.
type CustomPredicate struct {
	Alias    string
	Format   string
	Values   []interface{}
	Negative bool
}

// ToSQL marshals the CustomPredicate into an SQL query.
func (p CustomPredicate) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	if p.Negative {
		buf.WriteString("NOT ")
	}
	ExpandValues(buf, args, excludedTableQualifiers, p.Format, p.Values)
}

func Predicatef(format string, values ...interface{}) CustomPredicate {
	return CustomPredicate{
		Format: format,
		Values: values,
	}
}

func (p CustomPredicate) As(alias string) CustomPredicate {
	p.Alias = alias
	return p
}

// Not implements the Predicate interface.
func (p CustomPredicate) Not() Predicate {
	p.Negative = !p.Negative
	return p
}

// GetAlias implements the Field interface.
func (p CustomPredicate) GetAlias() string {
	return p.Alias
}

// GetName implements the Field interface.
func (p CustomPredicate) GetName() string {
	return ""
}

func Exists(query Query) CustomPredicate {
	return CustomPredicate{
		Format: "EXISTS(?)",
		Values: []interface{}{query},
	}
}

// VariadicPredicateOperator is an operator that can join a variadic number of
// Predicates together.
type VariadicPredicateOperator string

// Possible VariadicOperators
const (
	PredicateOr  VariadicPredicateOperator = "OR"
	PredicateAnd VariadicPredicateOperator = "AND"
)

// VariadicPredicate represents the "x AND y AND z..." or "x OR y OR z..." SQL
// construct.
type VariadicPredicate struct {
	// Toplevel indicates if the variadic predicate is the top level predicate
	// i.e. it does not need enclosing brackets
	Toplevel   bool
	Alias      string
	Operator   VariadicPredicateOperator
	Predicates []Predicate
	Negative   bool
}

// ToSQL marshals the VariadicPredicate into an SQL query and args as described
// in the VariadicPredicate struct description.
func (p VariadicPredicate) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	if p.Operator == "" {
		p.Operator = PredicateAnd
	}
	switch len(p.Predicates) {
	case 0: // nothing to do here
	case 1:
		if p.Negative {
			buf.WriteString("NOT ")
		}
		switch v := p.Predicates[0].(type) {
		case nil:
			buf.WriteString("NULL")
		case VariadicPredicate:
			if !p.Toplevel {
				buf.WriteString("(")
			}
			v.Toplevel = true
			v.AppendSQLExclude(buf, args, excludedTableQualifiers)
			if !p.Toplevel {
				buf.WriteString(")")
			}
		default:
			p.Predicates[0].AppendSQLExclude(buf, args, excludedTableQualifiers)
		}
	default:
		if p.Negative {
			buf.WriteString("NOT ")
		}
		if !p.Toplevel {
			buf.WriteString("(")
		}
		for i, predicate := range p.Predicates {
			if i > 0 {
				buf.WriteString(" ")
				buf.WriteString(string(p.Operator))
				buf.WriteString(" ")
			}
			if predicate == nil {
				buf.WriteString("NULL")
			} else {
				predicate.AppendSQLExclude(buf, args, excludedTableQualifiers)
			}
		}
		if !p.Toplevel {
			buf.WriteString(")")
		}
	}
}

// Not implements the Predicate interface.
func (p VariadicPredicate) Not() Predicate {
	p.Negative = !p.Negative
	return p
}

// GetAlias implements the Field interface.
func (p VariadicPredicate) GetAlias() string {
	return p.Alias
}

// GetName implements the Field interface.
func (p VariadicPredicate) GetName() string {
	return ""
}

func And(predicates ...Predicate) VariadicPredicate {
	return VariadicPredicate{
		Operator:   PredicateAnd,
		Predicates: predicates,
	}
}

func Or(predicates ...Predicate) VariadicPredicate {
	return VariadicPredicate{
		Operator:   PredicateOr,
		Predicates: predicates,
	}
}

func Eq(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f1, f2},
	}
}

func Ne(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f1, f2},
	}
}

func Gt(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f1, f2},
	}
}

func Ge(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f1, f2},
	}
}

func Lt(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f1, f2},
	}
}

func Le(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f1, f2},
	}
}
