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
	table2 := NewTable(
		"TABLE_B",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)
	tableJoined := table1.InnerJoin(table2, table1.C("test1").Eq(table2.C("id")))

	var cases = []statementTestCase{{
		stmt:   DropTable(table1),
		query:  `DROP TABLE "TABLE_A";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   DropTable(nil),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt:   DropTable(tableJoined),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is not natural table.",
	}}
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
