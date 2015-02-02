package sqlbuilder

type UpdateStatement struct {
	table Table
	set   []serializable
	where Condition
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

func (b *UpdateStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	bldr.Append("UPDATE ")
	bldr.AppendItem(b.table)

	bldr.Append(" SET ")
	bldr.AppendItems(b.set, ", ")

	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
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
