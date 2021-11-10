// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 10 Nov 21 16:41 UTC

package check

var (
	// Int provides checks on type int.
	Int = numberCheckerProvider[int]{}
	// Int8 provides checks on type int8.
	Int8 = numberCheckerProvider[int8]{}
	// Int16 provides checks on type int16.
	Int16 = numberCheckerProvider[int16]{}
	// Int32 provides checks on type int32.
	Int32 = numberCheckerProvider[int32]{}
	// Int64 provides checks on type int64.
	Int64 = numberCheckerProvider[int64]{}
	// Uint provides checks on type uint.
	Uint = numberCheckerProvider[uint]{}
	// Uint8 provides checks on type uint8.
	Uint8 = numberCheckerProvider[uint8]{}
	// Uint16 provides checks on type uint16.
	Uint16 = numberCheckerProvider[uint16]{}
	// Uint32 provides checks on type uint32.
	Uint32 = numberCheckerProvider[uint32]{}
	// Uint64 provides checks on type uint64.
	Uint64 = numberCheckerProvider[uint64]{}
	// Float32 provides checks on type float32.
	Float32 = numberCheckerProvider[float32]{}
	// Float64 provides checks on type float64.
	Float64 = numberCheckerProvider[float64]{}
	// Bool provides checks on type bool.
	Bool = boolCheckerProvider{}
	// Bytes provides checks on type []byte.
	Bytes = bytesCheckerProvider{}
	// Context provides checks on type context.Context.
	Context = contextCheckerProvider{}
	// Duration provides checks on type time.Duration.
	Duration = durationCheckerProvider{}
	// HTTPHeader provides checks on type http.Header.
	HTTPHeader = httpHeaderCheckerProvider{}
	// HTTPRequest provides checks on type *http.Request.
	HTTPRequest = httpRequestCheckerProvider{}
	// HTTPResponse provides checks on type *http.Response.
	HTTPResponse = httpResponseCheckerProvider{}
	// Map provides checks on kind map.
	Map = mapCheckerProvider{}
	// Slice provides checks on kind slice.
	Slice = sliceCheckerProvider{}
	// String provides checks on type string.
	String = stringCheckerProvider{}
	// Struct provides checks on kind struct.
	Struct = structCheckerProvider{}
	// Value provides checks on type any.
	Value = valueCheckerProvider{}
)
