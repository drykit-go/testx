package check

import (
	"fmt"
	"net/http"
)

// httpHeaderCheckerProvider provides checks on type http.Header.
type httpHeaderCheckerProvider struct{}

// KeySet checks the gotten http.Header has a spcific key set.
// The corresponding value is ignored, meaning an empty value
// for that key passes the check.
func (f httpHeaderCheckerProvider) KeySet(key string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return f.keySet(key, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to have key \"%s\" set, got %+v",
			label, key, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// KeyNotSet checks the gotten http.Header does not have
// a specific key set.
func (f httpHeaderCheckerProvider) KeyNotSet(key string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return !f.keySet(key, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to have key \"%s\" set, got %+v",
			label, key, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// ValueSet checks the gotten http.Heaser has any key
// witha a matching value.
func (f httpHeaderCheckerProvider) ValueSet(val string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return f.valueSet(val, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to have value \"%s\" set, got %+v",
			label, val, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// ValueNotSet checks the gotten http.Header does not have
// any key with a matching value.
func (f httpHeaderCheckerProvider) ValueNotSet(val string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return !f.valueSet(val, got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to have value \"%s\" set, got %+v",
			label, val, got,
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// ValueOf checks the gotten http.Header has a value
// for the matching key that passes the given StringChecker.
func (httpHeaderCheckerProvider) ValueOf(key string, c StringChecker) HTTPHeaderChecker {
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

// Helpers

func (f httpHeaderCheckerProvider) keySet(key string, header http.Header) bool {
	_, ok := header[key]
	return ok
}

func (f httpHeaderCheckerProvider) valueSet(val string, header http.Header) bool {
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
