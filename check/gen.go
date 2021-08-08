package check

//go:generate go run ../cmd/gen/main.go -t check.gotmpl -o check.go

// For every type {N,T} defined in ./gen/types.go, running go generate
// will create the following definitions:
//
//	NPassFunc func(got T) bool
//
//	NPasser interface{ Pass(got T) bool }
//
//	NChecker interface {
//		NPasser
//		Explainer
//	}
//
//	type nCheck struct {
//		passFunc NPassFunc
//		explFunc ExplainFunc
//	}
//
//	func (c nCheck) Pass(got T) bool { return c.passFunc(got) }
//
//	func (c nCheck) Explain(label string, got interface{}) string { return c.explFunc(label, got) }
//
//	func NewNCheck(passFunc NPassFunc, explainFunc ExplainFunc) NChecker {
//		return nCheck{passFunc: passFunc, explFunc: explainFunc}
//	}
