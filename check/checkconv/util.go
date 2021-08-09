package checkconv

import "reflect"

const (
	methodPass = "Pass"
	methodExpl = "Explain"
)

func IsChecker(in interface{}) bool {
	v := reflect.ValueOf(in)
	return isPasser(v) && isExplainer(v)
}

func isPasser(in reflect.Value) bool {
	m, ok := methodByName(in, methodPass)
	if !ok {
		return false
	}
	t := m.Type()
	return t.NumIn() == 1 && t.NumOut() == 1 && t.Out(0).Kind() == reflect.Bool
}

func isExplainer(in reflect.Value) bool {
	m, ok := methodByName(in, methodExpl)
	if !ok {
		return false
	}
	t := m.Type()
	return t.NumIn() == 2 && t.NumOut() == 1 && t.Out(0).Kind() == reflect.String
}

func methodByName(rval reflect.Value, name string) (reflect.Value, bool) {
	m := rval.MethodByName(name)
	return m, m.IsValid()
}
