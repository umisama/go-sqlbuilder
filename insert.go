package sqlbuilder

type InsertStatement struct {
	columns ColumnList
	values  []serializable
	into    Table
}

func Insert(into Table) *InsertStatement {
	return &InsertStatement{
		into: into,
	}
}

func (b *InsertStatement) Columns(columns ...Column) *InsertStatement {
	b.columns = ColumnList(columns)
	return b
}

func (b *InsertStatement) Values(values ...interface{}) *InsertStatement {
	sl := make([]serializable, len(values))
	for i := range values {
		sl[i] = toLiteral(values[i])
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
		bldr.AppendItem(b.columns)
		bldr.Append(" )")
	}

	// VALUES
	bldr.Append(" VALUES ( ")
	bldr.AppendItems(b.values, ", ")
	bldr.Append(" )")

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}
