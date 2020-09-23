package sq

import "strings"

// Subquery represents an SQL subquery.
type Subquery map[string]CustomField

// Subquery converts a SelectQuery into a Subquery.
func (q SelectQuery) Subquery(alias string) Subquery {
	subquery := Subquery{
		metadataQuery: {Values: []interface{}{q}},
		metadataAlias: {Values: []interface{}{alias}},
	}
	for _, field := range q.SelectFields {
		column := field.GetAlias()
		if column == "" {
			column = field.GetName()
		}
		subquery[column] = CustomField{Format: alias + "." + column}
	}
	return subquery
}

// Subquery converts a VariadicQuery into a Subquery.
func (vq VariadicQuery) Subquery(name string) Subquery {
	subquery := map[string]CustomField{
		metadataQuery: {Values: []interface{}{vq}},
		metadataAlias: {Values: []interface{}{name}},
	}
	if len(vq.Queries) > 0 {
		switch q := vq.Queries[0].(type) {
		case SelectQuery:
			for _, field := range q.SelectFields {
				column := getAliasOrName(field)
				subquery[column] = CustomField{Format: name + "." + column}
			}
		}
	}
	return subquery
}

// ToSQL marshals the Subquery into a query string and args slice.
func (subq Subquery) ToSQL() (string, []interface{}) {
	buf := &strings.Builder{}
	var args []interface{}
	subq.AppendSQL(buf, &args, nil)
	return buf.String(), args
}

// AppendSQL marshals the Subquery into a buffer and args slice.
func (subq Subquery) AppendSQL(buf *strings.Builder, args *[]interface{}, params map[string]int) {
	q := subq.GetQuery()
	if q == nil {
		return
	}
	q.NestThis().AppendSQL(buf, args, nil)
}

// GetQuery returns the Subquery's underlying Query.
func (subq Subquery) GetQuery() Query {
	field := subq[metadataQuery]
	if len(field.Values) > 0 {
		if q, ok := field.Values[0].(Query); ok {
			return q
		}
	}
	return nil
}

// GetName returns the name of the Subquery.
func (subq Subquery) GetName() string {
	return ""
}

// GetAlias returns the alias of the Subquery.
func (subq Subquery) GetAlias() string {
	field := subq[metadataAlias]
	if len(field.Values) > 0 {
		if alias, ok := field.Values[0].(string); ok {
			return alias
		}
	}
	return ""
}

// NestThis indicates to the Subquery that it is nested.
func (subq Subquery) NestThis() Query {
	return subq
}
