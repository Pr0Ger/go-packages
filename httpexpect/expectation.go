package httpexpect

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
)

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
}

type Expectation struct {
	t            TestingT
	handler      http.HandlerFunc
	middlewares  []func(http.Handler) http.Handler
	method       string
	target       string
	contentType  string
	extraHeaders map[string][]string
	payload      io.Reader
	context      context.Context

	recorder *httptest.ResponseRecorder

	require bool
}

func newExpectation(t TestingT, handler http.HandlerFunc, method string) Expectation {
	return Expectation{
		t:            t,
		handler:      handler,
		method:       method,
		extraHeaders: map[string][]string{},
	}
}

func (e *Expectation) performRequest() {
	if e.recorder != nil {
		// expectation is already executed; noop
		return
	}

	if e.target == "" {
		e.target = "/"
	}

	req := httptest.NewRequest(e.method, e.target, e.payload)
	if e.context != nil {
		req = req.WithContext(e.context)
	}
	if e.contentType != "" {
		req.Header.Add("Content-Type", e.contentType)
	}
	for key, values := range e.extraHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	e.recorder = httptest.NewRecorder()

	if len(e.middlewares) == 0 {
		e.handler.ServeHTTP(e.recorder, req)
	} else {
		h := e.middlewares[len(e.middlewares)-1](e.handler)
		for i := len(e.middlewares) - 2; i >= 0; i-- {
			h = e.middlewares[i](h)
		}

		h.ServeHTTP(e.recorder, req)
	}
}

func (e Expectation) Require() Expectation {
	e.require = true
	return e
}

func (e Expectation) errorf(format string, args ...interface{}) {
	e.t.Helper()

	e.t.Errorf(format, args...)
	if e.require {
		e.t.FailNow()
	}
}

func (e Expectation) fatalf(format string, args ...interface{}) {
	e.t.Helper()

	e.t.Errorf(format, args...)
	e.t.FailNow()
}

func Get(t TestingT, handler http.HandlerFunc) Expectation {
	return newExpectation(t, handler, http.MethodGet)
}

func Head(t TestingT, handler http.HandlerFunc) Expectation {
	return newExpectation(t, handler, http.MethodHead)
}

func Post(t TestingT, handler http.HandlerFunc) Expectation {
	return newExpectation(t, handler, http.MethodPost)
}

func Put(t TestingT, handler http.HandlerFunc) Expectation {
	return newExpectation(t, handler, http.MethodPut)
}

func Patch(t TestingT, handler http.HandlerFunc) Expectation {
	return newExpectation(t, handler, http.MethodPatch)
}

func Delete(t TestingT, handler http.HandlerFunc) Expectation {
	return newExpectation(t, handler, http.MethodDelete)
}
