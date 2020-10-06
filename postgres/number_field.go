package sq

import "strings"

// NumberField either represents a number column, a number expression or a
// literal number value.
type NumberField struct {
	// NumberField will be one of the following:

	// 1) Number expression
	// Examples of number expressions:
	// | query                  | args        |
	// |------------------------|-------------|
	// | ? / ?                  | 22, 7       |
	// | FLOOR(? + tbl.column)  | 5           |
	// | (ABS(?) + (? % ?)) - ? | -3, 5, 4, 8 |
	format *string
	values []interface{}

	// 2) Literal number value
	// Examples of literal number values:
	// | query | args    |
	// |-------|---------|
	// | ?     | 5       |
	// | ?     | 3.14159 |
	value interface{}

	// 3) Number column
	// Examples of number columns:
	// | query                    | args   |
	// |--------------------------|--------|
	// | users.uid                |        |
	// | uid                      |        |
	// | users.uid ASC NULLS LAST |        |
	alias      string
	table      Table
	name       string
	descending *bool
	nullsfirst *bool
}

// AppendSQLExclude marshals the NumberField into an SQL query and args as
// described in the NumberField internal struct comments.
func (f NumberField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	switch {
	case f.format != nil:
		// 1) Number expression
		expandValues(buf, args, excludedTableQualifiers, *f.format, f.values)
	case f.value != nil:
		// 2) Literal number value
		buf.WriteString("?")
		*args = append(*args, f.value)
	default:
		// 3) Number column
		tableQualifier := f.table.GetAlias()
		if tableQualifier == "" {
			tableQualifier = f.table.GetName()
		}
		for _, excludedTableQualifier := range excludedTableQualifiers {
			if tableQualifier == excludedTableQualifier {
				tableQualifier = ""
				break
			}
		}
		if tableQualifier != "" {
			if strings.ContainsAny(tableQualifier, " \t") {
				buf.WriteString(`"`)
				buf.WriteString(tableQualifier)
				buf.WriteString(`".`)
			} else {
				buf.WriteString(tableQualifier)
				buf.WriteString(".")
			}
		}
		if strings.ContainsAny(f.name, " \t") {
			buf.WriteString(`"`)
			buf.WriteString(f.name)
			buf.WriteString(`"`)
		} else {
			buf.WriteString(f.name)
		}
	}
	if f.descending != nil {
		if *f.descending {
			buf.WriteString(" DESC")
		} else {
			buf.WriteString(" ASC")
		}
	}
	if f.nullsfirst != nil {
		if *f.nullsfirst {
			buf.WriteString(" NULLS FIRST")
		} else {
			buf.WriteString(" NULLS LAST")
		}
	}
}

// NewNumberField returns a new NumberField representing a number TableInfo column.
func NewNumberField(name string, table Table) NumberField {
	return NumberField{
		name:  name,
		table: table,
	}
}

// Int returns a new NumberField representing a literal int value.
func Int(num int) NumberField {
	return NumberField{
		value: num,
	}
}

// Int64 returns a new NumberField representing a literal int64 value.
func Int64(num int64) NumberField {
	return NumberField{
		value: num,
	}
}

// Float64 returns a new NumberField representing a literal float64 value.
func Float64(num float64) NumberField {
	return NumberField{
		value: num,
	}
}

// Set returns a FieldAssignment associating the NumberField to the value i.e.
// 'field = value'.
func (f NumberField) Set(val interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: val,
	}
}

// SetInt returns a FieldAssignment associating the NumberField to the int value
// i.e. 'field = value'.
func (f NumberField) SetInt(num int) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: num,
	}
}

// SetInt64 returns a FieldAssignment associating the NumberField to the int64
// value i.e. 'field = value'.
func (f NumberField) SetInt64(num int64) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: num,
	}
}

// SetFloat64 returns a FieldAssignment associating the NumberField to the float64
// value i.e. 'field = value'.
func (f NumberField) SetFloat64(num float64) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: num,
	}
}

// As returns a new NumberField with the new field Alias i.e. 'field AS Alias'.
func (f NumberField) As(alias string) NumberField {
	f.alias = alias
	return f
}

