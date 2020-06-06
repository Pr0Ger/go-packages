package httpexpect

import (
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
	t       TestingT
	handler http.HandlerFunc

	recorder *httptest.ResponseRecorder

	require bool
}

func newExpectation(t TestingT, handler http.HandlerFunc, method string, payload io.Reader) Expectation {
	req := httptest.NewRequest(method, "/", payload)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	return Expectation{
		t:        t,
		handler:  handler,
		recorder: recorder,
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
	return newExpectation(t, handler, http.MethodGet, nil)
}

func Head(t TestingT, handler http.HandlerFunc) Expectation {
	return newExpectation(t, handler, http.MethodHead, nil)
}

func Post(t TestingT, handler http.HandlerFunc) ExpectationBuilder {
	return newExpectationBuilder(t, handler, http.MethodPost)
}

func Put(t TestingT, handler http.HandlerFunc) ExpectationBuilder {
	return newExpectationBuilder(t, handler, http.MethodPut)
}

func Patch(t TestingT, handler http.HandlerFunc) ExpectationBuilder {
	return newExpectationBuilder(t, handler, http.MethodPatch)
}

func Delete(t TestingT, handler http.HandlerFunc) ExpectationBuilder {
	return newExpectationBuilder(t, handler, http.MethodDelete)
}
