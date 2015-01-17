package sqlbuilder

type Column interface {
	Name() string
	NotNull() bool
	TableName() *Table
	SetTableName(*Table)
}

func IntColumn(name string, notnull bool) Column {
	return &intColumn{
		baseColumn{
			name:    name,
			notnull: notnull,
		},
	}
}

func StringColumn(name string, notnull bool) Column {
	return &stringColumn{
		baseColumn{
			name:    name,
			notnull: notnull,
		},
	}
}

func DateColumn(name string, notnull bool) Column {
	return &dateColumn{
		baseColumn{
			name:    name,
			notnull: notnull,
		},
	}
}

func FloatColumn(name string, notnull bool) Column {
	return &floatColumn{
		baseColumn{
			name:    name,
			notnull: notnull,
		},
	}
}

func BoolColumn(name string, notnull bool) Column {
	return &boolColumn{
		baseColumn{
			name:    name,
			notnull: notnull,
		},
	}
}

func BytesColumn(name string, notnull bool) Column {
	return &bytesColumn{
		baseColumn{
			name:    name,
			notnull: notnull,
		},
	}
}

type baseColumn struct {
	name    string
	notnull bool
	table   *Table
}

func (m *baseColumn) Name() string {
	return m.name
}

func (m *baseColumn) NotNull() bool {
	return m.notnull
}

func (m *baseColumn) TableName() *Table {
	return m.table
}

func (m *baseColumn) SetTableName(name *Table) {
	m.table = name
}

type intColumn struct {
	baseColumn
}

type stringColumn struct {
	baseColumn
}

type dateColumn struct {
	baseColumn
}

type floatColumn struct {
	baseColumn
}

type boolColumn struct {
	baseColumn
}

type bytesColumn struct {
	baseColumn
}
