package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDropTable(t *testing.T) {
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
		DropTable(table1),
		`DROP TABLE "TABLE_A";`,
		[]interface{}{},
		false,
	}, {
		DropTable(nil),
		``,
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
