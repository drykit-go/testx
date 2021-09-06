package check

import (
	"fmt"
	"net/http"
)

// httpHeaderCheckerProvider provides checks on type http.Header.
type httpHeaderCheckerProvider struct{ baseCheckerProvider }

// HasKey checks the gotten http.Header has a spcific key set.
// The corresponding value is ignored, meaning an empty value
// for that key passes the check.
func (p httpHeaderCheckerProvider) HasKey(key string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return p.hasKey(got, key) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to have key "+key, got)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// HasNotKey checks the gotten http.Header does not have
// a specific key set.
func (p httpHeaderCheckerProvider) HasNotKey(key string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return !p.hasKey(got, key) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to have key "+key, got)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// HasValue checks the gotten http.Heaser has any value equal to val.
// It only compares the first result for each key.
func (p httpHeaderCheckerProvider) HasValue(val string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return p.hasValue(got, val) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to have value "+val, got)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// HasNotValue checks the gotten http.Header does not have a value equal to val.
// It only compares the first result for each key.
func (p httpHeaderCheckerProvider) HasNotValue(val string) HTTPHeaderChecker {
	pass := func(got http.Header) bool { return !p.hasValue(got, val) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to have value "+val, got)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// CheckValue checks the gotten http.Header has a value for the matching key
// that passes the given StringChecker.
// It only checks the first result for the given key.
func (p httpHeaderCheckerProvider) CheckValue(key string, c StringChecker) HTTPHeaderChecker {
	var val string
	pass := func(got http.Header) bool {
		v, ok := p.get(got, key)
		if !ok {
			return false
		}
		val = v
		return c.Pass(v)
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			fmt.Sprintf("value for key %s to pass StringChecker", key),
			c.Explain("value", val),
		)
	}
	return NewHTTPHeaderChecker(pass, expl)
}

// Helpers

func (httpHeaderCheckerProvider) get(h http.Header, key string) (string, bool) {
	values, ok := h[key]
	if !ok || len(values) == 0 {
		return "", false
	}
	return values[0], true
}

func (httpHeaderCheckerProvider) hasKey(h http.Header, key string) bool {
	_, ok := h[key]
	return ok
}

func (httpHeaderCheckerProvider) hasValue(h http.Header, val string) bool {
	for _, values := range h {
		if len(values) == 0 {
			continue
		}
		if values[0] == val {
			return true
		}
	}
	return false
}
