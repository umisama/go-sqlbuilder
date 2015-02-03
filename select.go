package sqlbuilder

import (
	"errors"
)

type SelectStatement struct {
	columns    []serializable
	from       Table
	where      Condition
	distinct   bool
	groupBy    []serializable
	orderByAsc string
	orderBy    []serializable
	limit      int
	offset     int
	having     Condition
}

func Select(columns ...Column) *SelectStatement {
	ex_column := make([]serializable, len(columns))
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

func (b *SelectStatement) Distinct() *SelectStatement {
	b.distinct = true
	return b
}

func (b *SelectStatement) GroupBy(columns ...Column) *SelectStatement {
	ex_column := make([]serializable, len(columns))
	for i := range columns {
		ex_column[i] = columns[i]
	}
	b.groupBy = ex_column
	return b
}

func (b *SelectStatement) Having(cond Condition) *SelectStatement {
	b.having = cond
	return b
}

func (b *SelectStatement) OrderBy(desc bool, columns ...Column) *SelectStatement {
	ex_column := make([]serializable, len(columns))
	for i := range columns {
		ex_column[i] = columns[i]
	}
	b.orderBy = ex_column

	if desc {
		b.orderByAsc = " DESC"
	} else {
		b.orderByAsc = " ASC"
	}
	return b
}

func (b *SelectStatement) Limit(limit int) *SelectStatement {
	b.limit = limit
	return b
}

func (b *SelectStatement) Offset(offset int) *SelectStatement {
	b.offset = offset
	return b
}

func (b *SelectStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	// SELECT COLUMN
	bldr.Append("SELECT ")
	if b.distinct {
		bldr.Append("DISTINCT ")
	}

	bldr.AppendItems(b.columns, ", ")

	// FROM
	if b.from != nil {
		bldr.Append(" FROM ")
		bldr.AppendItem(b.from)
	} else {
		bldr.SetError(errors.New("from is not found"))
	}

	// WHERE
	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
	}

	// GROUP BY
	if b.groupBy != nil {
		bldr.Append(" GROUP BY ")
		bldr.AppendItems(b.groupBy, ",")
	}

	// HAVING
	if b.having != nil {
		if b.groupBy == nil {
			bldr.SetError(errors.New("GROUP BY by clause is not found"))
		}
		bldr.Append(" HAVING ")
		bldr.AppendItem(b.having)
	}

	// ORDER BY
	if b.orderBy != nil {
		bldr.Append(" ORDER BY ")
		bldr.AppendItems(b.orderBy, ",")
		bldr.Append(b.orderByAsc)
	}

	// LIMIT
	if b.limit != 0 {
		bldr.Append(" LIMIT ")
		bldr.AppendValue(b.limit)
	}

	// Offset
	if b.offset != 0 {
		bldr.Append(" OFFSET ")
		bldr.AppendValue(b.offset)
	}

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}
