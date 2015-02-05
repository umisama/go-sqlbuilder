package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
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
		Select(table1.C("test1"), table1.C("test2")).
			From(table1).
			Where(
			And(
				table1.C("id").Eq(1),
				table1.C("test1").Eq(2),
			)).
			Distinct().
			OrderBy(false, table1.C("id")).
			GroupBy(table1.C("id")).
			Having(table1.C("id").Eq(1)).
			Limit(10).
			Offset(20),
		`SELECT DISTINCT "TABLE_A"."test1", "TABLE_A"."test2" ` +
			`FROM "TABLE_A" WHERE "TABLE_A"."id"=? AND "TABLE_A"."test1"=? ` +
			`GROUP BY "TABLE_A"."id" HAVING "TABLE_A"."id"=? ORDER BY "TABLE_A"."id" ASC ` +
			`LIMIT ? OFFSET ?;`,
		[]interface{}{1, 2, 1, 10, 20},
		false,
	}, {
		Select(table1.C("test1"), table1.C("test2")).
			From(table1),
		`SELECT "TABLE_A"."test1", "TABLE_A"."test2" FROM "TABLE_A";`,
		[]interface{}{},
		false,
	}, {
		Select(table1.C("test1"), table1.C("test2")).
			From(nil),
		``,
		[]interface{}{},
		true,
	}, {
		Select(table1.C("test1"), table1.C("test2")).
			From(table1).Having(table1.C("id").Eq(1)),
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

func BenchmarkSelect(b *testing.B) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)

	for i := 0; i < b.N; i++ {
		Select(table1.C("test1"), table1.C("test2")).
			From(table1).
			Where(
			And(
				table1.C("id").Eq(1),
				table1.C("test1").Eq(2),
			)).
			ToSql()
	}
}
