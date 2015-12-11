package sqlbuilder

import (
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		StringColumn("str", &ColumnOption{
			Size: 255,
		}),
		BoolColumn("bool", nil),
		FloatColumn("float", nil),
		DateColumn("date", nil),
		BytesColumn("bytes", nil),
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
		stmt: Insert(table1).
			Columns(table1.C("str"), table1.C("bool"), table1.C("float"), table1.C("date"), table1.C("bytes")).
			Values("hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		query:  `INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		args:   []interface{}{"hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		errmsg: "",
	}, {
		stmt: Insert(table1).
			Set(table1.C("str"), "hoge").
			Set(table1.C("bool"), true).
			Set(table1.C("float"), 0.1).
			Set(table1.C("date"), time.Unix(0, 0).UTC()).
			Set(table1.C("bytes"), []byte{0x01}),
		query:  `INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		args:   []interface{}{"hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		errmsg: "",
	}, {
		stmt: Insert(table1).
			Set(table1.C("str"), "x").
			Returning(table1.C("id")),
		query:  `INSERT INTO "TABLE_A" ( "str" ) VALUES ( ? ) RETURNING "TABLE_A"."id";`,
		args:   []interface{}{"x"},
		errmsg: "",
	}, {
		stmt:   Insert(table1).Values(1, "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		query:  `INSERT INTO "TABLE_A" ( "id", "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ?, ? );`,
		args:   []interface{}{int64(1), "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		errmsg: "",
	}, {
		stmt:   Insert(table1).Columns(table1.C("id")).Values(1, 2, 3),
		query:  "",
		args:   []interface{}{},
		errmsg: "sqlbuilder: 1 values needed, but got 3.",
	}, {
		stmt:   Insert(nil).Columns(table1.C("id")).Values(1),
		query:  "",
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt:   Insert(table1).Columns(table1.C("str")).Values(1),
		query:  "",
		args:   []interface{}{},
		errmsg: "sqlbuilder: string column not accept int.",
	}, {
		stmt:   Insert(tableJoined).Columns(table1.C("str")).Values(1),
		query:  "",
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
