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
			MarshalLevel{
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
		case reflect.String:
			handleSimple(state, val)
			finishLevel(state)
		case reflect.Map:
			err := handleMap(state)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %T", current.Value)
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

func handleSimple(state *MarshalState, val reflect.Value) {
	switch val.Kind() {
	case reflect.String:
		state.Buf = strconv.AppendQuote(state.Buf, val.String())
	default:
		// no nothing
	}
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
