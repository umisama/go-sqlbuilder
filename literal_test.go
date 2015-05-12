package sqlbuilder

import (
	"reflect"
	"testing"
	"time"
)

func TestLiteralConvert(t *testing.T) {
	str := "makise-kurisu"
	var cases = []struct {
		lit    literal
		out    interface{}
		errmes string
	}{
		{
			lit:    toLiteral(int(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(int64(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(uint(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(uint64(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(float32(10)),
			out:    float64(10),
			errmes: "",
		}, {
			lit:    toLiteral(float64(10)),
			out:    float64(10),
			errmes: "",
		}, {
			lit:    toLiteral(bool(true)),
			out:    bool(true),
			errmes: "",
		}, {
			lit:    toLiteral([]byte{0x11}),
			out:    []byte{0x11},
			errmes: "",
		}, {
			lit:    toLiteral(string("makise-kurisu")),
			out:    string("makise-kurisu"),
			errmes: "",
		}, {
			lit:    toLiteral(&str),
			out:    str,
			errmes: "",
		}, {
			lit:    toLiteral((*string)(nil)),
			out:    nil,
			errmes: "",
		}, {
			lit:    toLiteral(time.Unix(0, 0)),
			out:    time.Unix(0, 0),
			errmes: "",
		}, {
			lit:    toLiteral(nil),
			out:    nil,
			errmes: "",
		}, {
			lit:    toLiteral(complex(0, 0)),
			out:    nil,
			errmes: "sqlbuilder: got complex128 type, but literal is not supporting this.",
		}}

	for num, c := range cases {
		val, err := c.lit.(*literalImpl).converted()
		if !reflect.DeepEqual(c.out, val) {
			t.Errorf("failed on %d", num)
		}
		if len(c.errmes) != 0 {
			if err == nil {
				t.Errorf("failed on %d", num)
			}
			if err.Error() != c.errmes {
				t.Errorf("failed on %d", num)
				panic(err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("failed on %d", num)
			}
		}
	}
}

func TestLiteralString(t *testing.T) {
	var cases = []struct {
		lit    literal
		out    string
		errmes string
	}{
		{
			lit:    toLiteral(int(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(int64(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(uint(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(uint64(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(float32(10)),
			out:    "10.0000000000",
			errmes: "",
		}, {
			lit:    toLiteral(float64(10)),
			out:    "10.0000000000",
			errmes: "",
		}, {
			lit:    toLiteral(bool(true)),
			out:    "true",
			errmes: "",
		}, {
			lit:    toLiteral([]byte{0x11}),
			out:    string([]byte{0x11}),
			errmes: "",
		}, {
			lit:    toLiteral(string("shibuya-rin")),
			out:    "shibuya-rin",
			errmes: "",
		}, {
			lit:    toLiteral(time.Unix(0, 0).UTC()),
			out:    "1970-01-01 00:00:00",
			errmes: "",
		}, {
			lit:    toLiteral(nil),
			out:    "NULL",
			errmes: "",
		}, {
			lit:    toLiteral(complex(0, 0)),
			out:    "",
			errmes: "aaa",
		}}

	for num, c := range cases {
		val := c.lit.(*literalImpl).string()
		if c.out != val {
			t.Error("failed on %d", num)
		}
	}
}

func TestLiteralIsNil(t *testing.T) {
	var cases = []struct {
		in  literal
		out bool
	}{
		{toLiteral(int(10)), false},
		{toLiteral([]byte{}), false},
		{toLiteral(nil), true},
		{toLiteral([]byte(nil)), true},
	}

	for num, c := range cases {
		isnil := c.in.IsNil()
		if c.out != isnil {
			t.Error("failed on %d", num)
		}
	}
}
