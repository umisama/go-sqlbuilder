package dialects

import (
	"errors"
	"fmt"
	sb "github.com/umisama/go-sqlbuilder"
	"strconv"
)

type Postgresql struct{}

func (m Postgresql) QuerySuffix() string {
	return ";"
}

func (m Postgresql) BindVar(i int) string {
	return "$" + strconv.Itoa(i)
}

func (m Postgresql) QuoteField(field string) string {
	return "\"" + field + "\""
}

func (m Postgresql) ColumnTypeToString(cc sb.ColumnConfig) (string, error) {
	if cc.Option().SqlType != "" {
		return cc.Option().SqlType, nil
	}

	typ := ""
	switch cc.Type() {
	case sb.ColumnTypeInt:
		if cc.Option().AutoIncrement {
			typ = "SERIAL"
		} else {
			typ = "BIGINT"
		}
	case sb.ColumnTypeString:
		typ = fmt.Sprintf("VARCHAR(%d)", cc.Option().Size)
	case sb.ColumnTypeDate:
		typ = "TIMESTAMP"
	case sb.ColumnTypeFloat:
		typ = "REAL"
	case sb.ColumnTypeBool:
		typ = "BOOLEAN"
	case sb.ColumnTypeBytes:
		typ = "BYTEA"
	}

	if typ == "" {
		return "", errors.New("dialects: unknown column type")
	} else {
		return typ, nil
	}
}

func (m Postgresql) ColumnOptionToString(co *sb.ColumnOption) (string, error) {
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
