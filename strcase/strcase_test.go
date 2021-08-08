package strcase_test

import (
	"testing"

	"github.com/drykit-go/testix"
	"github.com/drykit-go/testix/strcase"
)

func TestCamel(t *testing.T) {
	testix.Table(strcase.Camel, nil).Cases([]testix.Case{
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
