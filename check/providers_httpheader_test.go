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

	t.Run("HasKey pass", func(t *testing.T) {
		c := check.HTTPHeader.HasKey("API_KEY")
		assertPassHTTPHeaderChecker(t, "HasKey", c, h)
	})

	t.Run("HasKey fail", func(t *testing.T) {
		c := check.HTTPHeader.HasKey("password")
		assertFailHTTPHeaderChecker(t, "HasKey", c, h)
	})

	t.Run("HasNotKey pass", func(t *testing.T) {
		c := check.HTTPHeader.HasNotKey("password")
		assertPassHTTPHeaderChecker(t, "HasNotKey", c, h)
	})

	t.Run("HasNotKey fail", func(t *testing.T) {
		c := check.HTTPHeader.HasNotKey("API_KEY")
		assertFailHTTPHeaderChecker(t, "HasNotKey", c, h)
	})

	t.Run("HasValue pass", func(t *testing.T) {
		c := check.HTTPHeader.HasValue("42")
		assertPassHTTPHeaderChecker(t, "HasValue", c, h)
		c = check.HTTPHeader.HasValue("secret0")
		assertPassHTTPHeaderChecker(t, "HasValue", c, h)
	})

	t.Run("HasValue fail", func(t *testing.T) {
		c := check.HTTPHeader.HasValue("secret42")
		assertFailHTTPHeaderChecker(t, "HasValue", c, h)
		c = check.HTTPHeader.HasValue("secret1")
		assertFailHTTPHeaderChecker(t, "HasValue", c, h)
	})

	t.Run("HasNotValue pass", func(t *testing.T) {
		c := check.HTTPHeader.HasNotValue("secret42")
		assertPassHTTPHeaderChecker(t, "HasNotValue", c, h)
	})

	t.Run("HasNotValue fail", func(t *testing.T) {
		c := check.HTTPHeader.HasNotValue("42")
		assertFailHTTPHeaderChecker(t, "HasNotValue", c, h)
		c = check.HTTPHeader.HasNotValue("secret0")
		assertFailHTTPHeaderChecker(t, "HasNotValue", c, h)
	})

	t.Run("CheckValue pass", func(t *testing.T) {
		c := check.HTTPHeader.CheckValue("API_KEY", check.String.Is("secret0"))
		assertPassHTTPHeaderChecker(t, "CheckValue", c, h)
	})

	t.Run("CheckValue fail", func(t *testing.T) {
		c := check.HTTPHeader.CheckValue("API_KEY", check.String.Not("secret0"))
		assertFailHTTPHeaderChecker(t, "CheckValue", c, h)
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
