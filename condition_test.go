package sqlbuilder

import (
	"testing"
)

func TestBinaryCondition(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	var cases = []conditionTestCase{
		{
			table1.C("id").Eq(table1.C("test1")),
			`"TABLE_A"."id"="TABLE_A"."test1"`,
			[]interface{}{},
			nil,
		}, {
			table1.C("id").Eq(1),
			`"TABLE_A"."id"=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			table1.C("id").NotEq(1),
			`"TABLE_A"."id"<>?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			table1.C("id").Gt(1),
			`"TABLE_A"."id">?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			table1.C("id").GtEq(1),
			`"TABLE_A"."id">=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			table1.C("id").Lt(1),
			`"TABLE_A"."id"<?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			table1.C("id").LtEq(1),
			`"TABLE_A"."id"<=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			table1.C("id").Like("hoge"),
			`"TABLE_A"."id" LIKE ?`,
			[]interface{}{"hoge"},
			nil,
		}, {
			table1.C("id").Between(1, 2),
			`"TABLE_A"."id" BETWEEN ? AND ?`,
			[]interface{}{int64(1), int64(2)},
			nil,
		}, {
			table1.C("id").In(1, 2),
			`"TABLE_A"."id" IN ( ?, ? )`,
			[]interface{}{int64(1), int64(2)},
			nil,
		}, {
			table1.C("id").Eq(nil),
			`"TABLE_A"."id" IS NULL`,
			[]interface{}{},
			nil,
		}, {
			table1.C("id").NotEq([]byte(nil)),
			`"TABLE_A"."id" IS NOT NULL`,
			[]interface{}{},
			nil,
		}, {
			table1.C("id").Gt([]byte(nil)),
			`"TABLE_A"."id"`,
			[]interface{}{},
			newError("NULL can not be used with %s operator"),
		}, {
			// case for fail
			table1.C("id").In(NewTable("DUMMY TABLE", &TableOption{}, StringColumn("id", nil))),
			`"TABLE_A"."id" IN ( `,
			[]interface{}{},
			newError("got %T type, but literal is not supporting this"),
		},
	}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}

func TestBinaryConditionForSqlFunctions(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	var cases = []conditionTestCase{
		{
			Func("count", table1.C("id")).Eq(table1.C("test1")),
			`count("TABLE_A"."id")="TABLE_A"."test1"`,
			[]interface{}{},
			nil,
		}, {
			Func("count", table1.C("id")).Eq(1),
			`count("TABLE_A"."id")=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Func("count", table1.C("id")).NotEq(1),
			`count("TABLE_A"."id")<>?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Func("count", table1.C("id")).Gt(1),
			`count("TABLE_A"."id")>?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Func("count", table1.C("id")).GtEq(1),
			`count("TABLE_A"."id")>=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Func("count", table1.C("id")).Lt(1),
			`count("TABLE_A"."id")<?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Func("count", table1.C("id")).LtEq(1),
			`count("TABLE_A"."id")<=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Func("count", table1.C("id")).Like("hoge"),
			`count("TABLE_A"."id") LIKE ?`,
			[]interface{}{"hoge"},
			nil,
		}, {
			Func("count", table1.C("id")).Between(1, 2),
			`count("TABLE_A"."id") BETWEEN ? AND ?`,
			[]interface{}{int64(1), int64(2)},
			nil,
		}, {
			Func("count", table1.C("id")).In(1, 2),
			`count("TABLE_A"."id") IN ( ?, ? )`,
			[]interface{}{int64(1), int64(2)},
			nil,
		}, {
			// case for fail
			Func("count", table1.C("id")).In(NewTable("DUMMY TABLE", &TableOption{}, StringColumn("id", nil))),
			`count("TABLE_A"."id") IN ( `,
			[]interface{}{},
			newError("unsupported type"),
		},
	}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}

}

func TestConnectCondition(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	cases := []conditionTestCase{{
		And(
			table1.C("id").Eq(table1.C("test1")),
			table1.C("id").Eq(1),
			table1.C("id").Eq(2),
		),
		`"TABLE_A"."id"="TABLE_A"."test1" AND "TABLE_A"."id"=? AND "TABLE_A"."id"=?`,
		[]interface{}{int64(1), int64(2)},
		nil,
	}, {
		Or(
			table1.C("id").Eq(table1.C("test1")),
			table1.C("id").Eq(1),
		),
		`"TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=?`,
		[]interface{}{int64(1)},
		nil,
	}, {
		And(
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
		),
		`( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? ) AND ( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? )`,
		[]interface{}{int64(1), int64(1)},
		nil,
	}, {
		And(
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
			table1.C("id").Eq(table1.C("test1")),
		),
		`( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? ) AND "TABLE_A"."id"="TABLE_A"."test1"`,
		[]interface{}{int64(1)},
		nil,
	}}
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
