package sqlbuilder

import (
	"fmt"
	"strconv"
)

type Dialect interface {
	QuerySuffix() string
	BindVar(i int) string
	QuoteField(field string) string
	ColumnTypeToString(ColumnConfig) (string, error)
	ColumnOptionToString(*ColumnOption) (string, error)
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

func (m SqliteDialect) ColumnTypeToString(cc ColumnConfig) (string, error) {
	typ := ""
	switch cc.Type() {
	case columnTypeInt:
		typ = "INTEGER"
	case columnTypeString:
		typ = "TEXT"
	case columnTypeDate:
		typ = "DATE"
	case columnTypeFloat:
		typ = "REAL"
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

func (m SqliteDialect) ColumnOptionToString(co *ColumnOption) (string, error) {
	apnd := func(str, opt string) string {
		if len(str) != 0 {
			str += " "
		}
		str += opt
		return str
	}

	opt := ""
	if co.PrimaryKey {
		opt = apnd(opt, "PRIMARY KEY")
	}
	if co.AutoIncrement {
		opt = apnd(opt, "AUTOINCREMENT")
	}
	if co.NotNull {
		opt = apnd(opt, "NOT NULL")
	}
	if co.Unique {
		opt = apnd(opt, "UNIQUE")
	}

	return opt, nil
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

func (m MysqlDialect) ColumnTypeToString(cc ColumnConfig) (string, error) {
	typ := ""
	switch cc.Type() {
	case columnTypeInt:
		typ = "INTEGER"
	case columnTypeString:
		typ = fmt.Sprintf("VARCHAR(%d)", cc.Option().Size)
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

func (m MysqlDialect) ColumnOptionToString(co *ColumnOption) (string, error) {
	apnd := func(str, opt string) string {
		if len(str) != 0 {
			str += " "
		}
		str += opt
		return str
	}

	opt := ""
	if co.PrimaryKey {
		opt = apnd(opt, "PRIMARY KEY")
	}
	if co.AutoIncrement {
		opt = apnd(opt, "AUTO_INCREMENT")
	}
	if co.NotNull {
		opt = apnd(opt, "NOT NULL")
	}
	if co.Unique {
		opt = apnd(opt, "UNIQUE")
	}

	return opt, nil
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

func (m PostgresDialect) ColumnTypeToString(cc ColumnConfig) (string, error) {
	typ := ""
	switch cc.Type() {
	case columnTypeInt:
		if cc.Option().AutoIncrement {
			typ = "SERIAL"
		} else {
			typ = "BIGINT"
		}
	case columnTypeString:
		typ = fmt.Sprintf("VARCHAR(%d)", cc.Option().Size)
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

func (m PostgresDialect) ColumnOptionToString(co *ColumnOption) (string, error) {
	apnd := func(str, opt string) string {
		if len(str) != 0 {
			str += " "
		}
		str += opt
		return str
	}

	opt := ""
	if co.PrimaryKey {
		opt = apnd(opt, "PRIMARY KEY")
	}
	if co.AutoIncrement {
		// do nothing
	}
	if co.NotNull {
		opt = apnd(opt, "NOT NULL")
	}
	if co.Unique {
		opt = apnd(opt, "UNIQUE")
	}

	return opt, nil
}
