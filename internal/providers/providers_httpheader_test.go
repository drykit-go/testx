package providers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/internal/providers"
)

func TestHTTPHeaderCheckerProvider(t *testing.T) {
	checkHTTPHeader := providers.HTTPHeaderCheckerProvider{}

	h := http.Header{
		"Content-Length": []string{"42"},
		"API_KEY":        []string{"secret0", "secret1"},
	}
	hstr := fmt.Sprint(h)

	t.Run("HasKey pass", func(t *testing.T) {
		c := checkHTTPHeader.HasKey("API_KEY")
		assertPassChecker(t, "HTTPHeader.HasKey", c, h)
	})

	t.Run("HasKey fail", func(t *testing.T) {
		c := checkHTTPHeader.HasKey("password")
		assertFailChecker(t, "HTTPHeader.HasKey", c, h,
			makeExpl(`to have key "password"`, hstr),
		)
	})

	t.Run("HasNotKey pass", func(t *testing.T) {
		c := checkHTTPHeader.HasNotKey("password")
		assertPassChecker(t, "HTTPHeader.HasNotKey", c, h)
	})

	t.Run("HasNotKey fail", func(t *testing.T) {
		c := checkHTTPHeader.HasNotKey("API_KEY")
		assertFailChecker(t, "HTTPHeader.HasNotKey", c, h,
			makeExpl(`not to have key "API_KEY"`, hstr),
		)
	})

	t.Run("HasValue pass", func(t *testing.T) {
		c := checkHTTPHeader.HasValue("42")
		assertPassChecker(t, "HTTPHeader.HasValue", c, h)
		c = checkHTTPHeader.HasValue("secret0")
		assertPassChecker(t, "HTTPHeader.HasValue", c, h)
	})

	t.Run("HasValue fail", func(t *testing.T) {
		c := checkHTTPHeader.HasValue("secret42")
		assertFailChecker(t, "HTTPHeader.HasValue", c, h,
			makeExpl(`to have value secret42`, hstr),
		)

		c = checkHTTPHeader.HasValue("secret1")
		assertFailChecker(t, "HTTPHeader.HasValue", c, h,
			makeExpl(`to have value secret1`, hstr),
		)
	})

	t.Run("HasNotValue pass", func(t *testing.T) {
		c := checkHTTPHeader.HasNotValue("secret42")
		assertPassChecker(t, "HTTPHeader.HasNotValue", c, h)
	})

	t.Run("HasNotValue fail", func(t *testing.T) {
		c := checkHTTPHeader.HasNotValue("42")
		assertFailChecker(t, "HTTPHeader.HasNotValue", c, h,
			makeExpl(`not to have value 42`, hstr),
		)

		c = checkHTTPHeader.HasNotValue("secret0")
		assertFailChecker(t, "HTTPHeader.HasNotValue", c, h,
			makeExpl(`not to have value secret0`, hstr),
		)
	})

	t.Run("CheckValue pass", func(t *testing.T) {
		c := checkHTTPHeader.CheckValue("API_KEY", check.String.Is("secret0"))
		assertPassChecker(t, "HTTPHeader.CheckValue", c, h)
	})

	t.Run("CheckValue fail", func(t *testing.T) {
		c := checkHTTPHeader.CheckValue("API_KEY", check.String.Not("secret0"))
		assertFailChecker(t, "HTTPHeader.CheckValue", c, h, makeExpl(
			`value for key "API_KEY" to pass Checker[string]`,
			`explanation: http.Header["API_KEY"]:`+"\n"+makeExpl(
				"not secret0",
				"secret0",
			),
		))
	})
}
