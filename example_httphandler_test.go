package testx_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

func MyHTTPHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "42" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	time.Sleep(500 * time.Millisecond) // some processing
	w.Write([]byte("ok"))
}

func ExampleHTTPHandlerFunc() {
	t := &testing.T{} // ignore: emulating a testing context

	t.Run("good request", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/endpoint?id=42", nil)
		testx.HTTPHandlerFunc(MyHTTPHandler).WithRequest(r).
			Response(
				check.HTTPResponse.StatusCode(check.Int.InRange(200, 299)),
				check.HTTPResponse.Body(check.Bytes.Is([]byte("ok"))),
			).
			Run(t)
	})

	t.Run("bad request", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/endpoint?id=404", nil)
		testx.HTTPHandlerFunc(MyHTTPHandler).WithRequest(r).
			Response(check.HTTPResponse.Status(check.String.Contains("Not Found"))).
			Duration(check.Duration.Under(10 * time.Millisecond)).
			Run(t)
	})
}

func TestExampleHTTPHandlerFunc_middleware(t *testing.T) {
	var expCtxKey interface{} = "userID"

	// authenticate is a middleware that sets a userID to the request context.
	authenticate := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), expCtxKey, 42)
			next(w, r.WithContext(ctx))
		}
	}

	t.Run("middleware authenticate", func(t *testing.T) {
		testx.HTTPHandlerFunc(MyHTTPHandler, authenticate).
			Request(
				check.HTTPRequest.Context(check.Context.HasKeys(expCtxKey)),
			).
			Run(t)
	})
}

func ExampleHTTPHandlerFunc_dryRun() {
	handlerRunner := testx.HTTPHandlerFunc(MyHTTPHandler)

	goodRequest, _ := http.NewRequest("GET", "/endpoint?id=42", nil)
	goodRequestResults := handlerRunner.WithRequest(goodRequest).
		Response(
			check.HTTPResponse.StatusCode(check.Int.InRange(200, 299)),
			check.HTTPResponse.Body(check.Bytes.Is([]byte("ok"))),
		).
		DryRun()

	badRequest, _ := http.NewRequest("GET", "/endpoint?id=404", nil)
	badRequestResults := handlerRunner.WithRequest(badRequest).
		Response(check.HTTPResponse.Status(check.String.Contains("Not Found"))).
		Duration(check.Duration.Under(10 * time.Millisecond)).
		DryRun()

	fmt.Println(goodRequestResults.Passed())
	fmt.Println(goodRequestResults.ResponseStatus())
	fmt.Println(badRequestResults.Passed())
	fmt.Println(badRequestResults.ResponseStatus())

	// Output:
	// true
	// 200 OK
	// true
	// 404 Not Found
}
