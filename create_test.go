package sqlbuilder

import (
	"testing"
)

func TestCreate(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey:    true,
			AutoIncrement: true,
		}),
		IntColumn("test1", &ColumnOption{
			Unique: true,
		}),
		StringColumn("test2", &ColumnOption{
			Size: 255,
		}),
	)
	table2 := NewTable(
		"TABLE_B",
		&TableOption{},
		StringColumn("id", &ColumnOption{
			PrimaryKey:    true,
			AutoIncrement: true,
			SqlType:       "VARCHAR(255)",
		}),
		AnyColumn("test1", &ColumnOption{
			Unique:  true,
			SqlType: "INTEGER",
		}),
	)
	table3 := NewTable(
		"TABLE_C",
		&TableOption{
			Unique: [][]string{{"test1", "test2"}},
		},
		IntColumn("id", &ColumnOption{
			PrimaryKey:    true,
			AutoIncrement: true,
		}),
		IntColumn("test1", &ColumnOption{
			Unique: true,
		}),
		StringColumn("test2", &ColumnOption{
			Size: 255,
		}),
	)
	tableJoined := table1.InnerJoin(table2, table1.C("test1").Eq(table2.C("id")))
	tableZeroColumns := &table{
		name:    "ZERO_TABLE",
		columns: make([]Column, 0),
	}

	var cases = []statementTestCase{{
		CreateTable(table1).IfNotExists(),
		`CREATE TABLE IF NOT EXISTS "TABLE_A" ( "id" INTEGER PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE, "test2" TEXT );`,
		[]interface{}{},
		false,
	}, {
		CreateTable(table2).IfNotExists(),
		`CREATE TABLE IF NOT EXISTS "TABLE_B" ( "id" VARCHAR(255) PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE );`,
		[]interface{}{},
		false,
	}, {
		CreateTable(table3).IfNotExists(),
		`CREATE TABLE IF NOT EXISTS "TABLE_C" ( "id" INTEGER PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE, "test2" TEXT ) UNIQUE("test1", "test2");`,
		[]interface{}{},
		false,
	}, {
		CreateTable(tableZeroColumns),
		``,
		[]interface{}{},
		true,
	}, {
		CreateTable(nil),
		``,
		[]interface{}{},
		true,
	}, {
		CreateTable(tableJoined),
		``,
		[]interface{}{},
		true,
	}, {
		CreateIndex(table1).Name("I_TABLE_A").IfNotExists().Columns(table1.C("test1"), table1.C("test2")),
		`CREATE INDEX IF NOT EXISTS "I_TABLE_A" ON "TABLE_A" ( "test1", "test2" );`,
		[]interface{}{},
		false,
	}, {
		CreateIndex(table1).Columns(table1.C("test1"), table1.C("test2")),
		``,
		[]interface{}{},
		true,
	}, {
		CreateIndex(table1).Name("I_TABLE_A"),
		``,
		[]interface{}{},
		true,
	}, {
		CreateIndex(nil).Name("I_TABLE_A"),
		``,
		[]interface{}{},
		true,
	}, {
		CreateIndex(tableJoined).Name("I_TABLE_A"),
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
