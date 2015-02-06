package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlFunc(t *testing.T) {
	a := assert.New(t)
	b := newBuilder()
	table1 := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)

	Func("funcname", table1.C("id")).serialize(b)
	a.Equal(`funcname("TABLE_A"."id")`, b.Query())
	a.Equal([]interface{}{}, b.Args())
	a.NoError(b.Err())
}
