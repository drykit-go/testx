package check

import (
	"context"
	"net/http"

	"github.com/drykit-go/testx/internal/ioutil"
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
		body = ioutil.NopRead(&got.Body)
		return c.Pass(body)
	}
	return NewHTTPRequestChecker(pass, p.explainBodyFunc(c, body))
}

// Context checks the gotten *http.Request Context passes
// the input ContextChecker.
func (p httpRequestCheckerProvider) Context(c ContextChecker) HTTPRequestChecker {
	var ctx context.Context
	pass := func(got *http.Request) bool {
		ctx = got.Context()
		return c.Pass(ctx)
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"context to pass ContextChecker",
			c.Explain("context", ctx),
		)
	}
	return NewHTTPRequestChecker(pass, expl)
}
