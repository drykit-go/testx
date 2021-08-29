package check

import (
	"fmt"
	"net/http"
)

type httpHeaderCheckerFactory struct{}

func (httpHeaderCheckerFactory) KeySet(key string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return keySet(key, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to have key \"%s\" set, got %+v",
			label, key, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

func (httpHeaderCheckerFactory) KeyNotSet(key string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return !keySet(key, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to have key \"%s\" set, got %+v",
			label, key, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

func (httpHeaderCheckerFactory) ValueSet(val string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return valueSet(val, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to have value \"%s\" set, got %+v",
			label, val, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

func (httpHeaderCheckerFactory) ValueNotSet(val string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return !valueSet(val, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to have value \"%s\" set, got %+v",
			label, val, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

func (httpHeaderCheckerFactory) ValueOf(key string, c StringChecker) HTTPHeaderChecker {
	var val string
	pass := func(got http.Header) bool { val = got.Get(key); return c.Pass(val) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s value for key \"%s\" to pass StringChecker, got:\n\t%s",
			label, key, c.Explain("given string", val),
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

func keySet(key string, header http.Header) bool {
	_, ok := header[key]
	return ok
}

func valueSet(val string, header http.Header) bool {
	for _, v := range header {
		if len(v) == 0 {
			continue
		}
		for _, vv := range v {
			if vv == val {
				return true
			}
		}
	}
	return false
}
