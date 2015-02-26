package sqlbuilder

type AlterTableStatement struct {
	table          Table
	rename_to      string
	add_columns    []*alterTableAddColumn
	drop_columns   []Column
	change_columns []*alterTableChangeColumn
}

func AlterTable(table Table) *AlterTableStatement {
	return &AlterTableStatement{
		table:          table,
		add_columns:    make([]*alterTableAddColumn, 0),
		change_columns: make([]*alterTableChangeColumn, 0),
	}
}

func (b *AlterTableStatement) RenameTo(name string) *AlterTableStatement {
	b.rename_to = name
	return b
}

func (b *AlterTableStatement) AddColumn(col ColumnConfig) *AlterTableStatement {
	b.add_columns = append(b.add_columns, &alterTableAddColumn{
		column: col,
		first:  false,
		after:  nil,
	})
	return b
}

func (b *AlterTableStatement) AddColumnAfter(col ColumnConfig, after Column) *AlterTableStatement {
	b.add_columns = append(b.add_columns, &alterTableAddColumn{
		column: col,
		first:  false,
		after:  after,
	})
	return b
}

func (b *AlterTableStatement) AddColumnFirst(col ColumnConfig) *AlterTableStatement {
	b.add_columns = append(b.add_columns, &alterTableAddColumn{
		column: col,
		first:  true,
		after:  nil,
	})
	return b
}

func (b *AlterTableStatement) DropColumn(col Column) *AlterTableStatement {
	b.drop_columns = append(b.drop_columns, col)
	return b
}

func (b *AlterTableStatement) ChangeColumn(old_column Column, new_column ColumnConfig) *AlterTableStatement {
	b.change_columns = append(b.change_columns, &alterTableChangeColumn{
		old_column: old_column,
		new_column: new_column,
		first:      false,
		after:      nil,
	})
	return b
}

func (b *AlterTableStatement) ChangeColumnAfter(old_column Column, new_column ColumnConfig, after Column) *AlterTableStatement {
	b.change_columns = append(b.change_columns, &alterTableChangeColumn{
		old_column: old_column,
		new_column: new_column,
		first:      false,
		after:      after,
	})
	return b
}

func (b *AlterTableStatement) ChangeColumnFirst(old_column Column, new_column ColumnConfig) *AlterTableStatement {
	b.change_columns = append(b.change_columns, &alterTableChangeColumn{
		old_column: old_column,
		new_column: new_column,
		first:      true,
		after:      nil,
	})
	return b
}

func (b *AlterTableStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()

	bldr.Append("ALTER TABLE ")
	bldr.AppendItem(b.table)
	bldr.Append(" ")

	first := true
	for _, add_column := range b.add_columns {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		bldr.AppendItem(add_column)
	}
	for _, change_column := range b.change_columns {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		bldr.AppendItem(change_column)
	}
	for _, drop_column := range b.drop_columns {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		bldr.Append("DROP COLUMN ")
		if colname := drop_column.column_name(); len(colname) != 0 {
			bldr.Append(dialect().QuoteField(colname))
		} else {
			bldr.AppendItem(drop_column)
		}
	}
	if len(b.rename_to) != 0 {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		bldr.Append("RENAME TO ")
		bldr.Append(dialect().QuoteField(b.rename_to))
	}

	return "", nil, nil
}

type alterTableAddColumn struct {
	column ColumnConfig
	first  bool
	after  Column
}

func (b *alterTableAddColumn) serialize(bldr *builder) {
	bldr.Append("ADD COLUMN ")
	bldr.AppendItem(b.column)

	// SQL data name
	typ, err := dialect().ColumnTypeToString(b.column)
	if err != nil {
		bldr.SetError(err)
	} else if len(typ) == 0 {
		bldr.SetError(newError("Column type is required.(maybe, a bug is in implements of dialect.)"))
	} else {
		bldr.Append(" ")
		bldr.Append(typ)
	}

	opt, err := dialect().ColumnOptionToString(b.column.Option())
	if err != nil {
		bldr.SetError(err)
	} else if len(opt) != 0 {
		bldr.Append(" ")
		bldr.Append(opt)
	}

	if b.first {
		bldr.Append(" FIRST")
	} else if b.after != nil {
		bldr.Append(" AFTER ")
		if colname := b.after.column_name(); len(colname) != 0 {
			bldr.Append(dialect().QuoteField(colname))
		} else {
			bldr.AppendItem(b.after)
		}
	}
}

type alterTableChangeColumn struct {
	old_column Column
	new_column ColumnConfig
	first      bool
	after      Column
}

func (b *alterTableChangeColumn) serialize(bldr *builder) {
	bldr.Append("CHANGE COLUMN ")
	if colname := b.old_column.column_name(); len(colname) != 0 {
		bldr.Append(dialect().QuoteField(colname))
	} else {
		bldr.AppendItem(b.old_column)
	}
	bldr.Append(" ")
	bldr.AppendItem(b.new_column)

	typ, err := dialect().ColumnTypeToString(b.new_column)
	if err != nil {
		bldr.SetError(err)
	} else if len(typ) == 0 {
		bldr.SetError(newError("Column type is required.(maybe, a bug is in implements of dialect.)"))
	} else {
		bldr.Append(" ")
		bldr.Append(typ)
	}

	opt, err := dialect().ColumnOptionToString(b.new_column.Option())
	if err != nil {
		bldr.SetError(err)
	} else if len(opt) != 0 {
		bldr.Append(" ")
		bldr.Append(opt)
	}

	if b.first {
		bldr.Append(" FIRST")
	} else if b.after != nil {
		bldr.Append(" AFTER ")
		if colname := b.after.column_name(); len(colname) != 0 {
			bldr.Append(dialect().QuoteField(colname))
		} else {
			bldr.AppendItem(b.after)
		}
	}
}
