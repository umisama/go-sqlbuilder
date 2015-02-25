package dialects

import (
	"errors"
	"fmt"
	sb "github.com/umisama/go-sqlbuilder"
)

type MySql struct{}

func (m MySql) QuerySuffix() string {
	return ";"
}

func (m MySql) BindVar(i int) string {
	return "?"
}

func (m MySql) QuoteField(field string) string {
	return "`" + field + "`"
}

func (m MySql) ColumnTypeToString(cc sb.ColumnConfig) (string, error) {
	if cc.Option().SqlType != "" {
		return cc.Option().SqlType, nil
	}

	typ := ""
	switch cc.Type() {
	case sb.ColumnTypeInt:
		typ = "INTEGER"
	case sb.ColumnTypeString:
		typ = fmt.Sprintf("VARCHAR(%d)", cc.Option().Size)
	case sb.ColumnTypeDate:
		typ = "DATETIME"
	case sb.ColumnTypeFloat:
		typ = "FLOAT"
	case sb.ColumnTypeBool:
		typ = "BOOLEAN"
	case sb.ColumnTypeBytes:
		typ = "BLOB"
	}

	if typ == "" {
		return "", errors.New("dialects: unknown column type")
	} else {
		return typ, nil
	}
}

func (m MySql) ColumnOptionToString(co *sb.ColumnOption) (string, error) {
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
