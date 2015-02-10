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
