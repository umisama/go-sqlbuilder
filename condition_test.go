package sqlbuilder

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBinaryCondition(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)

	type testcase struct {
		cond  Condition
		query string
		attrs []interface{}
		err   error
	}

	var cases = []testcase{
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
			[]interface{}{1, 2},
			nil,
		}, {
			table1.C("id").In(1, 2),
			`"TABLE_A"."id" IN ( ?, ? )`,
			[]interface{}{1, 2},
			nil,
		}, {
			// case for fail
			table1.C("id").In(NewTable("DUMMY TABLE")),
			``,
			[]interface{}{},
			errors.New("sqlbuilder: unsupported type"),
		},
	}

	for _, c := range cases {
		bldr := newBuilder()
		c.cond.serialize(bldr)
		a.Equal(c.query, bldr.Query())
		a.Equal(c.attrs, bldr.Args())
		a.Equal(c.err, bldr.Err())
	}

}

func TestBinaryConditionForSqlFunctions(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)

	type testcase struct {
		cond  Condition
		query string
		attrs []interface{}
		err   error
	}

	var cases = []testcase{
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
			[]interface{}{1, 2},
			nil,
		}, {
			Func("count", table1.C("id")).In(1, 2),
			`count("TABLE_A"."id") IN ( ?, ? )`,
			[]interface{}{1, 2},
			nil,
		}, {
			// case for fail
			Func("count", table1.C("id")).In(NewTable("DUMMY TABLE")),
			``,
			[]interface{}{},
			errors.New("sqlbuilder: unsupported type"),
		},
	}

	for _, c := range cases {
		bldr := newBuilder()
		c.cond.serialize(bldr)
		a.Equal(c.query, bldr.Query())
		a.Equal(c.attrs, bldr.Args())
		a.Equal(c.err, bldr.Err())
	}

}

func TestAndCondition(t *testing.T) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)
	eq1 := table1.C("id").Eq(table1.C("test1"))
	eq2 := table1.C("id").Eq(1)
	eq3 := table1.C("id").Eq(2)

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
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)
	eq1 := table1.C("id").Eq(table1.C("test1"))
	eq2 := table1.C("id").Eq(1)

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
