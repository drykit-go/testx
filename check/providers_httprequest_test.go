package check_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestHTTPRequestCheckerProvider(t *testing.T) {
	newCtx := func(key, val interface{}) context.Context {
		return context.WithValue(context.Background(), key, val)
	}
	newReq := func() *http.Request {
		ctx := newCtx("userID", 42)
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
		assertPassHTTPRequestChecker(t, "ContentLength", c, newReq())
	})

	t.Run("ContentLength fail", func(t *testing.T) {
		c := check.HTTPRequest.ContentLength(check.Int.Not(expContentLength))
		assertFailHTTPRequestChecker(t, "ContentLength", c, newReq(), makeExpl(
			"content length to pass IntChecker",
			fmt.Sprintf(
				"explanation: content length:\nexp not %d\ngot %d",
				expContentLength, expContentLength,
			),
		))
	})

	t.Run("Header pass", func(t *testing.T) {
		c := check.HTTPRequest.Header(check.HTTPHeader.HasKey("Content-Type"))
		assertPassHTTPRequestChecker(t, "Header", c, newReq())
	})

	t.Run("Header fail", func(t *testing.T) {
		c := check.HTTPRequest.Header(check.HTTPHeader.HasNotKey("Content-Type"))
		r := newReq()
		assertFailHTTPRequestChecker(t, "Header", c, r, makeExpl(
			"header to pass HTTPHeaderChecker",
			fmt.Sprintf(
				"explanation: http.Header:\nexp not to have key \"Content-Type\"\ngot %v",
				r.Header,
			),
		))
	})

	t.Run("Body pass", func(t *testing.T) {
		c := check.HTTPRequest.Body(check.Bytes.Is(expBody))
		assertPassHTTPRequestChecker(t, "Body", c, newReq())
	})

	t.Run("Body fail", func(t *testing.T) {
		c := check.HTTPRequest.Body(check.Bytes.Not(expBody))
		assertFailHTTPRequestChecker(t, "Body", c, newReq(), makeExpl(
			"body to pass BytesChecker",
			"explanation: bytes:\n"+makeExpl(
				"not "+fmt.Sprint(expBody),
				fmt.Sprint(expBody),
			),
		))
	})

	t.Run("Context pass", func(t *testing.T) {
		c := check.HTTPRequest.Context(check.Context.Value(expCtxKey, check.Value.Is(expCtxVal)))
		assertPassHTTPRequestChecker(t, "Context", c, newReq())
	})

	t.Run("Context fail", func(t *testing.T) {
		c := check.HTTPRequest.Context(check.Context.Value(expCtxKey, check.Value.Not(expCtxVal)))
		assertFailHTTPRequestChecker(t, "Context", c, newReq(), makeExpl(
			"context to pass ContextChecker",
			"explanation: context:\n"+makeExpl(
				"value for key userID to pass ValueChecker",
				"explanation: value:\n"+makeExpl(
					"not 42",
					"42",
				),
			),
		))
	})
}

// Helpers

func assertPassHTTPRequestChecker(t *testing.T, method string, c check.HTTPRequestChecker, r *http.Request) {
	t.Helper()
	if !c.Pass(r) {
		failHTTPRequestCheckerTest(t, true, method, r, c.Explain)
	}
}

func assertFailHTTPRequestChecker(t *testing.T, method string, c check.HTTPRequestChecker, r *http.Request, expexpl string) {
	t.Helper()
	if c.Pass(r) {
		failHTTPRequestCheckerTest(t, false, method, r, c.Explain)
	}
	assertGoodExplain(t, c, r, expexpl)
}

func failHTTPRequestCheckerTest(t *testing.T, expPass bool, method string, r *http.Request, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "HTTPRequest."+method, explain("HTTPRequest value", r))
}
