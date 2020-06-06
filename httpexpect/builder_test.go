package httpexpect_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.pr0ger.dev/httpexpect"
)

func TestExpectationBuilder_WithoutBody(t *testing.T) {
	called := false
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Len(t, resp, 0)

		called = true
	})

	_ = httpexpect.Post(&testing.T{}, stubHandler).WithoutBody()

	assert.True(t, called)
}

func TestExpectationBuilder_WithPlainText(t *testing.T) {
	payload := []byte("test")

	called := false
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, payload, resp)

		called = true
	})

	_ = httpexpect.Post(&testing.T{}, stubHandler).WithPlainText(payload)

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
			stubHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, tt.serialized, resp)

				called = true
			})

			_ = httpexpect.Post(&testing.T{}, stubHandler).WithJSON(tt.data)

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

		_ = httpexpect.Post(&testing.T{}, func(http.ResponseWriter, *http.Request) {}).WithJSON(data)
	})
}
