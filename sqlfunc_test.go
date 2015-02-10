package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlFuncImplements(t *testing.T) {
	a := assert.New(t)
	a.Implements(new(Column), &columnImpl{})
}

func TestSqlFunc(t *testing.T) {
	a := assert.New(t)
	b := newBuilder()
	table1 := NewTable(
		"TABLE_A",
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	Func("funcname", table1.C("id")).serialize(b)
	a.Equal(`funcname("TABLE_A"."id")`, b.Query())
	a.Equal([]interface{}{}, b.Args())
	a.NoError(b.Err())
}
