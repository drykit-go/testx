package check

import (
	"regexp"
	"time"
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
