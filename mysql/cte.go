package sq

import (
	"strings"
)

// https://www.topster.net/text/utf-schriften.html serif italics
const (
	metadataQuery     = "𝑞𝑢𝑒𝑟𝑦"
	metadataRecursive = "𝑟𝑒𝑐𝑢𝑟𝑠𝑖𝑣𝑒"
	metadataName      = "𝑛𝑎𝑚𝑒"
	metadataAlias     = "𝑎𝑙𝑖𝑎𝑠"
	metadataColumns   = "𝑐𝑜𝑙𝑢𝑚𝑛𝑠"
)

// CTE represents an SQL CTE.
type CTE map[string]CustomField

func appendCTEs(buf *strings.Builder, args *[]interface{}, CTEs []CTE, fromTable Table, joinTables []JoinTable) {
	type TmpCTE struct {
		name    string
		columns []string
		query   Query
	}
	var tmpCTEs []TmpCTE
	cteNames := map[string]bool{} // track CTE names we have already seen; used to remove duplicates
	hasRecursiveCTE := false
	addTmpCTE := func(table Table) {
		cte, ok := table.(CTE)
		if !ok {
			return // not a CTE, skip
		}
		name := cte.GetName()
		if cteNames[name] {
			return // already seen this CTE, skip
		} else {
			cteNames[name] = true
		}
		if !hasRecursiveCTE && cte.IsRecursive() {
			hasRecursiveCTE = true
		}
		tmpCTEs = append(tmpCTEs, TmpCTE{
			name:    name,
			columns: cte.GetColumns(),
			query:   cte.GetQuery(),
		})
	}
	for _, cte := range CTEs {
		addTmpCTE(cte)
	}
	addTmpCTE(fromTable)
	for _, joinTable := range joinTables {
		addTmpCTE(joinTable.Table)
	}
	if len(tmpCTEs) == 0 {
		return // there were no CTEs in the list of tables, return
	}
	if hasRecursiveCTE {
		buf.WriteString("WITH RECURSIVE ")
	} else {
		buf.WriteString("WITH ")
	}
	for i, cte := range tmpCTEs {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(cte.name)
		if len(cte.columns) > 0 {
			buf.WriteString(" (")
			buf.WriteString(strings.Join(cte.columns, ", "))
			buf.WriteString(")")
		}
		buf.WriteString(" AS (")
		switch q := cte.query.(type) {
		case nil:
			buf.WriteString("NULL")
		case VariadicQuery:
			q.topLevel = true
			q.NestThis().AppendSQL(buf, args)
		default:
			q.NestThis().AppendSQL(buf, args)
		}
		buf.WriteString(")")
	}
	buf.WriteString(" ")
}

// CTE converts a SelectQuery into a CTE.
func (q SelectQuery) CTE(name string, columns ...string) CTE {
	cte := map[string]CustomField{
		metadataQuery:   {Values: []interface{}{q}},
		metadataName:    {Values: []interface{}{name}},
		metadataAlias:   {Values: []interface{}{""}},
		metadataColumns: {Values: []interface{}{columns}},
	}
	for _, field := range q.SelectFields {
		column := getAliasOrName(field)
		cte[column] = CustomField{Format: name + "." + column}
	}
	return cte
}

// CTE converts a VariadicQuery into a CTE.
func (vq VariadicQuery) CTE(name string, columns ...string) CTE {
	cte := map[string]CustomField{
		metadataQuery:   {Values: []interface{}{vq}},
		metadataName:    {Values: []interface{}{name}},
		metadataAlias:   {Values: []interface{}{""}},
		metadataColumns: {Values: []interface{}{columns}},
	}
	if len(columns) > 0 {
		for _, column := range columns {
			cte[column] = CustomField{Format: name + "." + column}
		}
		return cte
	}
	if len(vq.Queries) > 0 {
		switch q := vq.Queries[0].(type) {
		case SelectQuery:
			for _, field := range q.SelectFields {
				column := getAliasOrName(field)
				cte[column] = CustomField{Format: name + "." + column}
			}
		}
	}
	return cte
}

// As returns a new CTE with the alias i.e. 'CTE AS alias'.
func (cte CTE) As(alias string) CTE {
	newcte := map[string]CustomField{
		metadataQuery:   {Values: []interface{}{cte.GetQuery()}},
		metadataName:    {Values: []interface{}{cte.GetName()}},
		metadataAlias:   {Values: []interface{}{alias}},
		metadataColumns: {Values: []interface{}{cte.GetColumns()}},
	}
	for column := range cte {
		switch column {
		case metadataQuery, metadataName, metadataAlias, metadataColumns:
			continue
		}
		newcte[column] = CustomField{Format: alias + "." + column}
	}
	return newcte
}

