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
		Delete(table1).Where(table1.C("id").Eq(1)),
		`DELETE FROM "TABLE_A" WHERE "TABLE_A"."id"=?;`,
		[]interface{}{int64(1)},
		false,
	}, {
		Delete(nil).Where(table1.C("id").Eq(1)),
		``,
		[]interface{}{},
		true,
	}, {
		Delete(tableJoined).Where(table1.C("id").Eq(1)),
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
