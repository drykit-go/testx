package check

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/drykit-go/testx/internal/fmtexpl"
)

type baseCheckerProvider struct{}

// sameJSON returns true if x and y evaluate to the same JSON value,
// ignoring the keys order.
// xptr and yptr must be pointers, as their values are filled
// with unmarshaled x and y respectively.
//
// It panics on the first error encountered in the process.
func (p baseCheckerProvider) sameJSON(x, y []byte, xptr, yptr interface{}) bool {
	mustUnmarshal := func(b []byte, ptr interface{}) {
		if err := json.Unmarshal(b, &ptr); err != nil {
			panic(err)
		}
	}
	mustUnmarshal(x, xptr)
	mustUnmarshal(y, yptr)
	return p.deq(xptr, yptr)
}

// sameJSONProduced returns true if xdata and ydata result in the same JSON value,
// ignoring the keys order.
// xptr and yptr must be pointers, as their values are filled
// with marshaled+unmarshaled json from xdata and ydata respectively.
//
// It panics on the first error encountered in the process.
func (p baseCheckerProvider) sameJSONProduced(xdata, ydata, xptr, yptr interface{}) bool {
	mustMarshal := func(in interface{}) []byte {
		b, err := json.Marshal(in)
		if err != nil {
			panic(err)
		}
		return b
	}
	bx := mustMarshal(xdata)
	by := mustMarshal(ydata)
	return p.sameJSON(bx, by, &xptr, &yptr)
}

func (p baseCheckerProvider) deq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

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
