package sqlbuilder

import "fmt"

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
			bldr.Append(" " + c.connector + " ")
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
	left     serializable
	right    serializable
	operator string
	err      error
}

func newBinaryOperationCondition(left, right interface{}, operator string) *binaryOperationCondition {
	cond := &binaryOperationCondition{
		operator: operator,
	}
	column_exist := false
	switch t := left.(type) {
	case Column:
		column_exist = true
		cond.left = t
	default:
		cond.left = toLiteral(t)
	}
	switch t := right.(type) {
	case Column:
		column_exist = true
		cond.right = t
	default:
		cond.right = toLiteral(t)
	}
	if !column_exist {
		cond.err = fmt.Errorf("hello world")
	}

	return cond
}

func newBetweenCondition(left Column, low, high interface{}) Condition {
	low_literal := toLiteral(low)
	high_literal := toLiteral(high)

	return &betweenCondition{
		left:   left,
		lower:  low_literal,
		higher: high_literal,
	}
}

func (c *binaryOperationCondition) serialize(bldr *builder) {
	bldr.AppendItem(c.left)
	bldr.Append(c.operator)
	bldr.AppendItem(c.right)
	return
}

type betweenCondition struct {
	left   serializable
	lower  serializable
	higher serializable
}

func (c *betweenCondition) serialize(bldr *builder) {
	bldr.AppendItem(c.left)
	bldr.Append(" BETWEEN ")
	bldr.AppendItem(c.lower)
	bldr.Append(" AND ")
	bldr.AppendItem(c.higher)
	return
}
