package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	query, attrs, err := Update(table1).
		Set(table1.C("test1"), 10).
		Set(table1.C("test2"), 20).
		Where(table1.C("id").Eq(1)).
		ToSql()
	a.Equal(`UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=?;`, query)
	a.Equal([]interface{}{int64(10), int64(20), int64(1)}, attrs)
	a.Nil(err)

}
