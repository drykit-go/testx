package httpconv_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/drykit-go/testx/internal/httpconv"
)

func TestMerge(t *testing.T) {
	results := []int{}

	appendResultsMiddleware := func(n int) func(next http.HandlerFunc) http.HandlerFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				results = append(results, n)
				next(w, r)
			}
		}
	}

	mergedMiddleware := httpconv.Merge(
		appendResultsMiddleware(1),
		appendResultsMiddleware(2),
		appendResultsMiddleware(3),
	)

	executeMiddleware(mergedMiddleware)

	expResults := []int{1, 2, 3}
	if !reflect.DeepEqual(results, expResults) {
		t.Errorf("exp %v\ngot %v", expResults, results)
	}
}

func TestMiddlewareFuncs(t *testing.T) {
	results := []int{}

	appendResultsMiddleware := func(n int) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				results = append(results, n)
				next.ServeHTTP(w, r)
			})
		}
	}

	middlewareFuncs := httpconv.MiddlewareFuncs(
		appendResultsMiddleware(1),
		appendResultsMiddleware(2),
		appendResultsMiddleware(3),
	)

	for _, m := range middlewareFuncs {
		executeMiddleware(m)
	}

	expResults := []int{1, 2, 3}
	if !reflect.DeepEqual(results, expResults) {
		t.Errorf("exp %v\ngot %v", expResults, results)
	}
}

func executeMiddleware(m func(http.HandlerFunc) http.HandlerFunc) {
	rr, rq := httptest.NewRecorder(), httptest.NewRequest("", "/", nil)
	m(func(w http.ResponseWriter, r *http.Request) {})(rr, rq)
}
