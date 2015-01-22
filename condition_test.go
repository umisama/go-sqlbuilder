package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestBinaryCondition(t *testing.T) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	type testcase struct {
		cond  Condition
		query string
		attrs []interface{}
		err   error
	}

	var cases = []testcase{
		{
			Eq(table1.C("id"), table1.C("test1")),
			`"TABLE_A"."id"="TABLE_A"."test1"`,
			[]interface{}{},
			nil,
		}, {
			Eq(table1.C("id"), L(1)),
			`"TABLE_A"."id"=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			NotEq(table1.C("id"), L(1)),
			`"TABLE_A"."id"<>?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Gt(table1.C("id"), L(1)),
			`"TABLE_A"."id">?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Gte(table1.C("id"), L(1)),
			`"TABLE_A"."id">=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Lt(table1.C("id"), L(1)),
			`"TABLE_A"."id"<?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Lte(table1.C("id"), L(1)),
			`"TABLE_A"."id"<=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Lte(table1.C("id"), L(1)),
			`"TABLE_A"."id"<=?`,
			[]interface{}{int64(1)},
			nil,
		}, {
			Like(table1.C("id"), L("hoge")),
			`"TABLE_A"."id" LIKE ?`,
			[]interface{}{"hoge"},
			nil,
		},
	}

	for i, c := range cases {
		query, attrs, err := c.cond.serialize()
		if query != c.query {
			t.Error("got:", query, " case:", i)
		}
		if !reflect.DeepEqual(attrs, c.attrs) {
			t.Error("got:", attrs, " case:", i)
		}
		if err != c.err {
			t.Error("got:", err, " case:", i)
		}
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
	eq2 := Eq(table1.C("id"), L(1))
	eq3 := Eq(table1.C("id"), L(2))

	and := And(eq1, eq2, eq3)
	query, attrs, err := and.serialize()
	if query != `"TABLE_A"."id"="TABLE_A"."test1" AND "TABLE_A"."id"=? AND "TABLE_A"."id"=?` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{int64(1), int64(2)}) {
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
	eq2 := Eq(table1.C("id"), L(1))

	or := Or(eq1, eq2)
	query, attrs, err := or.serialize()
	if query != `"TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=?` {
		t.Error("got", query)
	}
	if !reflect.DeepEqual(attrs, []interface{}{int64(1)}) {
		t.Error("got", attrs)
	}
	if err != nil {
		t.Error("got", err)
	}
}
