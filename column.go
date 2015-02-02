package sqlbuilder

type ColumnConfig interface {
	toColumn(Table) Column
}

type Column interface {
	serializable
	serializableForColumnList

	column_name() string
	not_null() bool

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
}

func (m *columnConfigImpl) toColumn(table Table) Column {
	return &columnImpl{
		m, table,
	}
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

func (m *columnImpl) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.table.Name()) + "." + dialect.QuoteField(m.name))
	return
}

func (m *columnImpl) serializeForColumnList(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.name))
	return
}

func IntColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
	}
}

func StringColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
	}
}

func DateColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
	}
}

func FloatColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
	}
}

func BoolColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
	}
}

func BytesColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
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
