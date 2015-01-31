package sqlbuilder

import (
	"errors"
)

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

type Table interface {
	serializable

	C(name string) Column
	Name() string
	Columns() []Column

	InnerJoin(Table, Condition) Table
	LeftOuterJoin(Table, Condition) Table
	RightOuterJoin(Table, Condition) Table
	FullOuterJoin(Table, Condition) Table
}

func NewTable(name string, columns ...Column) (Table, error) {
	if len(columns) == 0 {
		return nil, errors.New("column is needed")
	}

	t := &table{
		name:    name,
		columns: columns,
	}

	for i := range t.columns {
		t.columns[i].setTableName(t)
	}

	return t, nil
}

func (m *table) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.name))
	return
}

func (m *table) C(name string) Column {
	for _, column := range m.columns {
		if column.Name() == name {
			return column
		}
	}

	return nil
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
	col := m.left.C(name)
	if col != nil {
		return col
	}
	return m.right.C(name)
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
