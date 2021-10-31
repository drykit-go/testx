package check

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/drykit-go/testx/internal/fmtexpl"
)

type baseCheckerProvider struct{}

// sameJSON returns true if x and y evaluate to the same JSON value,
// ignoring the keys order.
// xptr and yptr must be pointers, as their values are filled
// with unmarshaled x and y respectively.
//
// It panics on the first error encountered in the process.
func (p baseCheckerProvider) sameJSON(x, y []byte, xptr, yptr any) bool {
	mustUnmarshal := func(b []byte, ptr any) {
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
func (p baseCheckerProvider) sameJSONProduced(xdata, ydata, xptr, yptr any) bool {
	mustMarshal := func(in any) []byte {
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

func (p baseCheckerProvider) formatList(values []string) string {
	var b strings.Builder
	b.WriteByte('[')
	b.WriteString(strings.Join(values, ", "))
	b.WriteByte(']')
	return b.String()
}

func (p baseCheckerProvider) deq(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func (baseCheckerProvider) explain(label string, exp, got any) string {
	return fmtexpl.Default(label, exp, got)
}

func (p baseCheckerProvider) explainNot(label string, exp, got any) string {
	return p.explain(label, fmt.Sprintf("not %v", exp), got)
}

func (p baseCheckerProvider) explainCheck(label, expStr, gotExpl string) string {
	return fmtexpl.Checker(label, expStr, gotExpl)
}

type baseHTTPCheckerProvider struct{ baseCheckerProvider }

func (p baseHTTPCheckerProvider) explainContentLengthFunc(
	c Checker[int],
	got func() int,
) ExplainFunc {
	return func(label string, _ any) string {
		return p.explainCheck(label,
			"content length to pass Checker[int]",
			c.Explain("content length", got()),
		)
	}
}

func (p baseHTTPCheckerProvider) explainHeaderFunc(
	c Checker[http.Header],
	got func() http.Header,
) ExplainFunc {
	return func(label string, _ any) string {
		return p.explainCheck(label,
			"header to pass Checker[http.Header]",
			c.Explain("http.Header", got()),
		)
	}
}

func (p baseHTTPCheckerProvider) explainBodyFunc(
	c Checker[[]byte],
	got func() []byte,
) ExplainFunc {
	return func(label string, _ any) string {
		return p.explainCheck(label,
			"body to pass Checker[[]byte]",
			c.Explain("bytes", got()),
		)
	}
}
