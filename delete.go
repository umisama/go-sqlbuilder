package sqlbuilder

type DeleteStatement struct {
	from  Table
	where Condition
}

func Delete(from Table) *DeleteStatement {
	return &DeleteStatement{
		from: from,
	}
}

func (b *DeleteStatement) Where(cond Condition) *DeleteStatement {
	b.where = cond
	return b
}

func (b *DeleteStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	bldr.Append("DELETE FROM ")
	bldr.AppendItem(b.from)
	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
	}

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}
