package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestTable(t *testing.T) {
	var table1 Table
	var fnPanic = func(fn func()) (ok bool) {
		defer func() {
			if r := recover(); r != nil {
				ok = true
			} else {
				ok = false
			}
		}()
		fn()
		return
	}

	// be panic
	if !fnPanic(func() {
		table1 = NewTable(
			"TABLE_NAME",
			&TableOption{},
		)
	}) {
		t.Errorf("failed")
	}
	if table1 != nil {
		t.Errorf("failed")
	}

	// not panic
	if fnPanic(func() {
		table1 = NewTable(
			"TABLE_NAME",
			&TableOption{},
			IntColumn("id", nil),
		)
	}) {
		t.Errorf("failed")
	}
	if table1 == nil {
		t.Errorf("failed")
	}

	// not panic
	if fnPanic(func() {
		table1 = NewTable(
			"TABLE_NAME",
			nil,
			IntColumn("id", nil),
		)
	}) {
		t.Errorf("failed")
	}
	if table1 == nil {
		t.Errorf("failed")
	}

	// not panic
	if fnPanic(func() {
		table1 = NewTable(
			"TABLE_NAME",
			nil,
			IntColumn("id", nil),
		)
	}) {
		t.Errorf("failed")
	}
	if table1 == nil {
		t.Errorf("failed")
	}
}

func TestJoinTable(t *testing.T) {
	l_table := NewTable(
		"LEFT_TABLE",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("right_id", nil),
	)
	r_table := NewTable(
		"RIGHT_TABLE",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("value", nil),
	)
	rr_table := NewTable(
		"RIGHTRIGHT_TABLE",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)

	// inner join
	b := newBuilder()
	joinedTable := l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	if `"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"` != b.query.String() {
		t.Error("failed")
	}
	if b.err != nil {
		t.Error("failed")
	}
	if len(b.args) != 0 {
		t.Error("failed")
	}

	// left outer join
	b = newBuilder()
	joinedTable = l_table.LeftOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	if `"LEFT_TABLE" LEFT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"` != b.query.String() {
		t.Error("failed")
	}
	if b.err != nil {
		t.Error("failed")
	}
	if len(b.args) != 0 {
		t.Error("failed")
	}

	// right outer join
	b = newBuilder()
	joinedTable = l_table.RightOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	if `"LEFT_TABLE" RIGHT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"` != b.query.String() {
		t.Error("failed")
	}
	if b.err != nil {
		t.Error("failed")
	}
	if len(b.args) != 0 {
		t.Error("failed")
	}

	// full outer join
	b = newBuilder()
	joinedTable = l_table.FullOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	if `"LEFT_TABLE" FULL OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"` != b.query.String() {
		t.Error("failed")
	}
	if b.err != nil {
		t.Error("failed")
	}
	if len(b.args) != 0 {
		t.Error("failed")
	}

	// joined table column
	if !reflect.DeepEqual(l_table.C("right_id"), joinedTable.C("right_id")) {
		t.Error("failed")
	}
	if !reflect.DeepEqual(r_table.C("value"), joinedTable.C("value")) {
		t.Error("failed")
	}
	if _, ok := joinedTable.C("not_exist_column").(*errorColumn); !ok {
		t.Error("failed")
	}
	if _, ok := joinedTable.C("id").(*errorColumn); !ok {
		t.Error("failed")
	}

	// combination
	b = newBuilder()
	joinedTable = l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id"))).InnerJoin(rr_table, l_table.C("right_id").Eq(rr_table.C("id")))
	joinedTable.serialize(b)
	if `"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id" INNER JOIN "RIGHTRIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHTRIGHT_TABLE"."id"` != b.query.String() {
		t.Error("failed")
	}
	if b.err != nil {
		t.Error("failed")
	}
	if len(b.args) != 0 {
		t.Error("failed")
	}
}

func TestTableColumnOperation(t *testing.T) {
	var fnEqualColumnName = func(cols []Column, expect []string) bool {
		if len(cols) != len(expect) {
			return false
		}
		for i, col := range cols {
			if col.column_name() != expect[i] {
				return false
			}
		}
		return true
	}

	table1 := NewTable(
		"TABLE_NAME",
		nil,
		IntColumn("id", nil),
	).(*table)

	// initial check
	if !fnEqualColumnName(table1.Columns(), []string{"id"}) {
		t.Error("failed")
	}

	// AddColumnLast
	err := table1.AddColumnLast(IntColumn("test1", nil))
	if err != nil {
		t.Error("failed")
	}
	if !fnEqualColumnName(table1.Columns(), []string{"id", "test1"}) {
		t.Error("failed")
	}

	// AddColumnFirst
	err = table1.AddColumnFirst(IntColumn("first", nil))
	if err != nil {
		t.Error("failed")
	}
	if !fnEqualColumnName(table1.Columns(), []string{"first", "id", "test1"}) {
		t.Error("failed")
	}

	// AddColumnAfter
	err = table1.AddColumnAfter(IntColumn("second", nil), table1.C("first"))
	if err != nil {
		t.Error("failed")
	}
	err = table1.AddColumnAfter(IntColumn("aaa", nil), table1.C("invalid"))
	if err == nil {
		t.Error("failed")
	}
	if !fnEqualColumnName(table1.Columns(), []string{"first", "second", "id", "test1"}) {
		t.Error("failed")
	}

	// ChangeColumn
	err = table1.ChangeColumn(table1.C("id"), IntColumn("third", nil))
	if err != nil {
		t.Error("failed")
	}
	err = table1.ChangeColumn(table1.C("invalid"), IntColumn("third", nil))
	if err == nil {
		t.Error("failed")
	}
	if !fnEqualColumnName(table1.Columns(), []string{"first", "second", "third", "test1"}) {
		t.Error("failed")
	}

	// ChangeColumnFirst
	err = table1.ChangeColumnFirst(table1.C("test1"), IntColumn("new_first", nil))
	if err != nil {
		t.Error("failed")
	}
	err = table1.ChangeColumnFirst(table1.C("invalid"), IntColumn("new_first", nil))
	if err == nil {
		t.Error("failed")
	}
	if !fnEqualColumnName(table1.Columns(), []string{"new_first", "first", "second", "third"}) {
		t.Error("failed")
	}

	// ChangeColumnAfter
	err = table1.ChangeColumnAfter(table1.C("new_first"), IntColumn("fourth", nil), table1.C("third"))
	if err != nil {
		t.Error("failed")
	}
	err = table1.ChangeColumnAfter(table1.C("invalid"), IntColumn("fourth", nil), table1.C("third"))
	if err == nil {
		t.Error("failed")
	}
	err = table1.ChangeColumnAfter(table1.C("second"), IntColumn("fourth", nil), table1.C("invalid"))
	if err == nil {
		t.Error("failed")
	}
	if !fnEqualColumnName(table1.Columns(), []string{"first", "second", "third", "fourth"}) {
		t.Error("failed")
	}

	// ChangeColumnAfter
	err = table1.DropColumn(table1.C("fourth"))
	if err != nil {
		t.Error("failed")
	}
	err = table1.DropColumn(table1.C("invalid"))
	if err == nil {
		t.Error("failed")
	}
	if !fnEqualColumnName(table1.Columns(), []string{"first", "second", "third"}) {
		t.Error("failed")
	}

	// SetName
	table1.SetName("TABLE_MODIFIED")
	if "TABLE_MODIFIED" != table1.Name() {
		t.Error("failed")
	}
	return
}
