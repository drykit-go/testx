package check

import (
	"fmt"
	"net/http"

	"github.com/drykit-go/testx/internal/fmtexpl"
)

type baseCheckerProvider struct{}

func (baseCheckerProvider) explain(label string, exp, got interface{}) string {
	return fmtexpl.Default(label, exp, got)
}

func (p baseCheckerProvider) explainNot(label string, exp, got interface{}) string {
	return p.explain(label, fmt.Sprintf("not %v", exp), got)
}

func (p baseCheckerProvider) explainCheck(label, expStr, gotExpl string) string {
	return fmtexpl.Checker(label, expStr, gotExpl)
}

type baseHTTPCheckerProvider struct{ baseCheckerProvider }

func (p baseHTTPCheckerProvider) explainContentLengthFunc(c IntChecker, got int) ExplainFunc {
	return func(label string, _ interface{}) string {
		return p.explainCheck(label,
			"content length to pass IntChecker",
			c.Explain("content length", got),
		)
	}
}

func (p baseHTTPCheckerProvider) explainHeaderFunc(c HTTPHeaderChecker, got http.Header) ExplainFunc {
	return func(label string, _ interface{}) string {
		return p.explainCheck(label,
			"header to pass HTTPHeaderChecker",
			c.Explain("http.Header", got),
		)
	}
}

func (p baseHTTPCheckerProvider) explainBodyFunc(c BytesChecker, got []byte) ExplainFunc {
	return func(label string, _ interface{}) string {
		return p.explainCheck(label,
			"body to pass BytesChecker",
			c.Explain("bytes", got),
		)
	}
}
