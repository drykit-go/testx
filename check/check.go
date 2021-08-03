// Package check provides types to perform checks on values
// in a testing context.
package check

import (
	"net/http"
	"regexp"
	"time"
)

// TODO: naming:
// check -> check.Int.InRange()
// should -> should.Int.InRange()
// expect -> expect.Int.InRange() <-- cool
// xpect -> xpect.Int.InRange() <-- cool
// want -> want.Int.InRange() <-- cool but conflict w/ vars named 'want'
//
// TODO: new checks:
// http.Header (map[string][]string)
// map[string]interface{}
// []interface{}

type (
	BytesPassFunc      func(got []byte) bool
	StringPassFunc     func(got string) bool
	IntPassFunc        func(got int) bool
	DurationPassFunc   func(got time.Duration) bool
	HTTPHeaderPassFunc func(got http.Header) bool
	UntypedPassFunc    func(got interface{}) bool

	ExplainFunc func(label string, got interface{}) string
)

type (
	BytesPasser      interface{ Pass(got []byte) bool }
	StringPasser     interface{ Pass(got string) bool }
	IntPasser        interface{ Pass(got int) bool }
	DurationPasser   interface{ Pass(got time.Duration) bool }
	HTTPHeaderPasser interface{ Pass(got http.Header) bool }
	UntypedPasser    interface{ Pass(got interface{}) bool }

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

	HTTPHeaderChecker interface {
		HTTPHeaderPasser
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
		NotContains(tar string) StringChecker
		Match(rgx *regexp.Regexp) StringChecker
		NotMatch(rgx *regexp.Regexp) StringChecker
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

	HTTPHeaderNativeChecks interface {
		KeySet(key string) HTTPHeaderChecker
		KeyNotSet(key string) HTTPHeaderChecker
		ValueSet(val string) HTTPHeaderChecker
		ValueNotSet(val string) HTTPHeaderChecker
		ValueOf(key string, c StringChecker) HTTPHeaderChecker
	}

	UntypedNativeChecks interface {
		Custom(desc string, f UntypedPassFunc) UntypedChecker
	}
)

var (
	Bytes      BytesNativeChecks      = bytesValue{}      // Bytes provides checks on type []byte
	String     StringNativeChecks     = stringValue{}     // String provides checks on type []byte
	Int        IntNativeChecks        = intValue{}        // Int provides checks on type int
	Duration   DurationNativeChecks   = durationValue{}   // Duration provides checks on type time.Duration
	HTTPHeader HTTPHeaderNativeChecks = httpHeaderValue{} // HTTPHeader provides checks on type http.Header
	Untyped    UntypedNativeChecks    = untypedValue{}    // Untyped provides checks on untyped values
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

func NewHTTPHeaderCheck(passFunc HTTPHeaderPassFunc, explainFunc ExplainFunc) HTTPHeaderChecker {
	return httpHeaderCheck{passFunc: passFunc, explFunc: explainFunc}
}

func NewUntypedCheck(passFunc UntypedPassFunc, explainFunc ExplainFunc) UntypedChecker {
	return untypedCheck{passFunc: passFunc, explFunc: explainFunc}
}
