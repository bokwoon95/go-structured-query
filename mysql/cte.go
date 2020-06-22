package sq

import "strings"

// CTE represents an SQL Common Table Expression.
type CTE struct {
	Recursive bool
	Name      string
	Query     Query
	Columns   []string
}

// AppendSQL marshals the CTE name into a buffer and an args slice.
func (cte CTE) AppendSQL(buf Buffer, args *[]interface{}) {
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

// GetAlias returns the alias of the CTE, which is always an empty string. This
// is because CTEs do not have aliases, only AliasedCTEs do.
func (cte CTE) GetAlias() string {
	return ""
}

// GetName returns the name of the CTE.
func (cte CTE) GetName() string {
	return cte.Name
}

// Get returns a Field from the CTE, identified by fieldName.
func (cte CTE) Get(fieldName string) CustomField {
	return CustomField{
		Format: cte.Name + "." + fieldName,
	}
}

// CTEs represents a list of CTEs
type CTEs []CTE

// AppendSQL marshals the CTEs into a buffer and an args slice.
func (ctes CTEs) AppendSQL(buf Buffer, args *[]interface{}) {
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

// As aliases the CTE and returns an AliasedCTE.
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

// AppendSQL marshals the AliasedCTE name into a buffer and an args slice.
// There is no need for it to write the alias into the buffer, as the caller of
// AppendSQL should be responsible for checking for the alias as well.
func (cte AliasedCTE) AppendSQL(buf Buffer, _ *[]interface{}) {
	buf.WriteString(cte.Name)
}

// GetAlias returns the alias of the AliasedCTE.
func (cte AliasedCTE) GetAlias() string {
	return cte.Alias
}

// GetAlias returns the name of the parent CTE.
func (cte AliasedCTE) GetName() string {
	return cte.Name
}

// Get returns a Field from the AliasedCTE, identified by fieldName.
func (cte AliasedCTE) Get(fieldName string) CustomField {
	return CustomField{
		Format: cte.Alias + "." + fieldName,
	}
}
