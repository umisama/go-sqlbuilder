package sqlbuilder

type Condition interface {
	serializable
}

type connectCondition struct {
	connector string
	conds     []Condition
}

func (c *connectCondition) serialize(bldr *builder) {
	first := true
	for _, cond := range c.conds {
		if first {
			first = false
		} else {
			bldr.Append(" "+c.connector+" ", nil)
		}
		cond.serialize(bldr)
	}
	return
}

func And(conds ...Condition) Condition {
	return &connectCondition{
		connector: "AND",
		conds:     conds,
	}
}

func Or(conds ...Condition) Condition {
	return &connectCondition{
		connector: "OR",
		conds:     conds,
	}
}

type binaryOperationCondition struct {
	left     Expression
	right    Expression
	operator string
}

func Eq(left, right Expression) Condition {
	return &binaryOperationCondition{
		left:     left,
		right:    right,
		operator: "=",
	}
}

func NotEq(left, right Expression) Condition {
	return &binaryOperationCondition{
		left:     left,
		right:    right,
		operator: "<>",
	}
}

func Gt(left, right Expression) Condition {
	return &binaryOperationCondition{
		left:     left,
		right:    right,
		operator: ">",
	}
}

func Gte(left, right Expression) Condition {
	return &binaryOperationCondition{
		left:     left,
		right:    right,
		operator: ">=",
	}
}

func Lt(left, right Expression) Condition {
	return &binaryOperationCondition{
		left:     left,
		right:    right,
		operator: "<",
	}
}

func Lte(left, right Expression) Condition {
	return &binaryOperationCondition{
		left:     left,
		right:    right,
		operator: "<=",
	}
}

func Like(left, right Expression) Condition {
	return &binaryOperationCondition{
		left:     left,
		right:    right,
		operator: " LIKE ",
	}
}

func Between(left, low, high Expression) Condition {
	return &betweenCondition{
		left:   left,
		lower:  low,
		higher: high,
	}
}

func (c *binaryOperationCondition) serialize(bldr *builder) {
	bldr.AppendItem(c.left)
	bldr.Append(c.operator, nil)
	bldr.AppendItem(c.right)
	return
}

type betweenCondition struct {
	left   Expression
	lower  Expression
	higher Expression
}

func (c *betweenCondition) serialize(bldr *builder) {
	bldr.AppendItem(c.left)
	bldr.Append(" BETWEEN ", nil)
	bldr.AppendItem(c.lower)
	bldr.Append(" AND ", nil)
	bldr.AppendItem(c.higher)
	return
}
