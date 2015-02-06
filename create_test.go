package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	a := assert.New(t)
	table1 := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey, CO_AutoIncrement),
		IntColumn("test1", CO_Unique),
		IntColumn("test2"),
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
		`CREATE TABLE IF NOT EXISTS "TABLE_A" ( "id" INTEGER PRIMARY KEY AUTO INCREMENT, "test1" INTEGER UNIQUE, "test2" INTEGER );`,
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
			a.NotNil(err)
		} else {
			a.Nil(err)
		}
	}
}
