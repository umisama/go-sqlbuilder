package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlterTable(t *testing.T) {
	a := assert.New(t)
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	type testcase struct {
		stmt  Statement
		query string
		args  []interface{}
		err   bool
	}
	var cases = []testcase{{
		AlterTable(table1).
			RenameTo("TABLE_AAA").
			AddColumn(IntColumn("test3", nil)).
			AddColumn(IntColumn("test4", nil)).
			ChangeColumn(table1.C("test1"), IntColumn("test1a", nil)).
			DropColumn(table1.C("test1")),
		`ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER, ADD COLUMN "test4" INTEGER, CHANGE COLUMN "test1" "test1a" INTEGER, DROP COLUMN "test1", RENAME TO "TABLE_AAA";`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).RenameTo("TABLE_AAA"),
		`ALTER TABLE "TABLE_A" RENAME TO "TABLE_AAA";`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).AddColumn(IntColumn("test3", &ColumnOption{
			Unique: true,
		})),
		`ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER UNIQUE;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).AddColumnAfter(IntColumn("test0", nil), table1.C("id")),
		`ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER AFTER "id";`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).AddColumnFirst(IntColumn("test0", nil)),
		`ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER FIRST;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).ChangeColumn(table1.C("test1"), IntColumn("test1a", &ColumnOption{
			Unique: true,
		})),
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER UNIQUE;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).ChangeColumnFirst(table1.C("test1"), IntColumn("test1a", nil)),
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER FIRST;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).ChangeColumnAfter(table1.C("test1"), IntColumn("test1a", nil), table1.C("test2")),
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER AFTER "test2";`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).DropColumn(table1.C("test1")),
		`ALTER TABLE "TABLE_A" DROP COLUMN "test1";`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).DropColumn(table1.C("invalid")),
		``,
		[]interface{}{},
		true,
	}, {
		AlterTable(table1).ChangeColumnAfter(table1.C("invalid"), IntColumn("test1a", nil), table1.C("test2")),
		``,
		[]interface{}{},
		true,
	}, {
		AlterTable(table1).ChangeColumnAfter(table1.C("test1"), IntColumn("test1a", nil), table1.C("invalid")),
		``,
		[]interface{}{},
		true,
	}, {
		AlterTable(table1).AddColumnAfter(IntColumn("test0", nil), table1.C("invalid")),
		``,
		[]interface{}{},
		true,
	}}

	for _, c := range cases {
		query, args, err := c.stmt.ToSql()
		a.Equal(c.query, query)
		a.Equal(c.args, args)
		if c.err {
			a.Error(err)
		} else {
			a.NoError(err)
		}
	}
}
