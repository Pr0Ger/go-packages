package httpexpect_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.pr0ger.dev/x/httpexpect"
)

func TestExpectation_WithMiddlewares(t *testing.T) {
	calledMiddleware, calledHandler := false, false
	stubHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calledHandler = true
	})

	_ = httpexpect.Post(&testing.T{}, stubHandler).
		WithMiddlewares(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calledMiddleware = true
				next.ServeHTTP(w, r)
			})
		}).
		Status(http.StatusOK)

	assert.True(t, calledMiddleware)
	assert.True(t, calledHandler)
}

func TestExpectation_WithExtraHeader(t *testing.T) {
	called := false
	stubHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("X-Extra-Header")
		assert.Equal(t, "value", header)

		called = true
	})

	_ = httpexpect.Post(&testing.T{}, stubHandler).
		WithExtraHeader("X-Extra-Header", "value").
		Status(http.StatusOK)

	assert.True(t, called)
}

func TestExpectationBuilder_WithoutBody(t *testing.T) {
	called := false
	stubHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		resp, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Empty(t, resp)

		called = true
	})

	_ = httpexpect.Post(&testing.T{}, stubHandler).WithoutBody().Status(http.StatusOK)

	assert.True(t, called)
}

func TestExpectationBuilder_WithPlainText(t *testing.T) {
	payload := []byte("test")

	called := false
	stubHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		resp, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, payload, resp)

		called = true
	})

	_ = httpexpect.Post(&testing.T{}, stubHandler).WithPlainText(payload).Status(http.StatusOK)

	assert.True(t, called)
}

func TestExpectationBuilder_WithJSON(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		serialized []byte
	}{
		{"map", map[string]int{"key": 1}, []byte(`{"key":1}`)},
		{"list", []string{"key1", "key2"}, []byte(`["key1","key2"]`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			stubHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				resp, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.serialized, resp)

				called = true
			})

			_ = httpexpect.Post(&testing.T{}, stubHandler).WithJSON(tt.data).Status(http.StatusOK)

			assert.True(t, called)
		})
	}
}

type marshalError struct {
	P *marshalError
}

func TestExpectationBuilder_WithInvalidJSON(t *testing.T) {
	assert.Panics(t, func() {
		// json.Marshall can't handle recursive structs and will return an error therefore WithJSON will panic
		data := marshalError{}
		data.P = &data

		_ = httpexpect.
			Post(&testing.T{}, func(http.ResponseWriter, *http.Request) {}).
			WithJSON(data).
			Status(http.StatusOK)
	})
}
