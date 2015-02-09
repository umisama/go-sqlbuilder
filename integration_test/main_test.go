package sqlbuilder_integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	sb "github.com/umisama/go-sqlbuilder"
	_ "github.com/ziutek/mymysql/godrv"
)

var db *sql.DB

// Table for testing
var tbl_person = sb.NewTable(
	"PERSON",
	sb.IntColumn("id", sb.CO_PrimaryKey),
	sb.StringColumn("name", 255, sb.CO_Unique),
	sb.DateColumn("birth"),
)
var tbl_phone = sb.NewTable(
	"PHONE",
	sb.IntColumn("id", sb.CO_PrimaryKey, sb.CO_AutoIncrement),
	sb.IntColumn("person_id"),
	sb.StringColumn("number", 255),
)
var tbl_email = sb.NewTable(
	"EMAIL",
	sb.IntColumn("id", sb.CO_PrimaryKey, sb.CO_AutoIncrement),
	sb.IntColumn("person_id"),
	sb.StringColumn("address", 255),
)

// Data for testing
type Person struct {
	Id    int
	Name  string
	Birth time.Time
}

type Phone struct {
	PersonId int
	Number   string
}

type Email struct {
	PersonId int
	Address  string
}

var persons = []Person{{
	Id:    1,
	Name:  "Rintaro Okabe",
	Birth: time.Date(1991, time.December, 14, 0, 0, 0, 0, time.UTC),
}, {
	Id:    2,
	Name:  "Mayuri Shiina",
	Birth: time.Date(1994, time.February, 1, 0, 0, 0, 0, time.UTC),
}, {
	Id:    3,
	Name:  "Itaru Hashida",
	Birth: time.Date(1991, time.May, 19, 0, 0, 0, 0, time.UTC),
}}

var phones = []Phone{{
	PersonId: 1,
	Number:   "000-0000-0000",
}, {
	PersonId: 2,
	Number:   "111-1111-1111",
}, {
	PersonId: 2,
	Number:   "111-1111-2222",
}}

var emails = []Email{{
	PersonId: 1,
	Address:  "sg-epk@jtk93.x29.jp",
}, {
	PersonId: 1,
	Address:  "okarin@example.org",
}, {
	PersonId: 2,
	Address:  "mayusii@example.org",
}, {
	PersonId: 3,
	Address:  "hashida@example.org",
}}

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
		{"postgres", sb.PostgresDialect{}, "postgres", "user=postgres dbname=go_sqlbuilder_test sslmode=disable"},
	}

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
