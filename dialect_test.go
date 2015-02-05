package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// tests for dialects in ./integration_test

func TestDialectImplements(t *testing.T) {
	a := assert.New(t)
	a.Implements(new(Dialect), SqliteDialect{})
	a.Implements(new(Dialect), MysqlDialect{})
	a.Implements(new(Dialect), PostgresDialect{})
}
