package check_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestHTTPHeaderCheckerProvider(t *testing.T) {
	h := http.Header{
		"Content-Length": []string{"42"},
		"API_KEY":        []string{"secret0", "secret1"},
	}
	hstr := fmt.Sprint(h)

	t.Run("HasKey pass", func(t *testing.T) {
		c := check.HTTPHeader.HasKey("API_KEY")
		assertPassHTTPHeaderChecker(t, "HasKey", c, h)
	})

	t.Run("HasKey fail", func(t *testing.T) {
		c := check.HTTPHeader.HasKey("password")
		assertFailHTTPHeaderChecker(t, "HasKey", c, h,
			makeExpl(`to have key "password"`, hstr),
		)
	})

	t.Run("HasNotKey pass", func(t *testing.T) {
		c := check.HTTPHeader.HasNotKey("password")
		assertPassHTTPHeaderChecker(t, "HasNotKey", c, h)
	})

	t.Run("HasNotKey fail", func(t *testing.T) {
		c := check.HTTPHeader.HasNotKey("API_KEY")
		assertFailHTTPHeaderChecker(t, "HasNotKey", c, h,
			makeExpl(`not to have key "API_KEY"`, hstr),
		)
	})

	t.Run("HasValue pass", func(t *testing.T) {
		c := check.HTTPHeader.HasValue("42")
		assertPassHTTPHeaderChecker(t, "HasValue", c, h)
		c = check.HTTPHeader.HasValue("secret0")
		assertPassHTTPHeaderChecker(t, "HasValue", c, h)
	})

	t.Run("HasValue fail", func(t *testing.T) {
		c := check.HTTPHeader.HasValue("secret42")
		assertFailHTTPHeaderChecker(t, "HasValue", c, h,
			makeExpl(`to have value secret42`, hstr),
		)

		c = check.HTTPHeader.HasValue("secret1")
		assertFailHTTPHeaderChecker(t, "HasValue", c, h,
			makeExpl(`to have value secret1`, hstr),
		)
	})

	t.Run("HasNotValue pass", func(t *testing.T) {
		c := check.HTTPHeader.HasNotValue("secret42")
		assertPassHTTPHeaderChecker(t, "HasNotValue", c, h)
	})

	t.Run("HasNotValue fail", func(t *testing.T) {
		c := check.HTTPHeader.HasNotValue("42")
		assertFailHTTPHeaderChecker(t, "HasNotValue", c, h,
			makeExpl(`not to have value 42`, hstr),
		)

		c = check.HTTPHeader.HasNotValue("secret0")
		assertFailHTTPHeaderChecker(t, "HasNotValue", c, h,
			makeExpl(`not to have value secret0`, hstr),
		)
	})

	t.Run("CheckValue pass", func(t *testing.T) {
		c := check.HTTPHeader.CheckValue("API_KEY", check.String.Is("secret0"))
		assertPassHTTPHeaderChecker(t, "CheckValue", c, h)
	})

	t.Run("CheckValue fail", func(t *testing.T) {
		c := check.HTTPHeader.CheckValue("API_KEY", check.String.Not("secret0"))
		assertFailHTTPHeaderChecker(t, "CheckValue", c, h, makeExpl(
			`value for key "API_KEY" to pass StringChecker`,
			`explanation: http.Header["API_KEY"]:`+"\n"+makeExpl(
				"not secret0",
				"secret0",
			),
		))
	})
}

// Helpers

func assertPassHTTPHeaderChecker(t *testing.T, method string, c check.HTTPHeaderChecker, h http.Header) {
	t.Helper()
	if !c.Pass(h) {
		failHTTPHeaderCheckerTest(t, true, method, h, c.Explain)
	}
}

func assertFailHTTPHeaderChecker(t *testing.T, method string, c check.HTTPHeaderChecker, h http.Header, expexpl string) {
	t.Helper()
	if c.Pass(h) {
		failHTTPHeaderCheckerTest(t, false, method, h, c.Explain)
	}
	assertGoodExplain(t, c, h, expexpl)
}

func failHTTPHeaderCheckerTest(t *testing.T, expPass bool, method string, h http.Header, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "HTTPHeader."+method, explain("HTTPHeader value", h))
}
