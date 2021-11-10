package httpexpect

import (
	"bytes"
	"context"
	"encoding/json"
)

func (e Expectation) WithContext(ctx context.Context) Expectation {
	if e.recorder != nil {
		panic("handler is already invoked")
	}

	e.context = ctx

	return e
}

func (e Expectation) WithQuery(query string) Expectation {
	if e.recorder != nil {
		panic("handler is already invoked")
	}

	e.target = query

	return e
}

func (e Expectation) WithoutBody() Expectation {
	if e.recorder != nil {
		panic("handler is already invoked")
	}

	return e
}

func (e Expectation) WithPlainText(data []byte) Expectation {
	if e.recorder != nil {
		panic("handler is already invoked")
	}

	e.contentType = "text/plain"
	e.payload = bytes.NewBuffer(data)

	return e
}

func (e Expectation) WithJSON(data interface{}) Expectation {
	if e.recorder != nil {
		panic("handler is already invoked")
	}

	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	e.contentType = "application/json"
	e.payload = bytes.NewBuffer(body)
	return e
}
