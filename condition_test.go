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
		bldr := newBuilder()
		c.cond.serialize(bldr)
		if bldr.Query() != c.query {
			t.Error("got:", bldr.Query(), " case:", i)
		}
		if !reflect.DeepEqual(bldr.Args(), c.attrs) {
			t.Error("got:", bldr.Args(), " case:", i)
		}
		if bldr.Err() != c.err {
			t.Error("got:", bldr.Err(), " case:", i)
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

	bldr := newBuilder()
	and.serialize(bldr)
	if bldr.Query() != `"TABLE_A"."id"="TABLE_A"."test1" AND "TABLE_A"."id"=? AND "TABLE_A"."id"=?` {
		t.Error("got", bldr.Query())
	}
	if !reflect.DeepEqual(bldr.Args(), []interface{}{int64(1), int64(2)}) {
		t.Error("got", bldr.Args())
	}
	if bldr.Err() != nil {
		t.Error("got", bldr.Err())
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

	b := newBuilder()
	or.serialize(b)
	if b.Query() != `"TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=?` {
		t.Error("got", b.Query())
	}
	if !reflect.DeepEqual(b.Args(), []interface{}{int64(1)}) {
		t.Error("got", b.Args())
	}
	if b.Err() != nil {
		t.Error("got", b.Err())
	}
}
