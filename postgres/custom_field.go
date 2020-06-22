package sq

import "strings"

// CustomField is a Field that can render itself in an arbitrary way by calling
// ExpandValues on its Format and Values.
type CustomField struct {
	Alias        string
	Format       string
	Values       []interface{}
	IsDesc       *bool
	IsNullsFirst *bool
}

// AppendSQL marshals the CustomField into an SQL query and args as described in
// the CustomField struct description.
func (f CustomField) AppendSQLExclude(buf Buffer, args *[]interface{}, excludedTableQualifiers []string) {
	ExpandValues(buf, args, excludedTableQualifiers, f.Format, f.Values)
	if f.IsDesc != nil {
		if *f.IsDesc {
			buf.WriteString(" DESC")
		} else {
			buf.WriteString(" ASC")
		}
	}
	if f.IsNullsFirst != nil {
		if *f.IsNullsFirst {
			buf.WriteString(" NULLS FIRST")
		} else {
			buf.WriteString(" NULLS LAST")
		}
	}
}

// Fieldf is a CustomField constructor.
func Fieldf(format string, values ...interface{}) CustomField {
	return CustomField{
		Format: format,
		Values: values,
	}
}

// As returns a new CustomField with the new alias i.e. 'field AS Alias'.
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

// NullsFirst returns a new CustomField indicating that it should be ordered
// with nulls first i.e. 'ORDER BY field NULLS FIRST'.
func (f CustomField) NullsFirst() CustomField {
	isNullsFirst := true
	f.IsNullsFirst = &isNullsFirst
	return f
}

// NullsLast returns a new CustomField indicating that it should be ordered
// with nulls last i.e. 'ORDER BY field NULLS LAST'.
func (f CustomField) NullsLast() CustomField {
	isNullsFirst := false
	f.IsNullsFirst = &isNullsFirst
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
	var format string
	var values []interface{}
	switch v := v.(type) {
	case RowValue:
		format = "? IN ?"
		values = []interface{}{f, v}
	case Query:
		format = "? IN (?)"
		values = []interface{}{f, v.NestThis()}
	default:
		format = "? IN (?)"
		values = []interface{}{f, v}
	}
	return CustomPredicate{
		Format: format,
		Values: values,
	}
}

// String implements the fmt.Stringer interface. It returns the string
// representation of a CustomField.
func (f CustomField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil)
	return QuestionInterpolate(buf.String(), args...)
}

// GetAlias implements the Field interface. It returns the alias of thee
// CustomField.
func (f CustomField) GetAlias() string {
	return f.Alias
}

// GetName implements the Field interface. It returns the name of the
// CustomField.
func (f CustomField) GetName() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil)
	return buf.String()
}
