package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLiteralConvert(t *testing.T) {
	a := assert.New(t)
	type testcase struct {
		in  literal
		out interface{}
		err bool
	}
	str := "makise-kurisu"
	var cases = []testcase{
		{toLiteral(int(10)), int64(10), false},
		{toLiteral(int64(10)), int64(10), false},
		{toLiteral(uint(10)), int64(10), false},
		{toLiteral(uint64(10)), int64(10), false},
		{toLiteral(float32(10)), float64(10), false},
		{toLiteral(float64(10)), float64(10), false},
		{toLiteral(bool(true)), bool(true), false},
		{toLiteral([]byte{0x11}), []byte{0x11}, false},
		{toLiteral(string("makise-kurisu")), string("makise-kurisu"), false},
		{toLiteral(&str), str, false},
		{toLiteral((*string)(nil)), nil, false},
		{toLiteral(time.Unix(0, 0)), time.Unix(0, 0), false},
		{toLiteral(nil), nil, false},
		{toLiteral(complex(0, 0)), nil, true},
	}

	for _, c := range cases {
		val, err := c.in.(*literalImpl).converted()
		a.Equal(c.out, val)
		if c.err {
			a.Error(err)
		} else {
			a.NoError(err)
		}
	}
}

func TestLiteralString(t *testing.T) {
	a := assert.New(t)
	type testcase struct {
		in  literal
		out string
		err bool
	}
	var cases = []testcase{
		{toLiteral(int(10)), "10", false},
		{toLiteral(int64(10)), "10", false},
		{toLiteral(uint(10)), "10", false},
		{toLiteral(uint64(10)), "10", false},
		{toLiteral(float32(10)), "10.0000000000", false},
		{toLiteral(float64(10)), "10.0000000000", false},
		{toLiteral(bool(true)), "true", false},
		{toLiteral([]byte{0x11}), string([]byte{0x11}), false},
		{toLiteral(string("shibuya-rin")), "shibuya-rin", false},
		{toLiteral(time.Unix(0, 0).UTC()), "1970-01-01 00:00:00", false},
		{toLiteral(nil), "NULL", false},
		{toLiteral(complex(0, 0)), "", true},
	}

	for _, c := range cases {
		val := c.in.(*literalImpl).string()
		a.Equal(c.out, val)
	}
}

func TestLiteralIsNil(t *testing.T) {
	a := assert.New(t)
	type testcase struct {
		in  literal
		out bool
	}
	var cases = []testcase{
		{toLiteral(int(10)), false},
		{toLiteral([]byte{}), false},
		{toLiteral(nil), true},
		{toLiteral([]byte(nil)), true},
	}

	for _, c := range cases {
		isnil := c.in.IsNil()
		a.Equal(c.out, isnil)
	}
}
