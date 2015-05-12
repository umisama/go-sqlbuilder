package sqlbuilder_integration

import (
	"reflect"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	sb "github.com/umisama/go-sqlbuilder"
	_ "github.com/ziutek/mymysql/godrv"
)

func TestCreateTable(t *testing.T) {
	for _, table := range []sb.Table{tbl_person, tbl_phone, tbl_email} {
		query, args, err := sb.CreateTable(table).ToSql()
		if err != nil {
			t.Error(err.Error())
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestAlterTable(t *testing.T) {
	stmt := sb.AlterTable(tbl_person).AddColumn(sb.IntColumn("other_column", nil))
	query, args, err := stmt.ToSql()
	_, err = db.Exec(query, args...)
	if err != nil {
		t.Error(err.Error())
	}

	err = stmt.ApplyToTable()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCreateIndex(t *testing.T) {
	for _, table := range []sb.Table{tbl_phone, tbl_email} {
		query, args, err := sb.CreateIndex(table).Name("I_" + table.Name() + "_PERSONID").Columns(table.C("person_id")).
			ToSql()
		if err != nil {
			t.Error(err.Error())
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestInsertToPersonTable(t *testing.T) {
	for _, person := range persons {
		query, args, err := sb.Insert(tbl_person).
			Columns(tbl_person.C("id"), tbl_person.C("name"), tbl_person.C("birth")).
			Values(person.Id, person.Name, person.Birth).
			ToSql()
		if err != nil {
			t.Error(err.Error())
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestInsertToPhoneTable(t *testing.T) {
	for _, phone := range phones {
		query, args, err := sb.Insert(tbl_phone).
			Columns(tbl_phone.C("person_id"), tbl_phone.C("number")).
			Values(phone.PersonId, phone.Number).
			ToSql()
		if err != nil {
			t.Error(err.Error())
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestInsertToEmailTable(t *testing.T) {
	for _, email := range emails {
		query, args, err := sb.Insert(tbl_email).
			Columns(tbl_email.C("person_id"), tbl_email.C("address")).
			Values(email.PersonId, email.Address).
			ToSql()
		if err != nil {
			t.Error(err.Error())
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestSelectSimple(t *testing.T) {
	query, args, err := sb.Select(tbl_person).
		Columns(tbl_person.C("id"), tbl_person.C("name"), tbl_person.C("birth")).
		OrderBy(false, tbl_person.C("id")).
		ToSql()
	if err != nil {
		t.Error(err.Error())
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		t.Error(err.Error())
	}

	got_persons := make([]Person, 0)
	if rows != nil {
		for rows.Next() {
			id, name, birth := 0, "", time.Time{}
			err := rows.Scan(&id, &name, &birth)
			if err != nil {
				t.Error(err.Error())
			}
			got_persons = append(got_persons, Person{
				Id:    id,
				Name:  name,
				Birth: birth.UTC(),
			})
		}
		rows.Close()
	} else {
		t.Error("faild")
	}
	if !reflect.DeepEqual(got_persons, persons) {
		t.Error("faild")
	}
}

func TestSelectJoinedWithAlias(t *testing.T) {
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
	col_id := tbl_person.C("id").As("id")
	col_name := tbl_person.C("name").As("name")
	col_birth := tbl_person.C("birth").As("birth")
	col_email := tbl_email.C("address").As("email")
	query, args, err := sb.Select(tbl_person_email).
		Columns(col_id, col_name, col_birth, col_email).
		OrderBy(false, col_id).
		OrderBy(false, tbl_email.C("id")).
		ToSql()
	if err != nil {
		t.Error(err.Error())
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		t.Error(err.Error())
	}

	got_persons := make([]PersonEmail, 0)
	if rows != nil {
		for rows.Next() {
			id, name, birth, email := 0, "", time.Time{}, ""
			err := rows.Scan(&id, &name, &birth, &email)
			if err != nil {
				t.Error(err.Error())
			}
			got_persons = append(got_persons, PersonEmail{
				Id:    id,
				Name:  name,
				Birth: birth.UTC(),
				Email: email,
			})
		}
		rows.Close()
	} else {
		t.Error("failed")
	}
	if !reflect.DeepEqual(got_persons, expect) {
		t.Error("failed")
	}
}

func TestSelectJoined(t *testing.T) {
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
	query, args, err := sb.Select(tbl_person_email).
		Columns(tbl_person.C("id"), tbl_person.C("name"), tbl_person.C("birth"), tbl_email.C("address")).
		OrderBy(false, tbl_person.C("id")).
		OrderBy(false, tbl_email.C("id")).
		ToSql()
	if err != nil {
		t.Error(err.Error())
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		t.Error(err.Error())
	}

	got_persons := make([]PersonEmail, 0)
	if rows != nil {
		for rows.Next() {
			id, name, birth, email := 0, "", time.Time{}, ""
			err := rows.Scan(&id, &name, &birth, &email)
			if err != nil {
				t.Error(err.Error())
			}
			got_persons = append(got_persons, PersonEmail{
				Id:    id,
				Name:  name,
				Birth: birth.UTC(),
				Email: email,
			})
		}
		rows.Close()
	} else {
		t.Error("failed")
	}
	if !reflect.DeepEqual(got_persons, expect) {
		t.Error("failed")
	}
}

func TestDelete(t *testing.T) {
	query, args, err := sb.Delete(tbl_phone).Where(tbl_phone.C("id").Eq(1)).ToSql()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSqlFunction1(t *testing.T) {
	query, args, err := sb.Select(tbl_phone).
		Columns(sb.Func("count", tbl_phone.C("id"))).
		ToSql()
	if err != nil {
		t.Error(err.Error())
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		t.Error(err.Error())
	}
	if rows != nil {
		defer rows.Close()
		rows.Next()
		value := 0
		err := rows.Scan(&value)
		if err != nil {
			t.Error(err.Error())
		}
		if value != 2 {
			t.Error("failed")
		}
	} else {
		t.Error("failed")
	}
}

func TestDropTable(t *testing.T) {
	for _, table := range []sb.Table{tbl_person, tbl_phone, tbl_email} {
		query, args, err := sb.DropTable(table).ToSql()
		if err != nil {
			t.Error(err.Error())
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			t.Error(err.Error())
		}
	}
}
