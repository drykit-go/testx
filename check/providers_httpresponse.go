package check

import (
	"net/http"

	"github.com/drykit-go/testx/internal/ioutil"
)

// httpResponseCheckerProvider provides checks on type *http.Response.
type httpResponseCheckerProvider struct{ baseHTTPCheckerProvider }

// StatusCode checks the gotten *http.Response StatusCode passes
// the input Checker[int].
func (p httpResponseCheckerProvider) StatusCode(c Checker[int]) Checker[*http.Response] {
	var code int
	pass := func(got *http.Response) bool {
		code = got.StatusCode
		return c.Pass(code)
	}
	expl := func(label string, _ any) string {
		return p.explainCheck(label,
			"status code to pass Checker[int]",
			c.Explain("status code", code),
		)
	}
	return NewChecker(pass, expl)
}

// Status checks the gotten *http.Response Status passes
// the input Checker[string].
func (p httpResponseCheckerProvider) Status(c Checker[string]) Checker[*http.Response] {
	var status string
	pass := func(got *http.Response) bool {
		status = got.Status
		return c.Pass(status)
	}
	expl := func(label string, _ any) string {
		return p.explainCheck(label,
			"status to pass Checker[string]",
			c.Explain("status", status),
		)
	}
	return NewChecker(pass, expl)
}

// ContentLength checks the gotten *http.Response ContentLength passes
// the input Checker[int].
func (p httpResponseCheckerProvider) ContentLength(c Checker[int]) Checker[*http.Response] {
	var clen int
	pass := func(got *http.Response) bool {
		clen = int(got.ContentLength)
		return c.Pass(clen)
	}
	expl := func(label string, got any) string {
		return p.explainContentLengthFunc(c, func() int { return clen })(label, got)
	}
	return NewChecker(pass, expl)
}

// Header checks the gotten *http.Response Header passes
// the input Checker[http.Header].
func (p httpResponseCheckerProvider) Header(c Checker[http.Header]) Checker[*http.Response] {
	var header http.Header
	pass := func(got *http.Response) bool {
		header = got.Header
		return c.Pass(header)
	}
	expl := func(label string, got any) string {
		return p.explainHeaderFunc(c, func() http.Header { return header })(label, got)
	}
	return NewChecker(pass, expl)
}

// Body checks the gotten *http.Response Body passes the input Checker[[]byte].
// It should be used only once on a same *http.Response as it closes its body
// after reading it.
func (p httpResponseCheckerProvider) Body(c Checker[[]byte]) Checker[*http.Response] {
	var body []byte
	pass := func(got *http.Response) bool {
		body = ioutil.NopRead(&got.Body)
		return c.Pass(body)
	}
	expl := func(label string, got any) string {
		return p.explainBodyFunc(c, func() []byte { return body })(label, got)
	}
	return NewChecker(pass, expl)
}
