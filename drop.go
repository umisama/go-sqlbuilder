package sqlbuilder

type DropTableStatement struct {
	table Table
}

func DropTable(table Table) *DropTableStatement {
	return &DropTableStatement{
		table: table,
	}
}

func (b *DropTableStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	bldr.Append("DROP TABLE ")
	bldr.AppendItem(b.table)

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}
