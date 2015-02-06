package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
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
		Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20).
			OrderBy(true, table1.C("test1")).
			Limit(1).
			Offset(2),
		`UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=? ORDER BY "TABLE_A"."test1" DESC LIMIT ? OFFSET ?;`,
		[]interface{}{10, 20, 1, 1, 2},
		false,
	}, {
		Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20),
		`UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=?;`,
		[]interface{}{10, 20, 1},
		false,
	}, {
		Update(nil).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20),
		``,
		[]interface{}{},
		true,
	}, {
		Update(table1).Where(table1.C("id").Eq(1)),
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
