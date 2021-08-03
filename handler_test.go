package testix_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	testix "github.com/drykit-go/testix"
	"github.com/drykit-go/testix/check"
	"github.com/drykit-go/testix/check/expect"
)

type responseBody map[string]interface{}

// Main tests

func TestHandlerFunc(t *testing.T) {
	const (
		expCode = 200
		expBody = `{"message":"Hello World!"}`
	)

	h := handler(408, responseBody{"message": "Hello World!"})
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	testix.HandlerFunc(h, r).
		ResponseStatus(expect.String.Contains("Timeout")).
		ResponseCode(check.Int.Equal(408)).
		ResponseCode(check.Int.NotEqual(199)).
		ResponseCode(check.Int.NotInRange(201, 202)).
		ResponseCode(check.Int.InRange(400, 499)).
		Duration(check.Duration.Under(50 * time.Millisecond)).
		ResponseCode(isEven).
		ResponseBody(check.Bytes.Equal([]byte(expBody))).
		ResponseBody(
			check.Bytes.Len(check.Int.GreaterOrEqual(20)),
		).
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
	w.Header().Add("yo", "weeeesssshhhhhh")
	w.Write(b)
}

var isEven = check.NewIntCheck(
	func(got int) bool {
		return got%2 == 0
	},
	func(label string, got interface{}) string {
		return fmt.Sprintf("expect %s to be odd, got odd value %d", label, got)
	},
)
