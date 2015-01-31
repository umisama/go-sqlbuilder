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
	bldr.Append("INSERT ", nil)
	if len(b.columns) != 0 {
		bldr.Append("( ", nil)
		bldr.AppendItems(b.columns, " ")
		bldr.Append(" )", nil)
	}

	// INTO Table
	bldr.Append(" INTO ", nil)
	bldr.AppendItem(b.into)

	// VALUES
	bldr.Append(" VALUES ( ", nil)
	bldr.AppendItems(b.values, " ")
	bldr.Append(" )", nil)

	bldr.Append(dialect.QuerySuffix(), nil)
	return bldr.Query(), bldr.Args(), bldr.Err()
}
