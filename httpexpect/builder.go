package httpexpect

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ExpectationBuilder struct {
	t       TestingT
	handler http.HandlerFunc
	method  string
}

func newExpectationBuilder(t TestingT, handler http.HandlerFunc, method string) ExpectationBuilder {
	return ExpectationBuilder{
		t:       t,
		handler: handler,
		method:  method,
	}
}

func (b ExpectationBuilder) WithoutBody() Expectation {
	return newExpectation(b.t, b.handler, b.method, nil)
}

func (b ExpectationBuilder) WithPlainText(data []byte) Expectation {
	return newExpectation(b.t, b.handler, b.method, bytes.NewBuffer(data))
}

func (b ExpectationBuilder) WithJSON(data interface{}) Expectation {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return newExpectation(b.t, b.handler, b.method, bytes.NewBuffer(body))
}
