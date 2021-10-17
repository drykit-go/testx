package httpconv

import "net/http"

// Merge merges middlewares into a single one, wrapping each other
// starting from the last one.
//
// Example:
// 	Merge(m1, m2, m3) == m1(m2(m3))
func Merge(
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// MiddlewareFunc converts a http.Handler middleware to a http.HandlerFunc
// middleware.
func MiddlewareFunc(
	middleware func(http.Handler) http.Handler,
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(next).ServeHTTP
	}
}

// MiddlewareFuncs converts several http.Handler middlewares to a slice
// of http.HandlerFunc middlewares.
func MiddlewareFuncs(
	middlewares ...func(http.Handler) http.Handler,
) (middlewareFuncs []func(http.HandlerFunc) http.HandlerFunc) {
	for _, m := range middlewares {
		middlewareFuncs = append(middlewareFuncs, MiddlewareFunc(m))
	}
	return
}
