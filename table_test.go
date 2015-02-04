package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {
	a := assert.New(t)

	table1, err := NewTable(
		"TABLE_NAME",
	)
	a.Error(err)
	a.Nil(table1)

	table1, err = NewTable(
		"TABLE_NAME",
		IntColumn("id"),
	)
	a.Nil(err)
	a.NotNil(table1)
}

func TestJoinTable(t *testing.T) {
	a := assert.New(t)
	b := newBuilder()

	l_table, _ := NewTable(
		"LEFT_TABLE",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("right_id"),
	)
	r_table, _ := NewTable(
		"RIGHT_TABLE",
		IntColumn("id", CO_PrimaryKey),
	)
	rr_table, _ := NewTable(
		"RIGHTRIGHT_TABLE",
		IntColumn("id", CO_PrimaryKey),
	)

	joinedTable := l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id"))).InnerJoin(rr_table, l_table.C("right_id").Eq(rr_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(b.Query(), `"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id" INNER JOIN "RIGHTRIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHTRIGHT_TABLE"."id"`)
	a.Nil(b.Err())
	a.Empty(b.Args())
}
