package sqlbuilder

import (
	errs "errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetDialect(TestDialect{})
	os.Exit(m.Run())
}

func TestError(t *testing.T) {
	a := assert.New(t)
	err := newError("hogehogestring")
	a.Equal("sqlbuilder: hogehogestring", err.Error())
}

type TestDialect struct{}

func (m TestDialect) QuerySuffix() string {
	return ";"
}

func (m TestDialect) BindVar(i int) string {
	return "?"
}

func (m TestDialect) QuoteField(field string) string {
	return "\"" + field + "\""
}

func (m TestDialect) ColumnTypeToString(cc ColumnConfig) (string, error) {
	if cc.Option().SqlType != "" {
		return cc.Option().SqlType, nil
	}

	typ := ""
	switch cc.Type() {
	case ColumnTypeInt:
		typ = "INTEGER"
	case ColumnTypeString:
		typ = "TEXT"
	case ColumnTypeDate:
		typ = "DATE"
	case ColumnTypeFloat:
		typ = "REAL"
	case ColumnTypeBool:
		typ = "BOOLEAN"
	case ColumnTypeBytes:
		typ = "BLOB"
	}
	if typ == "" {
		return "", errs.New("dialects: unknown column type")
	} else {
		return typ, nil
	}
}

func (m TestDialect) ColumnOptionToString(co *ColumnOption) (string, []interface{}, error) {
	apnd := func(str, opt string) string {
		if len(str) != 0 {
			str += " "
		}
		str += opt
		return str
	}

	opt := ""
	args := make([]interface{}, 0)
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
	if co.Default == nil {
		opt = apnd(opt, "DEFAULT NULL")
	} else {
		opt = apnd(opt, "DEFAULT ?")
		args = append(args, co.Default)
	}

	return opt, args, nil
}

func (m TestDialect) TableOptionToString(to *TableOption) (string, []interface{}, error) {
	opt := ""
	args := make([]interface{}, 0)
	apnd := func(str, opt string) string {
		if len(str) != 0 {
			str += " "
		}
		str += opt
		return str
	}

	if to.Unique != nil {
		opt = apnd(opt, m.tableOptionUnique(to.Unique))
	}
	return opt, args, nil
}

func (m TestDialect) tableOptionUnique(op [][]string) string {
	opt := ""
	first_op := true
	for _, unique := range op {
		if first_op {
			first_op = false
		} else {
			opt += " "
		}

		opt += "UNIQUE("
		first := true
		for _, col := range unique {
			if first {
				first = false
			} else {
				opt += ", "
			}
			opt += m.QuoteField(col)
		}
		opt += ")"
	}
	return opt
}
