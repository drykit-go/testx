package check

import (
	"net/http"
)

// httpRequestCheckerProvider provides checks on type *http.Request.
type httpRequestCheckerProvider struct{ baseCheckerProvider }

// ContentLength checks the gotten *http.Request ContentLength passes
// the input IntChecker.
func (p httpRequestCheckerProvider) ContentLength(c IntChecker) HTTPRequestChecker {
	var clen int
	pass := func(got *http.Request) bool {
		clen = int(got.ContentLength)
		return c.Pass(clen)
	}
	expl := func(label string, _ interface{}) string {
		return p.explain(label, "content length to pass IntChecker", clen)
	}
	return NewHTTPRequestChecker(pass, expl)
}

// Header checks the gotten *http.Request Header passes
// the input HTTPHeaderChecker.
func (p httpRequestCheckerProvider) Header(c HTTPHeaderChecker) HTTPRequestChecker {
	var header http.Header
	pass := func(got *http.Request) bool {
		header = got.Header
		return c.Pass(header)
	}
	expl := func(label string, _ interface{}) string {
		return p.explain(label, "header to pass HTTPHeaderChecker", header)
	}
	return NewHTTPRequestChecker(pass, expl)
}
