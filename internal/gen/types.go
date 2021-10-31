package gen

// types is a slice of source types used to generate
// their subsequent definitions.
var types = []struct {
	// Name, underlying Type
	N, T string
}{
	{N: "Bool", T: "bool"},
	{N: "Bytes", T: "[]byte"},
	{N: "String", T: "string"},
	{N: "Int", T: "int"},
	{N: "Float64", T: "float64"},
	{N: "Duration", T: "time.Duration"},
	{N: "Context", T: "context.Context"},
	{N: "HTTPHeader", T: "http.Header"},
	{N: "HTTPRequest", T: "*http.Request"},
	{N: "HTTPResponse", T: "*http.Response"},
	{N: "Value", T: "any"},
}
