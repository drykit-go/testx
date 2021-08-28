package checkconv

import "reflect"

type signature struct {
	name    string
	in, out []reflect.Kind
}

func (s signature) match(v reflect.Value) bool {
	m := v.MethodByName(s.name)
	if !m.IsValid() {
		return false
	}
	t := m.Type()
	return s.matchIn(t) && s.matchOut(t)
}

func (s signature) matchIn(t reflect.Type) bool {
	return s.matchValues(t.NumIn(), t.In, s.in)
}

func (s signature) matchOut(t reflect.Type) bool {
	return s.matchValues(t.NumOut(), t.Out, s.out)
}

func (s signature) matchValues(numValues int, getIthVal func(int) reflect.Type, expKinds []reflect.Kind) bool {
	if numValues != len(expKinds) {
		return false
	}
	for i := 0; i < numValues; i++ {
		if !s.validKind(getIthVal(i).Kind(), expKinds[i]) {
			return false
		}
	}
	return true
}

func (s signature) validKind(gotk, expk reflect.Kind) bool {
	// reflect.Invalid is used as a generic wildcard
	return expk == reflect.Invalid || gotk == expk
}

var (
	signaturePass = signature{
		name: "Pass",
		in:   []reflect.Kind{reflect.Invalid}, // reflect.Invalid is used as a generic wildcard
		out:  []reflect.Kind{reflect.Bool},
	}

	signatureExpl = signature{
		name: "Explain",
		in:   []reflect.Kind{reflect.String, reflect.Interface},
		out:  []reflect.Kind{reflect.String},
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
	return signaturePass.match(v)
}

func isExplainer(v reflect.Value) bool {
	return signatureExpl.match(v)
}
