package sqlbuilder

type InsertStatement struct {
	columns []serializable
	values  []serializable
	into    Table
}

func Insert(into Table) *InsertStatement {
	return &InsertStatement{
		into: into,
	}
}

func (b *InsertStatement) Columns(columns ...Column) *InsertStatement {
	sl := make([]serializable, len(columns))
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

	// INSERT (COLUMN)
	bldr.Append("INSERT ")
	if len(b.columns) != 0 {
		bldr.Append("( ")
		bldr.AppendItems(b.columns, " ")
		bldr.Append(" )")
	}

	// INTO Table
	bldr.Append(" INTO ")
	bldr.AppendItem(b.into)

	// VALUES
	bldr.Append(" VALUES ( ")
	bldr.AppendItems(b.values, " ")
	bldr.Append(" )")

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}
