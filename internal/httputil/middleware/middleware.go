package middleware

import "net/http"

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

func AsFunc(
	middleware func(http.Handler) http.Handler,
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(next).ServeHTTP
	}
}

func AsFuncs(
	middlewares ...func(http.Handler) http.Handler,
) (middlewareFuncs []func(http.HandlerFunc) http.HandlerFunc) {
	for _, m := range middlewares {
		middlewareFuncs = append(middlewareFuncs, AsFunc(m))
	}
	return
}
