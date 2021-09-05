package check_test

import (
	"net/http"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestHTTPHeaderCheckerProvider(t *testing.T) {
	h := http.Header{
		"Content-Length": []string{"42"},
		"API_KEY":        []string{"secret0", "secret1"},
	}

	t.Run("KeySet pass", func(t *testing.T) {
		c := check.HTTPHeader.KeySet("API_KEY")
		assertPassHTTPHeaderChecker(t, "KeySet", c, h)
	})

	t.Run("KeySet fail", func(t *testing.T) {
		c := check.HTTPHeader.KeySet("password")
		assertFailHTTPHeaderChecker(t, "KeySet", c, h)
	})

	t.Run("KeyNotSet pass", func(t *testing.T) {
		c := check.HTTPHeader.KeyNotSet("password")
		assertPassHTTPHeaderChecker(t, "KeyNotSet", c, h)
	})

	t.Run("KeyNotSet fail", func(t *testing.T) {
		c := check.HTTPHeader.KeyNotSet("API_KEY")
		assertFailHTTPHeaderChecker(t, "KeyNotSet", c, h)
	})

	t.Run("ValueSet pass", func(t *testing.T) {
		c := check.HTTPHeader.ValueSet("42")
		assertPassHTTPHeaderChecker(t, "ValueSet", c, h)
		c = check.HTTPHeader.ValueSet("secret0")
		assertPassHTTPHeaderChecker(t, "ValueSet", c, h)
	})

	t.Run("ValueSet fail", func(t *testing.T) {
		c := check.HTTPHeader.ValueSet("secret42")
		assertFailHTTPHeaderChecker(t, "ValueSet", c, h)
		c = check.HTTPHeader.ValueSet("secret1")
		assertFailHTTPHeaderChecker(t, "ValueSet", c, h)
	})

	t.Run("ValueNotSet pass", func(t *testing.T) {
		c := check.HTTPHeader.ValueNotSet("secret42")
		assertPassHTTPHeaderChecker(t, "ValueNotSet", c, h)
	})

	t.Run("ValueNotSet fail", func(t *testing.T) {
		c := check.HTTPHeader.ValueNotSet("42")
		assertFailHTTPHeaderChecker(t, "ValueNotSet", c, h)
		c = check.HTTPHeader.ValueNotSet("secret0")
		assertFailHTTPHeaderChecker(t, "ValueNotSet", c, h)
	})

	t.Run("ValueOf pass", func(t *testing.T) {
		c := check.HTTPHeader.ValueOf("API_KEY", check.String.Is("secret0"))
		assertPassHTTPHeaderChecker(t, "ValueOf", c, h)
	})

	t.Run("ValueOf fail", func(t *testing.T) {
		c := check.HTTPHeader.ValueOf("API_KEY", check.String.Not("secret0"))
		assertFailHTTPHeaderChecker(t, "ValueOf", c, h)
	})
}

// Helpers

func assertPassHTTPHeaderChecker(t *testing.T, method string, c check.HTTPHeaderChecker, h http.Header) {
	t.Helper()
	if !c.Pass(h) {
		failHTTPHeaderCheckerTest(t, true, method, h, c.Explain)
	}
}

func assertFailHTTPHeaderChecker(t *testing.T, method string, c check.HTTPHeaderChecker, h http.Header) {
	t.Helper()
	if c.Pass(h) {
		failHTTPHeaderCheckerTest(t, false, method, h, c.Explain)
	}
}

func failHTTPHeaderCheckerTest(t *testing.T, expPass bool, method string, h http.Header, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "HTTPHeader."+method, explain("HTTPHeader value", h))
}
