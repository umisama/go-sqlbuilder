package sqlbuilder

// ColumnConfig represents a config for table's column.
// This has a name, data type and some options.
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

// ColumnList represents list of Column.
type ColumnList []Column

// Column represents a table column.
type Column interface {
	serializable

	column_name() string
	not_null() bool
	config() ColumnConfig

	// Eq creates Condition for "column==right".  Type for right is column's one or other Column.
	Eq(right interface{}) Condition

	// NotEq creates Condition for "column<>right".  Type for right is column's one or other Column.
	NotEq(right interface{}) Condition

	// GtEq creates Condition for "column>right".  Type for right is column's one or other Column.
	Gt(right interface{}) Condition

	// GtEq creates Condition for "column>=right".  Type for right is column's one or other Column.
	GtEq(right interface{}) Condition

	// Lt creates Condition for "column<right".  Type for right is column's one or other Column.
	Lt(right interface{}) Condition

	// LtEq creates Condition for "column<=right".  Type for right is column's one or other Column.
	LtEq(right interface{}) Condition

	// Like creates Condition for "column LIKE right".  Type for right is column's one or other Column.
	Like(right string) Condition

	// Between creates Condition for "column BETWEEN lower AND higher".  Type for lower/higher is int or time.Time.
	Between(lower, higher interface{}) Condition

	// In creates Condition for "column IN (values[0], values[1] ...)".  Type for values is column's one or other Column.
	In(values ...interface{}) Condition
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

// IntColumn creates config for INTEGER type column.
func IntColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeInt,
	}
}

// StringColumn creates config for TEXT or VARCHAR type column.
func StringColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeString,
	}
}

// DateColumn creates config for DATETIME type column.
func DateColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeDate,
	}
}

// FloatColumn creates config for REAL or FLOAT type column.
func FloatColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeFloat,
	}
}

// BoolColumn creates config for BOOLEAN type column.
func BoolColumn(name string, notnull bool) ColumnConfig {
	return &columnConfigImpl{
		name:    name,
		notnull: notnull,
		typ:     columnTypeBool,
	}
}

// BytesColumn creates config for BLOB type column.
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

func (left *columnImpl) In(val ...interface{}) Condition {
	return newInCondition(left, val...)
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
