package testx_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	testx "github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

func TestHTTPHandlerRunner(t *testing.T) {
	hf := func(w http.ResponseWriter, _ *http.Request) {
		b, _ := json.Marshal(map[string]interface{}{"message": "Hello, World!"})
		w.WriteHeader(200)
		w.Write(b)
	}
	r, _ := http.NewRequestWithContext(
		context.WithValue(context.Background(), "userID", 42), //nolint:revive,staticcheck // context-keys-type
		http.MethodGet, "/", nil,
	)

	expBody := []byte(`{"message":"Hello, World!"}`)

	t.Run("should pass", func(t *testing.T) {
		res := testx.HTTPHandlerFunc(hf).WithRequest(r).
			Duration(check.Duration.Under(100*time.Millisecond)).
			Request(check.HTTPRequest.Context(check.Context.HasKeys("userID"))).
			Response(
				check.HTTPResponse.StatusCode(check.Int.Is(200)),
				check.HTTPResponse.Body(check.Bytes.SameJSON(expBody)),
			).
			DryRun()

		exp := handlerResults{
			baseResults: baseResults{
				passed:  true,
				failed:  false,
				nPassed: 4,
				nFailed: 0,
				nChecks: 4,
				checks: []testx.CheckResult{
					{Passed: true, Reason: ""},
					{Passed: true, Reason: ""},
					{Passed: true, Reason: ""},
					{Passed: true, Reason: ""},
				},
			},
			header: http.Header{},
			status: "200 OK",
			code:   200,
			body:   expBody,
		}

		assertEqualHandlerResults(t, res, exp)
	})

	t.Run("should fail", func(t *testing.T) {
		res := testx.HTTPHandlerFunc(hf).WithRequest(r).
			Request(check.HTTPRequest.Context(check.Context.HasKeys("abc"))).
			Response(
				check.HTTPResponse.StatusCode(check.Int.Is(-1)),
				check.HTTPResponse.Body(check.Bytes.SameJSON(expBody)),
			).
			DryRun()

		exp := handlerResults{
			baseResults: baseResults{
				passed:  false,
				failed:  true,
				nPassed: 1,
				nFailed: 2,
				nChecks: 3,
				checks: []testx.CheckResult{
					{Passed: false, Reason: "http request:\nexp context to pass ContextChecker\ngot explanation: context:\nexp to have keys [abc]\ngot keys not set"},
					{Passed: false, Reason: "http response:\nexp status code to pass IntChecker\ngot explanation: status code:\nexp -1\ngot 200"},
					{Passed: true, Reason: ""},
				},
			},
			header: http.Header{},
			status: "200 OK",
			code:   200,
			body:   expBody,
		}

		assertEqualHandlerResults(t, res, exp)
	})
}

// Helpers

type handlerResults struct {
	baseResults
	header http.Header
	status string
	code   int
	body   []byte
	// duration time.Duration // cannot predict exact duration
}

func assertEqualHandlerResults(t *testing.T, res testx.HandlerResulter, exp handlerResults) {
	t.Helper()
	if got := toHandlerResults(res); !deq(got, exp) {
		failBadResults(t, "handlerResults", got, exp)
	}
}

func toHandlerResults(res testx.HandlerResulter) handlerResults {
	return handlerResults{
		baseResults: toBaseResults(res),
		header:      res.ResponseHeader(),
		status:      res.ResponseStatus(),
		code:        res.ResponseCode(),
		body:        res.ResponseBody(),
		// duration:    res.ResponseDuration(), // cannot predict exact duration
	}
}
