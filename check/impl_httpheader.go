package check

import (
	"fmt"
	"net/http"
)

type httpHeaderCheckFactory struct{}

func (httpHeaderCheckFactory) KeySet(key string) HTTPHeaderChecker {
	return httpHeaderCheck{
		passFunc: func(got http.Header) bool { return keySet(key, got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to have key \"%s\" set, got %+v",
				label, key, got,
			)
		},
	}
}

func (httpHeaderCheckFactory) KeyNotSet(key string) HTTPHeaderChecker {
	return httpHeaderCheck{
		passFunc: func(got http.Header) bool { return !keySet(key, got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to have key \"%s\" set, got %+v",
				label, key, got,
			)
		},
	}
}

func (httpHeaderCheckFactory) ValueSet(val string) HTTPHeaderChecker {
	return httpHeaderCheck{
		passFunc: func(got http.Header) bool { return valueSet(val, got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to have value \"%s\" set, got %+v",
				label, val, got,
			)
		},
	}
}

func (httpHeaderCheckFactory) ValueNotSet(val string) HTTPHeaderChecker {
	return httpHeaderCheck{
		passFunc: func(got http.Header) bool { return !valueSet(val, got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to have value \"%s\" set, got %+v",
				label, val, got,
			)
		},
	}
}

func (httpHeaderCheckFactory) ValueOf(key string, c StringChecker) HTTPHeaderChecker {
	var val string
	return httpHeaderCheck{
		passFunc: func(got http.Header) bool { val = got.Get(key); return c.Pass(val) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s value for key \"%s\" to pass StringChecker, got:\n\t%s",
				label, key, c.Explain("given string", val),
			)
		},
	}
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
