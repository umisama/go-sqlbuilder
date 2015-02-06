package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	a := assert.New(t)
	table1 := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		StringColumn("str"),
		BoolColumn("bool"),
		FloatColumn("float"),
		DateColumn("date"),
		BytesColumn("bytes"),
	)

	type testcase struct {
		stmt  Statement
		query string
		args  []interface{}
		err   bool
	}
	var cases = []testcase{{
		Insert(table1).
			Columns(table1.C("str"), table1.C("bool"), table1.C("float"), table1.C("date"), table1.C("bytes")).
			Values("hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		`INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		[]interface{}{"hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		false,
	}, {
		// all columns if Columns() was not setted.
		Insert(table1).Values(1, "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		`INSERT INTO "TABLE_A" ( "id", "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ?, ? );`,
		[]interface{}{1, "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		false,
	}, {
		// error if column's length and value's length are not eaual.
		Insert(table1).Columns(table1.C("id")).Values(1, 2, 3),
		"",
		[]interface{}{},
		true,
	}, {
		// error if into is nil.
		Insert(nil).Columns(table1.C("id")).Values(1),
		"",
		[]interface{}{},
		true,
	}, {
		// error if into is nil.
		Insert(table1).Columns(table1.C("str")).Values(1),
		"",
		[]interface{}{},
		true,
	}}

	for _, c := range cases {
		query, args, err := c.stmt.ToSql()
		a.Equal(c.query, query)
		a.Equal(c.args, args)
		if c.err {
			a.Error(err)
		} else {
			a.NoError(err)
		}
	}
}
