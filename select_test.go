package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	query, attrs, err := Select(table1.C("test1"), table1.C("test2")).
		From(table1).
		Where(
		And(
			Eq(table1.C("id"), L(1)),
			Eq(table1.C("test1"), L(2)),
		)).
		ToSql()
	if query != `SELECT "TABLE_A"."test1" "TABLE_A"."test2" FROM "TABLE_A" WHERE "TABLE_A"."id"=? AND "TABLE_A"."test1"=?;` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{int64(1), int64(2)}) {
		t.Error("got", attrs)
	}
	if err != nil {
		t.Error("got", err)
	}
}

func BenchmarkSelect(b *testing.B) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	for i := 0; i < b.N; i++ {
		Select(table1.C("test1"), table1.C("test2")).
			From(table1).
			Where(
			And(
				Eq(table1.C("id"), L(1)),
				Eq(table1.C("test1"), L(2)),
			)).
			ToSql()
	}
}
