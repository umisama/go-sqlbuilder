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
	table2 := NewTable(
		"TABLE_B",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)
	tableJoined := table1.InnerJoin(table2, table1.C("test1").Eq(table2.C("id")))

	var cases = []statementTestCase{{
		stmt: Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20).
			OrderBy(true, table1.C("test1")).
			Limit(1).
			Offset(2),
		query:  `UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=? ORDER BY "TABLE_A"."test1" DESC LIMIT ? OFFSET ?;`,
		args:   []interface{}{int64(10), int64(20), int64(1), 1, 2},
		errmsg: "",
	}, {
		stmt: Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20),
		query:  `UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=?;`,
		args:   []interface{}{int64(10), int64(20), int64(1)},
		errmsg: "",
	}, {
		stmt: Update(nil).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt:   Update(table1).Where(table1.C("id").Eq(1)),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: length of sets is 0.",
	}, {
		stmt: Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), "foo"),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: int column not accept string.",
	}, {
		stmt:   Update(tableJoined),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: length of sets is 0.",
	}}
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
