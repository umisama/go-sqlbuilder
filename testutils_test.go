package sqlbuilder

import (
	"reflect"
)

type statementTestCase struct {
	stmt  Statement
	query string
	args  []interface{}
	err   bool
}

func (testCase statementTestCase) Run() (message string, args []interface{}, ok bool) {
	query, args, err := testCase.stmt.ToSql()
	if testCase.err {
		if err == nil {
			return "error: expect returns error but got nil.", []interface{}{}, false
		}
	} else {
		if err != nil {
			return "error: expect returns no error got %s.", []interface{}{err.Error()}, false
		}
	}
	if testCase.query != query {
		return "expect returns query \n%s \nbut got\n%s.", []interface{}{testCase.query, query}, false
	}
	if !reflect.DeepEqual(testCase.args, args) {
		return "expect returns arguments \n%s \nbut got\n%s.", []interface{}{testCase.args, args}, false
	}
	return "", nil, true
}

type conditionTestCase struct {
	cond  Condition
	query string
	args  []interface{}
	err   error
}

func (testCase conditionTestCase) Run() (message string, args []interface{}, ok bool) {
	bldr := newBuilder()
	testCase.cond.serialize(bldr)
	if testCase.err != nil {
		if bldr.err == nil {
			return "error: expect returns error but got nil.", []interface{}{}, false
		}
	} else {
		if bldr.err != nil {
			return "error: expect returns no error got %s.", []interface{}{bldr.err.Error()}, false
		}
	}
	if bldr.query.String() != testCase.query {
		return "expect returns query \n%s \nbut got\n%s.", []interface{}{testCase.query, bldr.query.String()}, false
	}
	if !reflect.DeepEqual(bldr.args, testCase.args) {
		return "expect returns arguments \n%s \nbut got\n%s.", []interface{}{testCase.args, args}, false
	}
	return "", nil, true
}
