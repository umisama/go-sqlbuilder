package sqlbuilder

import (
	"errors"
)

type SelectStatement struct {
	columns []Expression
	from    Table
	where   Condition
}

func Select(columns ...Column) *SelectStatement {
	ex_column := make([]Expression, len(columns))
	for i := range columns {
		ex_column[i] = columns[i]
	}

	return &SelectStatement{
		columns: ex_column,
	}
}

func (b *SelectStatement) From(table Table) *SelectStatement {
	b.from = table
	return b
}

func (b *SelectStatement) Where(cond Condition) *SelectStatement {
	b.where = cond
	return b
}

func (b *SelectStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	// SELECT COLUMN
	bldr.Append("SELECT ", nil)
	bldr.AppendExpressions(b.columns, " ")

	// FROM
	if b.from != nil {
		bldr.Append(" FROM ", nil)
		bldr.AppendItem(b.from)
	} else {
		bldr.SetError(errors.New("from is not found"))
	}

	// WHERE
	if b.where != nil {
		bldr.Append(" WHERE ", nil)
		bldr.AppendItem(b.where)
	}

	bldr.Append(dialect.QuerySuffix(), nil)
	return bldr.Query(), bldr.Args(), bldr.Err()
}
