package sqlbuilder

import "strconv"

type UpdateStatement struct {
	table      Table
	set        []serializable
	where      Condition
	orderByAsc string
	orderBy    []serializable
	limit      int
	offset     int
}

func Update(table Table) *UpdateStatement {
	return &UpdateStatement{
		table: table,
		set:   make([]serializable, 0),
	}
}

func (b *UpdateStatement) Set(col Column, val interface{}) *UpdateStatement {
	b.set = append(b.set, newUpdateValue(col, val))
	return b
}

func (b *UpdateStatement) Where(cond Condition) *UpdateStatement {
	b.where = cond
	return b
}

func (b *UpdateStatement) Limit(limit int) *UpdateStatement {
	b.limit = limit
	return b
}

func (b *UpdateStatement) Offset(offset int) *UpdateStatement {
	b.offset = offset
	return b
}

func (b *UpdateStatement) OrderBy(desc bool, columns ...Column) *UpdateStatement {
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

func (b *UpdateStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	// UPDATE TABLE SET (COLUMN=VALUE)
	bldr.Append("UPDATE ")
	bldr.AppendItem(b.table)
	bldr.Append(" SET ")
	bldr.AppendItems(b.set, ", ")

	// WHERE
	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
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
		bldr.Append(strconv.Itoa(b.limit))
	}

	// Offset
	if b.offset != 0 {
		bldr.Append(" OFFSET ")
		bldr.Append(strconv.Itoa(b.offset))
	}

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}

type updateValue struct {
	col Column
	val literal
}

func newUpdateValue(col Column, val interface{}) updateValue {
	return updateValue{
		col: col,
		val: toLiteral(val),
	}
}

func (m updateValue) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.col.column_name()))
	bldr.Append("=")
	bldr.AppendItem(m.val)
}
