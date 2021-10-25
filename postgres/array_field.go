package sq

import (
	"fmt"
	"strings"
)

// ArrayField either represents an ARRAY column, or a literal slice value.
type ArrayField struct {
	// ArrayField will be one of the following:

	// 1) Literal slice value (only []bool, []float64, []int64 or []string
	// slices are supported.) Nested slices are also not supported even though
	// both Go and Postgres support nested slices/arrays because I'm not even
	// sure if it's possible to convert between the two with lib/pq.
	// Additionally []int is supported, but note that it only works when
	// converting from Go slices to Postgres arrays. When converting from
	// postgres arrays to Go slices, you have to use []int64 instead.
	// Examples of literal array values:
	// | query             | args                    |
	// |-------------------|-------------------------|
	// | ARRAY[?, ?, ?, ?] | 1, 2, 3, 4              |
	// | ARRAY[?, ?, ?]    | 22.7, 3.15, 4.0         |
	// | ARRAY[?, ?, ?]    | apple, banana, cucumber |
	value interface{}

	// 2) Array column
	// Examples of boolean columns:
	// | query                 | args |
	// |-----------------------|------|
	// | film.special_features |      |
	// | special_features      |      |
	alias      string
	table      Table
	name       string
	descending *bool
	nullsfirst *bool
}

// AppendSQLExclude marshals the ArrayField into a buffer and an args slice. It
// will not table qualify itself if its table qualifer appears in the
// excludedTableQualifiers list.
func (f ArrayField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	switch {
	case f.value != nil:
		// 1) Literal slice value
		switch array := f.value.(type) {
		case []bool:
			if len(array) == 0 {
				buf.WriteString("ARRAY[]::BOOLEAN[]")
			} else {
				buf.WriteString("ARRAY[?")
				buf.WriteString(strings.Repeat(", ?", len(array)-1))
				buf.WriteString("]")
				for _, arg := range array {
					*args = append(*args, arg)
				}
			}
		case []float64:
			if len(array) == 0 {
				buf.WriteString("ARRAY[]::FLOAT[]")
			} else {
				buf.WriteString("ARRAY[?")
				buf.WriteString(strings.Repeat(", ?", len(array)-1))
				buf.WriteString("]")
				for _, arg := range array {
					*args = append(*args, arg)
				}
			}
		case []int:
			if len(array) == 0 {
				buf.WriteString("ARRAY[]::INT[]")
			} else {
				buf.WriteString("ARRAY[?")
				buf.WriteString(strings.Repeat(", ?", len(array)-1))
				buf.WriteString("]")
				for _, arg := range array {
					*args = append(*args, arg)
				}
			}
		case []int64:
			if len(array) == 0 {
				buf.WriteString("ARRAY[]::BIGINT[]")
			} else {
				buf.WriteString("ARRAY[?")
				buf.WriteString(strings.Repeat(", ?", len(array)-1))
				buf.WriteString("]")
				for _, arg := range array {
					*args = append(*args, arg)
				}
			}
		case []string:
			if len(array) == 0 {
				buf.WriteString("ARRAY[]::TEXT[]")
			} else {
				buf.WriteString("ARRAY[?")
				buf.WriteString(strings.Repeat(", ?", len(array)-1))
				buf.WriteString("]")
				for _, arg := range array {
					*args = append(*args, arg)
				}
			}
		default:
			buf.WriteString(fmt.Sprintf("(unsupported type %#v: only []bool/[]float64/[]int64/[]string/[]int slices are supported.)", f.value))
		}
	default:
		// 2) Array column
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

// NewArrayField returns a new ArrayField representing an array column.
func NewArrayField(name string, table Table) ArrayField {
	return ArrayField{
		name:  name,
		table: table,
	}
}

// Array returns a new ArrayField representing a literal string value.
func Array(slice interface{}) ArrayField {
	return ArrayField{
		value: slice,
	}
}

// Set returns a FieldAssignment associating the ArrayField to the value i.e.
// 'field = value'.
func (f ArrayField) Set(value interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

// As returns a new ArrayField with the new field Alias i.e. 'field AS Alias'.
func (f ArrayField) As(alias string) ArrayField {
	f.alias = alias
	return f
}

// Asc returns a new ArrayField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f ArrayField) Asc() ArrayField {
	desc := false
	f.descending = &desc
	return f
}

// Desc returns a new ArrayField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f ArrayField) Desc() ArrayField {
	desc := true
	f.descending = &desc
	return f
}

// NullsFirst returns a new ArrayField indicating that it should be ordered
// with nulls first i.e. 'ORDER BY field NULLS FIRST'.
func (f ArrayField) NullsFirst() ArrayField {
	nullsfirst := true
	f.nullsfirst = &nullsfirst
	return f
}

// NullsLast returns a new ArrayField indicating that it should be ordered
// with nulls last i.e. 'ORDER BY field NULLS LAST'.
func (f ArrayField) NullsLast() ArrayField {
	nullsfirst := false
	f.nullsfirst = &nullsfirst
	return f
}

// IsNull returns an 'X IS NULL' Predicate.
func (f ArrayField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

// IsNotNull returns an 'X IS NOT NULL' Predicate.
func (f ArrayField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

// Eq returns an 'X = Y' Predicate. It only accepts ArrayField.
func (f ArrayField) Eq(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, field},
	}
}

// Ne returns an 'X <> Y' Predicate. It only accepts ArrayField.
func (f ArrayField) Ne(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, field},
	}
}

// Gt returns an 'X > Y' Predicate. It only accepts ArrayField.
func (f ArrayField) Gt(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? > ?",
		Values: []interface{}{f, field},
	}
}

// Ge returns an 'X >= Y' Predicate. It only accepts ArrayField.
func (f ArrayField) Ge(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? >= ?",
		Values: []interface{}{f, field},
	}
}

// Lt returns an 'X < Y' Predicate. It only accepts ArrayField.
func (f ArrayField) Lt(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? < ?",
		Values: []interface{}{f, field},
	}
}

// Le returns an 'X <= Y' Predicate. It only accepts ArrayField.
func (f ArrayField) Le(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? <= ?",
		Values: []interface{}{f, field},
	}
}

// Contains checks whether the subject ArrayField contains the object
// ArrayField.
func (f ArrayField) Contains(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? @> ?",
		Values: []interface{}{f, field},
	}
}

// ContainedBy checks whether the subject ArrayField is contained by the object
// ArrayField.
func (f ArrayField) ContainedBy(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? <@ ?",
		Values: []interface{}{f, field},
	}
}

// Overlaps checks whether the subject ArrayField and the object ArrayField
// have any values in common.
func (f ArrayField) Overlaps(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? && ?",
		Values: []interface{}{f, field},
	}
}

// Concat concatenates the object ArrayField to the subject ArrayField.
func (f ArrayField) Concat(field ArrayField) Field {
	return CustomField{
		Format: "? || ?",
		Values: []interface{}{f, field},
	}
}

// String implements the fmt.Stringer interface. It returns the string
// representation of an ArrayField.
func (f ArrayField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil, nil)
	return questionInterpolate(buf.String(), args...)
}

// GetAlias implements the Field interface. It returns the Alias of the
// ArrayField.
func (f ArrayField) GetAlias() string {
	return f.alias
}

// GetName implements the Field interface. It returns the Name of the
// ArrayField.
func (f ArrayField) GetName() string {
	return f.name
}
