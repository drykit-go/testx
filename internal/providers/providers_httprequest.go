package providers

import (
	"context"
	"net/http"

	check "github.com/drykit-go/testx/internal/checktypes"
	"github.com/drykit-go/testx/internal/ioutil"
)

// HTTPRequestCheckerProvider provides checks on type *http.Request.
type HTTPRequestCheckerProvider struct{ baseHTTPCheckerProvider }

// ContentLength checks the gotten *http.Request ContentLength passes
// the input Checker[int].
func (p HTTPRequestCheckerProvider) ContentLength(c check.Checker[int]) check.Checker[*http.Request] {
	var clen int
	pass := func(got *http.Request) bool {
		clen = int(got.ContentLength)
		return c.Pass(clen)
	}
	expl := func(label string, got any) string {
		return p.explainContentLengthFunc(c, func() int { return clen })(label, got)
	}
	return check.NewChecker(pass, expl)
}

// Header checks the gotten *http.Request Header passes
// the input Checker[http.Header].
func (p HTTPRequestCheckerProvider) Header(c check.Checker[http.Header]) check.Checker[*http.Request] {
	var header http.Header
	pass := func(got *http.Request) bool {
		header = got.Header
		return c.Pass(header)
	}
	expl := func(label string, got any) string {
		return p.explainHeaderFunc(c, func() http.Header { return header })(label, got)
	}
	return check.NewChecker(pass, expl)
}

// Body checks the gotten *http.Request Body passes the input Checker[[]byte].
// It should be used only once on a same *http.Request as it closes its body
// after reading it.
func (p HTTPRequestCheckerProvider) Body(c check.Checker[[]byte]) check.Checker[*http.Request] {
	var body []byte
	pass := func(got *http.Request) bool {
		body = ioutil.NopRead(&got.Body)
		return c.Pass(body)
	}
	expl := func(label string, got any) string {
		return p.explainBodyFunc(c, func() []byte { return body })(label, got)
	}
	return check.NewChecker(pass, expl)
}

// Context checks the gotten *http.Request Context passes
// the input Checker[context.Context].
func (p HTTPRequestCheckerProvider) Context(c check.Checker[context.Context]) check.Checker[*http.Request] {
	var ctx context.Context
	pass := func(got *http.Request) bool {
		ctx = got.Context()
		return c.Pass(ctx)
	}
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			"context to pass Checker[context.Context]",
			c.Explain("context", ctx),
		)
	}
	return check.NewChecker(pass, expl)
}
