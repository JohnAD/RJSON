package main

import (
	"testing"
)

func TestNativeValueMarshal(t *testing.T) {
	// Test simple values that are native to the std json library
	//
	// arrange
	//
	strDoc := "hello world"
	strExpected := "\"hello world\"\n"
	numDoc := 123.45
	numExpected := "123.45\n"
	trueDoc := true
	trueExpected := "true\n"
	falseDoc := false
	falseExpected := "false\n"
	nullDoc := interface{}(nil)
	nullExpected := "null\n"
	//
	// act
	//
	strResult, strErr := Marshal(strDoc)
	numResult, numErr := Marshal(numDoc)
	trueResult, trueErr := Marshal(trueDoc)
	falseResult, falseErr := Marshal(falseDoc)
	nullResult, nullErr := Marshal(nullDoc)
	//
	// assert
	//
	if strErr != nil {
		t.Error(strErr)
	}
	if numErr != nil {
		t.Error(numErr)
	}
	if trueErr != nil {
		t.Error(trueErr)
	}
	if falseErr != nil {
		t.Error(falseErr)
	}
	if nullErr != nil {
		t.Error(nullErr)
	}
	str := string(strResult)
	if str != strExpected {
		t.Errorf("got %s, want %s", str, strExpected)
	}
	num := string(numResult)
	if num != numExpected {
		t.Errorf("got %s, want %s", num, numExpected)
	}
	tser := string(trueResult)
	if tser != trueExpected {
		t.Errorf("got %s, want %s", tser, trueExpected)
	}
	fser := string(falseResult)
	if fser != falseExpected {
		t.Errorf("got %s, want %s", fser, falseExpected)
	}
	null := string(nullResult)
	if null != nullExpected {
		t.Errorf("got %s, want %s", null, nullExpected)
	}
}

func TestRJSONValueMarshal(t *testing.T) {
	// Test simple values that are part of the RJSON collection type
	//
	// arrange
	//
	strDoc := RString("hello world")
	strExpected := "\"hello world\"\n"
	numDoc := RNumber("123.45")
	numExpected := "123.45\n"
	trueDoc := RTrue()
	trueExpected := "true\n"
	falseDoc := RFalse()
	falseExpected := "false\n"
	nullDoc := RNull()
	nullExpected := "null\n"
	//
	// act
	//
	strResult, strErr := Marshal(strDoc)
	numResult, numErr := Marshal(numDoc)
	trueResult, trueErr := Marshal(trueDoc)
	falseResult, falseErr := Marshal(falseDoc)
	nullResult, nullErr := Marshal(nullDoc)
	//
	// assert
	//
	if strErr != nil {
		t.Error(strErr)
	}
	if numErr != nil {
		t.Error(numErr)
	}
	if trueErr != nil {
		t.Error(trueErr)
	}
	if falseErr != nil {
		t.Error(falseErr)
	}
	if nullErr != nil {
		t.Error(nullErr)
	}
	str := string(strResult)
	if str != strExpected {
		t.Errorf("got %s, want %s", str, strExpected)
	}
	num := string(numResult)
	if num != numExpected {
		t.Errorf("got %s, want %s", num, numExpected)
	}
	tser := string(trueResult)
	if tser != trueExpected {
		t.Errorf("got %s, want %s", tser, trueExpected)
	}
	fser := string(falseResult)
	if fser != falseExpected {
		t.Errorf("got %s, want %s", fser, falseExpected)
	}
	null := string(nullResult)
	if null != nullExpected {
		t.Errorf("got %s, want %s", null, nullExpected)
	}
}

func TestGeneralMarshal(t *testing.T) {
	//
	// arrange
	//
	doc := map[string]interface{}{
		"a_str_name": "foo",
		"d_obj_name": map[string]interface{}{
			"111": "x",
			"222": "y",
		},
	}
	expected := `{
  "a_str_name": "foo",
  "d_obj_name": {
    "111": "x",
    "222": "y"
  }
}
`
	//
	// act
	//
	result, err := Marshal(doc)
	//
	// assert
	//
	if err != nil {
		t.Error(err)
	}
	resultStr := string(result)
	if resultStr != expected {
		t.Errorf("marshalling failed, expected %s, got %s", expected, result)
	}
}

func TestObjectCombosMarshal(t *testing.T) {
	//
	// arrange
	//
	doc := map[string]interface{}{
		"c": map[string]interface{}{
			"111": "x",
			"222": "y",
		},
		"b": map[string]interface{}{},
		"a": map[string]interface{}{
			"111": "x",
			"aaa": map[string]interface{}{
				"111": "x",
				"222": "y",
			},
			"222": "y",
		},
	}
	expected := `{
  "a": {
    "111": "x",
    "222": "y",
    "aaa": {
      "111": "x",
      "222": "y"
    }
  },
  "b": {},
  "c": {
    "111": "x",
    "222": "y"
  }
}
`
	//
	// act
	//
	result, err := Marshal(doc)
	//
	// assert
	//
	if err != nil {
		t.Error(err)
	}
	resultStr := string(result)
	if resultStr != expected {
		t.Errorf("marshalling failed, expected %s, got %s", expected, result)
	}
}

// TODO: add other generic types next; then Decimal

//func TestStructMarshal(t *testing.T) {
//	//
//	// arrange
//	//
//	type Example struct {
//		Name string         `json:"name"`
//		Data map[string]int `json:"data"`
//	}
//	doc := Example{
//		Name: "test",
//		Data: map[string]int{"a": 1, "b": 2},
//	}
//	expected := `{
//  "name": "test",
//  "data": {
//    "a": 1,
//    "b": 2
//  }
//}
//`
//	//
//	// act
//	//
//	result, err := Marshal(doc)
//	//
//	// assert
//	//
//	if err != nil {
//		t.Error(err)
//	}
//	resultStr := string(result)
//	if resultStr != expected {
//		t.Errorf("marshalling failed, expected %s, got %s", expected, result)
//	}
//}
