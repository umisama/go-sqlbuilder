package sqlbuilder_integration

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	sb "github.com/umisama/go-sqlbuilder"
	_ "github.com/ziutek/mymysql/godrv"
)

func TestCreateTable(t *testing.T) {
	a := assert.New(t)

	query, args, err := sb.CreateTable(table1).ToSql()
	a.Nil(err)

	_, err = db.Exec(query, args...)
	a.Nil(err)
}

func TestInsert(t *testing.T) {
	a := assert.New(t)

	// data 1
	query, args, err := sb.Insert(table1).
		Columns(table1.C("id"), table1.C("value")).
		Values(1, 10).
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
		Values(2, 20).
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
		From(table1).
		Where(table1.C("id").Eq(1)).
		Limit(1).OrderBy(false, table1.C("id")).ToSql()
	a.Nil(err)

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

func TestDelete(t *testing.T) {
	a := assert.New(t)
	query, args, err := sb.Delete(table1).Where(table1.C("id").Eq(0)).ToSql()
	a.Nil(err)

	_, err = db.Exec(query, args...)
	a.Nil(err)
}

func TestDropTable(t *testing.T) {
	a := assert.New(t)
	query, args, err := sb.DropTable(table1).ToSql()
	a.Nil(err)

	_, err = db.Exec(query, args...)
	a.Nil(err)
}
