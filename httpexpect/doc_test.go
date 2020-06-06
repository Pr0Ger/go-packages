package httpexpect

import (
	"net/http"
	"testing"
)

func ExampleGet() {
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"bool":  true, "number": 1337, "string":  "str", "object": {}, "array": []}`))
	})

	t := &testing.T{}

	Get(t, stubHandler).Status(http.StatusOK).JSONObject().Bool("bool").True()
}
