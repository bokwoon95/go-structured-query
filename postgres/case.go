package sq

// PredicateCase represents a Predicate and the Result if the Predicate is
// true.
type PredicateCase struct {
	Condition Predicate
	Result    interface{}
}

// PredicateCases is the general form of the CASE expression.
type PredicateCases struct {
	Alias    string
	Cases    []PredicateCase
	Fallback interface{}
}

// AppendSQLExclude marshals the PredicateCases into a buffer and an args
// slice. It propagates the excludedTableQualifiers down to its child elements.
func (f PredicateCases) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	buf.WriteString("CASE")
	for i := range f.Cases {
		buf.WriteString(" WHEN ")
		AppendSQLValue(buf, args, excludedTableQualifiers, f.Cases[i].Condition)
		buf.WriteString(" THEN ")
		AppendSQLValue(buf, args, excludedTableQualifiers, f.Cases[i].Result)
	}
	if f.Fallback != nil {
		buf.WriteString(" ELSE ")
		AppendSQLValue(buf, args, excludedTableQualifiers, f.Fallback)
	}
	buf.WriteString(" END")
}

// CaseWhen creates a new PredicateCases i.e. CASE WHEN X THEN Y.
func CaseWhen(predicate Predicate, result interface{}) PredicateCases {
	return PredicateCases{
		Cases: []PredicateCase{{
			Condition: predicate,
			Result:    result,
		}},
	}
}

// When adds a new PredicateCase to the PredicateCases i.e. WHEN X THEN Y.
func (f PredicateCases) When(predicate Predicate, result interface{}) PredicateCases {
	f.Cases = append(f.Cases, PredicateCase{
		Condition: predicate,
		Result:    result,
	})
	return f
}

// Else adds the fallback value for the PredicateCases i.e. ELSE X.
func (f PredicateCases) Else(fallback interface{}) PredicateCases {
	f.Fallback = fallback
	return f
}

// As aliases the PredicateCases.
func (f PredicateCases) As(alias string) PredicateCases {
	f.Alias = alias
	return f
}

// GetAlias returns the alias of the PredicateCases.
func (f PredicateCases) GetAlias() string {
	return f.Alias
}

// GetName returns the name of the PredicateCases, which is always an empty
// string.
func (f PredicateCases) GetName() string {
	return ""
}

// SimpleCase represents a Value to be compared against and the Result if it
// matches.
type SimpleCase struct {
	Value  interface{}
	Result interface{}
}

// SimpleCases is the simple form of the CASE expression.
type SimpleCases struct {
	Alias      string
	Expression interface{}
	Cases      []SimpleCase
	Fallback   interface{}
}

// AppendSQLExclude marshals the SimpleCases into a buffer and an args slice.
// It propagates the excludedTableQualifiers down to its child elements.
func (f SimpleCases) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	buf.WriteString("CASE ")
	AppendSQLValue(buf, args, excludedTableQualifiers, f.Expression)
	for i := range f.Cases {
		buf.WriteString(" WHEN ")
		AppendSQLValue(buf, args, excludedTableQualifiers, f.Cases[i].Value)
		buf.WriteString(" THEN ")
		AppendSQLValue(buf, args, excludedTableQualifiers, f.Cases[i].Result)
	}
	if f.Fallback != nil {
		buf.WriteString(" ELSE ")
		AppendSQLValue(buf, args, excludedTableQualifiers, f.Fallback)
	}
	buf.WriteString(" END")
}

// Case creates a new SimpleCases i.e. CASE X
func Case(field Field) SimpleCases {
	return SimpleCases{
		Expression: field,
	}
}

// When adds a new SimpleCase to the SimpleCases i.e. WHEN X THEN Y.
func (f SimpleCases) When(field Field, result Field) SimpleCases {
	f.Cases = append(f.Cases, SimpleCase{
		Value:  field,
		Result: result,
	})
	return f
}

// Else adds the fallback value for the SimpleCases i.e. ELSE X.
func (f SimpleCases) Else(field Field) SimpleCases {
	f.Fallback = field
	return f
}

// As aliases the SimpleCases.
func (f SimpleCases) As(alias string) SimpleCases {
	f.Alias = alias
	return f
}

// GetAlias returns the alias of the SimpleCases.
func (f SimpleCases) GetAlias() string {
	return f.Alias
}

// GetName returns the name of the simple cases, which is always an empty
// string.
func (f SimpleCases) GetName() string {
	return ""
}
