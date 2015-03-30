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
		`ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER DEFAULT NULL, ADD COLUMN "test4" INTEGER DEFAULT NULL, CHANGE COLUMN "test1" "test1a" INTEGER DEFAULT NULL, DROP COLUMN "test1", RENAME TO "TABLE_AAA";`,
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
		`ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER UNIQUE DEFAULT NULL;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).AddColumnAfter(IntColumn("test0", nil), table1.C("id")),
		`ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER DEFAULT NULL AFTER "id";`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).AddColumnFirst(IntColumn("test0", nil)),
		`ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER DEFAULT NULL FIRST;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).ChangeColumn(table1.C("test1"), IntColumn("test1a", &ColumnOption{
			Unique: true,
		})),
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER UNIQUE DEFAULT NULL;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).ChangeColumnFirst(table1.C("test1"), IntColumn("test1a", nil)),
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER DEFAULT NULL FIRST;`,
		[]interface{}{},
		false,
	}, {
		AlterTable(table1).ChangeColumnAfter(table1.C("test1"), IntColumn("test1a", nil), table1.C("test2")),
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER DEFAULT NULL AFTER "test2";`,
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
	}, {
		AlterTable(table1.InnerJoin(table1, table1.C("id").Eq(table1.C("id")))).AddColumnAfter(IntColumn("test0", nil), table1.C("invalid")),
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

func TestAlterTableApplyToTable(t *testing.T) {
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

	stmt := AlterTable(table1).
		RenameTo("TABLE_AAA").
		AddColumn(IntColumn("test3", nil)).
		AddColumnFirst(IntColumn("test4", nil)).
		AddColumnAfter(IntColumn("test5", nil), table1.C("id")).
		ChangeColumn(table1.C("test1"), IntColumn("test1a", nil)).
		ChangeColumnFirst(table1.C("test2"), IntColumn("test2a", nil)).
		DropColumn(table1.C("id"))
	err := stmt.ApplyToTable()
	a.NoError(err)

	expect := []string{"test2a", "test4", "test5", "test1a", "test3"}
	a.Len(table1.Columns(), len(expect))
	for i, col := range table1.Columns() {
		a.Equal(expect[i], col.column_name())
	}
	a.Equal(table1.Name(), "TABLE_AAA")

	table1 = NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	stmt = AlterTable(table1).
		ChangeColumnAfter(table1.C("test1"), IntColumn("test1a", nil), table1.C("test2"))
	err = stmt.ApplyToTable()
	a.NoError(err)
	expect = []string{"id", "test2", "test1a"}
	a.Len(table1.Columns(), len(expect))
	for i, col := range table1.Columns() {
		a.Equal(expect[i], col.column_name())
	}
}
