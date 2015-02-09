package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {
	a := assert.New(t)

	var table1 Table
	a.Panics(func() {
		table1 = NewTable(
			"TABLE_NAME",
		)
	})

	a.NotPanics(func() {
		table1 = NewTable(
			"TABLE_NAME",
			IntColumn("id"),
		)
	})
	a.NotNil(table1)
}

func TestJoinTable(t *testing.T) {
	a := assert.New(t)

	l_table := NewTable(
		"LEFT_TABLE",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("right_id"),
	)
	r_table := NewTable(
		"RIGHT_TABLE",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("value"),
	)
	rr_table := NewTable(
		"RIGHTRIGHT_TABLE",
		IntColumn("id", CO_PrimaryKey),
	)

	// inner join
	b := newBuilder()
	joinedTable := l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.Query())
	a.NoError(b.Err())
	a.Empty(b.Args())

	// left outer join
	b = newBuilder()
	joinedTable = l_table.LeftOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" LEFT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.Query())
	a.NoError(b.Err())
	a.Empty(b.Args())

	// right outer join
	b = newBuilder()
	joinedTable = l_table.RightOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" RIGHT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.Query())
	a.NoError(b.Err())
	a.Empty(b.Args())

	// full outer join
	b = newBuilder()
	joinedTable = l_table.FullOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" FULL OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.Query())
	a.NoError(b.Err())
	a.Empty(b.Args())

	// joined table column
	a.Equal(l_table.C("right_id"), joinedTable.C("right_id"))
	a.Equal(r_table.C("value"), joinedTable.C("value"))
	a.IsType(&errorColumn{}, joinedTable.C("not_exist_column"))
	a.IsType(&errorColumn{}, joinedTable.C("id"))

	// combination
	b = newBuilder()
	joinedTable = l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id"))).InnerJoin(rr_table, l_table.C("right_id").Eq(rr_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id" INNER JOIN "RIGHTRIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHTRIGHT_TABLE"."id"`, b.Query())
	a.NoError(b.Err())
	a.Empty(b.Args())
}
