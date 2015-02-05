package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	a := assert.New(t)
	err := newError("hogehogestring")
	a.Equal("sqlbuilder: hogehogestring", err.Error())
}
