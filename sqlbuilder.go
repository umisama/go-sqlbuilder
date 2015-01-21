package sqlbuilder

var dialect Dialect

type Sqlizable interface {
	ToSql() (query string, attrs []interface{}, err error)
}

func appendListToQuery(parts []Sqlizable, query string, attrs []interface{}, separator string) (string, []interface{}, error) {
	first := true
	for _, part := range parts {
		if first {
			first = false
		} else {
			query += separator
		}

		var err error
		query, attrs, err = appendToQuery(part, query, attrs)
		if err != nil {
			return "", []interface{}{}, nil
		}
	}
	return query, attrs, nil
}

func appendToQuery(part Sqlizable, query string, attrs []interface{}) (string, []interface{}, error) {
	parts_query, parts_attrs, err := part.ToSql()
	if err != nil {
		return "", []interface{}{}, err
	}

	query += parts_query
	attrs = append(attrs, parts_attrs...)
	return query, attrs, nil
}

func SetDialect(opt Dialect) {
	dialect = opt
}

func init() {
	// initial setup
	SetDialect(&sqliteDialect{})
}
