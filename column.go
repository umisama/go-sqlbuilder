package sqlbuilder

// ColumnConfig represents a config for table's column.
// This has a name, data type and some options.
type ColumnConfig interface {
	serializable

	toColumn(Table) Column
	Name() string
	Type() columnType
	Options() []ColumnOption
	HasOption(ColumnOption) bool
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

// ColumnOption represents options for columns. ex: primary key.
// Use const CO_*
type ColumnOption int

const (
	CO_PrimaryKey ColumnOption = iota
	CO_NotNull
	CO_Unique
	CO_AutoIncrement
)

// ColumnList represents list of Column.
type ColumnList []Column

// Column represents a table column.
type Column interface {
	serializable

	column_name() string
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
	name string
	typ  columnType
	opts []ColumnOption
}

func (c *columnConfigImpl) Name() string {
	return c.name
}

func (c *columnConfigImpl) Type() columnType {
	return c.typ
}

func (c *columnConfigImpl) Options() []ColumnOption {
	return c.opts
}

func (c *columnConfigImpl) HasOption(trg ColumnOption) bool {
	for _, v := range c.opts {
		if v == trg {
			return true
		}
	}
	return false
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

func (m *columnImpl) config() ColumnConfig {
	return m.columnConfigImpl
}

func (m *columnImpl) serialize(bldr *builder) {
	bldr.Append(dialect.QuoteField(m.table.Name()) + "." + dialect.QuoteField(m.name))
	return
}

// IntColumn creates config for INTEGER type column.
func IntColumn(name string, opts ...ColumnOption) ColumnConfig {
	return &columnConfigImpl{
		name: name,
		typ:  columnTypeInt,
		opts: opts,
	}
}

// StringColumn creates config for TEXT or VARCHAR type column.
func StringColumn(name string, opts ...ColumnOption) ColumnConfig {
	return &columnConfigImpl{
		name: name,
		typ:  columnTypeString,
		opts: opts,
	}
}

// DateColumn creates config for DATETIME type column.
func DateColumn(name string, opts ...ColumnOption) ColumnConfig {
	return &columnConfigImpl{
		name: name,
		typ:  columnTypeDate,
		opts: opts,
	}
}

// FloatColumn creates config for REAL or FLOAT type column.
func FloatColumn(name string, opts ...ColumnOption) ColumnConfig {
	return &columnConfigImpl{
		name: name,
		typ:  columnTypeFloat,
		opts: opts,
	}
}

// BoolColumn creates config for BOOLEAN type column.
func BoolColumn(name string, opts ...ColumnOption) ColumnConfig {
	return &columnConfigImpl{
		name: name,
		typ:  columnTypeBool,
		opts: opts,
	}
}

// BytesColumn creates config for BLOB type column.
func BytesColumn(name string, opts ...ColumnOption) ColumnConfig {
	return &columnConfigImpl{
		name: name,
		typ:  columnTypeBytes,
		opts: opts,
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
