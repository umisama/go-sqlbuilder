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
	}}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
