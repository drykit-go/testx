package checkconv

//go:generate  ../bin/gen -kind types -name assert

// For every type {N,T} defined in ../gen/types.go, running go generate
// will create the following definitions:
//
//	func FromN(c check.NChecker) check.ValueChecker {
//		return check.NewValueCheck(
//			func(got any) bool { return c.Pass(got.(T) },
//			c.Explain,
//		)
//	}
//
// It will also add a new case in the switch statement of func Assert:
//
//	case check.NChecker:
//		return FromN(c)
