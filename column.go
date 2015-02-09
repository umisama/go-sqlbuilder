package sqlbuilder

import (
	"reflect"
	"time"
)

// ColumnConfig represents a config for table's column.
// This has a name, data type and some options.
type ColumnConfig interface {
	serializable

	toColumn(Table) Column
	Name() string
	Type() columnType
	Options() []ColumnOption
	HasOption(ColumnOption) bool
	Size() int
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

func (t columnType) String() string {
	switch t {
	case columnTypeInt:
		return "int"
	case columnTypeString:
		return "string"
	case columnTypeDate:
		return "date"
	case columnTypeFloat:
		return "float"
	case columnTypeBool:
		return "bool"
	case columnTypeBytes:
		return "bytes"
	}
	panic(newError("unknown columnType"))
}

func (t columnType) CapableTypes() []reflect.Type {
	switch t {
	case columnTypeInt:
		return []reflect.Type{
			reflect.TypeOf(int(0)),
			reflect.TypeOf(int8(0)),
			reflect.TypeOf(int16(0)),
			reflect.TypeOf(int32(0)),
			reflect.TypeOf(int64(0)),
			reflect.TypeOf(uint(0)),
			reflect.TypeOf(uint8(0)),
			reflect.TypeOf(uint16(0)),
			reflect.TypeOf(uint32(0)),
			reflect.TypeOf(uint64(0)),
		}
	case columnTypeString:
		return []reflect.Type{
			reflect.TypeOf(""),
		}
	case columnTypeDate:
		return []reflect.Type{
			reflect.TypeOf(time.Time{}),
		}
	case columnTypeFloat:
		return []reflect.Type{
			reflect.TypeOf(float32(0)),
			reflect.TypeOf(float64(0)),
		}
	case columnTypeBool:
		return []reflect.Type{
			reflect.TypeOf(bool(true)),
		}
	case columnTypeBytes:
		return []reflect.Type{
			reflect.TypeOf([]byte{}),
		}
	}
	return []reflect.Type{}
}

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
	acceptType(interface{}) bool

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
	size int // size for varchar column
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

func (m *columnImpl) acceptType(val interface{}) bool {
	lit, ok := val.(literal)
	if !ok || lit == nil {
		return false
	}
	if reflect.ValueOf(lit).IsNil() {
		return !m.HasOption(CO_NotNull)
	}

	valt := reflect.TypeOf(lit.Raw())
	for _, t := range m.typ.CapableTypes() {
		if t == valt {
			return true
		}
	}
	return false
}

func (m *columnConfigImpl) Size() int {
	return m.size
}

func (m *columnImpl) serialize(bldr *builder) {
	if m == Star {
		bldr.Append("*")
	} else {
		bldr.Append(dialect.QuoteField(m.table.Name()) + "." + dialect.QuoteField(m.name))
	}
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
func StringColumn(name string, size int, opts ...ColumnOption) ColumnConfig {
	return &columnConfigImpl{
		name: name,
		typ:  columnTypeString,
		opts: opts,
		size: size,
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
		if column == nil {
			bldr.SetError(newError("column is not found"))
			return
		}
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		bldr.Append(dialect.QuoteField(column.column_name()))
	}
	return
}
