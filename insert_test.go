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
		Insert(table1).
			Columns(table1.C("str"), table1.C("bool"), table1.C("float"), table1.C("date"), table1.C("bytes")).
			Values("hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		`INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		[]interface{}{"hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		false,
	}, {
		Insert(table1).
			Set(table1.C("str"), "hoge").
			Set(table1.C("bool"), true).
			Set(table1.C("float"), 0.1).
			Set(table1.C("date"), time.Unix(0, 0).UTC()).
			Set(table1.C("bytes"), []byte{0x01}),
		`INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		[]interface{}{"hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		false,
	}, {
		// all columns if Columns() was not setted.
		Insert(table1).Values(1, "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		`INSERT INTO "TABLE_A" ( "id", "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ?, ? );`,
		[]interface{}{int64(1), "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		false,
	}, {
		// error if column's length and value's length are not eaual.
		Insert(table1).Columns(table1.C("id")).Values(1, 2, 3),
		"",
		[]interface{}{},
		true,
	}, {
		// error if into is nil.
		Insert(nil).Columns(table1.C("id")).Values(1),
		"",
		[]interface{}{},
		true,
	}, {
		// error if value type is invalid.
		Insert(table1).Columns(table1.C("str")).Values(1),
		"",
		[]interface{}{},
		true,
	}, {
		// error if table is not natural.
		Insert(tableJoined).Columns(table1.C("str")).Values(1),
		"",
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
