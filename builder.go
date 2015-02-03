package sqlbuilder

import (
	"bytes"
)

type builder struct {
	query *bytes.Buffer
	args  []interface{}
	err   error
}

func newBuilder() *builder {
	return &builder{
		query: bytes.NewBuffer(make([]byte, 0, 256)),
		args:  make([]interface{}, 0, 8),
		err:   nil,
	}
}

func (b *builder) Err() error {
	if b.err != nil {
		return b.err
	}
	return nil
}

func (b *builder) Query() string {
	if b.err != nil {
		return ""
	}
	return b.query.String()
}

func (b *builder) Args() []interface{} {
	if b.err != nil {
		return []interface{}{}
	}
	return b.args
}

func (b *builder) SetError(err error) {
	if b.err != nil {
		return
	}
	b.err = err
	return
}

func (b *builder) Append(query string) {
	if b.err != nil {
		return
	}

	b.query.WriteString(query)
}

func (b *builder) AppendValue(val interface{}) {
	if b.err != nil {
		return
	}

	b.query.WriteString(dialect.BindVar(len(b.args) + 1))
	b.args = append(b.args, val)
	return
}

func (b *builder) AppendItems(parts []serializable, sep string) {
	if b.err != nil {
		return
	}

	first := true
	for _, part := range parts {
		if first {
			first = false
		} else {
			b.Append(sep)
		}
		part.serialize(b)
	}
	return
}

func (b *builder) AppendItem(part serializable) {
	if b.err != nil {
		return
	}
	part.serialize(b)
}
