package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestEqCondition(t *testing.T) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	// case Eq:  (column)=(column)
	eq := Eq(table1.C("id"), table1.C("test1"))
	query, attrs, err := eq.toSql()
	if query != `"TABLE_A"."id"="TABLE_A"."test1"` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{}) {
		t.Error("got", attrs)
	}
	if err != nil {
		t.Error("got", err)
	}

	// case EqL:  (column)=(literal)
	eqL := EqL(table1.C("id"), 1)
	query, attrs, err = eqL.toSql()
	if query != `"TABLE_A"."id"=?` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{1}) {
		t.Error("got", attrs)
	}
	if err != nil {
		t.Error("got", err)
	}
}

func TestAndCondition(t *testing.T) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)
	eq1 := Eq(table1.C("id"), table1.C("test1"))
	eq2 := EqL(table1.C("id"), 1)

	and := And(eq1, eq2)
	query, attrs, err := and.toSql()
	if query != `"TABLE_A"."id"="TABLE_A"."test1" AND "TABLE_A"."id"=?` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{1}) {
		t.Error("got", attrs)
	}
	if err != nil {
		t.Error("got", err)
	}
}

func TestOrCondition(t *testing.T) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)
	eq1 := Eq(table1.C("id"), table1.C("test1"))
	eq2 := EqL(table1.C("id"), 1)

	or := Or(eq1, eq2)
	query, attrs, err := or.toSql()
	if query != `"TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=?` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{1}) {
		t.Error("got", attrs)
	}
	if err != nil {
		t.Error("got", err)
	}
}
