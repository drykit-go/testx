package httpconv

import "net/http"

// NopHandler returns a http.HandlerFunc that does nothing.
func NopHandler() http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {}
}

// SafeHandler returns a no-op http.Hanldler if h == nil, else h.
func SafeHandler(h http.Handler) http.Handler {
	if h == nil {
		return NopHandler()
	}
	return h
}

// SafeHandlerFunc returns a no-op http.HandlerFunc if hf == nil, else hf.
func SafeHandlerFunc(hf http.HandlerFunc) http.HandlerFunc {
	if hf == nil {
		return NopHandler()
	}
	return hf
}
