package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLockClauseImplements(t *testing.T) {
	a := assert.New(t)
	a.Implements(new(LockClause), &lockClauseImpl{})
}

func TestLockClause(t *testing.T) {
	a := assert.New(t)
	b := newBuilder()
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)
	table2 := NewTable(
		"TABLE_B",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)

	Lock("UPDATE", table1, table2).serialize(b)
	a.Equal(`UPDATE OF "TABLE_A", "TABLE_B"`, b.query.String())
	a.Equal([]interface{}{}, b.Args())
	a.NoError(b.Err())
}
