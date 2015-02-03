package sqlbuilder

type ColumnConfig interface {
	serializable

	toColumn(Table) Column
	Name() string
	NotNull() bool
	Type() columnType
	Options() interface{}
}

type columnType int

const (
	columnTypeInt columnType = iota
	columnTypeString
	columnTypeDate
	columnTypeFloat
	columnTypeBool
	columnTypeBytes
)

type ColumnList []Column

type Column interface {
	serializable

	column_name() string
	not_null() bool
	config() ColumnConfig

	Eq(right interface{}) Condition
	NotEq(right interface{}) Condition
	Gt(right interface{}) Condition
	GtEq(right interface{}) Condition
	Lt(right interface{}) Condition
	LtEq(right interface{}) Condition
	Like(right string) Condition
	Between(lower, higher interface{}) Condition
}

type columnConfigImpl struct {
	name    string
	notnull bool
	typ     columnType
	opt     interface{}
}

func (c *columnConfigImpl) Name() string {
	return c.name
}

func (c *columnConfigImpl) NotNull() bool {
	return c.notnull
}

func (c *columnConfigImpl) Type() columnType {
	return c.typ
}

func (c *columnConfigImpl) Options() interface{} {
	return c.opt
}

func (m *columnConfigImpl) toColumn(table Table) Column {
	return &columnImpl{
		m, table,
	}
}

func (m *columnConfigImpl) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.name))
	return
}

type columnImpl struct {
	*columnConfigImpl
	table Table
}

func (m *columnImpl) column_name() string {
	return m.name
}

func (m *columnImpl) not_null() bool {
	return m.notnull
}

func (m *columnImpl) config() ColumnConfig {
	return m.columnConfigImpl
}

func (m *columnImpl) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.table.Name()) + "." + dialect.QuoteField(m.name))
	return
}

func IntColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeInt,
	}
}

func StringColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeString,
	}
}

func DateColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeDate,
	}
}

func FloatColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeFloat,
	}
}

func BoolColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeBool,
	}
}

func BytesColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeBytes,
	}
}

func (left *columnImpl) Eq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "=")
}

func (left *columnImpl) NotEq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "<>")
}

func (left *columnImpl) Gt(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, ">")
}

func (left *columnImpl) GtEq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, ">=")
}

func (left *columnImpl) Lt(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "<")
}

func (left *columnImpl) LtEq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "<=")
}

func (left *columnImpl) Like(right string) Condition {
	return newBinaryOperationCondition(left, right, " LIKE ")
}

func (left *columnImpl) Between(lower, higher interface{}) Condition {
	return newBetweenCondition(left, lower, higher)
}

func (b ColumnList) serialize(bldr *builder) {
	first := true
	for _, column := range b {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		bldr.Append(dialect.QuoteField(column.column_name()))
	}
	return
}
