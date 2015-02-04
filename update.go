package sqlbuilder

// UpdateStatement represents a UPDATE statement.
type UpdateStatement struct {
	table      Table
	set        []serializable
	where      Condition
	orderByAsc string
	orderBy    []serializable
	limit      int
	offset     int
}

// Update returns new UPDATE statement. The table is Table object to update.
func Update(table Table) *UpdateStatement {
	return &UpdateStatement{
		table: table,
		set:   make([]serializable, 0),
	}
}

// Set sets SETS clause like col=val.  Call many time for update multi columns.
func (b *UpdateStatement) Set(col Column, val interface{}) *UpdateStatement {
	b.set = append(b.set, newUpdateValue(col, val))
	return b
}

// Where sets WHERE clause.  The cond is filter condition.
func (b *UpdateStatement) Where(cond Condition) *UpdateStatement {
	b.where = cond
	return b
}

// Limit sets LIMIT clause.
func (b *UpdateStatement) Limit(limit int) *UpdateStatement {
	b.limit = limit
	return b
}

// Limit sets OFFSET clause.
func (b *UpdateStatement) Offset(offset int) *UpdateStatement {
	b.offset = offset
	return b
}

// OrderBy sets "ORDER BY" clause. Use descending order if the desc is true, by the columns.
func (b *UpdateStatement) OrderBy(desc bool, columns ...Column) *UpdateStatement {
	ex_column := make([]serializable, len(columns))
	for i := range columns {
		ex_column[i] = columns[i]
	}
	b.orderBy = ex_column

	if desc {
		b.orderByAsc = " DESC"
	} else {
		b.orderByAsc = " ASC"
	}
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *UpdateStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()

	// UPDATE TABLE SET (COLUMN=VALUE)
	bldr.Append("UPDATE ")
	bldr.AppendItem(b.table)
	bldr.Append(" SET ")
	bldr.AppendItems(b.set, ", ")

	// WHERE
	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
	}

	// ORDER BY
	if b.orderBy != nil {
		bldr.Append(" ORDER BY ")
		bldr.AppendItems(b.orderBy, ",")
		bldr.Append(b.orderByAsc)
	}

	// LIMIT
	if b.limit != 0 {
		bldr.Append(" LIMIT ")
		bldr.AppendValue(b.limit)
	}

	// Offset
	if b.offset != 0 {
		bldr.Append(" OFFSET ")
		bldr.AppendValue(b.offset)
	}

	bldr.Append(dialect.QuerySuffix())
	return bldr.Query(), bldr.Args(), bldr.Err()
}

type updateValue struct {
	col Column
	val literal
}

func newUpdateValue(col Column, val interface{}) updateValue {
	return updateValue{
		col: col,
		val: toLiteral(val),
	}
}

func (m updateValue) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.col.column_name()))
	bldr.Append("=")
	bldr.AppendItem(m.val)
}
