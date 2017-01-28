package test

import (
	"srcp-rs/handlers"
	"strings"
	"testing"
)

const update_cv_json = `{
  "type": "CV",
  "values": [
    42
  ]
}`

const sm_init = "INIT 2 SM NMRA"
const cv_update = "SET 2 SM 3 CV 4 42"
const sm_term = "TERM 2 SM"

func TestUpdateCV(t *testing.T) {
	var expectedRequests = []string{sm_init, cv_update, sm_term}
	var numberOfRequests int
	sendAndReceive := func(request string) string {
		if strings.Compare(expectedRequests[numberOfRequests], request) == 0 {
			numberOfRequests++
			return "0000000000.000 200 OK"
		}
		t.Fatalf("expected: %s got: %s", expectedRequests[numberOfRequests], request)
		return ""
	}
	handlers.GetStore().SaveConnection(1, &MockSrcpConnection{t, sendAndReceive})

	w := sendRequest("PUT", "/sessions/1/buses/2/gls/3/cvs/4", update_cv_json)

	if 200 != w.Code {
		t.Fatalf("expected: %d got: %d", 200, w.Code)
	}

	if numberOfRequests != 3 {
		t.Fatalf("expected: %d requests, got: %d", len(expectedRequests), numberOfRequests)
	}
}
