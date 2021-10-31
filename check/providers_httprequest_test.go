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
		assertPassChecker(t, "HTTPRequest.ContentLength", c, newReq())
	})

	t.Run("ContentLength fail", func(t *testing.T) {
		c := check.HTTPRequest.ContentLength(check.Int.Not(expContentLength))
		assertFailChecker(t, "HTTPRequest.ContentLength", c, newReq(), makeExpl(
			"content length to pass Checker[int]",
			fmt.Sprintf(
				"explanation: content length:\nexp not %d\ngot %d",
				expContentLength, expContentLength,
			),
		))
	})

	t.Run("Header pass", func(t *testing.T) {
		c := check.HTTPRequest.Header(check.HTTPHeader.HasKey("Content-Type"))
		assertPassChecker(t, "HTTPRequest.Header", c, newReq())
	})

	t.Run("Header fail", func(t *testing.T) {
		c := check.HTTPRequest.Header(check.HTTPHeader.HasNotKey("Content-Type"))
		r := newReq()
		assertFailChecker(t, "HTTPRequest.Header", c, r, makeExpl(
			"header to pass Checker[http.Header]",
			fmt.Sprintf(
				"explanation: http.Header:\nexp not to have key \"Content-Type\"\ngot %v",
				r.Header,
			),
		))
	})

	t.Run("Body pass", func(t *testing.T) {
		c := check.HTTPRequest.Body(check.Bytes.Is(expBody))
		assertPassChecker(t, "HTTPRequest.Body", c, newReq())
	})

	t.Run("Body fail", func(t *testing.T) {
		c := check.HTTPRequest.Body(check.Bytes.Not(expBody))
		assertFailChecker(t, "HTTPRequest.Body", c, newReq(), makeExpl(
			"body to pass Checker[[]byte]",
			"explanation: bytes:\n"+makeExpl(
				"not "+fmt.Sprint(expBody),
				fmt.Sprint(expBody),
			),
		))
	})

	t.Run("Context pass", func(t *testing.T) {
		c := check.HTTPRequest.Context(check.Context.Value(expCtxKey, check.Value.Is(expCtxVal)))
		assertPassChecker(t, "HTTPRequest.Context", c, newReq())
	})

	t.Run("Context fail", func(t *testing.T) {
		c := check.HTTPRequest.Context(check.Context.Value(expCtxKey, check.Value.Not(expCtxVal)))
		assertFailChecker(t, "HTTPRequest.Context", c, newReq(), makeExpl(
			"context to pass Checker[context.Context]",
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
