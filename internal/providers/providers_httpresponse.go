package providers

import (
	"net/http"

	check "github.com/drykit-go/testx/internal/checktypes"
	"github.com/drykit-go/testx/internal/ioutil"
)

// HTTPResponseCheckerProvider provides checks on type *http.Response.
type HTTPResponseCheckerProvider struct{ baseHTTPCheckerProvider }

// StatusCode checks the gotten *http.Response StatusCode passes
// the input Checker[int].
func (p HTTPResponseCheckerProvider) StatusCode(c check.Checker[int]) check.Checker[*http.Response] {
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
	return check.NewChecker(pass, expl)
}

// Status checks the gotten *http.Response Status passes
// the input Checker[string].
func (p HTTPResponseCheckerProvider) Status(c check.Checker[string]) check.Checker[*http.Response] {
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
	return check.NewChecker(pass, expl)
}

// ContentLength checks the gotten *http.Response ContentLength passes
// the input Checker[int].
func (p HTTPResponseCheckerProvider) ContentLength(c check.Checker[int]) check.Checker[*http.Response] {
	var clen int
	pass := func(got *http.Response) bool {
		clen = int(got.ContentLength)
		return c.Pass(clen)
	}
	expl := func(label string, got any) string {
		return p.explainContentLengthFunc(c, func() int { return clen })(label, got)
	}
	return check.NewChecker(pass, expl)
}

// Header checks the gotten *http.Response Header passes
// the input Checker[http.Header].
func (p HTTPResponseCheckerProvider) Header(c check.Checker[http.Header]) check.Checker[*http.Response] {
	var header http.Header
	pass := func(got *http.Response) bool {
		header = got.Header
		return c.Pass(header)
	}
	expl := func(label string, got any) string {
		return p.explainHeaderFunc(c, func() http.Header { return header })(label, got)
	}
	return check.NewChecker(pass, expl)
}

// Body checks the gotten *http.Response Body passes the input Checker[[]byte].
// It should be used only once on a same *http.Response as it closes its body
// after reading it.
func (p HTTPResponseCheckerProvider) Body(c check.Checker[[]byte]) check.Checker[*http.Response] {
	var body []byte
	pass := func(got *http.Response) bool {
		body = ioutil.NopRead(&got.Body)
		return c.Pass(body)
	}
	expl := func(label string, got any) string {
		return p.explainBodyFunc(c, func() []byte { return body })(label, got)
	}
	return check.NewChecker(pass, expl)
}
