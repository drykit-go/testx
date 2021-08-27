package check

import (
	"regexp"
	"time"
)

type (
	// BytesCheckerProvider provides checks on type []byte.
	BytesCheckerProvider interface {
		// Is checks the gotten []byte is equal to the target.
		Is(tar []byte) BytesChecker
		// SameJSON checks the gotten []byte and the target returns
		// the same JSON object.
		SameJSON(tar []byte) BytesChecker
		// Len checks the gotten []byte's length passes the provided
		// IntChecker.
		Len(c IntChecker) BytesChecker
	}

	// StringCheckerProvider provides checks on type string.
	StringCheckerProvider interface {
		// Is checks the gotten string is equal to the target.
		Is(tar string) StringChecker
		// Contains checks the gotten string contains the target substring.
		Contains(tar string) StringChecker
		// NotContains checks the gotten string do not contain the target
		// substring.
		NotContains(tar string) StringChecker
		// Match checks the gotten string matches the given regexp.
		Match(rgx *regexp.Regexp) StringChecker
		// NotMatch checks the gotten string do not match the given regexp.
		NotMatch(rgx *regexp.Regexp) StringChecker
		// Len checks the gotten string's length passes the given IntChecker.
		Len(c IntChecker) StringChecker
	}

	// IntCheckerProvider provides checks on type int.
	IntCheckerProvider interface {
		// InRange checks the gotten int is in the closed interval [lo:hi].
		InRange(lo, hi int) IntChecker
		// OutRange checks the gotten int is not in the closed interval [lo:hi].
		OutRange(lo, hi int) IntChecker
		// Is checks the gotten int is equal to the target.
		Is(tar int) IntChecker
		// Not checks the gotten int is not equal to the target.
		Not(tar int) IntChecker
		// GT checks the gotten int is greater than the target.
		GT(tar int) IntChecker
		// GTE checks the gotten int is greater or equal to the target.
		GTE(tar int) IntChecker
		// LT checks the gotten int is lesser than the target.
		LT(tar int) IntChecker
		// LTE checks the gotten int is lesser or equal to the target.
		LTE(tar int) IntChecker
	}

	// DurationCheckerProvider provides checks on type time.Duration.
	DurationCheckerProvider interface {
		// Over checks the gotten time.Duration is over the target duration.
		Over(tar time.Duration) DurationChecker
		// Under checks the gotten time.Duration is under the target duration.
		Under(tar time.Duration) DurationChecker
	}

	// HTTPHeaderCheckerProvider provides checks on type http.Header.
	HTTPHeaderCheckerProvider interface {
		// KeySet checks the gotten http.Header has a spcific key set.
		// The corresponding value is ignored, meaning an empty value
		// for that key passes the check.
		KeySet(key string) HTTPHeaderChecker
		// KeyNotSet checks the gotten http.Header does not have
		// a specific key set.
		KeyNotSet(key string) HTTPHeaderChecker
		// ValueSet checks the gotten http.Heaser has any key
		// witha a matching value.
		ValueSet(val string) HTTPHeaderChecker
		// ValueNotSet checks the gotten http.Header does not have
		// any key with a matching value.
		ValueNotSet(val string) HTTPHeaderChecker
		// ValueOf checks the gotten http.Header has a value
		// for the matching key that passes the given StringChecker.
		ValueOf(key string, c StringChecker) HTTPHeaderChecker
	}

	// UntypedCheckerProvider provides checks on type interface{}.
	UntypedCheckerProvider interface {
		// Custom checks the gotten value passes the given UntypedPassFunc.
		// The description should typically begin with keywords like "expect"
		// or "should" for intelligible output.
		// For instance, "expect odd number" would output:
		// 	> "expect odd number, got 42"
		Custom(desc string, f UntypedPassFunc) UntypedChecker
	}
)

var (
	// Bytes implements BytesCheckerProvider.
	Bytes BytesCheckerProvider = bytesCheckerFactory{}
	// String implements StringCheckerProvider.
	String StringCheckerProvider = stringCheckerFactory{}
	// Int implements IntCheckerProvider.
	Int IntCheckerProvider = intCheckerFactory{}
	// Duration implements DurationCheckerProvider.
	Duration DurationCheckerProvider = durationCheckerFactory{}
	// HTTPHeader implements HTTPHeaderCheckerProvider.
	HTTPHeader HTTPHeaderCheckerProvider = httpHeaderCheckerFactory{}
	// Untyped implements UntypedCheckerProvider.
	Untyped UntypedCheckerProvider = untypedCheckerFactory{}
)
