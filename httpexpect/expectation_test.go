package httpexpect_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pr0ger.dev/x/httpexpect"
)

func TestMethod(t *testing.T) {
	tests := []struct {
		builderFunc func(t httpexpect.TestingT, handler http.HandlerFunc) httpexpect.Expectation
		method      string
	}{
		{httpexpect.Get, http.MethodGet},
		{httpexpect.Head, http.MethodHead},
		{httpexpect.Post, http.MethodPost},
		{httpexpect.Put, http.MethodPut},
		{httpexpect.Patch, http.MethodPatch},
		{httpexpect.Delete, http.MethodDelete},
	}

	dumbT := &testing.T{}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			called := false
			stubHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.method, r.Method)
				called = true
			})

			_ = tt.builderFunc(dumbT, stubHandler).Status(http.StatusOK)

			require.True(t, called)
		})
	}
}
