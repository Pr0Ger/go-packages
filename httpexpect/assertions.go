package httpexpect

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (e Expectation) JSONArray() *JSONArray {
	e.t.Helper()

	var value interface{}
	if err := json.Unmarshal(e.recorder.Body.Bytes(), &value); err != nil {
		e.fatalf("json unmarshall failed: %v", err)
		return nil
	}

	if obj, ok := value.(JSONArrayType); ok {
		return &JSONArray{
			expectation: &e,
			path:        ".",
			value:       obj,
		}
	}
	e.fatalf("expected array, got: %#v", value)

	return nil
}

func (e Expectation) JSONObject() *JSONObject {
	e.t.Helper()

	var value interface{}
	if err := json.Unmarshal(e.recorder.Body.Bytes(), &value); err != nil {
		e.fatalf("json unmarshall failed: %v", err)
		return nil
	}

	if obj, ok := value.(JSONObjectType); ok {
		return &JSONObject{
			expectation: &e,
			path:        ".",
			value:       obj,
		}
	}
	e.fatalf("expected object, got: %#v", value)

	return nil
}

func (e Expectation) NoContent() Expectation {
	e.t.Helper()

	if e.recorder.Body.Len() != 0 {
		e.errorf("Body is not empty")
	}

	return e
}

func (e Expectation) Status(status int) Expectation {
	e.t.Helper()

	if e.recorder.Code != status {
		e.errorf(fmt.Sprintf("\n"+
			"Not equal: \n"+
			" expected: %s\n"+
			"   actual: %s",
			http.StatusText(status),
			http.StatusText(e.recorder.Code)),
		)
	}

	return e
}
