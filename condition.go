package sqlbuilder

type Condition interface {
	toSql() (string, []interface{}, error)
}

type connectCondition struct {
	connector string
	conds     []Condition
}

func (c *connectCondition) toSql() (string, []interface{}, error) {
	query, attrs := "", []interface{}{}

	first := true
	for _, cond := range c.conds {
		if first {
			first = false
		} else {
			query += " " + c.connector + " "
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
