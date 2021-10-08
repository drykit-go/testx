package check

import (
	"net/http"
)

// httpRequestCheckerProvider provides checks on type *http.Request.
type httpRequestCheckerProvider struct{ baseHTTPCheckerProvider }

// ContentLength checks the gotten *http.Request ContentLength passes
// the input IntChecker.
func (p httpRequestCheckerProvider) ContentLength(c IntChecker) HTTPRequestChecker {
	var clen int
	pass := func(got *http.Request) bool {
		clen = int(got.ContentLength)
		return c.Pass(clen)
	}
	return NewHTTPRequestChecker(pass, p.explainContentLengthFunc(c, clen))
}

// Header checks the gotten *http.Request Header passes
// the input HTTPHeaderChecker.
func (p httpRequestCheckerProvider) Header(c HTTPHeaderChecker) HTTPRequestChecker {
	var header http.Header
	pass := func(got *http.Request) bool {
		header = got.Header
		return c.Pass(header)
	}
	return NewHTTPRequestChecker(pass, p.explainHeaderFunc(c, header))
}

// Body checks the gotten *http.Request Body passes the input BytesChecker.
// It should be used only once on a same *http.Request as it closes its body
// after reading it.
func (p httpRequestCheckerProvider) Body(c BytesChecker) HTTPRequestChecker {
	var body []byte
	pass := func(got *http.Request) bool {
		body = mustReadIO("check.HTTPRequest.Body", got.Body)
		return c.Pass(body)
	}
	return NewHTTPRequestChecker(pass, p.explainBodyFunc(c, body))
}
