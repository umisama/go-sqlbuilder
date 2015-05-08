package sqlbuilder

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	var cases = []statementTestCase{{
		Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20).
			OrderBy(true, table1.C("test1")).
			Limit(1).
			Offset(2),
		`UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=? ORDER BY "TABLE_A"."test1" DESC LIMIT ? OFFSET ?;`,
		[]interface{}{int64(10), int64(20), int64(1), 1, 2},
		false,
	}, {
		Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20),
		`UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=?;`,
		[]interface{}{int64(10), int64(20), int64(1)},
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
	}, {
		Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), "foo"),
		``,
		[]interface{}{},
		true,
	}}
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}

}
