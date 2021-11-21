// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 21 Nov 21 13:15 UTC

package check

import (
	"github.com/drykit-go/testx/internal/providers"
)

var (
	// Int provides checks on type int.
	Int = providers.NumberCheckerProvider[int]{}
	// Int8 provides checks on type int8.
	Int8 = providers.NumberCheckerProvider[int8]{}
	// Int16 provides checks on type int16.
	Int16 = providers.NumberCheckerProvider[int16]{}
	// Int32 provides checks on type int32.
	Int32 = providers.NumberCheckerProvider[int32]{}
	// Int64 provides checks on type int64.
	Int64 = providers.NumberCheckerProvider[int64]{}
	// Uint provides checks on type uint.
	Uint = providers.NumberCheckerProvider[uint]{}
	// Uint8 provides checks on type uint8.
	Uint8 = providers.NumberCheckerProvider[uint8]{}
	// Uint16 provides checks on type uint16.
	Uint16 = providers.NumberCheckerProvider[uint16]{}
	// Uint32 provides checks on type uint32.
	Uint32 = providers.NumberCheckerProvider[uint32]{}
	// Uint64 provides checks on type uint64.
	Uint64 = providers.NumberCheckerProvider[uint64]{}
	// Float32 provides checks on type float32.
	Float32 = providers.NumberCheckerProvider[float32]{}
	// Float64 provides checks on type float64.
	Float64 = providers.NumberCheckerProvider[float64]{}
	// Bytes provides checks on type []byte.
	Bytes = providers.BytesCheckerProvider{}
	// Context provides checks on type context.Context.
	Context = providers.ContextCheckerProvider{}
	// Duration provides checks on type time.Duration.
	Duration = providers.DurationCheckerProvider{}
	// HTTPHeader provides checks on type http.Header.
	HTTPHeader = providers.HTTPHeaderCheckerProvider{}
	// HTTPRequest provides checks on type *http.Request.
	HTTPRequest = providers.HTTPRequestCheckerProvider{}
	// HTTPResponse provides checks on type *http.Response.
	HTTPResponse = providers.HTTPResponseCheckerProvider{}
	// String provides checks on type string.
	String = providers.StringCheckerProvider{}
	// Struct provides checks on kind struct.
	Struct = providers.StructCheckerProvider{}
)

// Map provides checks on type map[Key]Val.
func Map[Key comparable, Val any]() providers.MapCheckerProvider[Key, Val] {
	return providers.MapCheckerProvider[Key, Val]{}
}

// Slice provides checks on type []Elem.
func Slice[Elem any]() providers.SliceCheckerProvider[Elem] {
	return providers.SliceCheckerProvider[Elem]{}
}

// ValueCheckerProvider provides generic checks on any type.
func Value[T any]() providers.ValueCheckerProvider[T] {
	return providers.ValueCheckerProvider[T]{}
}
