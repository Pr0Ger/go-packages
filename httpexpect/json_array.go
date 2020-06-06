package httpexpect

import (
	"fmt"
)

type JSONArrayType = []interface{}

type JSONArray struct {
	expectation *Expectation

	path  string
	value []interface{}
}

func (a *JSONArray) Len() *JSONNumber {
	return &JSONNumber{
		expectation: a.expectation,
		path:        fmt.Sprintf("%s.$len", a.path),
		value:       float64(len(a.value)),
	}
}

func (a *JSONArray) Number(idx int) *JSONNumber {
	if idx < len(a.value) {
		if number, ok := a.value[idx].(float64); ok {
			return &JSONNumber{
				expectation: a.expectation,
				path:        fmt.Sprintf("%s.$%d", a.path, idx),
				value:       number,
			}
		}
		a.expectation.fatalf(`element at %d is not a number`, idx)
		return nil
	}
	a.expectation.fatalf(`index %d is out of bounds (len=%d)`, idx, len(a.value))
	return nil
}

func (a *JSONArray) String(idx int) *JSONString {
	if idx < len(a.value) {
		if str, ok := a.value[idx].(string); ok {
			return &JSONString{
				expectation: a.expectation,
				path:        fmt.Sprintf("%s.$%d", a.path, idx),
				value:       str,
			}
		}
		a.expectation.fatalf(`element at %d is not a string`, idx)
		return nil
	}
	a.expectation.fatalf(`index %d is out of bounds (len=%d)`, idx, len(a.value))
	return nil
}

func (a *JSONArray) Bool(idx int) *JSONBool {
	if idx < len(a.value) {
		if b, ok := a.value[idx].(bool); ok {
			return &JSONBool{
				expectation: a.expectation,
				path:        fmt.Sprintf("%s.$%d", a.path, idx),
				value:       b,
			}
		}
		a.expectation.fatalf(`element at %d is not a bool`, idx)
		return nil
	}
	a.expectation.fatalf(`index %d is out of bounds (len=%d)`, idx, len(a.value))
	return nil
}

func (a *JSONArray) Array(idx int) *JSONArray {
	if idx < len(a.value) {
		if array, ok := a.value[idx].(JSONArrayType); ok {
			return &JSONArray{
				expectation: a.expectation,
				path:        fmt.Sprintf("%s.$%d", a.path, idx),
				value:       array,
			}
		}
		a.expectation.fatalf(`element at %d is not an array`, idx)
		return nil
	}
	a.expectation.fatalf(`index %d is out of bounds (len=%d)`, idx, len(a.value))
	return nil
}

func (a *JSONArray) Object(idx int) *JSONObject {
	if idx < len(a.value) {
		if obj, ok := a.value[idx].(JSONObjectType); ok {
			return &JSONObject{
				expectation: a.expectation,
				path:        fmt.Sprintf("%s.$%d", a.path, idx),
				value:       obj,
			}
		}
		a.expectation.fatalf(`element at %d is not an object`, idx)
		return nil
	}
	a.expectation.fatalf(`index %d is out of bounds (len=%d)`, idx, len(a.value))
	return nil
}

func (a *JSONArray) Null(idx int) {
	if idx < len(a.value) {
		if a.value[idx] != nil {
			a.expectation.fatalf(`element at %d is not a null`, idx)
			return
		}
		return
	}
	a.expectation.fatalf(`index %d is out of bounds (len=%d)`, idx, len(a.value))
}
