package httpconv_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/drykit-go/testx/internal/httpconv"
)

func TestNopHandler(t *testing.T) {
	nopHandler := func(_ http.ResponseWriter, _ *http.Request) {}
	assertSameHandlers(t, http.HandlerFunc(nopHandler), httpconv.NopHandler())
}

func TestSafeHandler(t *testing.T) {
	t.Run("nil input returns nop handler", func(t *testing.T) {
		assertSameHandlers(t, httpconv.SafeHandler(nil), httpconv.NopHandler())
		assertSameHandlers(t, httpconv.SafeHandlerFunc(nil), httpconv.NopHandler())
	})

	t.Run("non-nil input returns identity", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})
		assertSameHandlers(t, httpconv.SafeHandler(h), h)
		assertSameHandlers(t, httpconv.SafeHandlerFunc(h), h)
	})
}

func assertSameHandlers(t *testing.T, h1, h2 http.Handler) {
	t.Helper()

	rr1, rq1 := defaultRecReq()
	rr2, rq2 := defaultRecReq()

	h1.ServeHTTP(rr1, rq1)
	h2.ServeHTTP(rr2, rq2)

	if !reflect.DeepEqual(rr1, rr2) {
		t.Errorf("different response recorders:\n%#v\n%#v", rr1, rr2)
	}
	if !reflect.DeepEqual(rq1, rq2) {
		t.Errorf("different requests:\n%#v\n%#v", rq1, rq2)
	}
}

func defaultRecReq() (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)
}
