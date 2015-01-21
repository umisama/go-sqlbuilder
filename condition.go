package sqlbuilder

type Condition interface {
	toSql() (string, []interface{}, error)
}

type andCondition struct {
	conds []Condition
}

func (c *andCondition) toSql() (string, []interface{}, error) {
	query, attrs := "", []interface{}{}

	first := true
	for _, cond := range c.conds {
		if first {
			first = false
		} else {
			query += " AND "
		}

		q, a, err := cond.toSql()
		if err != nil {
			return "", []interface{}{}, nil
		}

		query += q
		attrs = append(attrs, a...)
	}

	return query, attrs, nil
}

type orCondition struct {
	conds []Condition
}

func (c *orCondition) toSql() (string, []interface{}, error) {
	query, attrs := "", []interface{}{}

	first := true
	for _, cond := range c.conds {
		if first {
			first = false
		} else {
			query += " OR "
		}

		q, a, err := cond.toSql()
		if err != nil {
			return "", []interface{}{}, nil
		}

		query += q
		attrs = append(attrs, a...)
	}

	return query, attrs, nil
}

func And(conds ...Condition) Condition {
	return &andCondition{
		conds: conds,
	}
}

func Or(conds ...Condition) Condition {
	return &orCondition{
		conds: conds,
	}
}

type eqCondition struct {
	left  interface{}
	right interface{}
}

func Eq(left, right Column) Condition {
	return &eqCondition{
		left:  left,
		right: right,
	}
}

func EqL(left Column, right interface{}) Condition {
	return &eqCondition{
		left:  left,
		right: right,
	}
}

func (c *eqCondition) toSql() (string, []interface{}, error) {
	query, attrs := "", []interface{}{}
	switch l := c.left.(type) {
	case Column:
		n, _, err := l.toSql()
		if err != nil {
			return "", []interface{}{}, err
		}
		query += n
	default:
		query += "?"
		attrs = append(attrs, l)
	}

	query += "="

	switch r := c.right.(type) {
	case Column:
		n, _, err := r.toSql()
		if err != nil {
			return "", []interface{}{}, err
		}
		query += n
	default:
		query += "?"
		attrs = append(attrs, r)
	}

	return query, attrs, nil
}
