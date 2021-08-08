package strcase_test

import (
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/strcase"
)

func TestCamel(t *testing.T) {
	testx.Table(strcase.Camel, nil).Cases([]testx.Case{
		{In: "A", Exp: "a"},
		{In: "ABC", Exp: "abc"},
		{In: "Name", Exp: "name"},
		{In: "AAAAAa", Exp: "aaaaAa"},
		{In: "SomeName", Exp: "someName"},
		{In: "HTTPHeader", Exp: "httpHeader"},
		{In: "MyHTTPHeader", Exp: "myHTTPHeader"},
		{In: "my_HTTP_Header", Exp: "myHTTPHeader"},
	}).Run(t)
}
