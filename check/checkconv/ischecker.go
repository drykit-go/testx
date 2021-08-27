package checkconv

import "reflect"

type methodKey string

const (
	keyPass methodKey = "Pass"
	keyExpl methodKey = "Explain"
)

type signature struct {
	name    string
	in, out []reflect.Kind
}

var (
	signaturePass = signature{
		name: string(keyPass),
		in:   []reflect.Kind{reflect.Invalid}, // we use reflect.Invalid as a "any" wildcard
		out:  []reflect.Kind{reflect.Bool},
	}
	signatureExpl = signature{
		name: string(keyExpl),
		in:   []reflect.Kind{reflect.String, reflect.Interface},
		out:  []reflect.Kind{reflect.String},
	}

	checkerSignatures = map[methodKey]signature{
		keyPass: signaturePass,
		keyExpl: signatureExpl,
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
	return isCheckerMethod(v, keyPass)
}

func isExplainer(v reflect.Value) bool {
	return isCheckerMethod(v, keyExpl)
}

func isCheckerMethod(v reflect.Value, k methodKey) bool {
	s, ok := checkerSignatures[k]
	return ok && matchMethod(v, s)
}

func matchMethod(v reflect.Value, s signature) bool {
	m := v.MethodByName(s.name)
	if !m.IsValid() {
		return false
	}
	t := m.Type()
	return matchIn(t, s.in) && matchOut(t, s.out)
}

func matchIn(t reflect.Type, in []reflect.Kind) bool {
	return matchValuesKind(t.NumIn(), t.In, in)
}

func matchOut(t reflect.Type, out []reflect.Kind) bool {
	return matchValuesKind(t.NumOut(), t.Out, out)
}

func matchValuesKind(gotLen int, getIthVal func(int) reflect.Type, expKinds []reflect.Kind) bool {
	if gotLen != len(expKinds) {
		return false
	}
	for i := 0; i < gotLen; i++ {
		if gotk, expk := getIthVal(i).Kind(), expKinds[i]; expk != reflect.Invalid && gotk != expk {
			return false
		}
	}
	return true
}
