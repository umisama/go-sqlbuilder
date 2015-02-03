package sqlbuilder

import (
	"errors"
)

type Dialect interface {
	QuerySuffix() string
	BindVar(i int) string
	QuoteField(field string) string
	SqlType(ColumnConfig) (string, error)
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

	return "", errors.New("unknown column type")
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

	return "", errors.New("unknown column type")
}
