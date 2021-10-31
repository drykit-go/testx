package check

//go:generate ../bin/gen -kind types -name check
//go:generate ../bin/gen -kind types -name checkers
// //go:generate ../bin/gen -kind interfaces -name providers

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
//	type nChecker struct {
//		passFunc NPassFunc
//		explFunc ExplainFunc
//	}
//
//	func (c nChecker) Pass(got T) bool { return c.passFunc(got) }
//
//	func (c nChecker) Explain(label string, got interface{}) string { return c.explFunc(label, got) }
//
//	func NewNChecker(passFunc NPassFunc, explainFunc ExplainFunc) NChecker {
//		return nChecker{passFunc: passFunc, explFunc: explainFunc}
//	}
