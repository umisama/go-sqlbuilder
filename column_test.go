package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnImplements(t *testing.T) {
	a := assert.New(t)
	a.Implements(new(Column), &columnImpl{})
	a.Implements(new(Column), &errorColumn{})
	a.Implements(new(Column), &aliasColumn{})
}

func TestColumnOptionImpl(t *testing.T) {
	a := assert.New(t)

	a.Equal(&columnConfigImpl{
		name: "name",
		typ:  ColumnTypeBytes,
		opt: &ColumnOption{
			Unique: true,
		},
	}, newColumnConfigImpl("name", ColumnTypeBytes, &ColumnOption{
		Unique: true,
	}))
	a.Equal(&columnConfigImpl{
		name: "name",
		typ:  ColumnTypeBytes,
		opt:  &ColumnOption{},
	}, newColumnConfigImpl("name", ColumnTypeBytes, nil))
}
