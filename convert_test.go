// Copyright 2013 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

// File convert_test.go tests the conversion
// to and from go data structures of a variety of types.

package zoom

import (
	"github.com/garyburd/redigo/redis"
	"reflect"
	"testing"
)

func TestPrimativeTypes(t *testing.T) {
	construct := func() (interface{}, error) {
		ms, err := newPrimativeTypesModels(1)
		return ms[0], err
	}
	testConvertType(reflect.TypeOf(primativeTypesModel{}), construct, t)
}

func TestPointerToPrimativeTypes(t *testing.T) {
	construct := func() (interface{}, error) {
		ms, err := newPointersToPrimativeTypesModels(1)
		return ms[0], err
	}
	testConvertType(reflect.TypeOf(pointersToPrimativeTypesModel{}), construct, t)
}

func TestInconvertibleTypes(t *testing.T) {
	construct := func() (interface{}, error) {
		ms, err := newInconvertibleTypesModels(1)
		return ms[0], err
	}
	testConvertType(reflect.TypeOf(inconvertibleTypesModel{}), construct, t)
}

func TestModelWithList(t *testing.T) {
	// use the generic tester to make sure what we get out is the same
	// as what we put in
	construct := func() (interface{}, error) {
		return &modelWithList{
			List: []string{"one", "two", "three"},
		}, nil
	}
	testConvertType(reflect.TypeOf(modelWithList{}), construct, t)

	// test to make sure the field is saved as a redis list type
	testingSetUp()
	defer testingTearDown()

	m1 := &modelWithList{
		List: []string{"one", "two", "three"},
	}
	m2 := &modelWithList{} // make sure empty slice doesn't throw an error
	if err := Save(m1, m2); err != nil {
		t.Error(err)
	}

	conn := GetConn()
	defer conn.Close()

	listKey := "modelWithList:" + m1.Id + ":List"
	result, err := redis.Strings(conn.Do("LRANGE", listKey, "0", "-1"))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(m1.List, result) {
		t.Errorf("List was not saved correctly.\nExpected: %+v\nGot: %+v\n", m1.List, result)
	}
}

func TestModelWithSet(t *testing.T) {
	// use the generic tester to make sure what we get out is the same
	// as what we put in
	construct := func() (interface{}, error) {
		return &modelWithSet{
			Set: []string{"one", "two", "three"},
		}, nil
	}
	testConvertType(reflect.TypeOf(modelWithSet{}), construct, t)

	// test to make sure the field is saved as a redis set type
	testingSetUp()
	defer testingTearDown()

	m1 := &modelWithSet{
		Set: []string{"one", "two", "three"},
	}
	m2 := &modelWithSet{} // make sure empty slice doesn't throw an error
	if err := Save(m1, m2); err != nil {
		t.Error(err)
	}

	conn := GetConn()
	defer conn.Close()

	setKey := "modelWithSet:" + m1.Id + ":Set"
	result, err := redis.Strings(conn.Do("SMEMBERS", setKey))
	if err != nil {
		t.Error(err)
	}
	if equal, msg := compareAsStringSet(m1.Set, result); !equal {
		t.Errorf("Set was not saved correctly.\nExpected: %+v\nGot: %+v\n%s\n", m1.Set, result, msg)
	}
}

// TODO:
//	- ModelWithHash

func TestEmbeddedStruct(t *testing.T) {
	construct := func() (interface{}, error) {
		return &embeddedStructModel{
			embed: embed{
				Int: 42,
			},
		}, nil
	}
	testConvertType(reflect.TypeOf(embeddedStructModel{}), construct, t)
}

func TestEmbeddedPointerToStruct(t *testing.T) {
	construct := func() (interface{}, error) {
		return &embeddedPointerToStructModel{
			embed: &embed{
				Int: 42,
			},
		}, nil
	}
	testConvertType(reflect.TypeOf(embeddedPointerToStructModel{}), construct, t)
}

// a general test that uses reflection
func testConvertType(typ reflect.Type, construct func() (in interface{}, err error), t *testing.T) {
	testingSetUp()
	defer testingTearDown()

	// construct a model using the construct function
	modelInterface, err := construct()
	if err != nil {
		t.Error(err)
	}
	model, ok := modelInterface.(Model)
	if !ok {
		t.Errorf("couldn't convert type %T to Model", modelInterface)
	}
	if err := Save(model); err != nil {
		t.Error(err)
	}

	// create a copy of the same type and use ScanById
	modelCopyInterface := reflect.New(typ).Interface()
	modelCopy, ok := modelCopyInterface.(Model)
	id := model.getId()
	if !ok {
		t.Errorf("couldn't convert type %T to Model", modelCopyInterface)
	}
	if err := ScanById(id, modelCopy); err != nil {
		t.Error(err)
	}

	// make sure the copy equals the original
	equal, err := looseEquals(model, modelCopy)
	if err != nil {
		t.Error(err)
	}
	if !equal {
		t.Errorf("model was not saved/retrieved correctly.\nExpected: %+v\nGot %+v\n", model, modelCopy)
	}
}
