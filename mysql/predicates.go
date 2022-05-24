package sq

import "strings"

// Not inverts the Predicate i.e. 'NOT Predicate'.
func Not(predicate Predicate) Predicate {
	return predicate.Not()
}

// CustomPredicate is a Query that can render itself in an arbitrary way by
// calling expandValues on its Format and Values.
type CustomPredicate struct {
	Alias    string
	Format   string
	Values   []interface{}
	Negative bool
}

// AppendSQLExclude marshals the CustomPredicate into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (p CustomPredicate) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	if p.Negative {
		buf.WriteString("NOT ")
	}
	expandValues(buf, args, excludedTableQualifiers, p.Format, p.Values)
}

// Predicatef creates a new CustomPredicate.
func Predicatef(format string, values ...interface{}) CustomPredicate {
	return CustomPredicate{
		Format: format,
		Values: values,
	}
}

// As aliases the CustomPredicate.
func (p CustomPredicate) As(alias string) CustomPredicate {
	p.Alias = alias
	return p
}

// Not inverts the CustomPredicate i.e. 'NOT CustomPredicate'.
func (p CustomPredicate) Not() Predicate {
	p.Negative = !p.Negative
	return p
}

// GetAlias returns the alias of the CustomPredicate.
func (p CustomPredicate) GetAlias() string {
	return p.Alias
}

// GetName returns the name of the CustomPredicate, which is always an empty
// string.
func (p CustomPredicate) GetName() string {
	return ""
}

// Exists represents the EXISTS() predicate.
func Exists(query Query) CustomPredicate {
	return CustomPredicate{
		Format: "EXISTS(?)",
		Values: []interface{}{query},
	}
}

// VariadicPredicateOperator is an operator that can join a variadic number of
// Predicates together.
type VariadicPredicateOperator string

// VariadicPredicateOperators
const (
	PredicateOr  VariadicPredicateOperator = "OR"
	PredicateAnd VariadicPredicateOperator = "AND"
)

// VariadicPredicate represents the "x AND y AND z..." or "x OR y OR z..." SQL
// construct.
type VariadicPredicate struct {
	// toplevel indicates if the variadic predicate is the top level predicate
	// i.e. it does not need enclosing brackets
	toplevel   bool
	Alias      string
	Operator   VariadicPredicateOperator
	Predicates []Predicate
	Negative   bool
}

// AppendSQLExclude marshals the VariadicPredicate into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (p VariadicPredicate) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	if p.Operator == "" {
		p.Operator = PredicateAnd
	}
	switch len(p.Predicates) {
	case 0: // no-op
	case 1:
		if p.Negative {
			buf.WriteString("NOT ")
		}
		switch v := p.Predicates[0].(type) {
		case nil:
			buf.WriteString("NULL")
		case VariadicPredicate:
			if !p.toplevel {
				buf.WriteString("(")
			}
			v.toplevel = true
			v.AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
			if !p.toplevel {
				buf.WriteString(")")
			}
		default:
			p.Predicates[0].AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
		}
	default:
		if p.Negative {
			buf.WriteString("NOT ")
		}
		if !p.toplevel {
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
				predicate.AppendSQLExclude(buf, args, nil, excludedTableQualifiers)
			}
		}
		if !p.toplevel {
			buf.WriteString(")")
		}
	}
}

// Not inverts the VariadicPredicate i.e. 'NOT VariadicPredicate'.
func (p VariadicPredicate) Not() Predicate {
	p.Negative = !p.Negative
	return p
}

// GetAlias returns the alias of the VariadicPredicate.
func (p VariadicPredicate) GetAlias() string {
	return p.Alias
}

// GetName returns the name of the VariadicPredicate, which is always an empty
// string.
func (p VariadicPredicate) GetName() string {
	return ""
}

// And joins the list of predicates together with the AND operator.
func And(predicates ...Predicate) VariadicPredicate {
	return VariadicPredicate{
		Operator:   PredicateAnd,
		Predicates: predicates,
	}
}

// Or joins the list of predicates together with the OR operator.
func Or(predicates ...Predicate) VariadicPredicate {
	return VariadicPredicate{
		Operator:   PredicateOr,
		Predicates: predicates,
	}
}

// Eq returns an 'X = Y' Predicate.
func Eq(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f1, f2},
	}
}

// In returns an 'X IN (Y)' Predicate.
func In(f1 interface{}, f2 []interface{}) Predicate {
	return CustomPredicate{
		Format: "? IN (?)",
		Values: []interface{}{f1, f2},
	}
}

// Ne returns an 'X <> Y' Predicate.
func Ne(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f1, f2},
	}
}

// Gt returns an 'X > Y' Predicate.
func Gt(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f1, f2},
	}
}

// Ge returns an 'X >= Y' Predicate.
func Ge(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f1, f2},
	}
}

// Lt returns an 'X < Y' Predicate.
func Lt(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f1, f2},
	}
}

// Le returns an 'X <= Y' Predicate.
func Le(f1, f2 interface{}) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f1, f2},
	}
}
