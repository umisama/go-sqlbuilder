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
			cond:   table1.C("id").Eq(table1.C("test1")),
			query:  `"TABLE_A"."id"="TABLE_A"."test1"`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   table1.C("id").Eq(1),
			query:  `"TABLE_A"."id"=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").NotEq(1),
			query:  `"TABLE_A"."id"<>?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Gt(1),
			query:  `"TABLE_A"."id">?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").GtEq(1),
			query:  `"TABLE_A"."id">=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Lt(1),
			query:  `"TABLE_A"."id"<?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").LtEq(1),
			query:  `"TABLE_A"."id"<=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Like("hoge"),
			query:  `"TABLE_A"."id" LIKE ?`,
			args:   []interface{}{"hoge"},
			errmsg: "",
		}, {
			cond:   table1.C("id").Between(1, 2),
			query:  `"TABLE_A"."id" BETWEEN ? AND ?`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   table1.C("id").In(1, 2),
			query:  `"TABLE_A"."id" IN ( ?, ? )`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Eq(nil),
			query:  `"TABLE_A"."id" IS NULL`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   table1.C("id").NotEq([]byte(nil)),
			query:  `"TABLE_A"."id" IS NOT NULL`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   table1.C("id").Gt([]byte(nil)),
			query:  `"TABLE_A"."id"`,
			args:   []interface{}{},
			errmsg: "sqlbuilder: NULL can not be used with > operator.",
		}, {
			// case for fail
			cond:   table1.C("id").In(NewTable("DUMMY TABLE", &TableOption{}, StringColumn("id", nil))),
			query:  `"TABLE_A"."id" IN ( `,
			args:   []interface{}{},
			errmsg: "sqlbuilder: got sqlbuilder.table type, but literal is not supporting this.",
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
			cond:   Func("count", table1.C("id")).Eq(table1.C("test1")),
			query:  `count("TABLE_A"."id")="TABLE_A"."test1"`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Eq(1),
			query:  `count("TABLE_A"."id")=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).NotEq(1),
			query:  `count("TABLE_A"."id")<>?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Gt(1),
			query:  `count("TABLE_A"."id")>?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).GtEq(1),
			query:  `count("TABLE_A"."id")>=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Lt(1),
			query:  `count("TABLE_A"."id")<?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).LtEq(1),
			query:  `count("TABLE_A"."id")<=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Like("hoge"),
			query:  `count("TABLE_A"."id") LIKE ?`,
			args:   []interface{}{"hoge"},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Between(1, 2),
			query:  `count("TABLE_A"."id") BETWEEN ? AND ?`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).In(1, 2),
			query:  `count("TABLE_A"."id") IN ( ?, ? )`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).In(NewTable("DUMMY TABLE", &TableOption{}, StringColumn("id", nil))),
			query:  `count("TABLE_A"."id") IN ( `,
			args:   []interface{}{},
			errmsg: "sqlbuilder: got sqlbuilder.table type, but literal is not supporting this.",
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
		cond: And(
			table1.C("id").Eq(table1.C("test1")),
			table1.C("id").Eq(1),
			table1.C("id").Eq(2),
		),
		query:  `"TABLE_A"."id"="TABLE_A"."test1" AND "TABLE_A"."id"=? AND "TABLE_A"."id"=?`,
		args:   []interface{}{int64(1), int64(2)},
		errmsg: "",
	}, {
		cond: Or(
			table1.C("id").Eq(table1.C("test1")),
			table1.C("id").Eq(1),
		),
		query:  `"TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=?`,
		args:   []interface{}{int64(1)},
		errmsg: "",
	}, {
		cond: And(
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
		),
		query:  `( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? ) AND ( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? )`,
		args:   []interface{}{int64(1), int64(1)},
		errmsg: "",
	}, {
		cond: And(
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
			table1.C("id").Eq(table1.C("test1")),
		),
		query:  `( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? ) AND "TABLE_A"."id"="TABLE_A"."test1"`,
		args:   []interface{}{int64(1)},
		errmsg: "",
	}}
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
