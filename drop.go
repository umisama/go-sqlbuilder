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

	bldr.Append("DROP TABLE ")
	bldr.AppendItem(b.table)

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}
