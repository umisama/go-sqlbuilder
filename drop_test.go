package sqlbuilder

import (
	"testing"
)

func TestDropTable(t *testing.T) {
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
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
