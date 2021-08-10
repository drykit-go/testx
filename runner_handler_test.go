package testx_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	testx "github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

type responseBody map[string]interface{}

// Main tests

func TestHandlerRunner(t *testing.T) {
	const (
		expCode = 200
		expBody = `{   "message"  :  "Hello World!"   }`
	)

	h := handler(408, responseBody{"message": "Hello World!"})
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	testx.HandlerFunc(h, r).
		ResponseStatus(
			check.String.Contains("Timeout"),
			check.String.NotContains("OK"),
		).
		ResponseCode(
			check.Int.Equal(408),
			check.Int.NotEqual(200),
			check.Int.NotInRange(200, 299),
			check.Int.InRange(400, 499),
			checkIntIsEven,
		).
		ResponseBody(
			check.Bytes.EqualJSON([]byte(expBody)),
			check.Bytes.Len(check.Int.GreaterOrEqual(20)),
		).
		ResponseHeader(
			check.HTTPHeader.ValueOf("marcel", check.String.Equal("patulacci")),
			check.HTTPHeader.ValueNotSet("secret"),
			check.HTTPHeader.KeySet("API_KEY"),
			check.HTTPHeader.KeyNotSet("password"),
		).
		Duration(check.Duration.Under(50 * time.Millisecond)).
		Run(t)
}

// Helpers

func handler(respCode int, respBody interface{}, funcs ...func()) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		for _, f := range funcs {
			f()
		}
		respond(w, respCode, respBody)
	}
}

func respond(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	b, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Add("marcel", "patulacci")
	w.Header()["API_KEY"] = []string{"abc"}
	w.Write(b)
}

var checkIntIsEven = check.NewIntCheck(
	func(got int) bool {
		return got%2 == 0
	},
	func(label string, got interface{}) string {
		return fmt.Sprintf("expect %s to be odd, got odd value %d", label, got)
	},
)
