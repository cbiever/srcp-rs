package test

import (
	"net/http/httptest"
	"strings"
)

func sendRequest(method string, path string, payload string) *httptest.ResponseRecorder {
	createRouter()

	r := httptest.NewRequest(method, path, strings.NewReader(payload))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	return w
}
