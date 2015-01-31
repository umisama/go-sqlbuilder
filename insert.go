package sqlbuilder

type InsertStatement struct {
	columns []serializableForColumnList
	values  []serializable
	into    Table
}

func Insert(into Table) *InsertStatement {
	return &InsertStatement{
		into: into,
	}
}

func (b *InsertStatement) Columns(columns ...Column) *InsertStatement {
	sl := make([]serializableForColumnList, len(columns))
	for i := range columns {
		sl[i] = columns[i]
	}
	b.columns = sl
	return b
}

func (b *InsertStatement) Values(values ...Literal) *InsertStatement {
	sl := make([]serializable, len(values))
	for i := range values {
		sl[i] = values[i]
	}
	b.values = sl
	return b
}

func (b *InsertStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	// INSERT
	bldr.Append("INSERT")

	// INTO Table(COLUMN)
	bldr.Append(" INTO ")
	bldr.AppendItem(b.into)
	if len(b.columns) != 0 {
		bldr.Append(" ( ")
		bldr.AppendItemsForColumnList(b.columns, ", ")
		bldr.Append(" )")
	}

	// VALUES
	bldr.Append(" VALUES ( ")
	bldr.AppendItems(b.values, ", ")
	bldr.Append(" )")

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}
