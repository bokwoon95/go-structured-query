package sq

import "strings"

type Window struct {
	WindowName        string
	RenderName        bool
	PartitionByFields Fields
	OrderByFields     Fields
	FrameDefinition   string
}

func (w Window) AppendSQL(buf *strings.Builder, args *[]interface{}) {
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

func (w Window) As(name string) Window {
	w.WindowName = name
	return w
}

func (w Window) Name() Window {
	if w.WindowName == "" {
		w.WindowName = RandomString(8)
	}
	w.RenderName = true
	return w
}

func PartitionBy(fields ...Field) Window {
	return Window{
		PartitionByFields: fields,
	}
}

func OrderBy(fields ...Field) Window {
	return Window{
		OrderByFields: fields,
	}
}

func (w Window) PartitionBy(fields ...Field) Window {
	w.PartitionByFields = fields
	return w
}

func (w Window) OrderBy(fields ...Field) Window {
	w.OrderByFields = fields
	return w
}

func (w Window) Frame(frameDefinition string) Window {
	w.FrameDefinition = frameDefinition
	return w
}

type Windows []Window

func (ws Windows) AppendSQL(buf *strings.Builder, args *[]interface{}) {
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
