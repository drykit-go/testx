package check

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// bytesCheckerProvider provides checks on type []byte.
type bytesCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten []byte is equal to the target.
func (p bytesCheckerProvider) Is(tar []byte) BytesChecker {
	pass := func(got []byte) bool { return p.eq(got, tar) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, tar, got)
	}
	return NewBytesChecker(pass, expl)
}

// Not checks the gotten []byte is not equal to the target.
func (p bytesCheckerProvider) Not(values ...[]byte) BytesChecker {
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
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, match, got)
	}
	return NewBytesChecker(pass, expl)
}

// SameJSON checks the gotten []byte and the target read as the same
// JSON value, ignoring formatting and keys order.
func (p bytesCheckerProvider) SameJSON(tar []byte) BytesChecker {
	var decGot, decTar interface{}
	pass := func(got []byte) bool {
		return sameJSON(got, tar, &decGot, &decTar)
	}
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("json data: %v", decTar),
			fmt.Sprintf("json data: %v", decGot),
		)
	}
	return NewBytesChecker(pass, expl)
}

// Len checks the gotten []byte's length passes the provided
// IntChecker.
func (p bytesCheckerProvider) Len(c IntChecker) BytesChecker {
	pass := func(got []byte) bool { return c.Pass(len(got)) }
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"length to pass IntChecker",
			c.Explain("length", len(got.([]byte))),
		)
	}
	return NewBytesChecker(pass, expl)
}

// Contains checks the gotten []byte contains a spcific subslice.
func (p bytesCheckerProvider) Contains(subslice []byte) BytesChecker {
	pass := func(got []byte) bool { return bytes.Contains(got, subslice) }
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			fmt.Sprintf("to contain subslice %v", subslice),
			got,
		)
	}
	return NewBytesChecker(pass, expl)
}

// NotContains checks the gotten []byte contains a spcific subslice.
func (p bytesCheckerProvider) NotContains(subslice []byte) BytesChecker {
	pass := func(got []byte) bool { return !bytes.Contains(got, subslice) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label,
			fmt.Sprintf("to contain subslice %v", subslice),
			got,
		)
	}
	return NewBytesChecker(pass, expl)
}

// AsMap checks the gotten []byte passes the given mapChecker
// once json-unmarshaled to a map[string]interface{}.
// It panics if it is not a valid JSON.
func (p bytesCheckerProvider) AsMap(mapChecker ValueChecker) BytesChecker {
	var m map[string]interface{}
	pass := func(got []byte) bool {
		err := json.NewDecoder(bytes.NewReader(got)).Decode(&m)
		if err != nil {
			panic(fmt.Sprintf("Bytes.AsMap: marshaling error: %s", err))
		}
		return mapChecker.Pass(m)
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"to pass MapChecker",
			mapChecker.Explain("unmarshaled json", m),
		)
	}
	return NewBytesChecker(pass, expl)
}

// AsString checks the gotten []byte passes the given StringChecker
// once converted to a string.
func (p bytesCheckerProvider) AsString(c StringChecker) BytesChecker {
	var s string
	pass := func(got []byte) bool {
		s = string(got)
		return c.Pass(s)
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"to pass StringChecker",
			c.Explain("converted bytes", s),
		)
	}
	return NewBytesChecker(pass, expl)
}

func (bytesCheckerProvider) eq(a, b []byte) bool {
	return bytes.Equal(a, b)
}
