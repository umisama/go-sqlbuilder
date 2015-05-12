# umisama/go-sqlbuilder
**go-sqlbuilder** is a SQL-query builder for golang.  This supports you using relational database with more readable and flexible code than raw SQL query string.

[![Build Status](https://travis-ci.org/umisama/go-sqlbuilder.svg?branch=master)](https://travis-ci.org/umisama/go-sqlbuilder)
[![Coverage Status](https://coveralls.io/repos/umisama/go-sqlbuilder/badge.svg)](https://coveralls.io/r/umisama/go-sqlbuilder)

## Status
!!!SUPER ALPHA!!!

## Support
 * Generate SQL query programmatically.
   * fluent flexibility! yeah!!
 * Basic SQL statements
   * SELECT/INSERT/UPDATE/DELETE/DROP/CREATE TABLE/CREATE INDEX
 * Strict error checking
 * Some database server
   * Sqlite3([mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))
   * MySQL([ziutek/mymysql](https://github.com/ziutek/mymysql))
   * MySQL([go-sql-driver/mysql](https://github.com/go-sql-driver/mysql))
   * PostgresSQL([lib/pq](https://github.com/lib/pq))
 * Subquery in SELECT FROM clause

## TODO
 * Support union

## Quick start
working on

## Examples
### Initialize
off course, go getable.

```shell-script
$ go get github.com/umisama/go-sqlbuilder
```

I recomended to set "sb" as sqlbuilder's shorthand.

```go
import sb "github.com/umisama/go-sqlbuilder"

// First, you set dialect for your DB
func init (
	sb.SetDialect(sb.SqliteDialect{})
)
```

### Define a table
Sqlbuilder needs table definition to strict query generating.

```go
table1  := sb.NewTable(
	"TABLE_A",
	sb.IntColumn("id", &sb.ColumnOption{
		PrimaryKey: true,
	}),
	sb.StrColumn("name", nil)
	sb.IntColumn("age", nil),
)
```

### CRATE TABLE statement
Sqlbuilder can generate CREATE TABLE statement from table object.  
Statement objects have ```ToSql()``` method.  ```Statement.ToSql()``` returns query(string), placeholder arguments([]interface{}) and error(error).

```go
query, args, err := sb.CreateTable(table1).ToSql()
if err != nil {
	panic(err)
}
// query == `CREATE TABLE "TABLE_A" ( "id" INTEGER PRIMARY KEY, "value" INTEGER );`
// args  == []interface{}{}
// err   == nil
```

You can exec with ```database/sql``` package or Table-struct mapper(for example, gorp).  
here is example,

```go
db, err := sql.Open("sqlite3", ":memory:")
if err != nil {
	panic(err)
}
_, err = db.Exec(query, args...)
if err != nil {
	panic(err)
}
```

### INSERT statement
Sqlbuilder can generate INSERT statement.  You can checkout a column with ```*Table.C([column_name])``` method.

```go
query, args, err := sb.Insert(table1).
	Columns(table1.C("id"), table1.C("value")).
	Values(1, 10).
	ToSql()
// query == `INSERT INTO "TABLE_A" ( "id", "value" ) VALUES ( ?, ? );`
// args  == []interface{}{1, 10}
// err   == nil
```

### SELECT statement
Sqlbuilder can generate SELECT statement with readable interfaces.  Condition object is generated from column object.

```go
query, args, err := sb.Select(table1.C("id"), table1.C("value")).
	From(table1).
	Where(
		table1.C("id").Eq(10),
	).
	Limit(1).OrderBy(false, table1.C("id")).
	ToSql()
// query == `SELECT "TABLE_A"."id", "TABLE_A"."value" FROM "TABLE_A" WHERE "TABLE_A"."id"=? ORDER BY "TABLE_A"."id" ASC LIMIT ?;`
// args  == []interface{}{10, 1}
// err   == nil
```

### Condition clause
You can define condition with Condition objects.  Condition object create from ```Column```'s method.

| example operation                     |  output means              |
|:-------------------------------------:|:--------------------------:|
|```table1.C("id").Eq(10)```              | "TABLE1"."id"=10           |
|```table1.C("id").Eq(table2.C("id"))```    | "TABLE1"."id"="TABLE2"."id"|

More than one condition can combine with AND & OR operator.

| example operation                     |  output means              |
|:-------------------------------------:|:--------------------------:|
|```And(table1.C("id").Eq(1), table2.C("id").Eq(2)``` | "TABLE1"."id"=1 AND "TABLE2"."id"=1 |
|```Or(table1.C("id").Eq(1), table2.C("id").Eq(2)```  | "TABLE1"."id"=1 OR "TABLE2"."id"=1 |

Sqlbuilder is supporting most common condition operators.  
Here is now supporting...

| columns method        |SQL operator|      SQL results     |
|:---------------------:|:----------:|:--------------------:|
|Eq(Column or value)    | ```=```    | "TABLE"."id" = 10    |
|NotEq(Column or value) | ```<>```   | "TABLE"."id" <> 10   |
|Gt(Column or value)    | ```>```    | "TABLE"."id" > 10    |
|GtEq(Column or value)  | ```>=```   | "TABLE"."id" >= 10   |
|Lt(Column or value)    | ```<```    | "TABLE"."id" < 10    |
|LtEq(Column or value)  | ```<=```   | "TABLE"."id" <= 10   |
|Like(string)           | ```LIKE``` | "TABLE"."id" LIKE "%hoge%"   |
|In(values array)       | ```IN``` | "TABLE"."id" IN ( 1, 2, 3 ) |
|Between(loewer, higher int | ```BETWEEN``` | "TABLE"."id" BETWEEN 10 AND 20)|

Document for all: [godoc(Column)](http://godoc.org/github.com/umisama/go-sqlbuilder#Column)

## More documents
[godoc.org](http://godoc.org/github.com/umisama/go-sqlbuilder)

## License
under the MIT license
