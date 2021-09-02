// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 01 Sep 21 22:21 UTC

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
		// Len checks the gotten []byte's length passes the provided
		// IntChecker.
		Len(c IntChecker) BytesChecker
		// SameJSON checks the gotten []byte and the target returns
		// the same JSON object.
		SameJSON(tar []byte) BytesChecker
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
		// KeyNotSet checks the gotten http.Header does not have
		// a specific key set.
		KeyNotSet(key string) HTTPHeaderChecker
		// KeySet checks the gotten http.Header has a spcific key set.
		// The corresponding value is ignored, meaning an empty value
		// for that key passes the check.
		KeySet(key string) HTTPHeaderChecker
		// ValueNotSet checks the gotten http.Header does not have
		// any key with a matching value.
		ValueNotSet(val string) HTTPHeaderChecker
		// ValueOf checks the gotten http.Header has a value
		// for the matching key that passes the given StringChecker.
		ValueOf(key string, c StringChecker) HTTPHeaderChecker
		// ValueSet checks the gotten http.Heaser has any key
		// witha a matching value.
		ValueSet(val string) HTTPHeaderChecker
	}

	// IntCheckerProvider provides checks on type int.
	IntCheckerProvider interface {
		// GT checks the gotten int is greater than the target.
		GT(tar int) IntChecker
		// GTE checks the gotten int is greater or equal to the target.
		GTE(tar int) IntChecker
		// InRange checks the gotten int is in the closed interval [lo:hi].
		InRange(lo, hi int) IntChecker
		// Is checks the gotten int is equal to the target.
		Is(tar int) IntChecker
		// LT checks the gotten int is lesser than the target.
		LT(tar int) IntChecker
		// LTE checks the gotten int is lesser or equal to the target.
		LTE(tar int) IntChecker
		// Not checks the gotten int is not equal to the target.
		Not(tar int) IntChecker
		// OutRange checks the gotten int is not in the closed interval [lo:hi].
		OutRange(lo, hi int) IntChecker
	}

	// StringCheckerProvider provides checks on type string.
	StringCheckerProvider interface {
		// Contains checks the gotten string contains the target substring.
		Contains(sub string) StringChecker
		// Is checks the gotten string is equal to the target.
		Is(tar string) StringChecker
		// Len checks the gotten string's length passes the given IntChecker.
		Len(c IntChecker) StringChecker
		// Match checks the gotten string matches the given regexp.
		Match(rgx *regexp.Regexp) StringChecker
		// NotContains checks the gotten string do not contain the target
		// substring.
		NotContains(sub string) StringChecker
		// NotMatch checks the gotten string do not match the given regexp.
		NotMatch(rgx *regexp.Regexp) StringChecker
	}

	// ValueCheckerProvider provides checks on type interface{}.
	ValueCheckerProvider interface {
		// Custom checks the gotten value passes the given ValuePassFunc.
		// The description should typically begin with keywords like "expect"
		// or "should" for intelligible output.
		// For instance, "expect odd number" would output:
		// 	> "expect odd number, got 42"
		Custom(desc string, f ValuePassFunc) ValueChecker
	}
)

var (
	// Bytes implements BytesCheckerProvider.
	Bytes BytesCheckerProvider = bytesCheckerProvider{}
	// Duration implements DurationCheckerProvider.
	Duration DurationCheckerProvider = durationCheckerProvider{}
	// HTTPHeader implements HTTPHeaderCheckerProvider.
	HTTPHeader HTTPHeaderCheckerProvider = httpHeaderCheckerProvider{}
	// Int implements IntCheckerProvider.
	Int IntCheckerProvider = intCheckerProvider{}
	// String implements StringCheckerProvider.
	String StringCheckerProvider = stringCheckerProvider{}
	// Value implements ValueCheckerProvider.
	Value ValueCheckerProvider = valueCheckerProvider{}
)