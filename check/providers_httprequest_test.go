package check_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestHTTPRequestCheckerProvider(t *testing.T) {
	newContext := func(key, val interface{}) context.Context {
		return context.WithValue(context.Background(), key, val)
	}
	newRequest := func() *http.Request {
		ctx := newContext("userID", 42)
		body, _ := json.Marshal(map[string]interface{}{"answer": 42})
		r, _ := http.NewRequestWithContext(ctx, "GET", "/endpoint?id=42", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		return r
	}

	var (
		expContentLength = 13
		expBody          = []byte(`{"answer":42}`)
		expCtxKey        = "userID"
		expCtxVal        = 42
	)

	t.Run("ContentLength pass", func(t *testing.T) {
		c := check.HTTPRequest.ContentLength(check.Int.Is(expContentLength))
		assertPassHTTPRequestChecker(t, "ContentLength", c, newRequest())
	})

	t.Run("ContentLength fail", func(t *testing.T) {
		c := check.HTTPRequest.ContentLength(check.Int.Not(expContentLength))
		assertFailHTTPRequestChecker(t, "ContentLength", c, newRequest())
	})

	t.Run("Header pass", func(t *testing.T) {
		c := check.HTTPRequest.Header(check.HTTPHeader.HasKey("Content-Type"))
		assertPassHTTPRequestChecker(t, "Header", c, newRequest())
	})

	t.Run("Header fail", func(t *testing.T) {
		c := check.HTTPRequest.Header(check.HTTPHeader.HasNotKey("Content-Type"))
		assertFailHTTPRequestChecker(t, "Header", c, newRequest())
	})

	t.Run("Body pass", func(t *testing.T) {
		c := check.HTTPRequest.Body(check.Bytes.Is(expBody))
		assertPassHTTPRequestChecker(t, "Body", c, newRequest())
	})

	t.Run("Body fail", func(t *testing.T) {
		c := check.HTTPRequest.Body(check.Bytes.Not(expBody))
		assertFailHTTPRequestChecker(t, "Body", c, newRequest())
	})

	t.Run("Context pass", func(t *testing.T) {
		c := check.HTTPRequest.Context(check.Context.Value(expCtxKey, check.Value.Is(expCtxVal)))
		assertPassHTTPRequestChecker(t, "Context", c, newRequest())
	})

	t.Run("Context fail", func(t *testing.T) {
		c := check.HTTPRequest.Context(check.Context.Value(expCtxKey, check.Value.Not(expCtxVal)))
		assertFailHTTPRequestChecker(t, "Context", c, newRequest())
	})
}

// Helpers

func assertPassHTTPRequestChecker(t *testing.T, method string, c check.HTTPRequestChecker, r *http.Request) {
	t.Helper()
	if !c.Pass(r) {
		failHTTPRequestCheckerTest(t, true, method, r, c.Explain)
	}
}

func assertFailHTTPRequestChecker(t *testing.T, method string, c check.HTTPRequestChecker, r *http.Request) {
	t.Helper()
	if c.Pass(r) {
		failHTTPRequestCheckerTest(t, false, method, r, c.Explain)
	}
}

func failHTTPRequestCheckerTest(t *testing.T, expPass bool, method string, r *http.Request, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "HTTPRequest."+method, explain("HTTPRequest value", r))
}
