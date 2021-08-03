// Package check provides types to perform checks on values
// in a testing context.
package check

import (
	"regexp"
	"time"
)

// TODO: naming:
// check -> check.Int.InRange()
// shouldbe -> shouldbe.Int.InRange()
// should -> should.Int.InRange()
// expect -> expect.Int.InRange() <-- cool
// xpect -> xpect.Int.InRange() <-- cool
//
// TODO: new checks:
// string
// http.Header (map[string][]string)
// map[string]interface{}
// []interface{}

type (
	BytesPassFunc    func(got []byte) bool
	StringPassFunc   func(got string) bool
	IntPassFunc      func(got int) bool
	DurationPassFunc func(got time.Duration) bool
	UntypedPassFunc  func(got interface{}) bool

	ExplainFunc func(label string, got interface{}) string
)

type (
	BytesPasser    interface{ Pass(got []byte) bool }
	StringPasser   interface{ Pass(got string) bool }
	IntPasser      interface{ Pass(got int) bool }
	DurationPasser interface{ Pass(got time.Duration) bool }
	UntypedPasser  interface{ Pass(got interface{}) bool }

	Explainer interface {
		Explain(label string, got interface{}) string
	}
)

type (
	BytesChecker interface {
		BytesPasser
		Explainer
	}

	StringChecker interface {
		StringPasser
		Explainer
	}

	IntChecker interface {
		IntPasser
		Explainer
	}

	DurationChecker interface {
		DurationPasser
		Explainer
	}

	UntypedChecker interface {
		UntypedPasser
		Explainer
	}
)

type (
	BytesNativeChecks interface {
		Equal(tar []byte) BytesChecker
		EqualJSON(tar []byte) BytesChecker
		Len(c IntChecker) BytesChecker
	}

	StringNativeChecks interface {
		Equal(tar string) StringChecker
		Contains(tar string) StringChecker
		Match(rgx *regexp.Regexp) StringChecker
		Len(c IntChecker) StringChecker
	}

	IntNativeChecks interface {
		InRange(lo, hi int) IntChecker
		NotInRange(lo, hi int) IntChecker
		Equal(tar int) IntChecker
		NotEqual(tar int) IntChecker
		GreaterThan(tar int) IntChecker
		GreaterOrEqual(tar int) IntChecker
		LesserThan(tar int) IntChecker
		LesserOrEqual(tar int) IntChecker
	}

	DurationNativeChecks interface {
		Over(tar time.Duration) DurationChecker
		Under(tar time.Duration) DurationChecker
	}

	UntypedNativeChecks interface {
		Custom(desc string, f UntypedPassFunc) UntypedChecker
	}
)

var (
	Bytes    BytesNativeChecks    = bytesValue{}    // Bytes provides checks on type []byte
	String   StringNativeChecks   = stringValue{}   // String provides checks on type []byte
	Int      IntNativeChecks      = intValue{}      // Int provides checks on type int
	Duration DurationNativeChecks = durationValue{} // Duration provides checks on type time.Duration
	Untyped  UntypedNativeChecks  = untypedValue{}  // Untyped provides checks on untyped values
)

func NewBytesCheck(passFunc BytesPassFunc, explainFunc ExplainFunc) BytesChecker {
	return bytesCheck{passFunc: passFunc, explFunc: explainFunc}
}

func NewStringCheck(passFunc StringPassFunc, explainFunc ExplainFunc) StringChecker {
	return stringCheck{passFunc: passFunc, explFunc: explainFunc}
}

func NewIntCheck(passFunc IntPassFunc, explainFunc ExplainFunc) IntChecker {
	return intCheck{passFunc: passFunc, explFunc: explainFunc}
}

func NewDurationCheck(passFunc DurationPassFunc, explainFunc ExplainFunc) DurationChecker {
	return durationCheck{passFunc: passFunc, explFunc: explainFunc}
}

func NewUntypedCheck(passFunc UntypedPassFunc, explainFunc ExplainFunc) UntypedChecker {
	return untypedCheck{passFunc: passFunc, explFunc: explainFunc}
}
