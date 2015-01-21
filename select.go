package sqlbuilder

import (
	"errors"
)

type SelectBuilder struct {
	columns []sqlizable
	from    *Table
	where   Condition
	err     error
}

func Select(columns ...Column) *SelectBuilder {
	sqlizable_column := make([]sqlizable, len(columns))
	for i := range columns {
		sqlizable_column[i] = columns[i]
	}

	return &SelectBuilder{
		columns: sqlizable_column,
	}
}

func (b *SelectBuilder) From(table *Table) *SelectBuilder {
	if b.err != nil {
		return b
	}

	b.from = table
	return b
}

func (b *SelectBuilder) Where(cond Condition) *SelectBuilder {
	if b.err != nil {
		return b
	}

	b.where = cond
	return b
}

func (b *SelectBuilder) Error() error {
	return b.err
}

func (b *SelectBuilder) ToSql() (query string, attrs []interface{}, err error) {
	if b.err != nil {
		return "", []interface{}{}, b.err
	}

	query, attrs, err = "", []interface{}{}, nil
	defer func() {
		query += dialect.QuerySuffix()
	}()

	// SELECT COLUMN
	query += "SELECT "
	query, attrs, err = appendListToQuery(b.columns, query, attrs, " ")
	if err != nil {
		return "", []interface{}{}, err
	}

	// FROM
	if b.from != nil {
		query += " FROM "
		query, attrs, err = appendToQuery(b.from, query, attrs)
		if err != nil {
			return "", []interface{}{}, err
		}
	} else {
		return "", []interface{}{}, errors.New("from is not found")
	}

	// WHERE
	if b.where != nil {
		query += " WHERE "
		q, a, err := b.where.toSql()
		if err != nil {
			return "", []interface{}{}, err
		}
		query += q
		attrs = append(attrs, a...)
	}

	return query, attrs, nil
}
