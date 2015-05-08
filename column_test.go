package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestColumnImplements(t *testing.T) {
	fnImplColumn := func(i interface{}) bool {
		return reflect.TypeOf(i).Implements(reflect.TypeOf(new(Column)).Elem())
	}
	if !fnImplColumn(&columnImpl{}) {
		t.Errorf("fail")
	}
	if !fnImplColumn(&errorColumn{}) {
		t.Errorf("fail")
	}
	if !fnImplColumn(&aliasColumn{}) {
		t.Errorf("fail")
	}
}

func TestColumnOptionImpl(t *testing.T) {
	if !reflect.DeepEqual(&columnConfigImpl{
		name: "name",
		typ:  ColumnTypeBytes,
		opt: &ColumnOption{
			Unique: true,
		}}, newColumnConfigImpl("name", ColumnTypeBytes, &ColumnOption{Unique: true})) {
		t.Errorf("fail")
	}
	if !reflect.DeepEqual(&columnConfigImpl{
		name: "name",
		typ:  ColumnTypeBytes,
		opt:  &ColumnOption{},
	}, newColumnConfigImpl("name", ColumnTypeBytes, nil)) {
		t.Errorf("fail")
	}
}
