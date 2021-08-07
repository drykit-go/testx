package strcase_test

import (
	"testing"

	"github.com/drykit-go/testix/strcase"
)

func TestCamel(t *testing.T) {
	cases := []struct {
		in, exp string
	}{
		{in: "A", exp: "a"},
		{in: "ABC", exp: "abc"},
		{in: "Name", exp: "name"},
		{in: "AAAAAa", exp: "aaaaAa"},
		{in: "SomeName", exp: "someName"},
		{in: "HTTPHeader", exp: "httpHeader"},
		{in: "MyHTTPHeader", exp: "myHTTPHeader"},
		{in: "my_HTTP_Header", exp: "myHTTPHeader"},
	}

	for _, c := range cases {
		if got := strcase.Camel(c.in); got != c.exp {
			t.Errorf(
				`strcase.Camel("%s") -> want "%s", got "%s"`,
				c.in, c.exp, got,
			)
		}
	}
}
