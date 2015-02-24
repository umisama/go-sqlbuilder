package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	a := assert.New(t)
	table1 := NewTable(
		"TABLE_A",
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
	table_zero_columns := &table{
		name:    "ZERO_TABLE",
		columns: make([]Column, 0),
	}

	type testcase struct {
		stmt  Statement
		query string
		args  []interface{}
		err   bool
	}
	var cases = []testcase{{
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
		CreateTable(table_zero_columns),
		``,
		[]interface{}{},
		true,
	}, {
		CreateTable(nil),
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
