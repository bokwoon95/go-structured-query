package sq

import "strings"

// TableInfo is struct that implements the Table interface, containing all the
// information needed to call itself a Table. It is meant to be embedded in
// arbitrary structs to also transform them into valid Tables.
type TableInfo struct {
	Schema string
	Name   string
	Alias  string
}

// AppendSQLExclude marshals the TableInfo into a buffer and an args slice.
func (tbl *TableInfo) AppendSQL(buf Buffer, args *[]interface{}) {
	if tbl == nil {
		return
	}
	if tbl.Schema != "" {
		if strings.ContainsAny(tbl.Schema, " \t") {
			buf.WriteString("`")
			buf.WriteString(tbl.Schema)
			buf.WriteString("`.")
		} else {
			buf.WriteString(tbl.Schema)
			buf.WriteString(".")
		}
	}
	if strings.ContainsAny(tbl.Name, " \t") {
		buf.WriteString("`")
		buf.WriteString(tbl.Name)
		buf.WriteString("`")
	} else {
		buf.WriteString(tbl.Name)
	}
}

// GetAlias returns the alias of the TableInfo.
func (tbl *TableInfo) GetAlias() string {
	if tbl == nil {
		return ""
	}
	return tbl.Alias
}

// GetName returns the name of the TableInfo.
func (tbl *TableInfo) GetName() string {
	if tbl == nil {
		return ""
	}
	return tbl.Name
}

// AssertBaseTable implements the BaseTable interface.
func (tbl *TableInfo) AssertBaseTable() {}
