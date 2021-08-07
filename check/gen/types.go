package gen

// types is a slice of source types used to generate
// their subsequent definitions.
var types = []struct {
	// Name, underlying Type
	N, T string
}{
	{N: "Bytes", T: "[]byte"},
	{N: "String", T: "string"},
	{N: "Int", T: "int"},
	{N: "Duration", T: "time.Duration"},
	{N: "HTTPHeader", T: "http.Header"},
	{N: "Untyped", T: "interface{}"},
}
