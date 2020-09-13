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

type CTE map[string]CustomField

func AppendCTEs(buf *strings.Builder, args *[]interface{}, CTEs []CTE, fromTable Table, joinTables []JoinTable) {
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
			q.TopLevel = true
			q.NestThis().AppendSQL(buf, args)
		default:
			q.NestThis().AppendSQL(buf, args)
		}
		buf.WriteString(")")
	}
	buf.WriteString(" ")
}

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

func (cte CTE) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	buf.WriteString(cte.GetName())
}

func (cte CTE) IsRecursive() bool {
	field := cte[metadataRecursive]
	if len(field.Values) > 0 {
		if recursive, ok := field.Values[0].(bool); ok {
			return recursive
		}
	}
	return false
}

func (cte CTE) GetQuery() Query {
	field := cte[metadataQuery]
	if len(field.Values) > 0 {
		if q, ok := field.Values[0].(Query); ok {
			return q
		}
	}
	return nil
}

func (cte CTE) GetColumns() []string {
	field := cte[metadataColumns]
	if len(field.Values) > 0 {
		if columns, ok := field.Values[0].([]string); ok {
			return columns
		}
	}
	return nil
}

func (cte CTE) GetName() string {
	field := cte[metadataName]
	if len(field.Values) > 0 {
		if name, ok := field.Values[0].(string); ok {
			return name
		}
	}
	return ""
}

func (cte CTE) GetAlias() string {
	field := cte[metadataAlias]
	if len(field.Values) > 0 {
		if alias, ok := field.Values[0].(string); ok {
			return alias
		}
	}
	return ""
}

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

type intermediateCTE map[string]CustomField

func (cte CTE) Initial(query Query) intermediateCTE {
	if !cte.IsRecursive() {
		return intermediateCTE(cte)
	}
	if cte == nil {
		cte = map[string]CustomField{}
	}
	cte[metadataQuery] = CustomField{Values: []interface{}{query}}
	name := cte.GetName()
	columns := cte.GetColumns()
	if len(columns) > 0 {
		return intermediateCTE(cte)
	}
	switch q := query.(type) {
	case SelectQuery:
		for _, field := range q.SelectFields {
			column := getAliasOrName(field)
			cte[column] = CustomField{Format: name + "." + column}
		}
	}
	return intermediateCTE(cte)
}

func (cte intermediateCTE) Union(queries ...Query) CTE {
	if !CTE(cte).IsRecursive() {
		return CTE(cte)
	}
	return cte.union(queries, QueryUnion)
}

func (cte intermediateCTE) UnionAll(queries ...Query) CTE {
	if !CTE(cte).IsRecursive() {
		return CTE(cte)
	}
	return cte.union(queries, QueryUnionAll)
}

func (cte intermediateCTE) union(queries []Query, operator VariadicQueryOperator) CTE {
	if cte == nil {
		cte = map[string]CustomField{}
	}
	initialQuery := CTE(cte).GetQuery()
	cte[metadataQuery] = CustomField{Values: []interface{}{VariadicQuery{
		Operator: operator,
		Queries:  append([]Query{initialQuery}, queries...),
	}}}
	return CTE(cte)
}
