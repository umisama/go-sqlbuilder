package sqlbuilder

type Dialect interface {
	QuerySuffix() string
	BindVar(i int) string
	QuoteField(field string) string
}

type SqliteDialect struct{}

func (m SqliteDialect) QuerySuffix() string {
	return ";"
}

func (m SqliteDialect) BindVar(i int) string {
	return "?"
}

func (m SqliteDialect) QuoteField(field string) string {
	return "\"" + field + "\""
}

type MysqlDialect struct{}

func (m MysqlDialect) QuerySuffix() string {
	return ";"
}

func (m MysqlDialect) BindVar(i int) string {
	return "?"
}

func (m MysqlDialect) QuoteField(field string) string {
	return "`" + field + "`"
}