// Asc returns a new NumberField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f NumberField) Asc() NumberField {
	desc := false
	f.descending = &desc
	return f
}

// Desc returns a new NumberField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f NumberField) Desc() NumberField {
	desc := true
	f.descending = &desc
	return f
}

// NullsFirst returns a new NumberField indicating that it should be ordered
// with nulls first i.e. 'ORDER BY field NULLS FIRST'.
func (f NumberField) NullsFirst() NumberField {
	nullsfirst := true
	f.nullsfirst = &nullsfirst
	return f
}

// NullsLast returns a new NumberField indicating that it should be ordered
// with nulls last i.e. 'ORDER BY field NULLS LAST'.
func (f NumberField) NullsLast() NumberField {
	nullsfirst := false
	f.nullsfirst = &nullsfirst
	return f
}

// IsNull returns an 'X IS NULL' Predicate.
func (f NumberField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f NumberField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// Eq returns an 'X = Y' Predicate. It only accepts NumberField.
func (f NumberField) Eq(field NumberField) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, field},
	}
}

// EqFloat64 returns an 'X = Y' Predicate. It only accepts float64.
func (f NumberField) EqFloat64(num float64) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, num},
	}
}

// EqInt returns an 'X = Y' Predicate. It only accepts int.
func (f NumberField) EqInt(num int) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, num},
	}
}

// Ne returns an 'X <> Y' Predicate. It only accepts NumberField.
func (f NumberField) Ne(field NumberField) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, field},
	}
}

// NeFloat64 returns an 'X <> Y' Predicate. It only accepts float64.
func (f NumberField) NeFloat64(num float64) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, num},
	}
}

// NeInt returns an 'X <> Y' Predicate. It only accepts int.
func (f NumberField) NeInt(num int) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, num},
	}
}

// Gt returns an 'X > Y' Predicate. It only accepts NumberField.
func (f NumberField) Gt(field NumberField) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, field},
	}
}

// GtFloat64 returns an 'X > Y' Predicate. It only accepts float64.
func (f NumberField) GtFloat64(num float64) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, num},
	}
}

// GtInt returns an 'X > Y' Predicate. It only accepts int.
func (f NumberField) GtInt(num int) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, num},
	}
}

// Ge returns an 'X >= Y' Predicate. It only accepts NumberField.
func (f NumberField) Ge(field NumberField) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, field},
	}
}

// GeFloat64 returns an 'X >= Y' Predicate. It only accepts float64.
func (f NumberField) GeFloat64(num float64) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, num},
	}
}

// GeInt returns an 'X >= Y' Predicate. It only accepts int.
func (f NumberField) GeInt(num int) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, num},
	}
}

// Lt returns an 'X < Y' Predicate. It only accepts NumberField.
func (f NumberField) Lt(field NumberField) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, field},
	}
}

// LtFloat64 returns an 'X < Y' Predicate. It only accepts float64.
func (f NumberField) LtFloat64(num float64) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, num},
	}
}

// LtInt returns an 'X < Y' Predicate. It only accepts int.
func (f NumberField) LtInt(num int) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, num},
	}
}

// Le returns an 'X <= Y' Predicate. It only accepts NumberField.
func (f NumberField) Le(field NumberField) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, field},
	}
}

// LeFloat64 returns an 'X <= Y' Predicate. It only accepts float64.
func (f NumberField) LeFloat64(num float64) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, num},
	}
}

// LeInt returns an 'X <= Y' Predicate. It only accepts int.
func (f NumberField) LeInt(num int) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, num},
	}
}

// In returns an 'X IN (Y)' Predicate.
func (f NumberField) In(v interface{}) Predicate {
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
// representation of a NumberField.
func (f NumberField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil, nil)
	return questionInterpolate(buf.String(), args...)
}

// GetAlias implements the Field interface. It returns the Alias of the
// NumberField.
func (f NumberField) GetAlias() string {
	return f.alias
}

// GetName implements the Field interface. It returns the Name of the
// NumberField.
func (f NumberField) GetName() string {
	return f.name
}

// NumberFieldf creates a new number expression.
func NumberFieldf(format string, values ...interface{}) NumberField {
	return NumberField{
		format: &format,
		values: values,
	}
}
