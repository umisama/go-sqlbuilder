package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	query, attrs, err := Select(table1.C("test1"), table1.C("test2")).
		From(table1).
		Where(
		And(
			Eq(table1.C("id"), L(1)),
			Eq(table1.C("test1"), L(2)),
		)).
		Distinct().
		OrderBy(false, table1.C("id")).
		GroupBy(table1.C("id")).
		Having(Eq(table1.C("id"), L(1))).
		Limit(10).
		Offset(20).
		ToSql()
	a.Equal(`SELECT DISTINCT "TABLE_A"."test1", "TABLE_A"."test2" `+
		`FROM "TABLE_A" WHERE "TABLE_A"."id"=? AND "TABLE_A"."test1"=? `+
		`GROUP BY "TABLE_A"."id" HAVING "TABLE_A"."id"=? ORDER BY "TABLE_A"."id" ASC `+
		`LIMIT 10 OFFSET 20;`, query)
	a.Equal([]interface{}{int64(1), int64(2), int64(1)}, attrs)
	a.Nil(err)
}

func BenchmarkSelect(b *testing.B) {
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	for i := 0; i < b.N; i++ {
		Select(table1.C("test1"), table1.C("test2")).
			From(table1).
			Where(
			And(
				Eq(table1.C("id"), L(1)),
				Eq(table1.C("test1"), L(2)),
			)).
			ToSql()
	}
}
