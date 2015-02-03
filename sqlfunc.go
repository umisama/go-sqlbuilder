package sqlbuilder

type SqlFunc interface {
	Column
}

type sqlFuncImpl struct {
	name    string
	columns []serializable
}

func Func(name string, columns ...Column) SqlFunc {
	cl := make([]serializable, 0, len(columns))
	for _, c := range columns {
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
