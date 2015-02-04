package sqlbuilder_integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
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
		sb.IntColumn("id", sb.CO_PrimaryKey),
		sb.IntColumn("value"),
	)

	for _, c := range cases {
		fmt.Println("START unit test for", c.name)

		var err error
		db, err = sql.Open(c.driver, c.dsn)
		if err != nil {
			fmt.Println(err.Error())
		}
		sb.SetDialect(c.dialect)

		results[c.name] = m.Run()
	}

	for _, v := range results {
		if v != 0 {
			os.Exit(v)
		}
	}
	os.Exit(0)
}
