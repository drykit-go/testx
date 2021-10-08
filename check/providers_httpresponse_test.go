package check_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drykit-go/testx/check"
)

//nolint:bodyclose
func TestHTTPResponseCheckerProvider(t *testing.T) {
	newResponse := func() *http.Response {
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
		assertPassHTTPResponseChecker(t, "StatusCode", c, newResponse())
	})

	t.Run("StatusCode fail", func(t *testing.T) {
		c := check.HTTPResponse.StatusCode(check.Int.Not(expStatusCode))
		assertFailHTTPResponseChecker(t, "StatusCode", c, newResponse())
	})

	t.Run("Status pass", func(t *testing.T) {
		c := check.HTTPResponse.Status(check.String.Is(expStatus))
		assertPassHTTPResponseChecker(t, "Status", c, newResponse())
	})

	t.Run("Status fail", func(t *testing.T) {
		c := check.HTTPResponse.Status(check.String.Not(expStatus))
		assertFailHTTPResponseChecker(t, "Status", c, newResponse())
	})

	t.Run("ContentLength pass", func(t *testing.T) {
		c := check.HTTPResponse.ContentLength(check.Int.Is(expContentLength))
		assertPassHTTPResponseChecker(t, "ContentLength", c, newResponse())
	})

	t.Run("ContentLength fail", func(t *testing.T) {
		c := check.HTTPResponse.ContentLength(check.Int.Not(expContentLength))
		assertFailHTTPResponseChecker(t, "ContentLength", c, newResponse())
	})

	t.Run("Header pass", func(t *testing.T) {
		c := check.HTTPResponse.Header(check.HTTPHeader.HasKey("Content-Type"))
		assertPassHTTPResponseChecker(t, "Header", c, newResponse())
	})

	t.Run("Header fail", func(t *testing.T) {
		c := check.HTTPResponse.Header(check.HTTPHeader.HasNotKey("Content-Type"))
		assertFailHTTPResponseChecker(t, "Header", c, newResponse())
	})

	t.Run("Body pass", func(t *testing.T) {
		c := check.HTTPResponse.Body(check.Bytes.Is(expBody))
		assertPassHTTPResponseChecker(t, "Body", c, newResponse())
	})

	t.Run("Body fail", func(t *testing.T) {
		c := check.HTTPResponse.Body(check.Bytes.Not(expBody))
		assertFailHTTPResponseChecker(t, "Body", c, newResponse())
	})
}

// Helpers

func assertPassHTTPResponseChecker(t *testing.T, method string, c check.HTTPResponseChecker, resp *http.Response) {
	t.Helper()
	if !c.Pass(resp) {
		failHTTPResponseCheckerTest(t, true, method, resp, c.Explain)
	}
}

func assertFailHTTPResponseChecker(t *testing.T, method string, c check.HTTPResponseChecker, resp *http.Response) {
	t.Helper()
	if c.Pass(resp) {
		failHTTPResponseCheckerTest(t, false, method, resp, c.Explain)
	}
}

func failHTTPResponseCheckerTest(t *testing.T, expPass bool, method string, resp *http.Response, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "HTTPResponse."+method, explain("HTTPResponse value", resp))
}
