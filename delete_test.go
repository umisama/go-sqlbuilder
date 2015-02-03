package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	query, args, err := Delete(table1).Where(table1.C("id").Eq(1)).ToSql()
	a.Equal(`DELETE FROM "TABLE_A" WHERE "TABLE_A"."id"=?;`, query)
	a.Equal([]interface{}{1}, args)
	a.Nil(err)
}
