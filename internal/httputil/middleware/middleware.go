package middleware

import "net/http"

// MergeRight merges middlewares into a single one, wrapping each other
// starting from the last one.
//
// Example:
// 	MergeRight(m1, m2, m3) == m1(m2(m3))
func MergeRight(
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// AsFunc converts a http.Handler middleware to a http.HandlerFunc middleware.
func AsFunc(
	middleware func(http.Handler) http.Handler,
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(next).ServeHTTP
	}
}

// AsFuncs converts several http.Handler middlewares to a slice
// of http.HandlerFunc middlewares.
func AsFuncs(
	middlewares ...func(http.Handler) http.Handler,
) (middlewareFuncs []func(http.HandlerFunc) http.HandlerFunc) {
	for _, m := range middlewares {
		middlewareFuncs = append(middlewareFuncs, AsFunc(m))
	}
	return
}
