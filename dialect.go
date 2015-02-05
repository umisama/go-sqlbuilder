package sqlbuilder

import (
	"strconv"
)

type Dialect interface {
	QuerySuffix() string
	BindVar(i int) string
	QuoteField(field string) string
	SqlType(ColumnConfig) (string, error)
	ColumnOptionToString(ColumnOption) (string, error)
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

func (m SqliteDialect) SqlType(cc ColumnConfig) (string, error) {
	switch cc.Type() {
	case columnTypeInt:
		return "INTEGER", nil
	case columnTypeString:
		return "TEXT", nil
	case columnTypeDate:
		return "DATE", nil
	case columnTypeFloat:
		return "REAL", nil
	case columnTypeBool:
		return "BOOLEAN", nil
	case columnTypeBytes:
		return "BLOB", nil
	}

	return "", newError("unknown column type")
}

func (m SqliteDialect) ColumnOptionToString(co ColumnOption) (string, error) {
	switch co {
	case CO_PrimaryKey:
		return "PRIMARY KEY", nil
	case CO_AutoIncrement:
		return "AUTO INCREMENT", nil
	case CO_NotNull:
		return "NOT NULL", nil
	case CO_Unique:
		return "UNIQUE", nil
	}

	return "", newError("unknown column option")
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

func (m MysqlDialect) SqlType(cc ColumnConfig) (string, error) {
	switch cc.Type() {
	case columnTypeInt:
		return "INTEGER", nil
	case columnTypeString:
		return "VARCHAR(255)", nil // FIXME:
	case columnTypeDate:
		return "DATETIME", nil
	case columnTypeFloat:
		return "FLOAT", nil
	case columnTypeBool:
		return "BOOLEAN", nil
	case columnTypeBytes:
		return "BLOB", nil
	}

	return "", newError("unknown column type")
}

func (m MysqlDialect) ColumnOptionToString(co ColumnOption) (string, error) {
	switch co {
	case CO_PrimaryKey:
		return "PRIMARY KEY", nil
	case CO_AutoIncrement:
		return "AUTO INCREMENT", nil
	case CO_NotNull:
		return "NOT NULL", nil
	case CO_Unique:
		return "UNIQUE", nil
	}

	return "", newError("unknown column option")
}

type PostgresDialect struct{}

func (m PostgresDialect) QuerySuffix() string {
	return ";"
}

func (m PostgresDialect) BindVar(i int) string {
	return "$" + strconv.Itoa(i)
}

func (m PostgresDialect) QuoteField(field string) string {
	return "\"" + field + "\""
}

func (m PostgresDialect) SqlType(cc ColumnConfig) (string, error) {
	switch cc.Type() {
	case columnTypeInt:
		if cc.HasOption(CO_AutoIncrement) {
			return "SERIAL", nil
		}
		return "BIGINT", nil
	case columnTypeString:
		return "VARCHAR(255)", nil // FIXME:
	case columnTypeDate:
		return "TIMESTAMP", nil
	case columnTypeFloat:
		return "REAL", nil
	case columnTypeBool:
		return "BOOLEAN", nil
	case columnTypeBytes:
		return "BYTEA", nil
	}

	return "", newError("unknown column type")
}

func (m PostgresDialect) ColumnOptionToString(co ColumnOption) (string, error) {
	switch co {
	case CO_PrimaryKey:
		return "PRIMARY KEY", nil
	case CO_AutoIncrement:
		return "", nil
	case CO_NotNull:
		return "NOT NULL", nil
	case CO_Unique:
		return "UNIQUE", nil
	}

	return "", newError("unknown column option")
}
