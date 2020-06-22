package sq

import "strings"

// CustomField is a Field that can render itself in an arbitrary way by calling
// ExpandValues on its Format and Values.
type CustomField struct {
	Alias  string
	Format string
	Values []interface{}
	IsDesc *bool
}

// AppendSQLExclude marshals the CustomField into a buffer and an args slice.
// It propagates the excludedTableQualifiers down to its child elements.
func (f CustomField) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	ExpandValues(buf, args, excludedTableQualifiers, f.Format, f.Values)
	if f.IsDesc != nil {
		if *f.IsDesc {
			buf.WriteString(" DESC")
		} else {
			buf.WriteString(" ASC")
		}
	}
}

// Fieldf creates a new CustomField.
func Fieldf(format string, values ...interface{}) CustomField {
	return CustomField{
		Format: format,
		Values: values,
	}
}

// Aliases the CustomField i.e. 'field AS Alias'.
func (f CustomField) As(alias string) CustomField {
	f.Alias = alias
	return f
}

// Asc returns a new CustomField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f CustomField) Asc() CustomField {
	isDesc := false
	f.IsDesc = &isDesc
	return f
}

// Desc returns a new CustomField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f CustomField) Desc() CustomField {
	isDesc := true
	f.IsDesc = &isDesc
	return f
}

// IsNull returns an 'X IS NULL' Predicate.
func (f CustomField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f CustomField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// Eq returns an 'X = Y' Predicate.
func (f CustomField) Eq(v interface{}) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, v},
	}
}

// Ne returns an 'X <> Y' Predicate.
func (f CustomField) Ne(v interface{}) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, v},
	}
}

// Gt returns an 'X > Y' Predicate.
func (f CustomField) Gt(v interface{}) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, v},
	}
}

// Ge returns an 'X >= Y' Predicate.
func (f CustomField) Ge(v interface{}) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, v},
	}
}

// Lt returns an 'X < Y' Predicate.
func (f CustomField) Lt(v interface{}) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, v},
	}
}

// Le returns an 'X <= Y' Predicate.
func (f CustomField) Le(v interface{}) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, v},
	}
}

// In returns an 'X IN (Y)' Predicate.
func (f CustomField) In(v interface{}) Predicate {
	switch v.(type) {
	case RowValue:
		return CustomPredicate{
			Format: "? IN ?",
			Values: []interface{}{f, v},
		}
	default:
		return CustomPredicate{
			Format: "? IN (?)",
			Values: []interface{}{f, v},
		}
	}
}

// String returns the string representation of the CustomField.
func (f CustomField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil)
	return QuestionInterpolate(buf.String(), args...)
}

// GetAlias returns the alias of the CustomField.
func (f CustomField) GetAlias() string {
	return f.Alias
}

// GetName returns the name of the CustomField.
func (f CustomField) GetName() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil)
	return buf.String()
}
