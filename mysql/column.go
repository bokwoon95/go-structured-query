package sq

import "time"

type colmode int

const (
	colmodeInsert colmode = iota
	colmodeUpdate
)

type Column struct {
	// mode determines if INSERT or UPDATE
	mode colmode
	// INSERT
	rowStart      bool
	rowEnd        bool
	firstField    string
	insertColumns Fields
	rowValues     RowValues
	// UPDATE
	assignments Assignments
}

func (col *Column) Set(field Field, value interface{}) {
	if field == nil {
		// should I panic with an error here instead?
		return
	}
	switch col.mode {
	case colmodeUpdate:
		col.assignments = append(col.assignments, FieldAssignment{
			Field: field,
			Value: value,
		})
	case colmodeInsert:
		fallthrough
	default:
		name := field.GetName()
		if !col.rowStart {
			col.rowStart = true
			col.firstField = name
			col.insertColumns = append(col.insertColumns, field)
			col.rowValues = append(col.rowValues, RowValue{value})
			return
		}
		switch name {
		case col.firstField: // Start a new RowValue
			if !col.rowEnd {
				col.rowEnd = true
			}
			col.rowValues = append(col.rowValues, RowValue{value})
		default: // Append to last RowValue
			if !col.rowEnd {
				col.insertColumns = append(col.insertColumns, field)
			}
			last := len(col.rowValues) - 1
			col.rowValues[last] = append(col.rowValues[last], value)
		}
	}
}

func (col *Column) SetBool(field BooleanField, value bool) {
	col.Set(field, value)
}

func (col *Column) SetFloat64(field NumberField, value float64) {
	col.Set(field, value)
}

func (col *Column) SetInt(field NumberField, value int) {
	col.Set(field, value)
}

func (col *Column) SetInt64(field NumberField, value int64) {
	col.Set(field, value)
}

func (col *Column) SetString(field StringField, value string) {
	col.Set(field, value)
}

func (col *Column) SetTime(field TimeField, value time.Time) {
	col.Set(field, value)
}
