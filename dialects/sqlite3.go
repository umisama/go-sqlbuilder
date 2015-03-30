package dialects

import (
	"errors"
	sb "github.com/umisama/go-sqlbuilder"
)

type Sqlite struct{}

func (m Sqlite) QuerySuffix() string {
	return ";"
}

func (m Sqlite) BindVar(i int) string {
	return "?"
}

func (m Sqlite) QuoteField(field string) string {
	return "\"" + field + "\""
}

func (m Sqlite) ColumnTypeToString(cc sb.ColumnConfig) (string, error) {
	if cc.Option().SqlType != "" {
		return cc.Option().SqlType, nil
	}

	typ := ""
	switch cc.Type() {
	case sb.ColumnTypeInt:
		typ = "INTEGER"
	case sb.ColumnTypeString:
		typ = "TEXT"
	case sb.ColumnTypeDate:
		typ = "DATE"
	case sb.ColumnTypeFloat:
		typ = "REAL"
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

func (m Sqlite) ColumnOptionToString(co *sb.ColumnOption) (string, []interface{}, error) {
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

func (m Sqlite) TableOptionToString(to *sb.TableOption) (string, []interface{}, error) {
	opt := ""
	args := make([]interface{}, 0)
	if to.Unique != nil {
		opt = str_append(opt, m.tableOptionUnique(to.Unique))
	}

	return "", args, nil
}

func (m Sqlite) tableOptionUnique(op [][]string) string {
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
