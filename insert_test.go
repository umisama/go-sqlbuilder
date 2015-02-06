package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	a := assert.New(t)
	table1 := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)

	type testcase struct {
		stmt  Statement
		query string
		args  []interface{}
		err   bool
	}
	var cases = []testcase{{
		Insert(table1).
			Columns(table1.C("test1"), table1.C("test2")).
			Values(2, 3),
		`INSERT INTO "TABLE_A" ( "test1", "test2" ) VALUES ( ?, ? );`,
		[]interface{}{2, 3},
		false,
	}, {
		// all columns if Columns() was not setted.
		Insert(table1).Values(1, 2, 3),
		`INSERT INTO "TABLE_A" ( "id", "test1", "test2" ) VALUES ( ?, ?, ? );`,
		[]interface{}{1, 2, 3},
		false,
	}, {
		// error if column's length and value's length are not eaual.
		Insert(table1).Columns(table1.C("test1")).Values(1, 2, 3),
		"",
		[]interface{}{},
		true,
	}, {
		// error if into is nil.
		Insert(nil).Columns(table1.C("test1")).Values(1),
		"",
		[]interface{}{},
		true,
	}}

	for _, c := range cases {
		query, args, err := c.stmt.ToSql()
		a.Equal(c.query, query)
		a.Equal(c.args, args)
		if c.err {
			a.NotNil(err)
		} else {
			a.Nil(err)
		}
	}
}
