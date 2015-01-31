package sqlbuilder_integration_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	sb "github.com/umisama/go-sqlbuilder"
	_ "github.com/ziutek/mymysql/godrv"
)

var table1 sb.Table
var db *sql.DB

func TestMain(m *testing.M) {
	results := make(map[string]int)
	type testcase struct {
		name    string
		dialect sb.Dialect
		driver  string
		dsn     string
	}
	var cases = []testcase{
		{"sqlite", sb.SqliteDialect{}, "sqlite3", ":memory:"},
		{"mymysql", sb.MysqlDialect{}, "mymysql", "go_sqlbuilder_test/root/"},
	}

	table1, _ = sb.NewTable(
		"TABLE_A",
		sb.IntColumn("id", true),
		sb.IntColumn("value", true),
	)

	for _, c := range cases {
		fmt.Println("START unit test for", c.name)

		var err error
		db, err = sql.Open(c.driver, c.dsn)
		if err != nil {
			fmt.Println(err.Error())
		}
		sb.SetDialect(c.dialect)
		_, err = db.Exec("CREATE TABLE `TABLE_A` (`id` integer primary key, `value` integer)")
		if err != nil {
			fmt.Println(err.Error())
		}
		results[c.name] = m.Run()
		_, err = db.Exec("DROP TABLE `TABLE_A`")
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for _, v := range results {
		if v != 0 {
			os.Exit(v)
		}
	}
	os.Exit(0)
}

func TestInsert(t *testing.T) {
	a := assert.New(t)

	// data 1
	query, args, err := sb.Insert(table1).
		Columns(table1.C("id"), table1.C("value")).
		Values(sb.L(1), sb.L(10)).
		ToSql()
	a.Nil(err)

	result, err := db.Exec(query, args...)
	a.Nil(err)

	if a.NotNil(result) {
		rows_affected, err := result.RowsAffected()
		a.Equal(1, rows_affected)
		a.Nil(err)
	}

	// data 2
	query, args, err = sb.Insert(table1).
		Columns(table1.C("id"), table1.C("value")).
		Values(sb.L(2), sb.L(20)).
		ToSql()
	a.Nil(err)

	result, err = db.Exec(query, args...)
	a.Nil(err)

	if a.NotNil(result) {
		rows_affected, err := result.RowsAffected()
		a.Equal(1, rows_affected)
		a.Nil(err)
	}
}

func TestSelect(t *testing.T) {
	a := assert.New(t)
	query, args, err := sb.Select(table1.C("id"), table1.C("value")).
		From(table1).Where(
		sb.Eq(
			table1.C("id"), sb.L(1),
		),
	).Limit(1).OrderBy(false, table1.C("id")).ToSql()

	rows, err := db.Query(query, args...)
	a.Nil(err)

	if a.NotNil(rows) {
		count := 0
		for rows.Next() {
			id, value := 0, 0
			err := rows.Scan(&id, &value)
			a.Nil(err)
			a.Equal(1, id)
			a.Equal(10, value)

			count += 1
		}
		a.Equal(1, count)

		rows.Close()
	}
}
