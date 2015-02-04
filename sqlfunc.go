package sqlbuilder

// SqlFunc represents function on SQL(ex:count(*)).  This can be use in the same way as Column.
type SqlFunc interface {
	Column
}

type sqlFuncImpl struct {
	name    string
	columns []serializable
}

// Func returns new SQL function.  The name is function name, and the args is arguments of function
func Func(name string, args ...Column) SqlFunc {
	cl := make([]serializable, 0, len(args))
	for _, c := range args {
		cl = append(cl, c)
	}
	return &sqlFuncImpl{
		name:    name,
		columns: cl,
	}
}

func (m *sqlFuncImpl) column_name() string {
	return m.name
}

func (m *sqlFuncImpl) not_null() bool {
	return true
}

func (m *sqlFuncImpl) config() ColumnConfig {
	return nil
}

func (m *sqlFuncImpl) serialize(bldr *builder) {
	bldr.Append(m.name)
	bldr.Append("(")
	bldr.AppendItems(m.columns, ", ")
	bldr.Append(")")
}

func (left *sqlFuncImpl) Eq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "=")
}

func (left *sqlFuncImpl) NotEq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "<>")
}

func (left *sqlFuncImpl) Gt(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, ">")
}

func (left *sqlFuncImpl) GtEq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, ">=")
}

func (left *sqlFuncImpl) Lt(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "<")
}

func (left *sqlFuncImpl) LtEq(right interface{}) Condition {
	return newBinaryOperationCondition(left, right, "<=")
}

func (left *sqlFuncImpl) Like(right string) Condition {
	return newBinaryOperationCondition(left, right, " LIKE ")
}

func (left *sqlFuncImpl) Between(lower, higher interface{}) Condition {
	return newBetweenCondition(left, lower, higher)
}

func (left *sqlFuncImpl) In(vals ...interface{}) Condition {
	return newInCondition(left, vals...)
}
