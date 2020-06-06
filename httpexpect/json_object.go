package httpexpect

import (
	"fmt"
)

type JSONObjectType = map[string]interface{}

type JSONObject struct {
	expectation *Expectation

	path  string
	value JSONObjectType
}

func (o *JSONObject) HasKey(key string) *JSONObject {
	if _, ok := o.value[key]; !ok {
		o.expectation.errorf(`key "%s" do not exists`, key)
	}
	return o
}

func (o *JSONObject) Number(key string) *JSONNumber {
	if value, ok := o.value[key]; ok {
		if number, ok := value.(float64); ok {
			return &JSONNumber{
				expectation: o.expectation,
				path:        fmt.Sprintf("%s.%s", o.path, key),
				value:       number,
			}
		}
		o.expectation.fatalf(`key "%s" is not an number`, key)
		return nil
	}
	o.expectation.fatalf(`key "%s" do not exists`, key)
	return nil
}

func (o *JSONObject) String(key string) *JSONString {
	if value, ok := o.value[key]; ok {
		if str, ok := value.(string); ok {
			return &JSONString{
				expectation: o.expectation,
				path:        fmt.Sprintf("%s.%s", o.path, key),
				value:       str,
			}
		}
		o.expectation.fatalf(`key "%s" is not a string`, key)
		return nil
	}
	o.expectation.fatalf(`key "%s" do not exists`, key)
	return nil
}

func (o *JSONObject) Bool(key string) *JSONBool {
	if value, ok := o.value[key]; ok {
		if b, ok := value.(bool); ok {
			return &JSONBool{
				expectation: o.expectation,
				path:        fmt.Sprintf("%s.%s", o.path, key),
				value:       b,
			}
		}
		o.expectation.fatalf(`key "%s" is not a bool`, key)
		return nil
	}
	o.expectation.fatalf(`key "%s" do not exists`, key)
	return nil
}

func (o *JSONObject) Array(key string) *JSONArray {
	if value, ok := o.value[key]; ok {
		if array, ok := value.(JSONArrayType); ok {
			return &JSONArray{
				expectation: o.expectation,
				path:        fmt.Sprintf("%s.%s", o.path, key),
				value:       array,
			}
		}
		o.expectation.fatalf(`key "%s" is not an array`, key)
		return nil
	}
	o.expectation.fatalf(`key "%s" do not exists`, key)
	return nil
}

func (o *JSONObject) Object(key string) *JSONObject {
	if value, ok := o.value[key]; ok {
		if obj, ok := value.(JSONObjectType); ok {
			return &JSONObject{
				expectation: o.expectation,
				path:        fmt.Sprintf("%s.%s", o.path, key),
				value:       obj,
			}
		}
		o.expectation.fatalf(`key "%s" is not an object`, key)
		return nil
	}
	o.expectation.fatalf(`key "%s" do not exists`, key)
	return nil
}

func (o *JSONObject) Null(key string) {
	if value, ok := o.value[key]; ok {
		if value != nil {
			o.expectation.fatalf(`key "%s" is not null`, key)
			return
		}
		return
	}
	o.expectation.fatalf(`key "%s" do not exists`, key)
}
