package sqlbuilder

type Dialect interface {
	QuerySuffix() string
	BindVar(i int) string
	QuoteField(field string) string
}

type sqliteDialect struct{}

func (m *sqliteDialect) QuerySuffix() string {
	return ";"
}

func (m *sqliteDialect) BindVar(i int) string {
	return "?"
}

func (m *sqliteDialect) QuoteField(field string) string {
	return "\"" + field + "\""
}
