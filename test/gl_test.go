package test

import (
	"encoding/json"
	"io/ioutil"
	"srcp-rs/handlers"
	"strings"
	"testing"
)

const create_gl_json = `{
   "data": {
      "id": "1-3",
      "attributes": {
         "name": null,
         "address": 3,
         "protocol": "N",
         "protocol-version": 1,
         "decoder-speed-steps": 28,
         "number-of-decoder-functions": 4,
         "drivemode": 1,
         "v": 0,
         "v-max": 28,
         "functions": [
         ]
      },
      "relationships": {
         "bus": {
            "data": {
               "type": "buses",
               "id": "1"
            }
         }
      },
      "type": "gls"
   }
}`

const gl_init = "INIT 2 GL 3 N 1 28 4"

func Test1CreatGL(t *testing.T) {
	sendAndReceive := func(request string) string {
		if strings.Compare(gl_init, request) != 0 {
			t.Fatalf("expected: %s got: %s", gl_init, request)
		}
		return "0000000000.000 200 OK"
	}
	handlers.GetStore().SaveConnection(1, &MockSrcpConnection{t, sendAndReceive})

	w := sendRequest("POST", "/sessions/1/buses/2/gls", create_gl_json)

	if 200 != w.Code {
		t.Fatalf("expected: %d got: %d", 200, w.Code)
	}
}

func Test2CreatGL(t *testing.T) {
	sendAndReceive := func(request string) string {
		if strings.Compare(gl_init, request) != 0 {
			t.Fatalf("expected: %s got: %s", gl_init, request)
		}
		return "0000000000.000 400 ERROR unknown"
	}
	handlers.GetStore().SaveConnection(1, &MockSrcpConnection{t, sendAndReceive})

	w := sendRequest("POST", "/sessions/1/buses/2/gls", create_gl_json)

	if 400 != w.Code {
		t.Fatalf("expected: %d got: %d", 400, w.Code)
	}

	body, _ := ioutil.ReadAll(w.Body)
	var srcpError handlers.SrcpError
	_ = json.Unmarshal(body, &srcpError)

	if 400 != srcpError.Code || strings.Compare("ERROR", srcpError.Status) != 0 || strings.Compare("unknown", srcpError.Text) != 0 {
		t.Fatal("Srcp error is not as expected")
	}
}
