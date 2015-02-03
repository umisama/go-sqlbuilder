package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	query, args, err := CreateTable(table1).IfNotExists().ToSql()
	a.Equal(`CREATE TABLE IF NOT EXISTS "TABLE_A" ( "id" INTEGER, "test1" INTEGER, "test2" INTEGER );`, query)
	a.Equal([]interface{}{}, args)
	a.Nil(err)
}

func TestCreateIndex(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	query, args, err := CreateIndex(table1).Name("I_TABLE_A").IfNotExists().Columns(table1.C("test1"), table1.C("test2")).ToSql()
	a.Equal(`CREATE INDEX IF NOT EXISTS "I_TABLE_A" ON "TABLE_A" ( "test1", "test2" );`, query)
	a.Equal([]interface{}{}, args)
	a.Nil(err)
}
