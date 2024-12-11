package RJSON

import (
	"testing"
)

func TestValueMarshal(t *testing.T) {
	//
	// arrange
	//
	doc := "hello world"
	expected := "\"hello world\"\n"
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
		t.Errorf("got %s, want %s", resultStr, expected)
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
