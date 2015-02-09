package sqlbuilder

type joinType int

const (
	inner_join joinType = iota
	left_outer_join
	right_outer_join
	full_outer_join
)

type table struct {
	name    string
	columns []Column
}

type joinTable struct {
	typ   joinType
	left  Table
	right Table
	on    Condition
}

// Table represents a table.
type Table interface {
	serializable

	// C returns table's column by the name.
	C(name string) Column

	// Name returns table' name.
	Name() string

	// Columns returns all columns.
	Columns() []Column

	// InnerJoin returns a joined table use with "INNER JOIN" clause.
	// The joined table can be handled in same way as single table.
	InnerJoin(Table, Condition) Table

	// LeftOuterJoin returns a joined table use with "LEFT OUTER JOIN" clause.
	// The joined table can be handled in same way as single table.
	LeftOuterJoin(Table, Condition) Table

	// RightOuterJoin returns a joined table use with "RIGHT OUTER JOIN" clause.
	// The joined table can be handled in same way as single table.
	RightOuterJoin(Table, Condition) Table

	// FullOuterJoin returns a joined table use with "FULL OUTER JOIN" clause.
	// The joined table can be handled in same way as single table.
	FullOuterJoin(Table, Condition) Table
}

// NewTable returns a new table named by the name.  Specify table columns by the column_config.
// Panic if column is empty.
func NewTable(name string, column_configs ...ColumnConfig) Table {
	if len(column_configs) == 0 {
		panic(newError("column is needed"))
	}

	t := &table{
		name: name,
	}

	columns := make([]Column, 0, len(column_configs))
	for _, column_config := range column_configs {
		columns = append(columns, column_config.toColumn(t))
	}
	t.columns = columns

	return t
}

func (m *table) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.name))
	return
}

func (m *table) C(name string) Column {
	for _, column := range m.columns {
		if column.column_name() == name {
			return column
		}
	}

	return newErrorColumn(newError("column %s.%s was not found.", m.name, name))
}

func (m *table) Name() string {
	return m.name
}

func (m *table) Columns() []Column {
	return m.columns
}

func (m *table) InnerJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   inner_join,
		on:    on,
	}
}

func (m *table) LeftOuterJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   left_outer_join,
		on:    on,
	}
}

func (m *table) RightOuterJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   right_outer_join,
		on:    on,
	}
}

func (m *table) FullOuterJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   full_outer_join,
		on:    on,
	}
}

func (m *joinTable) C(name string) Column {
	l_col := m.left.C(name)
	r_col := m.right.C(name)

	_, l_err := l_col.(*errorColumn)
	_, r_err := r_col.(*errorColumn)

	switch {
	case l_err && r_err:
		return newErrorColumn(newError("column %s was not found.", name))
	case l_err && !r_err:
		return r_col
	case !l_err && r_err:
		return l_col
	default:
		return newErrorColumn(newError("column %s was duplicated.", name))
	}
}

func (m *joinTable) Name() string {
	return ""
}

func (m *joinTable) Columns() []Column {
	return append(m.left.Columns(), m.right.Columns()...)
}

func (m *joinTable) InnerJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   inner_join,
		on:    on,
	}
}

func (m *joinTable) LeftOuterJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   left_outer_join,
		on:    on,
	}
}

func (m *joinTable) RightOuterJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   right_outer_join,
		on:    on,
	}
}

func (m *joinTable) FullOuterJoin(right Table, on Condition) Table {
	return &joinTable{
		left:  m,
		right: right,
		typ:   full_outer_join,
		on:    on,
	}
}

func (m *joinTable) serialize(bldr *builder) {
	bldr.AppendItem(m.left)
	switch m.typ {
	case inner_join:
		bldr.Append(" INNER JOIN ")
	case left_outer_join:
		bldr.Append(" LEFT OUTER JOIN ")
	case right_outer_join:
		bldr.Append(" RIGHT OUTER JOIN ")
	case full_outer_join:
		bldr.Append(" FULL OUTER JOIN ")
	}
	bldr.AppendItem(m.right)
	bldr.Append(" ON ")
	bldr.AppendItem(m.on)
	return
}
