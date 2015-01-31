package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	a := assert.New(t)

	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	query, attrs, err := Insert(table1).
		Columns(table1.C("id"), table1.C("test1"), table1.C("test2")).
		Values(L(1), L(2), L(3)).ToSql()
	a.Equal(`INSERT ( "TABLE_A"."id" "TABLE_A"."test1" "TABLE_A"."test2" ) `+
		`INTO "TABLE_A" VALUES ( ? ? ? );`, query)
	a.Equal([]interface{}{1, 2, 3}, attrs)
	a.Nil(err)

}
