package sqlbuilder_integration

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	sb "github.com/umisama/go-sqlbuilder"
	_ "github.com/ziutek/mymysql/godrv"
)

func TestCreateTable(t *testing.T) {
	a := assert.New(t)

	for _, table := range []sb.Table{tbl_person, tbl_phone, tbl_email} {
		query, args, err := sb.CreateTable(table).ToSql()
		a.NoError(err)

		_, err = db.Exec(query, args...)
		a.NoError(err)
	}
}

func TestCreateIndex(t *testing.T) {
	a := assert.New(t)
	for _, table := range []sb.Table{tbl_phone, tbl_email} {
		query, args, err := sb.CreateIndex(table).Name("I_" + table.Name() + "_PERSONID").Columns(table.C("person_id")).
			IfNotExists().
			ToSql()
		a.NoError(err)

		_, err = db.Exec(query, args...)
		a.NoError(err)
	}
}

func TestInsertToPersonTable(t *testing.T) {
	a := assert.New(t)
	for _, person := range persons {
		query, args, err := sb.Insert(tbl_person).
			Columns(tbl_person.C("id"), tbl_person.C("name"), tbl_person.C("birth")).
			Values(person.Id, person.Name, person.Birth).
			ToSql()
		a.NoError(err)

		_, err = db.Exec(query, args...)
		a.NoError(err)
	}
}

func TestInsertToPhoneTable(t *testing.T) {
	a := assert.New(t)
	for _, phone := range phones {
		query, args, err := sb.Insert(tbl_phone).
			Columns(tbl_phone.C("person_id"), tbl_phone.C("number")).
			Values(phone.PersonId, phone.Number).
			ToSql()
		a.NoError(err)

		_, err = db.Exec(query, args...)
		a.NoError(err)
	}
}

func TestInsertToEmailTable(t *testing.T) {
	a := assert.New(t)
	for _, email := range emails {
		query, args, err := sb.Insert(tbl_email).
			Columns(tbl_email.C("person_id"), tbl_email.C("address")).
			Values(email.PersonId, email.Address).
			ToSql()
		a.NoError(err)

		_, err = db.Exec(query, args...)
		a.NoError(err)
	}
}

func TestSelectSimple(t *testing.T) {
	a := assert.New(t)

	query, args, err := sb.Select(
		tbl_person.C("id"), tbl_person.C("name"), tbl_person.C("birth")).
		From(tbl_person).
		OrderBy(false, tbl_person.C("id")).
		ToSql()
	a.NoError(err)
	rows, err := db.Query(query, args...)
	a.NoError(err)

	got_persons := make([]Person, 0)
	if a.NotNil(rows) {
		for rows.Next() {
			id, name, birth := 0, "", time.Time{}
			err := rows.Scan(&id, &name, &birth)
			a.NoError(err)
			got_persons = append(got_persons, Person{
				Id:    id,
				Name:  name,
				Birth: birth.UTC(),
			})
		}
		rows.Close()
	}
	a.Equal(got_persons, persons)
}

func TestSelectJoined(t *testing.T) {
	a := assert.New(t)
	type PersonEmail struct {
		Id    int
		Name  string
		Birth time.Time
		Email string
	}
	expect := []PersonEmail{{
		Id:    1,
		Name:  "Rintaro Okabe",
		Birth: time.Date(1991, time.December, 14, 0, 0, 0, 0, time.UTC),
		Email: "sg-epk@jtk93.x29.jp",
	}, {
		Id:    1,
		Name:  "Rintaro Okabe",
		Birth: time.Date(1991, time.December, 14, 0, 0, 0, 0, time.UTC),
		Email: "okarin@example.org",
	}, {
		Id:    2,
		Name:  "Mayuri Shiina",
		Birth: time.Date(1994, time.February, 1, 0, 0, 0, 0, time.UTC),
		Email: "mayusii@example.org",
	}, {
		Id:    3,
		Name:  "Itaru Hashida",
		Birth: time.Date(1991, time.May, 19, 0, 0, 0, 0, time.UTC),
		Email: "hashida@example.org",
	}}

	tbl_person_email := tbl_person.LeftOuterJoin(tbl_email, tbl_email.C("person_id").Eq(tbl_person.C("id")))
	query, args, err := sb.Select(
		tbl_person.C("id"), tbl_person.C("name"), tbl_person.C("birth"), tbl_email.C("address")).
		From(tbl_person_email).
		OrderBy(false, tbl_person.C("id")).
		OrderBy(false, tbl_email.C("id")).
		ToSql()
	a.NoError(err)
	rows, err := db.Query(query, args...)
	a.NoError(err)

	got_persons := make([]PersonEmail, 0)
	if a.NotNil(rows) {
		for rows.Next() {
			id, name, birth, email := 0, "", time.Time{}, ""
			err := rows.Scan(&id, &name, &birth, &email)
			a.NoError(err)
			got_persons = append(got_persons, PersonEmail{
				Id:    id,
				Name:  name,
				Birth: birth.UTC(),
				Email: email,
			})
		}
		rows.Close()
	}
	a.Equal(expect, got_persons)
}

func TestDelete(t *testing.T) {
	a := assert.New(t)
	query, args, err := sb.Delete(tbl_phone).Where(tbl_phone.C("id").Eq(1)).ToSql()
	a.NoError(err)

	_, err = db.Exec(query, args...)
	a.NoError(err)
}

func TestSqlFunction1(t *testing.T) {
	a := assert.New(t)
	query, args, err := sb.Select(
		sb.Func("count", tbl_phone.C("id")),
	).From(tbl_phone).ToSql()
	a.NoError(err)

	rows, err := db.Query(query, args...)
	a.NoError(err)
	if a.NotNil(rows) {
		defer rows.Close()
		rows.Next()
		value := 0
		err := rows.Scan(&value)
		a.NoError(err)
		a.Equal(2, value)
	}
}

func TestDropTable(t *testing.T) {
	a := assert.New(t)
	for _, table := range []sb.Table{tbl_person, tbl_phone, tbl_email} {
		query, args, err := sb.DropTable(table).ToSql()
		a.NoError(err)

		_, err = db.Exec(query, args...)
		a.NoError(err)
	}
}
