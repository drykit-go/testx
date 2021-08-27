package checkconv

//go:generate go run ../../cmd/gen/main.go -t assert.gotmpl -o assert.go

// For every type {N,T} defined in ../gen/types.go, running go generate
// will create the following definitions:
//
//	func FromN(c check.NChecker) check.UntypedChecker {
//		return check.NewUntypedCheck(
//			func(got interface{}) bool { return c.Pass(got.(T) },
//			c.Explain,
//		)
//	}
//
// It will also add a new case in the switch statement of func Assert:
//
//	case check.NChecker:
//		return FromN(c)
