package sqlbuilder

import (
	sqldriver "database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Literal interface {
	serializable
	Raw() interface{}
}

type literal struct {
	raw         interface{}
	placeholder bool
}

func L(v interface{}) Literal {
	return &literal{
		raw:         v,
		placeholder: true,
	}
}

func (l *literal) serialize() (string, []interface{}, error) {
	val, err := l.converted()
	if err != nil {
		return "", []interface{}{}, err
	}

	if l.placeholder {
		return "?", []interface{}{val}, nil
	} else {
		return l.string(), []interface{}{}, nil
	}
}

func (l *literal) converted() (interface{}, error) {
	switch t := l.raw.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(t).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(t).Uint()), nil
	case float32, float64:
		return reflect.ValueOf(l.raw).Float(), nil
	case bool:
		return t, nil
	case []byte:
		return t, nil
	case string:
		return t, nil
	case time.Time:
		return t, nil
	case sqldriver.Valuer:
		return t, nil
	}

	return nil, errors.New("sqlbuilder: unsupported type")
}

func (l *literal) string() string {
	switch t := l.raw.(type) {
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', 10, 64)
	case bool:
		return strconv.FormatBool(t)
	case string:
		return t
	case time.Time:
		return t.Format(time.ANSIC)
	case fmt.Stringer:
		return t.String()
	}
	return ""
}

func (l *literal) Raw() interface{} {
	return l.raw
}
