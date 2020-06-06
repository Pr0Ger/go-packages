package httpexpect_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pr0ger.dev/httpexpect"
)

func TestMethod(t *testing.T) {
	tests := []struct {
		simpleFunk  func(t httpexpect.TestingT, handler http.HandlerFunc) httpexpect.Expectation
		builderFunk func(t httpexpect.TestingT, handler http.HandlerFunc) httpexpect.ExpectationBuilder
		method      string
	}{
		{httpexpect.Get, nil, http.MethodGet},
		{httpexpect.Head, nil, http.MethodHead},
		{nil, httpexpect.Post, http.MethodPost},
		{nil, httpexpect.Put, http.MethodPut},
		{nil, httpexpect.Patch, http.MethodPatch},
		{nil, httpexpect.Delete, http.MethodDelete},
	}

	dumbT := &testing.T{}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			called := false
			stubHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.method, r.Method)
				called = true
			})

			if tt.simpleFunk != nil {
				_ = tt.simpleFunk(dumbT, stubHandler)
			}
			if tt.builderFunk != nil {
				_ = tt.builderFunk(dumbT, stubHandler).WithoutBody()
			}

			require.True(t, called)
		})
	}
}
