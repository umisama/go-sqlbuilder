package sqlbuilder

type CreateIndexStatement struct {
	table       Table
	columns     []Column
	name        string
	ifNotExists bool
}

type CreateTableStatement struct {
	table       Table
	ifNotExists bool
}

func CreateTable(table Table) *CreateTableStatement {
	return &CreateTableStatement{
		table: table,
	}
}

func (b *CreateTableStatement) IfNotExists() *CreateTableStatement {
	b.ifNotExists = true
	return b
}

func CreateIndex(table Table) *CreateIndexStatement {
	return &CreateIndexStatement{
		table: table,
	}
}

func (b *CreateTableStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	bldr.Append("CREATE TABLE ")
	if b.ifNotExists {
		bldr.Append("IF NOT EXISTS ")
	}
	bldr.AppendItem(b.table)

	bldr.Append(" ( ")
	bldr.AppendItem(createTableColumnList(b.table.Columns()))
	bldr.Append(" )")

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}

func (b *CreateIndexStatement) IfNotExists() *CreateIndexStatement {
	b.ifNotExists = true
	return b
}

func (b *CreateIndexStatement) Columns(columns ...Column) *CreateIndexStatement {
	b.columns = columns
	return b
}

func (b *CreateIndexStatement) Name(name string) *CreateIndexStatement {
	b.name = name
	return b
}

func (b *CreateIndexStatement) ToSql() (string, []interface{}, error) {
	bldr := newBuilder()

	bldr.Append("CREATE INDEX ")
	if b.ifNotExists {
		bldr.Append("IF NOT EXISTS ")
	}

	bldr.Append(dialect.QuoteField(b.name))
	bldr.Append(" ON ")

	bldr.AppendItem(b.table)
	bldr.Append(" ( ")
	bldr.AppendItem(createIndexColumnList(b.columns))
	bldr.Append(" )")

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}

type createTableColumnList []Column

func (m createTableColumnList) serialize(bldr *builder) {
	first := true
	for _, column := range m {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		cc := column.config()
		bldr.AppendItem(cc)
		bldr.Append(" ")
		str, err := dialect.SqlType(cc)
		if err != nil {
			bldr.SetError(err)
		}
		bldr.Append(str)
	}
}

type createIndexColumnList []Column

func (m createIndexColumnList) serialize(bldr *builder) {
	first := true
	for _, column := range m {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		cc := column.config()
		bldr.AppendItem(cc)
	}
}
