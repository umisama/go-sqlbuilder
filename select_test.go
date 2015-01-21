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

	query, attrs, err := Select(table1.C("test1"), table1.C("test2")).From(table1).ToSql()
	if query != `SELECT "TABLE_A"."test1" "TABLE_A"."test2" FROM "TABLE_A";` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{}) {
		t.Error("got", attrs)
	}
	if err != nil {
		t.Error("got", err)
	}
}
