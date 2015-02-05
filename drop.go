package sqlbuilder

// DeleteTableStatement represents a "DROP TABLE" statement.
type DropTableStatement struct {
	table Table
}

// DropTable returns new "DROP TABLE" statement. The table is Table object to drop.
func DropTable(table Table) *DropTableStatement {
	return &DropTableStatement{
		table: table,
	}
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *DropTableStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	defer func() {
		bldr.Append(dialect.QuerySuffix())
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()

	bldr.Append("DROP TABLE ")
	if b.table != nil {
		bldr.AppendItem(b.table)
	} else {
		bldr.SetError(newError("table is nil"))
		return
	}

	return
}
