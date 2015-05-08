package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestSqlFuncImplements(t *testing.T) {
	fnImplColumn := func(i interface{}) bool {
		return reflect.TypeOf(i).Implements(reflect.TypeOf(new(Column)).Elem())
	}
	fnImplColumn(&columnImpl{})
}

func TestSqlFunc(t *testing.T) {
	b := newBuilder()
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	Func("funcname", table1.C("id")).serialize(b)
	if `funcname("TABLE_A"."id")` != b.query.String() {
		t.Errorf("failed")
	}
	if len(b.Args()) != 0 {
		t.Errorf("failed")
	}
	if b.Err() != nil {
		t.Errorf("failed")
	}
}
