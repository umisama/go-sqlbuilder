package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	acol_id := table1.C("id").As("tbl1id")

	var cases = []statementTestCase{{
		Select(table1).
			Columns(table1.C("test1"), table1.C("test2")).
			Where(
			And(
				table1.C("id").Eq(1),
				table1.C("test1").Eq(2),
			)).
			Distinct().
			OrderBy(false, table1.C("id")).
			GroupBy(table1.C("id")).
			Having(table1.C("id").Eq(1)).
			Limit(10).
			Offset(20),
		`SELECT DISTINCT "TABLE_A"."test1", "TABLE_A"."test2" ` +
			`FROM "TABLE_A" WHERE "TABLE_A"."id"=? AND "TABLE_A"."test1"=? ` +
			`GROUP BY "TABLE_A"."id" HAVING "TABLE_A"."id"=? ORDER BY "TABLE_A"."id" ASC ` +
			`LIMIT ? OFFSET ?;`,
		[]interface{}{int64(1), int64(2), int64(1), 10, 20},
		false,
	}, {
		Select(table1).
			Columns(table1.C("test1"), table1.C("test2")),
		`SELECT "TABLE_A"."test1", "TABLE_A"."test2" FROM "TABLE_A";`,
		[]interface{}{},
		false,
	}, {
		Select(table1).
			Columns(acol_id).
			Where(acol_id.Eq(1)),
		`SELECT "TABLE_A"."id" AS "tbl1id" FROM "TABLE_A" WHERE "tbl1id"=?;`,
		[]interface{}{int64(1)},
		false,
	}, {
		Select(table1).
			Columns(acol_id).
			Where(acol_id.Eq(1)).
			OrderBy(false, table1.C("test1")).
			OrderBy(true, table1.C("test2")),
		`SELECT "TABLE_A"."id" AS "tbl1id" FROM "TABLE_A" WHERE "tbl1id"=? ORDER BY "TABLE_A"."test1" ASC, "TABLE_A"."test2" DESC;`,
		[]interface{}{int64(1)},
		false,
	}, {
		Select(table1).
			Columns(Star),
		`SELECT * FROM "TABLE_A";`,
		[]interface{}{},
		false,
	}, {
		Select(table1),
		`SELECT * FROM "TABLE_A";`,
		[]interface{}{},
		false,
	}, {
		Select(nil).
			Columns(table1.C("test1"), table1.C("test2")),
		``,
		[]interface{}{},
		true,
	}, {
		Select(table1).
			Columns(table1.C("test1"), table1.C("test2")).
			Having(table1.C("id").Eq(1)),
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

func TestSubquery(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	subquery := Select(table1).Columns(table1.C("id")).ToSubquery("SQ1")
	query, attrs, err := Select(subquery).
		Columns(subquery.C("id")).
		Where(subquery.C("id").Eq(1)).ToSql()

	if `SELECT "SQ1"."id" FROM ( SELECT "TABLE_A"."id" FROM "TABLE_A" ) AS SQ1 WHERE "SQ1"."id"=?;` != query {
		t.Error("failed \ngot %s", query)
	}
	if !reflect.DeepEqual([]interface{}{int64(1)}, attrs) {
		t.Error("failed \ngot %#v", attrs)
	}
	if err != nil {
		t.Error("failed \ngot %#v", err.Error())
	}
}

func BenchmarkSelect(b *testing.B) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	for i := 0; i < b.N; i++ {
		Select(table1).
			Columns(table1.C("test1"), table1.C("test2")).
			Where(
			And(
				table1.C("id").Eq(1),
				table1.C("test1").Eq(2),
			)).
			ToSql()
	}
}
