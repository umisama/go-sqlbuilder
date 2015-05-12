package sqlbuilder

import (
	"testing"
)

func TestAlterTable(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	table2 := NewTable(
		"TABLE_B",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)
	tableJoined := table1.InnerJoin(table2, table1.C("test1").Eq(table2.C("id")))

	var cases = []statementTestCase{{
		stmt: AlterTable(table1).
			RenameTo("TABLE_AAA").
			AddColumn(IntColumn("test3", nil)).
			AddColumn(IntColumn("test4", nil)).
			ChangeColumn(table1.C("test1"), IntColumn("test1a", nil)).
			DropColumn(table1.C("test1")),
		query:  `ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER, ADD COLUMN "test4" INTEGER, CHANGE COLUMN "test1" "test1a" INTEGER, DROP COLUMN "test1", RENAME TO "TABLE_AAA";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   AlterTable(table1).RenameTo("TABLE_AAA"),
		query:  `ALTER TABLE "TABLE_A" RENAME TO "TABLE_AAA";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt: AlterTable(table1).AddColumn(IntColumn("test3", &ColumnOption{
			Unique: true,
		})),
		query:  `ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER UNIQUE;`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   AlterTable(table1).AddColumnAfter(IntColumn("test0", nil), table1.C("id")),
		query:  `ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER AFTER "id";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   AlterTable(table1).AddColumnFirst(IntColumn("test0", nil)),
		query:  `ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER FIRST;`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt: AlterTable(table1).ChangeColumn(table1.C("test1"), IntColumn("test1a", &ColumnOption{
			Unique: true,
		})),
		query:  `ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER UNIQUE;`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   AlterTable(table1).ChangeColumnFirst(table1.C("test1"), IntColumn("test1a", nil)),
		query:  `ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER FIRST;`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   AlterTable(table1).ChangeColumnAfter(table1.C("test1"), IntColumn("test1a", nil), table1.C("test2")),
		query:  `ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER AFTER "test2";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   AlterTable(table1).DropColumn(table1.C("test1")),
		query:  `ALTER TABLE "TABLE_A" DROP COLUMN "test1";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   AlterTable(table1).DropColumn(table1.C("invalid")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: column TABLE_A.invalid was not found.",
	}, {
		stmt:   AlterTable(table1).ChangeColumnAfter(table1.C("invalid"), IntColumn("test1a", nil), table1.C("test2")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: column TABLE_A.invalid was not found.",
	}, {
		stmt:   AlterTable(table1).ChangeColumnAfter(table1.C("test1"), IntColumn("test1a", nil), table1.C("invalid")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: column TABLE_A.invalid was not found.",
	}, {
		stmt:   AlterTable(table1).AddColumnAfter(IntColumn("test0", nil), table1.C("invalid")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: column TABLE_A.invalid was not found.",
	}, {
		stmt:   AlterTable(nil).DropColumn(table1.C("invalid")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt:   AlterTable(tableJoined).DropColumn(table1.C("id")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: AlterTable can use only natural table.",
	}}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}

func TestAlterTableApplyToTable(t *testing.T) {
	var cases = []struct {
		stmt           func(Table) *AlterTableStatement
		expect_columns []string
		expect_name    string
	}{{
		stmt: func(t Table) *AlterTableStatement {
			return AlterTable(t).
				RenameTo("TABLE_AAA").
				AddColumn(IntColumn("test3", nil)).
				AddColumnFirst(IntColumn("test4", nil)).
				AddColumnAfter(IntColumn("test5", nil), t.C("id")).
				ChangeColumn(t.C("test1"), IntColumn("test1a", nil)).
				ChangeColumnFirst(t.C("test2"), IntColumn("test2a", nil)).
				DropColumn(t.C("id"))
		},
		expect_columns: []string{"test2a", "test4", "test5", "test1a", "test3"},
		expect_name:    "TABLE_AAA",
	}, {
		stmt: func(t Table) *AlterTableStatement {
			return AlterTable(t).
				ChangeColumnAfter(t.C("test1"), IntColumn("test1a", nil), t.C("test2"))
		},
		expect_columns: []string{"id", "test2", "test1a"},
		expect_name:    "TABLE_A",
	}}

	for num, c := range cases {
		table1 := NewTable(
			"TABLE_A",
			&TableOption{},
			IntColumn("id", &ColumnOption{
				PrimaryKey: true,
			}),
			IntColumn("test1", nil),
			IntColumn("test2", nil),
		)

		err := c.stmt(table1).ApplyToTable()
		if err != nil {
			t.Errorf("failed on %d", num)
		}
		if len(table1.Columns()) != len(c.expect_columns) {
			t.Errorf("failed on %d", num)
		}
		for i, col := range table1.Columns() {
			if c.expect_columns[i] != col.column_name() {
				t.Errorf("failed on %d", num)
				break
			}
		}
		if table1.Name() != c.expect_name {
			t.Errorf("failed on %d", num)
		}
	}
}
