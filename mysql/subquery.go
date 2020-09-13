package sq

import "strings"

type Subquery map[string]CustomField

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

func (subq Subquery) ToSQL() (string, []interface{}) {
	buf := &strings.Builder{}
	var args []interface{}
	subq.AppendSQL(buf, &args)
	return buf.String(), args
}

func (subq Subquery) AppendSQL(buf *strings.Builder, args *[]interface{}) {
	q := subq.GetQuery()
	if q == nil {
		return
	}
	q.NestThis().AppendSQL(buf, args)
}

func (subq Subquery) GetQuery() Query {
	field := subq[metadataQuery]
	if len(field.Values) > 0 {
		if q, ok := field.Values[0].(Query); ok {
			return q
		}
	}
	return nil
}

func (subq Subquery) GetName() string {
	return ""
}

func (subq Subquery) GetAlias() string {
	field := subq[metadataAlias]
	if len(field.Values) > 0 {
		if alias, ok := field.Values[0].(string); ok {
			return alias
		}
	}
	return ""
}

func (subq Subquery) NestThis() Query {
	return subq
}
