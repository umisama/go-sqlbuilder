package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {
	a := assert.New(t)

	table1, err := NewTable(
		"TABLE_NAME",
	)
	a.Error(err)
	a.Nil(table1)

	table1, err = NewTable(
		"TABLE_NAME",
		IntColumn("id", true),
	)
	a.Nil(err)
	a.NotNil(table1)
}
