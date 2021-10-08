package check

import (
	"net/http"
)

// httpResponseCheckerProvider provides checks on type *http.Response.
type httpResponseCheckerProvider struct{ baseHTTPCheckerProvider }

// StatusCode checks the gotten *http.Response StatusCode passes
// the input IntChecker.
func (p httpResponseCheckerProvider) StatusCode(c IntChecker) HTTPResponseChecker {
	var code int
	pass := func(got *http.Response) bool {
		code = got.StatusCode
		return c.Pass(code)
	}
	expl := func(label string, _ interface{}) string {
		return p.explain(label, "status code to pass IntChecker", code)
	}
	return NewHTTPResponseChecker(pass, expl)
}

// Status checks the gotten *http.Response Status passes
// the input StringChecker.
func (p httpResponseCheckerProvider) Status(c StringChecker) HTTPResponseChecker {
	var status string
	pass := func(got *http.Response) bool {
		status = got.Status
		return c.Pass(status)
	}
	expl := func(label string, _ interface{}) string {
		return p.explain(label, "status to pass StringChecker", status)
	}
	return NewHTTPResponseChecker(pass, expl)
}

// ContentLength checks the gotten *http.Response ContentLength passes
// the input IntChecker.
func (p httpResponseCheckerProvider) ContentLength(c IntChecker) HTTPResponseChecker {
	var clen int
	pass := func(got *http.Response) bool {
		clen = int(got.ContentLength)
		return c.Pass(clen)
	}
	return NewHTTPResponseChecker(pass, p.explainContentLengthFunc(c, clen))
}

// Header checks the gotten *http.Response Header passes
// the input HTTPHeaderChecker.
func (p httpResponseCheckerProvider) Header(c HTTPHeaderChecker) HTTPResponseChecker {
	var header http.Header
	pass := func(got *http.Response) bool {
		header = got.Header
		return c.Pass(header)
	}
	return NewHTTPResponseChecker(pass, p.explainHeaderFunc(c, header))
}

// Body checks the gotten *http.Response Body passes the input BytesChecker.
// It should be used only once on a same *http.Response as it closes its body
// after reading it.
func (p httpResponseCheckerProvider) Body(c BytesChecker) HTTPResponseChecker {
	var body []byte
	pass := func(got *http.Response) bool {
		body = mustReadIO("check.HTTPResponse.Body", got.Body)
		return c.Pass(body)
	}
	return NewHTTPResponseChecker(pass, p.explainBodyFunc(c, body))
}
