package sqlbuilder

type Condition interface {
	serialize() (string, []interface{}, error)
}

type connectCondition struct {
	connector string
	conds     []Condition
}

func (c *connectCondition) serialize() (string, []interface{}, error) {
	query, attrs := "", []interface{}{}

	first := true
	for _, cond := range c.conds {
		if first {
			first = false
		} else {
			query += " " + c.connector + " "
		}

		q, a, err := cond.serialize()
		if err != nil {
			return "", []interface{}{}, nil
		}

		query += q
		attrs = append(attrs, a...)
	}

	return query, attrs, nil
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

func (c *binaryOperationCondition) serialize() (string, []interface{}, error) {
	query, attrs := "", []interface{}{}

	// left hand side
	var err error
	query, attrs, err = appendItemToQuery(c.left, query, attrs)
	if err != nil {
		return "", []interface{}{}, err
	}

	// operator
	query += c.operator

	// right hand side
	query, attrs, err = appendItemToQuery(c.right, query, attrs)
	if err != nil {
		return "", []interface{}{}, err
	}

	return query, attrs, nil
}

type betweenCondition struct {
	left   Expression
	lower  Expression
	higher Expression
}

func (c *betweenCondition) serialize() (string, []interface{}, error) {
	query, attrs := "", []interface{}{}

	// left hand side
	var err error
	query, attrs, err = appendItemToQuery(c.left, query, attrs)
	if err != nil {
		return "", []interface{}{}, err
	}

	// operator
	query += " BETWEEN "

	// right hand side(lower)
	query, attrs, err = appendItemToQuery(c.lower, query, attrs)
	if err != nil {
		return "", []interface{}{}, err
	}

	query += " AND "

	// right hand side(higher)
	query, attrs, err = appendItemToQuery(c.higher, query, attrs)
	if err != nil {
		return "", []interface{}{}, err
	}

	return query, attrs, nil
}
