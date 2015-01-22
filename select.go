package sqlbuilder

import (
	"errors"
)

type SelectStatement struct {
	columns []Expression
	from    *Table
	where   Condition
	err     error
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

func (b *SelectStatement) From(table *Table) *SelectStatement {
	if b.err != nil {
		return b
	}

	b.from = table
	return b
}

func (b *SelectStatement) Where(cond Condition) *SelectStatement {
	if b.err != nil {
		return b
	}

	b.where = cond
	return b
}

func (b *SelectStatement) Error() error {
	return b.err
}

func (b *SelectStatement) ToSql() (query string, attrs []interface{}, err error) {
	if b.err != nil {
		return "", []interface{}{}, b.err
	}

	query, attrs, err = "", []interface{}{}, nil
	defer func() {
		query += dialect.QuerySuffix()
	}()

	// SELECT COLUMN
	query += "SELECT "
	query, attrs, err = appendExpressionsToQuery(b.columns, query, attrs, " ")
	if err != nil {
		return "", []interface{}{}, err
	}

	// FROM
	if b.from != nil {
		query += " FROM "
		query, attrs, err = appendItemToQuery(b.from, query, attrs)
		if err != nil {
			return "", []interface{}{}, err
		}
	} else {
		return "", []interface{}{}, errors.New("from is not found")
	}

	// WHERE
	if b.where != nil {
		query += " WHERE "
		query, attrs, err = appendItemToQuery(b.where, query, attrs)
		if err != nil {
			return "", []interface{}{}, err
		}
	}

	return
}
