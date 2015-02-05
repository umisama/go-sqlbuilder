package sqlbuilder

// DeleteStatement represents a DELETE statement.
type DeleteStatement struct {
	from  Table
	where Condition
}

// Delete returns new DELETE statement. The table is Table object to delete from.
func Delete(from Table) *DeleteStatement {
	return &DeleteStatement{
		from: from,
	}
}

// Where sets WHERE clause. cond is filter condition.
func (b *DeleteStatement) Where(cond Condition) *DeleteStatement {
	b.where = cond
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *DeleteStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	defer func() {
		bldr.Append(dialect.QuerySuffix())
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()

	bldr.Append("DELETE FROM ")
	if b.from != nil {
		bldr.AppendItem(b.from)
	} else {
		bldr.SetError(newError("table is nil"))
		return
	}

	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
	}
	return
}
