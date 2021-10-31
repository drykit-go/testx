package check

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// bytesCheckerProvider provides checks on type []byte.
type bytesCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten []byte is equal to the target.
func (p bytesCheckerProvider) Is(tar []byte) Checker[[]byte] {
	pass := func(got []byte) bool { return p.eq(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}

// Not checks the gotten []byte is not equal to the target.
func (p bytesCheckerProvider) Not(values ...[]byte) Checker[[]byte] {
	match := []byte{}
	pass := func(got []byte) bool {
		for _, v := range values {
			if p.eq(got, v) {
				match = v
				return false
			}
		}
		return true
	}
	expl := func(label string, got any) string {
		return p.explainNot(label, match, got)
	}
	return NewChecker(pass, expl)
}

// SameJSON checks the gotten []byte and the target read as the same
// JSON value, ignoring formatting and keys order.
func (p bytesCheckerProvider) SameJSON(tar []byte) Checker[[]byte] {
	var decGot, decTar any
	pass := func(got []byte) bool {
		return p.sameJSON(got, tar, &decGot, &decTar)
	}
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("json data: %v", decTar),
			fmt.Sprintf("json data: %v", decGot),
		)
	}
	return NewChecker(pass, expl)
}

// Len checks the gotten []byte's length passes the provided
// Checker[int].
func (p bytesCheckerProvider) Len(c Checker[int]) Checker[[]byte] {
	pass := func(got []byte) bool { return c.Pass(len(got)) }
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			"length to pass Checker[int]",
			c.Explain("length", len(got.([]byte))),
		)
	}
	return NewChecker(pass, expl)
}

// Contains checks the gotten []byte contains a specific subslice.
func (p bytesCheckerProvider) Contains(subslice []byte) Checker[[]byte] {
	pass := func(got []byte) bool { return bytes.Contains(got, subslice) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("to contain subslice %v", subslice),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// NotContains checks the gotten []byte contains a specific subslice.
func (p bytesCheckerProvider) NotContains(subslice []byte) Checker[[]byte] {
	pass := func(got []byte) bool { return !bytes.Contains(got, subslice) }
	expl := func(label string, got any) string {
		return p.explainNot(label,
			fmt.Sprintf("to contain subslice %v", subslice),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// AsMap checks the gotten []byte passes the given mapChecker
// once json-unmarshaled to a map[string]any.
// It fails if it is not a valid JSON.
func (p bytesCheckerProvider) AsMap(mapChecker Checker[any]) Checker[[]byte] {
	var m map[string]any
	var goterr error
	pass := func(got []byte) bool {
		goterr = json.NewDecoder(bytes.NewReader(got)).Decode(&m)
		return goterr == nil && mapChecker.Pass(m)
	}
	expl := func(label string, _ any) string {
		if goterr != nil {
			return p.explain(label,
				"to pass MapChecker",
				fmt.Sprintf("error: %s", goterr),
			)
		}
		return p.explainCheck(label,
			"to pass MapChecker",
			mapChecker.Explain("json map", m),
		)
	}
	return NewChecker(pass, expl)
}

// AsString checks the gotten []byte passes the given Checker[string]
// once converted to a string.
func (p bytesCheckerProvider) AsString(c Checker[string]) Checker[[]byte] {
	var s string
	pass := func(got []byte) bool {
		s = string(got)
		return c.Pass(s)
	}
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			"to pass Checker[string]",
			c.Explain("converted bytes", s),
		)
	}
	return NewChecker(pass, expl)
}

func (bytesCheckerProvider) eq(a, b []byte) bool {
	return bytes.Equal(a, b)
}
