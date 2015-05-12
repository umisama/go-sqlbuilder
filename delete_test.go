package sqlbuilder

import (
	"testing"
)

func TestDelete(t *testing.T) {
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
		stmt:   Delete(table1).Where(table1.C("id").Eq(1)),
		query:  `DELETE FROM "TABLE_A" WHERE "TABLE_A"."id"=?;`,
		args:   []interface{}{int64(1)},
		errmsg: "",
	}, {
		stmt:   Delete(nil).Where(table1.C("id").Eq(1)),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: from is nil.",
	}, {
		stmt:   Delete(tableJoined).Where(table1.C("id").Eq(1)),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: CreateTable can use only natural table.",
	}}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
