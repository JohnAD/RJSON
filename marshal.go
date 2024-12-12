package RJSON

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

const (
	MaxDepth = 8096
	Indent   = "  "
)

type CollectionType uint

const (
	UNKNOWN CollectionType = iota
	VALUE
	MAP
	ARRAY
)

type MarshalLevel struct {
	Value        interface{}
	Type         CollectionType
	Keys         []string // used with Maps, Structs, and Sets
	CurrentIndex int      // used with Maps, Structs, Arrays, and Sets
}

type MarshalState struct {
	Levels       [MaxDepth]MarshalLevel
	CurrentLevel int
	Buf          []byte
}

func appendIndent(state *MarshalState, offset int) {
	for i := 0; i < (state.CurrentLevel + offset); i++ {
		state.Buf = append(state.Buf, Indent...)
	}
}

func Marshal(v interface{}) ([]byte, error) {
	state := MarshalState{
		Levels: [MaxDepth]MarshalLevel{
			{
				Value:        v,
				Type:         UNKNOWN,
				Keys:         []string{},
				CurrentIndex: 0,
			},
		},
		CurrentLevel: 0,
		Buf:          make([]byte, 0, 1024),
	}
	err := nonRecursiveMarshal(&state)
	return state.Buf, err
}

func nonRecursiveMarshal(state *MarshalState) error {
	for state.CurrentLevel >= 0 {
		current := state.Levels[state.CurrentLevel]
		val := reflect.ValueOf(current.Value)
		switch val.Kind() {
		case reflect.String, reflect.Float32, reflect.Float64, reflect.Bool, reflect.Interface, reflect.Invalid:
			err := handleSimple(state, val)
			if err != nil {
				return err
			}
			finishLevel(state)
		case reflect.Map:
			err := handleMap(state)
			if err != nil {
				return err
			}
		case reflect.Struct:
			switch current.Value.(type) {
			case RJsonElement:
				err := handleRJSON(state)
				if err != nil {
					return err
				}
				finishLevel(state)
			default:
				return fmt.Errorf("unsupported type %v", val.Type())
			}
		default:
			return fmt.Errorf("unsupported type: %T  [%v]", current.Value, val.Kind())
		}
	}
	state.Buf = append(state.Buf, "\n"...)
	return nil
}

func addLevel(state *MarshalState, value interface{}) {
	state.CurrentLevel++
	state.Levels[state.CurrentLevel].Type = UNKNOWN
	state.Levels[state.CurrentLevel].Value = value
	state.Levels[state.CurrentLevel].CurrentIndex = 0
	state.Levels[state.CurrentLevel].Keys = []string{}
}

func finishLevel(state *MarshalState) {
	state.Levels[state.CurrentLevel].Type = UNKNOWN
	state.CurrentLevel--
}

// handleSimple handles the "simple generic" case of JSON values. That is, that which is
// Unmarshalled on a generic interface{} by the `json` lib AND is not a structure (like array
// or map).
func handleSimple(state *MarshalState, val reflect.Value) error {
	switch val.Kind() {
	case reflect.String:
		state.Buf = strconv.AppendQuote(state.Buf, val.String())
	case reflect.Float32, reflect.Float64:
		state.Buf = strconv.AppendFloat(state.Buf, val.Float(), 'f', -1, 64)
	case reflect.Bool:
		state.Buf = strconv.AppendBool(state.Buf, val.Bool())
	case reflect.Interface:
		if val.IsNil() {
			state.Buf = append(state.Buf, "null"...)
		}
		return fmt.Errorf("unsupported non-nil type: %T", val)
	case reflect.Invalid:
		state.Buf = append(state.Buf, "null"...)
	default:
		return fmt.Errorf("unsupported non-nil type: %T", val)
	}
	return nil
}

func handleRJSON(state *MarshalState) error {
	level := &state.Levels[state.CurrentLevel]
	elem := level.Value.(RJsonElement)
	switch elem.kind {
	// TODO: add maps and arrays next
	case String:
		state.Buf = strconv.AppendQuote(state.Buf, elem.str)
	case Number:
		state.Buf = append(state.Buf, elem.num.String()...)
	case True:
		state.Buf = append(state.Buf, "true"...)
	case False:
		state.Buf = append(state.Buf, "false"...)
	case Null:
		state.Buf = append(state.Buf, "null"...)
	default:
		return fmt.Errorf("unsupported RJSON type %s", elem.kind.String())
	}
	return nil
}

func handleMap(state *MarshalState) error {
	level := &state.Levels[state.CurrentLevel]
	target, ok := level.Value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot map level %d", state.CurrentLevel)
	}
	length := len(target)
	if length == 0 {
		state.Buf = append(state.Buf, "{}"...)
		finishLevel(state)
		return nil
	}
	//
	// sort the keys (if not done)
	//
	if level.Type == UNKNOWN {
		level.Keys = make([]string, 0, length)
		for k := range target {
			level.Keys = append(level.Keys, k)
		}
		sort.Sort(sort.StringSlice(level.Keys))
		level.Type = MAP
		level.CurrentIndex = 0
		state.Buf = append(state.Buf, "{\n"...)
	} else {
		// else we are already underway, so point to the next key
		if level.CurrentIndex != length-1 {
			state.Buf = append(state.Buf, ',')
		}
		state.Buf = append(state.Buf, "\n"...)
		level.CurrentIndex++
	}
	//
	// iterate
	//
	for level.CurrentIndex < len(level.Keys) {
		key := level.Keys[level.CurrentIndex]
		value := target[key]
		appendIndent(state, 1)
		state.Buf = strconv.AppendQuote(state.Buf, key)
		state.Buf = append(state.Buf, ": "...)
		rVal := reflect.ValueOf(value)
		switch rVal.Kind() {
		case reflect.String:
			handleSimple(state, rVal)
		case reflect.Map:
			addLevel(state, value)
			return nil
		default:
			return fmt.Errorf("unsupported type: %T", rVal)
		}
		if level.CurrentIndex != length-1 {
			state.Buf = append(state.Buf, ',')
		}
		state.Buf = append(state.Buf, "\n"...)
		level.CurrentIndex++
	}
	appendIndent(state, 0)
	state.Buf = append(state.Buf, "}"...)
	finishLevel(state)
	return nil
}
