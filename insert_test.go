package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)

	query, attrs, err := Insert(table1).
		Columns(table1.C("id"), table1.C("test1"), table1.C("test2")).
		Values(1, 2, 3).ToSql()
	a.Equal(`INSERT INTO "TABLE_A" ( "id", "test1", "test2" ) VALUES ( ?, ?, ? );`, query)
	a.Equal([]interface{}{1, 2, 3}, attrs)
	a.Nil(err)

}
