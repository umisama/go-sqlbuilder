# umisama/go-sqlbuilder
[![Build Status](https://travis-ci.org/umisama/go-sqlbuilder.svg?branch=master)](https://travis-ci.org/umisama/go-sqlbuilder)

## Status
development...

## Motivation
 * Build SQL query in less limitations.
 * Flexible API.
 * (will be my ORM base...)

## Concepts
### initialize

```go
import sb "github.com/umisama/go-sqlbuilder"

func init() {
	sb.Setup(SqliteDialect{})
}
```

### define table

```go
table1 := sb.NewTable(
	"TABLE_NAME",
	sb.IntColumn("id", false),
	sb.StrColumn("name", []sb.ColumnOptions{
		sb.UTF8,
		sb.UTF8CaseInsensitive,
		sb.StrColumnSize(255),
		},
		false)
)

query, attrs, err := table1.ToSql()
// err == nil
// attrs == []interface{}{}
// query == `CREATE TABLE "TABLE_NAME" ("id" INTEGER, "name" VARCHAR(255))`
```

### query

```go
query, attrs, err := Select("*").From(table1).Where(
	sb.Eq(table1.Column("id"), 1),
	).ToSql()
// err == nil
// attrs == []interface{}{1}
// query == `SELECT * FROM "TABLE_NAME" WHERE "TABLE_NAME"."id" = ?`
```
