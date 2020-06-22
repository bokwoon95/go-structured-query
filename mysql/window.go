package sq

// Window represents a window usable in a window function.
type Window struct {
	WindowName        string
	RenderName        bool
	PartitionByFields Fields
	OrderByFields     Fields
	FrameDefinition   string
}

// AppendSQL marshals the Window into a buffer and args slice.
func (w Window) AppendSQL(buf Buffer, args *[]interface{}) {
	if w.RenderName {
		buf.WriteString(w.WindowName)
		return
	}
	buf.WriteString("(")
	var written bool
	if len(w.PartitionByFields) > 0 {
		buf.WriteString("PARTITION BY ")
		w.PartitionByFields.AppendSQLExclude(buf, args, nil)
		written = true
	}
	if len(w.OrderByFields) > 0 {
		if written {
			buf.WriteString(" ")
		}
		buf.WriteString("ORDER BY ")
		w.OrderByFields.AppendSQLExclude(buf, args, nil)
		written = true
	}
	if w.FrameDefinition != "" {
		if written {
			buf.WriteString(" ")
		}
		buf.WriteString(w.FrameDefinition)
	}
	buf.WriteString(")")
}

// As aliases the VariadicQuery i.e. 'query AS alias'.
func (w Window) As(name string) Window {
	w.WindowName = name
	return w
}

// Name returns the name of the Window.
func (w Window) Name() Window {
	if w.WindowName == "" {
		w.WindowName = RandomString(8)
	}
	w.RenderName = true
	return w
}

// PartitionBy creates a new Window.
func PartitionBy(fields ...Field) Window {
	return Window{
		PartitionByFields: fields,
	}
}

// OrderBy creates a new Window.
func OrderBy(fields ...Field) Window {
	return Window{
		OrderByFields: fields,
	}
}

// PartitionBy sets the fields of the window's PARTITION BY clause.
func (w Window) PartitionBy(fields ...Field) Window {
	w.PartitionByFields = fields
	return w
}

// OrderBy sets the fields of the window's ORDER BY clause.
func (w Window) OrderBy(fields ...Field) Window {
	w.OrderByFields = fields
	return w
}

// Frame sets the frame definition of the window e.g. RANGE BETWEEN 5 PRECEDING
// AND 10 FOLLOWING.
func (w Window) Frame(frameDefinition string) Window {
	w.FrameDefinition = frameDefinition
	return w
}

// Windows is a list of Windows.
type Windows []Window

// AppendSQL marshals the Windows into a buffer and args slice.
func (ws Windows) AppendSQL(buf Buffer, args *[]interface{}) {
	for i := range ws {
		if i > 0 {
			buf.WriteString(", ")
		}
		if ws[i].WindowName != "" {
			buf.WriteString(ws[i].WindowName)
		} else {
			buf.WriteString(RandomString(8))
		}
		buf.WriteString(" AS ")
		ws[i].AppendSQL(buf, args)
	}
}

// RowNumberOver represents the ROW_NUMBER() OVER window function.
func RowNumberOver(window Window) NumberField {
	format := "ROW_NUMBER() OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{window},
	}
}

// RankOver represents the RANK() OVER window function.
func RankOver(window Window) NumberField {
	format := "RANK() OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{window},
	}
}

// DenseRankOver represents the DENSE_RANK() OVER window function.
func DenseRankOver(window Window) NumberField {
	format := "DENSE_RANK() OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{window},
	}
}

// PercentRankOver represents the PERCENT_RANK() OVER window function.
func PercentRankOver(window Window) NumberField {
	format := "PERCENT_RANK() OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{window},
	}
}

// CumeDistOver represents the CUME_DIST() OVER window function.
func CumeDistOver(window Window) NumberField {
	format := "CUME_DIST() OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{window},
	}
}

// LeadOver represents the LEAD(field, offset, fallback) OVER window function.
func LeadOver(field interface{}, offset interface{}, fallback interface{}, window Window) CustomField {
	if offset == nil {
		offset = 1
	}
	return CustomField{
		Format: "LEAD(?, ?, ?) OVER ?",
		Values: []interface{}{field, offset, fallback, window},
	}
}

// LagOver represents the LAG(field, offset, fallback) OVER window function.
func LagOver(field interface{}, offset interface{}, fallback interface{}, window Window) CustomField {
	if offset == nil {
		offset = 1
	}
	return CustomField{
		Format: "LAG(?, ?, ?) OVER ?",
		Values: []interface{}{field, offset, fallback, window},
	}
}

// NtileOver represents the NTILE(n) OVER window function.
func NtileOver(n int, window Window) NumberField {
	format := "NTILE(?) OVER ?"
	return NumberField{
		format: &format,
		values: []interface{}{n, window},
	}
}

// FirstValueOver represents the FIRST_VALUE(field) OVER window function.
func FirstValueOver(field interface{}, window Window) CustomField {
	return CustomField{
		Format: "FIRST_VALUE(?) OVER ?",
		Values: []interface{}{field, window},
	}
}

// LastValueOver represents the LAST_VALUE(field) OVER window function.
func LastValueOver(field interface{}, window Window) CustomField {
	return CustomField{
		Format: "LAST_VALUE(?) OVER ?",
		Values: []interface{}{field, window},
	}
}

// NthValueOver represents the NTH_VALUE(field, n) OVER window function.
func NthValueOver(field interface{}, n int, window Window) CustomField {
	return CustomField{
		Format: "NTH_VALUE(?, ?) OVER ?",
		Values: []interface{}{field, n, window},
	}
}
