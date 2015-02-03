package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlFunc(t *testing.T) {
	a := assert.New(t)
	b := newBuilder()
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", false),
		IntColumn("test1", false),
		IntColumn("test2", false),
	)

	Func("funcname", table1.C("id")).serialize(b)
	a.Equal(`funcname("TABLE_A"."id")`, b.Query())
	a.Equal([]interface{}{}, b.Args())
	a.Nil(b.Err())
}
