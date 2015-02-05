package sqlbuilder

// InsertStatement represents a INSERT statement.
type InsertStatement struct {
	columns ColumnList
	values  []serializable
	into    Table
}

// Insert returns new INSERT statement. The table is Table object for into.
func Insert(into Table) *InsertStatement {
	return &InsertStatement{
		into: into,
	}
}

// Columns sets columns for insert.
// If not set this, get error on ToSql().
func (b *InsertStatement) Columns(columns ...Column) *InsertStatement {
	b.columns = ColumnList(columns)
	return b
}

// Values sets VALUES clause.  The values is data for insert.
func (b *InsertStatement) Values(values ...interface{}) *InsertStatement {
	sl := make([]serializable, len(values))
	for i := range values {
		sl[i] = toLiteral(values[i])
	}
	b.values = sl
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *InsertStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	defer func() {
		bldr.Append(dialect.QuerySuffix())
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()

	// INSERT
	bldr.Append("INSERT")

	// INTO Table
	bldr.Append(" INTO ")
	if b.into != nil {
		bldr.AppendItem(b.into)
	} else {
		bldr.SetError(newError("into is nil"))
		return
	}

	// (COLUMN)
	if len(b.columns) == 0 {
		b.columns = b.into.Columns()
	}
	bldr.Append(" ( ")
	bldr.AppendItem(b.columns)
	bldr.Append(" )")

	// VALUES
	if len(b.columns) != len(b.values) {
		bldr.SetError(newError("%d values needed, but got %d", len(b.columns), len(b.values)))
		return
	}
	bldr.Append(" VALUES ( ")
	bldr.AppendItems(b.values, ", ")
	bldr.Append(" )")

	return
}