// AppendSQL marshals the CTE into a buffer and args slice.
func (cte CTE) AppendSQL(buf *strings.Builder, _ *[]interface{}) {
	buf.WriteString(cte.GetName())
}

// IsRecursive checks if the CTE is recursive.
func (cte CTE) IsRecursive() bool {
	field := cte[metadataRecursive]
	if len(field.Values) > 0 {
		if recursive, ok := field.Values[0].(bool); ok {
			return recursive
		}
	}
	return false
}

// GetQuery returns the CTE's underlying Query.
func (cte CTE) GetQuery() Query {
	field := cte[metadataQuery]
	if len(field.Values) > 0 {
		if q, ok := field.Values[0].(Query); ok {
			return q
		}
	}
	return nil
}

// GetQuery returns the CTE's columns.
func (cte CTE) GetColumns() []string {
	field := cte[metadataColumns]
	if len(field.Values) > 0 {
		if columns, ok := field.Values[0].([]string); ok {
			return columns
		}
	}
	return nil
}

// GetName returns the name of the CTE.
func (cte CTE) GetName() string {
	field := cte[metadataName]
	if len(field.Values) > 0 {
		if name, ok := field.Values[0].(string); ok {
			return name
		}
	}
	return ""
}

// GetAlias returns the alias of the CTE.
func (cte CTE) GetAlias() string {
	field := cte[metadataAlias]
	if len(field.Values) > 0 {
		if alias, ok := field.Values[0].(string); ok {
			return alias
		}
	}
	return ""
}

// RecursiveCTE constructs a new recursive CTE.
func RecursiveCTE(name string, columns ...string) CTE {
	cte := map[string]CustomField{
		metadataRecursive: {Values: []interface{}{true}},
		metadataName:      {Values: []interface{}{name}},
		metadataAlias:     {Values: []interface{}{""}},
	}
	if len(columns) > 0 {
		cte[metadataColumns] = CustomField{Values: []interface{}{columns}}
		for _, column := range columns {
			cte[column] = CustomField{Format: name + "." + column}
		}
	}
	return cte
}

// IntermediateCTE is a CTE used to hold the intermediate state of a recursive
// CTE just after the CTE's initial query is declared. It can only be converted
// back into a CTE by adding the recursive queries that UNION into the CTE.
type IntermediateCTE map[string]CustomField

// Initial specifies recursive CTE's initial query. If the CTE is not
// recursive, this operation is a no-op.
func (cte *CTE) Initial(query Query) IntermediateCTE {
	if !cte.IsRecursive() {
		return IntermediateCTE(*cte)
	}
	if *cte == nil {
		*cte = map[string]CustomField{}
	}
	(*cte)[metadataQuery] = CustomField{Values: []interface{}{query}}
	name := cte.GetName()
	columns := cte.GetColumns()
	if len(columns) > 0 {
		return IntermediateCTE(*cte)
	}
	switch q := query.(type) {
	case SelectQuery:
		for _, field := range q.SelectFields {
			column := getAliasOrName(field)
			(*cte)[column] = CustomField{Format: name + "." + column}
		}
	}
	return IntermediateCTE(*cte)
}

// Union specifies the queries to be UNIONed into the CTE. If the CTE is not
// recursive, this operation is a no-op.
func (cte IntermediateCTE) Union(queries ...Query) CTE {
	if !CTE(cte).IsRecursive() {
		return CTE(cte)
	}
	return cte.union(queries, QueryUnion)
}

// Union specifies the queries to be UNION-ALLed into the CTE. If the CTE is
// not recursive, this operation is a no-op.
func (cte IntermediateCTE) UnionAll(queries ...Query) CTE {
	if !CTE(cte).IsRecursive() {
		return CTE(cte)
	}
	return cte.union(queries, QueryUnionAll)
}

func (cte *IntermediateCTE) union(queries []Query, operator VariadicQueryOperator) CTE {
	if *cte == nil {
		*cte = map[string]CustomField{}
	}
	initialQuery := CTE(*cte).GetQuery()
	(*cte)[metadataQuery] = CustomField{Values: []interface{}{VariadicQuery{
		Operator: operator,
		Queries:  append([]Query{initialQuery}, queries...),
	}}}
	return CTE(*cte)
}
