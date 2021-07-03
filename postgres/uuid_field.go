package sq

import (
	"strings"

	"github.com/google/uuid"
)

// UUIDField represents a UUID column or a literal UUID value.
type UUIDField struct {
	// 1) Literal UUID value
	// Examples of literal UUID values:
	// | query | args |
	// |-------|---------|
	// | ?     | [16]byte(123e4567-e89b-12d3-a456-426614174000) |
	value *uuid.UUID

	// 2) UUID column
	alias string
	table Table
	name string
	descending *bool
}

// AppendSQLExclude marshals the UUIDField into a buffer and an args slice.
// It will not table qualify itself if its table qualifer appears in the
// excludedTableQualifiers list.
func (f UUIDField) AppendSQLExclude(buf *strings.Builder, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) {
	switch {
	case f.value != nil:
		// 1) Literal UUID value
		buf.WriteString("?")
		*args = append(*args, *f.value)
	default:
		// 2) UUID column
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
}

func NewUUIDField(name string, table Table) UUIDField {
	return UUIDField{
		name: name,
		table: table,
	}
}

func UUID(u uuid.UUID) UUIDField {
	return UUIDField{
		value: &u,
	}
}

func (f UUIDField) Set(value interface{}) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: value,
	}
}

func (f UUIDField) SetUUID(u uuid.UUID) FieldAssignment {
	return FieldAssignment{
		Field: f,
		Value: u,
	}
}

func (f UUIDField) IsNull() Predicate {
	return CustomPredicate{
		Format: "? IS NULL",
		Values: []interface{}{f},
	}
}

func (f UUIDField) IsNotNull() Predicate {
	return CustomPredicate{
		Format: "? IS NOT NULL",
		Values: []interface{}{f},
	}
}

func (f UUIDField) Eq(field UUIDField) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, field},
	}
}

func (f UUIDField) Ne(field UUIDField) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, field},
	}
}

func (f UUIDField) EqUUID(u uuid.UUID) Predicate {
	return CustomPredicate{
		Format: "? = ?",
		Values: []interface{}{f, u},
	}
}

func (f UUIDField) NeUUID(u uuid.UUID) Predicate {
	return CustomPredicate{
		Format: "? <> ?",
		Values: []interface{}{f, u},
	}
}

func (f UUIDField) In(v interface{}) Predicate {
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

func (f UUIDField) String() string {
	buf := &strings.Builder{}
	var args []interface{}
	f.AppendSQLExclude(buf, &args, nil, nil)
	return questionInterpolate(buf.String(), args...)
}

func (f UUIDField) GetAlias() string {
	return f.alias
}

func (f UUIDField) GetName() string {
	return f.name
}
