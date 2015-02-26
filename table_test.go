package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {
	a := assert.New(t)

	var table1 Table
	a.Panics(func() {
		table1 = NewTable(
			"TABLE_NAME",
			&TableOption{},
		)
	})

	a.NotPanics(func() {
		table1 = NewTable(
			"TABLE_NAME",
			&TableOption{},
			IntColumn("id", nil),
		)
	})
	a.NotNil(table1)

	a.NotPanics(func() {
		table1 = NewTable(
			"TABLE_NAME",
			nil,
			IntColumn("id", nil),
		)
	})
	a.NotNil(table1)

	a.NotPanics(func() {
		table1 = NewTable(
			"TABLE_NAME",
			nil,
			IntColumn("id", nil),
		)
	})
	a.NotNil(table1)
}

func TestJoinTable(t *testing.T) {
	a := assert.New(t)

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
	a.Equal(`"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.query.String())
	a.NoError(b.err)
	a.Empty(b.args)

	// left outer join
	b = newBuilder()
	joinedTable = l_table.LeftOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" LEFT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.query.String())
	a.NoError(b.err)
	a.Empty(b.args)

	// right outer join
	b = newBuilder()
	joinedTable = l_table.RightOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" RIGHT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.query.String())
	a.NoError(b.err)
	a.Empty(b.args)

	// full outer join
	b = newBuilder()
	joinedTable = l_table.FullOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" FULL OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`, b.query.String())
	a.NoError(b.err)
	a.Empty(b.args)

	// joined table column
	a.Equal(l_table.C("right_id"), joinedTable.C("right_id"))
	a.Equal(r_table.C("value"), joinedTable.C("value"))
	a.IsType(&errorColumn{}, joinedTable.C("not_exist_column"))
	a.IsType(&errorColumn{}, joinedTable.C("id"))

	// combination
	b = newBuilder()
	joinedTable = l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id"))).InnerJoin(rr_table, l_table.C("right_id").Eq(rr_table.C("id")))
	joinedTable.serialize(b)
	a.Equal(`"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id" INNER JOIN "RIGHTRIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHTRIGHT_TABLE"."id"`, b.query.String())
	a.NoError(b.err)
	a.Empty(b.args)
}

func TestTableColumnOperation(t *testing.T) {
	a := assert.New(t)
	table1 := NewTable(
		"TABLE_NAME",
		nil,
		IntColumn("id", nil),
	).(*table)

	// initial check
	if a.Len(table1.Columns(), 1) {
		expect := []string{"id"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	// AddColumnLast
	err := table1.AddColumnLast(IntColumn("test1", nil))
	a.NoError(err)
	if a.Len(table1.Columns(), 2) {
		expect := []string{"id", "test1"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	// AddColumnFirst
	err = table1.AddColumnFirst(IntColumn("first", nil))
	a.NoError(err)
	if a.Len(table1.Columns(), 3) {
		expect := []string{"first", "id", "test1"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	// AddColumnAfter
	err = table1.AddColumnAfter(IntColumn("second", nil), table1.C("first"))
	a.NoError(err)
	err = table1.AddColumnAfter(IntColumn("aaa", nil), table1.C("invalid"))
	a.Error(err)
	if a.Len(table1.Columns(), 4) {
		expect := []string{"first", "second", "id", "test1"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	// ChangeColumn
	err = table1.ChangeColumn(table1.C("id"), IntColumn("third", nil))
	a.NoError(err)
	err = table1.ChangeColumn(table1.C("invalid"), IntColumn("third", nil))
	a.Error(err)
	if a.Len(table1.Columns(), 4) {
		expect := []string{"first", "second", "third", "test1"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	// ChangeColumnFirst
	err = table1.ChangeColumnFirst(table1.C("test1"), IntColumn("new_first", nil))
	a.NoError(err)
	err = table1.ChangeColumnFirst(table1.C("invalid"), IntColumn("new_first", nil))
	a.Error(err)
	if a.Len(table1.Columns(), 4) {
		expect := []string{"new_first", "first", "second", "third"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	// ChangeColumnAfter
	err = table1.ChangeColumnAfter(table1.C("new_first"), IntColumn("fourth", nil), table1.C("third"))
	a.NoError(err)
	err = table1.ChangeColumnAfter(table1.C("invalid"), IntColumn("fourth", nil), table1.C("third"))
	a.Error(err)
	err = table1.ChangeColumnAfter(table1.C("second"), IntColumn("fourth", nil), table1.C("invalid"))
	a.Error(err)
	if a.Len(table1.Columns(), 4) {
		expect := []string{"first", "second", "third", "fourth"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	// ChangeColumnAfter
	err = table1.DropColumn(table1.C("fourth"))
	a.NoError(err)
	err = table1.DropColumn(table1.C("invalid"))
	a.Error(err)
	if a.Len(table1.Columns(), 3) {
		expect := []string{"first", "second", "third"}
		for i, col := range table1.Columns() {
			a.Equal(col.column_name(), expect[i])
		}
	}

	table1.SetName("TABLE_MODIFIED")
	a.Equal("TABLE_MODIFIED", table1.Name())
	return
}
