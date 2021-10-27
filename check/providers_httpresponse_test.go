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
		body, _ := json.Marshal(map[string]interface{}{"answer": 42})
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
		assertPassHTTPResponseChecker(t, "StatusCode", c, newResp())
	})

	t.Run("StatusCode fail", func(t *testing.T) {
		c := check.HTTPResponse.StatusCode(check.Int.Not(expStatusCode))
		assertFailHTTPResponseChecker(t, "StatusCode", c, newResp(), makeExpl(
			"status code to pass IntChecker",
			"explanation: status code:\n"+makeExpl(
				"not "+fmt.Sprint(expStatusCode),
				fmt.Sprint(expStatusCode),
			),
		))
	})

	t.Run("Status pass", func(t *testing.T) {
		c := check.HTTPResponse.Status(check.String.Is(expStatus))
		assertPassHTTPResponseChecker(t, "Status", c, newResp())
	})

	t.Run("Status fail", func(t *testing.T) {
		c := check.HTTPResponse.Status(check.String.Not(expStatus))
		assertFailHTTPResponseChecker(t, "Status", c, newResp(), makeExpl(
			"status to pass StringChecker",
			"explanation: status:\n"+makeExpl(
				"not "+expStatus,
				expStatus,
			),
		))
	})

	t.Run("ContentLength pass", func(t *testing.T) {
		c := check.HTTPResponse.ContentLength(check.Int.Is(expContentLength))
		assertPassHTTPResponseChecker(t, "ContentLength", c, newResp())
	})

	t.Run("ContentLength fail", func(t *testing.T) {
		c := check.HTTPResponse.ContentLength(check.Int.Not(expContentLength))
		assertFailHTTPResponseChecker(t, "ContentLength", c, newResp(), makeExpl(
			"content length to pass IntChecker",
			"explanation: content length:\n"+makeExpl(
				"not "+fmt.Sprint(expContentLength),
				fmt.Sprint(expContentLength),
			),
		))
	})

	t.Run("Header pass", func(t *testing.T) {
		c := check.HTTPResponse.Header(check.HTTPHeader.HasKey("Content-Type"))
		assertPassHTTPResponseChecker(t, "Header", c, newResp())
	})

	t.Run("Header fail", func(t *testing.T) {
		c := check.HTTPResponse.Header(check.HTTPHeader.HasNotKey("Content-Type"))
		resp := newResp()
		assertFailHTTPResponseChecker(t, "Header", c, resp, makeExpl(
			"header to pass HTTPHeaderChecker",
			"explanation: http.Header:\n"+makeExpl(
				`not to have key "Content-Type"`,
				fmt.Sprint(resp.Header),
			),
		))
	})

	t.Run("Body pass", func(t *testing.T) {
		c := check.HTTPResponse.Body(check.Bytes.Is(expBody))
		assertPassHTTPResponseChecker(t, "Body", c, newResp())
	})

	t.Run("Body fail", func(t *testing.T) {
		c := check.HTTPResponse.Body(check.Bytes.Not(expBody))
		assertFailHTTPResponseChecker(t, "Body", c, newResp(), makeExpl(
			"body to pass BytesChecker",
			"explanation: bytes:\n"+makeExpl(
				"not "+fmt.Sprint(expBody),
				fmt.Sprint(expBody),
			),
		))
	})
}

// Helpers

func assertPassHTTPResponseChecker(t *testing.T, method string, c check.HTTPResponseChecker, resp *http.Response) {
	t.Helper()
	if !c.Pass(resp) {
		failHTTPResponseCheckerTest(t, true, method, resp, c.Explain)
	}
}

func assertFailHTTPResponseChecker(t *testing.T, method string, c check.HTTPResponseChecker, resp *http.Response, expexpl string) {
	t.Helper()
	if c.Pass(resp) {
		failHTTPResponseCheckerTest(t, false, method, resp, c.Explain)
	}
	assertGoodExplain(t, c, resp, expexpl)
}

func failHTTPResponseCheckerTest(t *testing.T, expPass bool, method string, resp *http.Response, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "HTTPResponse."+method, explain("HTTPResponse value", resp))
}
