package sq

import (
	"strings"
	"time"
)

// TimeField either represents a time column or a literal time.Time value.
type TimeField struct {
	// TimeField will be one of the following:

	// 1) Literal time.Time value
	// Examples of literal string values:
	// | query | args       |
	// |-------|------------|
	// | ?     | time.Now() |
	value *time.Time

	// 2) Time column
	// Examples of time columns:
	// | query            | args |
	// |------------------|------|
	// | users.created_at |      |
	// | created_at       |      |
	// | events.start_at  |      |
	alias      string
	table      Table
	name       string
	descending *bool
}

// AppendSQLExclude marshals the TimeField into an SQL query and args as described
// in the TimeField internal struct comments.
func (f TimeField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, excludedTableQualifiers []string) {
	switch {
	case f.value != nil:
		// 1) Literal time.Time value
		buf.WriteString("?")
		*args = append(*args, *f.value)
	default:
		// 2) Time column
		tableQualifier := f.table.GetAlias()
		if tableQualifier == "" {
			tableQualifier = f.table.GetName()
		}
		for i := range excludedTableQualifiers {
			if tableQualifier == excludedTableQualifiers[i] {
				tableQualifier = ""
				break
			}
		}
		if tableQualifier != "" {
			if strings.ContainsAny(tableQualifier, " \t") {
				buf.WriteString("`")
				buf.WriteString(tableQualifier)
				buf.WriteString("`.")
			} else {
				buf.WriteString(tableQualifier)
				buf.WriteString(".")
			}
		}
		if strings.ContainsAny(f.name, " \t") {
			buf.WriteString("`")
			buf.WriteString(f.name)
			buf.WriteString("`")
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
}

// NewTimeField returns a new TimeField representing a time column.
func NewTimeField(name string, table Table) TimeField {
	return TimeField{
		name:  name,
		table: table,
	}
}

// Time returns a new TimeField representing a literal time.Time value.
func Time(t time.Time) TimeField {
	return TimeField{
		value: &t,
	}
}

// Set returns a FieldAssignment associating the TimeField to the value i.e.
// 'field = value'.
func (f TimeField) Set(value interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

// SetTime returns a FieldAssignment associating the TimeField to the time.Time
// value i.e. 'field = value'.
func (f TimeField) SetTime(value time.Time) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

// As returns a new TimeField with the new field Alias i.e. 'field AS Alias'.
func (f TimeField) As(alias string) TimeField {
	f.alias = alias
	return f
}

// Asc returns a new TimeField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f TimeField) Asc() TimeField {
	desc := false
	f.descending = &desc
	return f
}

// Desc returns a new TimeField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f TimeField) Desc() TimeField {
	desc := true
	f.descending = &desc
	return f
}

// IsNull returns an 'X IS NULL' Predicate.
func (f TimeField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f TimeField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// Eq returns an 'X = Y' Predicate. It only accepts TimeField.
func (f TimeField) Eq(field TimeField) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, field},
	}
}

// EqTime returns an 'X = Y' Predicate. It only accepts time.Time.
func (f TimeField) EqTime(t time.Time) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, t},
	}
}

// Ne returns an 'X <> Y' Predicate. It only accepts TimeField.
func (f TimeField) Ne(field TimeField) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, field},
	}
}

// NeTime returns an 'X <> Y' Predicate. It only accepts time.Time.
func (f TimeField) NeTime(t time.Time) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, t},
	}
}

// Gt returns an 'X > Y' Predicate. It only accepts TimeField.
func (f TimeField) Gt(field TimeField) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, field},
	}
}

// GtTime returns an 'X > Y' Predicate. It only accepts time.Time.
func (f TimeField) GtTime(t time.Time) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, t},
	}
}

// Ge returns an 'X >= Y' Predicate. It only accepts TimeField.
func (f TimeField) Ge(field TimeField) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, field},
	}
}

// GeTime returns an 'X >= Y' Predicate. It only accepts time.Time.
func (f TimeField) GeTime(t time.Time) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, t},
	}
}

// Lt returns an 'X < Y' Predicate. It only accepts TimeField.
func (f TimeField) Lt(field TimeField) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, field},
	}
}

// LtTime returns an 'X < Y' Predicate. It only accepts time.Time.
func (f TimeField) LtTime(t time.Time) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, t},
	}
}

// Le returns an 'X <= Y' Predicate. It only accepts TimeField.
func (f TimeField) Le(field TimeField) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, field},
	}
}

// LeTime returns an 'X <= Y' Predicate. It only accepts time.Time.
func (f TimeField) LeTime(t time.Time) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, t},
	}
}

// Between returns an 'X BETWEEN Y AND Z' Predicate. It only accepts TimeField.
func (f TimeField) Between(start, end TimeField) Predicate {
	return CustomPredicate{
		Format: "? BETWEEN ? AND ?",
		Values: []interface{}{f, start, end},
	}
}

// BetweenTime returns an 'X BETWEEN Y AND Z' Predicate. It only accepts
// time.Time.
func (f TimeField) BetweenTime(start, end time.Time) Predicate {
	return CustomPredicate{
		Format: "? BETWEEN ? AND ?",
		Values: []interface{}{f, start, end},
	}
}

// NotBetween returns an 'X NOT BETWEEN Y AND Z' Predicate. It only accepts
// TimeField.
func (f TimeField) NotBetween(start, end TimeField) Predicate {
	return CustomPredicate{
		Format: "? NOT BETWEEN ? AND ?",
		Values: []interface{}{f, start, end},
	}
}

// NotBetweenTime returns an 'X NOT BETWEEN Y AND Z' Predicate. It only accepts
// time.Time.
func (f TimeField) NotBetweenTime(start, end time.Time) Predicate {
	return CustomPredicate{
		Format: "? NOT BETWEEN ? AND ?",
		Values: []interface{}{f, start, end},
	}
}

// String returns the string representation of the TimeField.
func (f TimeField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil)
	return QuestionInterpolate(buf.String(), args...)
}

// GetAlias returns the alias of the TimeField.
func (f TimeField) GetAlias() string {
	return f.alias
}

// GetName returns the name of the TimeField.
func (f TimeField) GetName() string {
	return f.name
}
