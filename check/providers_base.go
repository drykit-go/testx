package check

import "fmt"

type baseCheckerProvider struct{}

func (baseCheckerProvider) explain(label string, expStr, gotStr interface{}) string {
	return fmt.Sprintf("%s:\n    exp %v\n    got %v", label, expStr, gotStr)
}

func (p baseCheckerProvider) explainNot(label string, expStr, gotStr interface{}) string {
	return p.explain(label, fmt.Sprintf("not %v", expStr), gotStr)
}

func (p baseCheckerProvider) explainCheck(label, expStr, gotExpl string) string {
	return p.explain(label, expStr, "explanation: "+gotExpl)
}
