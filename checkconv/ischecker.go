package checkconv

import (
	"reflect"

	"github.com/drykit-go/testx/internal/reflectutil"
)

var (
	signaturePass = reflectutil.FuncSignature{
		Name: "Pass",
		In:   []reflect.Kind{reflectutil.AnyKind},
		Out:  []reflect.Kind{reflect.Bool},
	}

	signatureExpl = reflectutil.FuncSignature{
		Name: "Explain",
		In:   []reflect.Kind{reflect.String, reflect.Interface},
		Out:  []reflect.Kind{reflect.String},
	}
)

// IsChecker returns true if the provided value is a valid checker.
// A valid checker is any type exposing two methods:
// 	- Pass(got T) bool
// 	- Explain(label string, got interface{}) string
// Any custom implementation is considered valid whether or not it uses
// the package check.
//
// Note: The nature of Pass(got T) method, whose signature depend on the
// type of the tested value, prevents using a regular interface to identify
// a checker, hence the need of this helper.
func IsChecker(in interface{}) bool {
	v := reflect.ValueOf(in)
	return isPasser(v) && isExplainer(v)
}

func isPasser(v reflect.Value) bool {
	return signaturePass.ImplementedBy(v)
}

func isExplainer(v reflect.Value) bool {
	return signatureExpl.ImplementedBy(v)
}
