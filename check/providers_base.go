package check

import (
	"fmt"

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
