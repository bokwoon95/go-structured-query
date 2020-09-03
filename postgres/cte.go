package sq

import "strings"

// CTE represents an SQL Common Table Expression.
type CTE struct {
	Recursive bool
	Name      string
	Query     Query
	Columns   []string
}

// ToSQL simply returns the name of the CTE.
func (cte CTE) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	buf.WriteString(cte.Name)
}

// NewCTE creates a new CTE.
func NewCTE(name string, query Query, columns ...string) CTE {
	return CTE{
		Name:    name,
		Query:   query,
		Columns: columns,
	}
}

// NewRecursiveCTE creates a new recursive CTE.
func NewRecursiveCTE(name string, query Query, columns ...string) CTE {
	return CTE{
		Recursive: true,
		Name:      name,
		Query:     query,
		Columns:   columns,
	}
}

// GetAlias implements the Table interface. It always returns an empty string,
// because CTEs do not have aliases (only AliasedCTEs do).
func (cte CTE) GetAlias() string {
	return ""
}

// GetAlias implements the Table interface. It returns the name of the CTE.
func (cte CTE) GetName() string {
	return cte.Name
}

// Get returns a Field from the CTE identified by fieldName. No checks are done
// to see if the fieldName really exists in the CTE at all, CTE simply prepends
// its own name to the fieldName.
func (cte CTE) Get(fieldName string) CustomField {
	return CustomField{
		Format: cte.Name + "." + fieldName,
	}
}

// CTEs represents a list of CTEs
type CTEs []CTE

// AppendSQL will write the CTE clause into the buffer and args. If there are no
// CTEs to be written, it will simply write nothing. It returns a flag
// indicating whether it wrote anything into the buffer.
func (ctes CTEs) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	var hasRecursiveCTE bool
	for i := range ctes {
		if ctes[i].Recursive {
			hasRecursiveCTE = true
			break
		}
	}
	if hasRecursiveCTE {
		buf.WriteString("WITH RECURSIVE ")
	} else {
		buf.WriteString("WITH ")
	}
	for i := range ctes {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ctes[i].Name)
		if len(ctes[i].Columns) > 0 {
			buf.WriteString(" (")
			buf.WriteString(strings.Join(ctes[i].Columns, ", "))
			buf.WriteString(")")
		}
		buf.WriteString(" AS (")
		switch ctes[i].Query.(type) {
		case nil:
			buf.WriteString("NULL")
		default:
			ctes[i].Query.NestThis().AppendSQL(buf, args)
		}
		buf.WriteString(")")
	}
}

// As returns a an Aliased CTE derived from the parent CTE that it was called
// on.
func (cte CTE) As(alias string) AliasedCTE {
	return AliasedCTE{
		Name:  cte.Name,
		Alias: alias,
	}
}

// AliasedCTE is an aliased version of a CTE derived from a parent CTE.
type AliasedCTE struct {
	Name  string
	Alias string
}

// ToSQL returns the name of the parent CTE the AliasedCTE was derived from.
// There is no need to provide the alias, as the caller of ToSQL() should be
// responsible for calling GetAlias() as well.
func (cte AliasedCTE) AppendSQL(buf *strings.Builder, _ *[]interface{}) {
	buf.WriteString(cte.Name)
}

// GetAlias implements the Table interface. It returns the alias of the
// AliasedCTE.
func (cte AliasedCTE) GetAlias() string {
	return cte.Alias
}

// GetAlias implements the Table interface. It returns the name of the parent
// CTE.
func (cte AliasedCTE) GetName() string {
	return cte.Name
}

// Get returns a Field from the AliasedCTE identified by fieldName. No checks
// are done to see if the fieldName really exists in the AliasedCTE at all,
// AliasedCTE simply prepends its own alias to the fieldName.
func (cte AliasedCTE) Get(fieldName string) CustomField {
	return CustomField{
		Format: cte.Alias + "." + fieldName,
	}
}
