package sqlbuilder

import (
	"errors"
)

type Table struct {
	name    string
	columns []Column
}

func NewTable(name string, columns ...Column) (*Table, error) {
	if len(columns) == 0 {
		return nil, errors.New("column is needed")
	}

	t := &Table{
		name:    name,
		columns: columns,
	}

	for i := range t.columns {
		t.columns[i].setTableName(t)
	}

	return t, nil
}

func (m *Table) toSql() (string, []interface{}, error) {
	return dialect.QuoteField(m.name), []interface{}{}, nil
}

func (m *Table) C(name string) Column {
	for _, column := range m.columns {
		if column.Name() == name {
			return column
		}
	}

	return nil
}

func (m *Table) Name() string {
	return m.name
}

func (m *Table) Columns() []Column {
	return m.columns
}
