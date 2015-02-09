package sqlbuilder

import (
	"fmt"
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
		return "AUTOINCREMENT", nil
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
	typ := ""
	switch cc.Type() {
	case columnTypeInt:
		typ = "INTEGER"
	case columnTypeString:
		typ = fmt.Sprintf("VARCHAR(%d)", cc.Size())
	case columnTypeDate:
		typ = "DATETIME"
	case columnTypeFloat:
		typ = "FLOAT"
	case columnTypeBool:
		typ = "BOOLEAN"
	case columnTypeBytes:
		typ = "BLOB"
	}

	if typ == "" {
		return "", newError("unknown column type")
	} else {
		return typ, nil
	}
}

func (m MysqlDialect) ColumnOptionToString(co ColumnOption) (string, error) {
	switch co {
	case CO_PrimaryKey:
		return "PRIMARY KEY", nil
	case CO_AutoIncrement:
		return "AUTO_INCREMENT", nil
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
	typ := ""
	switch cc.Type() {
	case columnTypeInt:
		if cc.HasOption(CO_AutoIncrement) {
			typ = "SERIAL"
		}
		typ = "BIGINT"
	case columnTypeString:
		typ = fmt.Sprintf("VARCHAR(%d)", cc.Size())
	case columnTypeDate:
		typ = "TIMESTAMP"
	case columnTypeFloat:
		typ = "REAL"
	case columnTypeBool:
		typ = "BOOLEAN"
	case columnTypeBytes:
		typ = "BYTEA"
	}

	if typ == "" {
		return "", newError("unknown column type")
	} else {
		return typ, nil
	}
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
