package sqlbuilder

var dialect Dialect

type Statement interface {
	ToSql() (query string, attrs []interface{}, err error)
	Error() error
}

type serializable interface {
	serialize() (parts string, attrs []interface{}, err error)
}

type Clause interface {
	serializable
}

type Expression interface {
	serializable
}

func appendExpressionsToQuery(parts []Expression, query string, attrs []interface{}, separator string) (string, []interface{}, error) {
	first := true
	for _, part := range parts {
		if first {
			first = false
		} else {
			query += separator
		}

		var err error
		query, attrs, err = appendItemToQuery(part, query, attrs)
		if err != nil {
			return "", []interface{}{}, nil
		}
	}
	return query, attrs, nil
}

func appendItemsToQuery(parts []serializable, query string, attrs []interface{}, separator string) (string, []interface{}, error) {
	first := true
	for _, part := range parts {
		if first {
			first = false
		} else {
			query += separator
		}

		var err error
		query, attrs, err = appendItemToQuery(part, query, attrs)
		if err != nil {
			return "", []interface{}{}, nil
		}
	}
	return query, attrs, nil
}

func appendItemToQuery(part serializable, query string, attrs []interface{}) (string, []interface{}, error) {
	parts_query, parts_attrs, err := part.serialize()
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
