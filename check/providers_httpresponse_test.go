package check_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drykit-go/testx/check"
)

//nolint:bodyclose
func TestHTTPResponseCheckerProvider(t *testing.T) {
	newResp := func() *http.Response {
		rr := httptest.NewRecorder()
		rq := httptest.NewRequest("GET", "/", nil)
		body, _ := json.Marshal(map[string]any{"answer": 42})
		clen := len(body)

		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTeapot)
			w.Header().Set("Content-Type", "application/json")
			w.Write(body)
		}(rr, rq)

		resp := rr.Result()
		resp.Header = rr.Header()
		resp.ContentLength = int64(clen)
		return resp
	}

	var (
		expStatusCode    = 418
		expStatus        = "418 I'm a teapot"
		expContentLength = 13
		expBody          = []byte(`{"answer":42}`)
	)

	t.Run("StatusCode pass", func(t *testing.T) {
		c := check.HTTPResponse.StatusCode(check.Int.Is(expStatusCode))
		assertPassChecker(t, "HTTPResponse.StatusCode", c, newResp())
	})

	t.Run("StatusCode fail", func(t *testing.T) {
		c := check.HTTPResponse.StatusCode(check.Int.Not(expStatusCode))
		assertFailChecker(t, "HTTPResponse.StatusCode", c, newResp(), makeExpl(
			"status code to pass Checker[int]",
			"explanation: status code:\n"+makeExpl(
				"not "+fmt.Sprint(expStatusCode),
				fmt.Sprint(expStatusCode),
			),
		))
	})

	t.Run("Status pass", func(t *testing.T) {
		c := check.HTTPResponse.Status(check.String.Is(expStatus))
		assertPassChecker(t, "HTTPResponse.Status", c, newResp())
	})

	t.Run("Status fail", func(t *testing.T) {
		c := check.HTTPResponse.Status(check.String.Not(expStatus))
		assertFailChecker(t, "HTTPResponse.Status", c, newResp(), makeExpl(
			"status to pass Checker[string]",
			"explanation: status:\n"+makeExpl(
				"not "+expStatus,
				expStatus,
			),
		))
	})

	t.Run("ContentLength pass", func(t *testing.T) {
		c := check.HTTPResponse.ContentLength(check.Int.Is(expContentLength))
		assertPassChecker(t, "HTTPResponse.ContentLength", c, newResp())
	})

	t.Run("ContentLength fail", func(t *testing.T) {
		c := check.HTTPResponse.ContentLength(check.Int.Not(expContentLength))
		assertFailChecker(t, "HTTPResponse.ContentLength", c, newResp(), makeExpl(
			"content length to pass Checker[int]",
			"explanation: content length:\n"+makeExpl(
				"not "+fmt.Sprint(expContentLength),
				fmt.Sprint(expContentLength),
			),
		))
	})

	t.Run("Header pass", func(t *testing.T) {
		c := check.HTTPResponse.Header(check.HTTPHeader.HasKey("Content-Type"))
		assertPassChecker(t, "HTTPResponse.Header", c, newResp())
	})

	t.Run("Header fail", func(t *testing.T) {
		c := check.HTTPResponse.Header(check.HTTPHeader.HasNotKey("Content-Type"))
		resp := newResp()
		assertFailChecker(t, "HTTPResponse.Header", c, resp, makeExpl(
			"header to pass Checker[http.Header]",
			"explanation: http.Header:\n"+makeExpl(
				`not to have key "Content-Type"`,
				fmt.Sprint(resp.Header),
			),
		))
	})

	t.Run("Body pass", func(t *testing.T) {
		c := check.HTTPResponse.Body(check.Bytes.Is(expBody))
		assertPassChecker(t, "HTTPResponse.Body", c, newResp())
	})

	t.Run("Body fail", func(t *testing.T) {
		c := check.HTTPResponse.Body(check.Bytes.Not(expBody))
		assertFailChecker(t, "HTTPResponse.Body", c, newResp(), makeExpl(
			"body to pass Checker[[]byte]",
			"explanation: bytes:\n"+makeExpl(
				"not "+fmt.Sprint(expBody),
				fmt.Sprint(expBody),
			),
		))
	})
}
