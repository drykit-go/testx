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

func ExampleHTTPHandlerFunc_middleware() {
	// withLongProcess middleware processes something for 100 milliseconds.
	withLongProcess := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			next.ServeHTTP(w, r)
		})
	}

	// withContextValue middleware attaches the input key-val pair
	// to the request context.
	withContextValue := func(key, val interface{}) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), key, val)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}
	}

	// withContentType middleware sets the response header Content-Type
	// to contentType.
	withContentType := func(contentType string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", contentType)
				next.ServeHTTP(w, r)
			})
		}
	}

	results := testx.HTTPHandler(testx.NopHandler, // We can use NopHandler in this context.
		withLongProcess,
		withContextValue("userID", 42),
		withContentType("application/json"),
	).
		Duration(check.Duration.Over(100 * time.Millisecond)).
		Request(check.HTTPRequest.Context(check.Context.HasKeys("userID"))).
		Response(check.HTTPResponse.Header(check.HTTPHeader.HasValue("application/json"))).
		DryRun()

	fmt.Println(results.Passed())

	// Output:
	// true
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
