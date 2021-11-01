// Package slices provides functions to perform operations on generic
// slices.
//
// TODO: temporary package, move to an external one.
package slices

func Map[T, U any](src []T, f func(T) U) []U {
	out := make([]U, len(src))
	for i, v := range src {
		out[i] = f(v)
	}
	return out
}
