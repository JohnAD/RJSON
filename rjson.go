package main

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
)

type RJsonElementType uint8

const (
	Void RJsonElementType = iota
	Object
	Array
	String
	Number
	True
	False
	Null
)

func (t RJsonElementType) String() string {
	switch t {
	case Void:
		return "Void"
	case Object:
		return "Object"
	case Array:
		return "Array"
	case String:
		return "String"
	case Number:
		return "Number"
	case True:
		return "True"
	case False:
		return "False"
	case Null:
		return "Null"
	default:
		return "Invalid JSON Element Type"
	}
}

type RJsonElement struct {
	kind RJsonElementType
	obj  map[string]RJsonElement
	arr  []RJsonElement
	str  string
	num  decimal.Decimal
}

func (j RJsonElement) ToNativeInterface() interface{} {
	// this function is mostly for use internally in the library
	switch j.kind {
	case Void:
		return nil // THIS SHOULD NEVER HAPPEN
	case Object:
		return j.ToNativeMap()
	case Array:
		return j.ToNativeArray()
	case String:
		return j.str
	case Number:
		return j.num
	case True:
		return true
	case False:
		return false
	case Null:
		return nil
	}
	return nil
}

func JsonVoid() RJsonElement {
	// Void means "undefined" or "not-applicable"; NOT "unknown". For "unknown" call RNull() for the `null` type
	return RJsonElement{kind: Void}
}

func (j RJsonElement) IsVoid() bool {
	return j.kind == Void
}

func JsonObject() RJsonElement {
	newObject := RJsonElement{kind: Object, obj: make(map[string]RJsonElement)}
	return newObject
}

func (j RJsonElement) IsObject() bool {
	return j.kind == Object
}

func (j RJsonElement) ToNativeMap() map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range j.obj {
		if v.kind != Void {
			newMap[k] = v.ToNativeInterface()
		}
	}
	return newMap
}

func (j RJsonElement) Set(name string, value RJsonElement) {
	// Note: setting an existing value to Void deletes the entry
	if j.kind == Object {
		if value.kind == Void {
			if _, ok := j.obj[name]; ok {
				delete(j.obj, name)
			}
		} else {
			j.obj[name] = value
		}
	}
}

func (j RJsonElement) GetWithDefault(name string, defaultValue RJsonElement) RJsonElement {
	if j.kind == Object {
		if v, ok := j.obj[name]; ok {
			return v
		}
	}
	return defaultValue
}

func (j RJsonElement) Get(name string) RJsonElement {
	if j.kind == Object {
		if v, ok := j.obj[name]; ok {
			return v
		}
	}
	return RJsonElement{kind: Void}
}

func (j RJsonElement) HasKey(name string) bool {
	if j.kind == Object {
		if v, ok := j.obj[name]; ok {
			if v.kind != Void {
				return true
			}
		}
	}
	return false
}

func (j RJsonElement) Delete(name string) {
	if j.kind == Object {
		delete(j.obj, name)
	}
}

func JsonArray() RJsonElement {
	newArray := RJsonElement{kind: Array, arr: make([]RJsonElement, 0)}
	return newArray
}

func (j RJsonElement) ToNativeArray() []interface{} {
	var newArray []interface{}
	for i := range j.arr {
		v := j.arr[i]
		if v.kind != Void {
			newArray = append(newArray, v.ToNativeInterface())
		}
	}
	return newArray
}

func (j RJsonElement) Append(value RJsonElement) {
	if j.kind == Array {
		if value.kind != Void {
			j.arr = append(j.arr, value)
		}
	}
}

func (j RJsonElement) IsArray() bool {
	return j.kind == Array
}

func RString(str string) RJsonElement {
	newString := RJsonElement{kind: String, str: str}
	return newString
}

func (j RJsonElement) IsString() bool {
	return j.kind == String
}

func (j RJsonElement) ToStringWithDefault(defaultValue string) string {
	if j.kind == String {
		return j.str
	}
	return defaultValue
}

func (j RJsonElement) ToString() (string, error) {
	if j.kind == String {
		return j.str, nil
	}
	return "", errors.New(fmt.Sprintf("JSON value of kind %s cannot be converted to string", j.kind))
}

func RNumber(params ...interface{}) RJsonElement {
	newNumber := RJsonElement{kind: Number, num: decimal.Zero}
	if len(params) > 0 {
		value := params[0]
		if value, ok := value.(string); ok {
			newDec, err := decimal.NewFromString(value)
			if err == nil {
				newNumber.num = newDec
			}
		}
	}
	return newNumber
}

func (j RJsonElement) IsNumber() bool {
	return j.kind == Number
}

func (j RJsonElement) ToDecimalWithDefault(defaultValue decimal.Decimal) decimal.Decimal {
	if j.kind == Number {
		return j.num
	}
	return defaultValue
}

func (j RJsonElement) ToDecimal() (decimal.Decimal, error) {
	if j.kind == Number {
		return j.num, nil
	}
	return decimal.Zero, errors.New(fmt.Sprintf("JSON value of kind %s cannot be converted to Decimal", j.kind))
}

func (j RJsonElement) ToFloat64() (float64, error) {
	if j.kind == Number {
		newFloat, _ := j.num.Float64()
		return newFloat, nil
	}
	return 0.0, errors.New(fmt.Sprintf("JSON value of kind %s cannot be converted to float64", j.kind))
}

func (j RJsonElement) ToInt() (int, error) {
	if j.kind == Number {
		newFloat := j.num.InexactFloat64()
		newInt := int(math.Round(newFloat))
		return newInt, nil
	}
	return 0.0, errors.New(fmt.Sprintf("JSON value of kind %s cannot be converted to int", j.kind))
}

func RTrue() RJsonElement {
	return RJsonElement{kind: True}
}

func (j RJsonElement) IsTrue() bool {
	return j.kind == True
}

func RFalse() RJsonElement {
	return RJsonElement{kind: False}
}

func (j RJsonElement) IsFalse() bool {
	return j.kind == False
}

func RBool(bool bool) RJsonElement {
	if bool {
		return RJsonElement{kind: True}
	} else {
		return RJsonElement{kind: False}
	}
}

func (j RJsonElement) IsBool() bool {
	return j.kind == True || j.kind == False
}

func (j RJsonElement) ToBoolWithDefault(defaultValue bool) bool {
	switch j.kind {
	case True:
		return true
	case False:
		return false
	default:
		return defaultValue
	}
}

func (j RJsonElement) ToBool() (bool, error) {
	switch j.kind {
	case True:
		return true, nil
	case False:
		return false, nil
	default:
		return false, errors.New(fmt.Sprintf("JSON value of kind %s cannot be converted to bool", j.kind))
	}
}

func RNull() RJsonElement {
	return RJsonElement{kind: Null}
}

func (j RJsonElement) IsNull() bool {
	return j.kind == Null
}
