package sqlbuilder

type literal struct {
	val interface{}
}

func Literal(v interface{}) Expression {
	return &literal{v}
}

func (l *literal) serialize() (string, []interface{}, error) {
	return "?", []interface{}{l.val}, nil
}
