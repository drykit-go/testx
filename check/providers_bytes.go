package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

// bytesCheckerProvider provides checks on type []byte.
type bytesCheckerProvider struct{}

// Is checks the gotten []byte is equal to the target.
func (p bytesCheckerProvider) Is(tar []byte) BytesChecker {
	pass := func(got []byte) bool { return p.eq(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to equal %v, got %v",
			label, tar, got,
		)
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
		return fmt.Sprintf(
			"expect %s not to equal %v, got %v",
			label, match, got,
		)
	}
	return NewBytesChecker(pass, expl)
}

// SameJSON checks the gotten []byte and the target returns
// the same JSON object.
func (bytesCheckerProvider) SameJSON(tar []byte) BytesChecker {
	var decGot, decTar interface{}
	pass := func(got []byte) bool {
		if err := json.Unmarshal(got, &decGot); err != nil {
			log.Panic(err)
		}
		if err := json.Unmarshal(tar, &decTar); err != nil {
			log.Panic(err)
		}
		return reflect.DeepEqual(decGot, decTar)
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expected json encoded value to equal %#v, got %#v",
			decTar, decGot,
		)
	}
	return NewBytesChecker(pass, expl)
}

// Len checks the gotten []byte's length passes the provided
// IntChecker.
func (bytesCheckerProvider) Len(c IntChecker) BytesChecker {
	pass := func(got []byte) bool { return c.Pass(len(got)) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"unexpected %s length: %s",
			label, c.Explain(label, len(got.([]byte))),
		)
	}
	return NewBytesChecker(pass, expl)
}

func (bytesCheckerProvider) eq(a, b []byte) bool {
	return bytes.Equal(a, b)
}
