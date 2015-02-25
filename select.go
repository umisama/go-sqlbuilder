package sqlbuilder

// SelectStatement represents a SELECT statement.
type SelectStatement struct {
	columns  selectColumnList
	from     Table
	where    Condition
	distinct bool
	groupBy  []serializable
	orderBy  []serializable
	limit    int
	offset   int
	having   Condition
}

// Select returns new SELECT statement with from as FROM clause.
func Select(from Table) *SelectStatement {
	return &SelectStatement{
		from: from,
	}
}

// Columns set columns for select.
// Get all columns (use *) if it is not setted.
func (b *SelectStatement) Columns(columns ...Column) *SelectStatement {
	b.columns = selectColumnList(columns)
	return b

}

// Where sets WHERE clause.  The cond is filter condition.
func (b *SelectStatement) Where(cond Condition) *SelectStatement {
	b.where = cond
	return b
}

// Distinct sets DISTINCT clause.
func (b *SelectStatement) Distinct() *SelectStatement {
	b.distinct = true
	return b
}

// GroupBy sets "GROUP BY" clause by the columns.
func (b *SelectStatement) GroupBy(columns ...Column) *SelectStatement {
	ex_column := make([]serializable, len(columns))
	for i := range columns {
		ex_column[i] = columns[i]
	}
	b.groupBy = ex_column
	return b
}

// GroupBy sets "HAVING" clause with the cond.
func (b *SelectStatement) Having(cond Condition) *SelectStatement {
	b.having = cond
	return b
}

// OrderBy sets "ORDER BY" clause. Use descending order if the desc is true, by the columns.
func (b *SelectStatement) OrderBy(desc bool, columns ...Column) *SelectStatement {
	if b.orderBy == nil {
		b.orderBy = make([]serializable, 0)
	}

	for _, c := range columns {
		b.orderBy = append(b.orderBy, newOrderBy(desc, c))
	}
	return b
}

// Limit sets LIMIT clause.
func (b *SelectStatement) Limit(limit int) *SelectStatement {
	b.limit = limit
	return b
}

// Offset sets OFFSET clause.
func (b *SelectStatement) Offset(offset int) *SelectStatement {
	b.offset = offset
	return b
}

func (b *SelectStatement) serialize(bldr *builder) {
	// SELECT COLUMN
	bldr.Append("SELECT ")
	if b.distinct {
		bldr.Append("DISTINCT ")
	}

	bldr.AppendItem(b.columns)

	// FROM
	if b.from != nil {
		bldr.Append(" FROM ")
		bldr.AppendItem(b.from)
	} else {
		bldr.SetError(newError("from is nil"))
	}

	// WHERE
	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
	}

	// GROUP BY
	if b.groupBy != nil {
		bldr.Append(" GROUP BY ")
		bldr.AppendItems(b.groupBy, ",")
	}

	// HAVING
	if b.having != nil {
		if b.groupBy == nil {
			bldr.SetError(newError("GROUP BY by clause is not found"))
		}
		bldr.Append(" HAVING ")
		bldr.AppendItem(b.having)
	}

	// ORDER BY
	if b.orderBy != nil {
		bldr.Append(" ORDER BY ")
		bldr.AppendItems(b.orderBy, ", ")
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
	return
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *SelectStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	bldr.AppendItem(b)
	return bldr.Query(), bldr.Args(), bldr.Err()
}

func (m *SelectStatement) ToSubquery(alias string) Table {
	return newSubquery(m, alias)
}

type subquery struct {
	stat  *SelectStatement
	alias string
	err   error
}

func newSubquery(s *SelectStatement, alias string) *subquery {
	m := &subquery{
		stat:  s,
		alias: alias,
	}

	if len(alias) == 0 {
		m.err = newError("alias is empty")
	}
	return m
}

func (m *subquery) serialize(bldr *builder) {
	if m.err != nil {
		bldr.SetError(m.err)
	}

	bldr.Append("( ")
	bldr.AppendItem(m.stat)
	bldr.Append(" ) AS " + m.alias)
	return
}

func (m *subquery) Name() string {
	return m.alias
}

func (m *subquery) C(name string) Column {
	for _, col := range m.stat.columns {
		if ac, ok := col.(aliasedColumn); ok {
			if ac.column_alias() == name {
				return col.config().toColumn(m)
			}
		}
		if col.column_name() == name {
			return col.config().toColumn(m)
		}
	}
	return newErrorColumn(newError("column %s was not found.", name))
}

func (m *subquery) Columns() []Column {
	l := make([]Column, len(m.stat.columns))
	for _, col := range m.stat.columns {
		if _, ok := col.(aliasedColumn); ok {
			l = append(l, col.config().toColumn(m))
		}
		l = append(l, col.config().toColumn(m))
	}
	return nil
}

func (m *subquery) InnerJoin(Table, Condition) Table {
	m.err = newError("Subquery can not join")
	return m
}

func (m *subquery) LeftOuterJoin(Table, Condition) Table {
	m.err = newError("Subquery can not join")
	return m
}

func (m *subquery) RightOuterJoin(Table, Condition) Table {
	m.err = newError("Subquery can not join")
	return m
}

func (m *subquery) FullOuterJoin(Table, Condition) Table {
	m.err = newError("Subquery can not join")
	return m
}

type selectColumnList []Column

func (l selectColumnList) serialize(bldr *builder) {
	if len(l) == 0 {
		bldr.AppendItem(Star)
		return
	}

	first := true
	for _, col := range l {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		if ac, ok := col.(aliasedColumn); ok {
			bldr.AppendItem(ac.source())
			bldr.Append(" AS ")
			bldr.AppendItem(col)
		} else {
			bldr.AppendItem(col)
		}
	}
}

type orderBy struct {
	column Column
	desc   bool
}

func newOrderBy(desc bool, column Column) *orderBy {
	return &orderBy{
		column: column,
		desc:   desc,
	}
}

func (m *orderBy) serialize(bldr *builder) {
	bldr.AppendItem(m.column)
	if m.desc {
		bldr.Append(" DESC")
	} else {
		bldr.Append(" ASC")
	}
}
