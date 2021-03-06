// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 30 Oct 21 12:10 UTC

package check

import (
	"regexp"
	"time"
)

type (
	// BoolCheckerProvider provides checks on type bool.
	BoolCheckerProvider interface {
		// Is checks the gotten bool is equal to the target.
		Is(tar bool) BoolChecker
	}

	// BytesCheckerProvider provides checks on type []byte.
	BytesCheckerProvider interface {
		// AsMap checks the gotten []byte passes the given mapChecker
		// once json-unmarshaled to a map[string]interface{}.
		// It fails if it is not a valid JSON.
		AsMap(mapChecker ValueChecker) BytesChecker
		// AsString checks the gotten []byte passes the given StringChecker
		// once converted to a string.
		AsString(c StringChecker) BytesChecker
		// Contains checks the gotten []byte contains a specific subslice.
		Contains(subslice []byte) BytesChecker
		// Is checks the gotten []byte is equal to the target.
		Is(tar []byte) BytesChecker
		// Len checks the gotten []byte's length passes the provided
		// IntChecker.
		Len(c IntChecker) BytesChecker
		// Not checks the gotten []byte is not equal to the target.
		Not(values ...[]byte) BytesChecker
		// NotContains checks the gotten []byte contains a specific subslice.
		NotContains(subslice []byte) BytesChecker
		// SameJSON checks the gotten []byte and the target read as the same
		// JSON value, ignoring formatting and keys order.
		SameJSON(tar []byte) BytesChecker
	}

	// ContextCheckerProvider provides checks on type context.Context.
	ContextCheckerProvider interface {
		// Done checks the gotten context is done.
		Done(expectDone bool) ContextChecker
		// HasKeys checks the gotten context has the given keys set.
		HasKeys(keys ...interface{}) ContextChecker
		// Value checks the gotten context's value for the given key passes
		// the given ValueChecker. It fails if value is nil.
		//
		// Examples:
		// 	Context.Value("userID", Value.Is("abcde"))
		// 	Context.Value("userID", checkconv.Assert(String.Contains("abc")))
		Value(key interface{}, c ValueChecker) ContextChecker
	}

	// DurationCheckerProvider provides checks on type time.Duration.
	DurationCheckerProvider interface {
		// InRange checks the gotten time.Duration is in range [lo:hi]
		InRange(lo, hi time.Duration) DurationChecker
		// OutRange checks the gotten time.Duration is not in range [lo:hi]
		OutRange(lo, hi time.Duration) DurationChecker
		// Over checks the gotten time.Duration is over the target duration.
		Over(tar time.Duration) DurationChecker
		// Under checks the gotten time.Duration is under the target duration.
		Under(tar time.Duration) DurationChecker
	}

	// Float64CheckerProvider provides checks on type float64.
	Float64CheckerProvider interface {
		// GT checks the gotten float64 is greater than the target.
		GT(tar float64) Float64Checker
		// GTE checks the gotten float64 is greater or equal to the target.
		GTE(tar float64) Float64Checker
		// InRange checks the gotten float64 is in the closed interval [lo:hi].
		InRange(lo, hi float64) Float64Checker
		// Is checks the gotten float64 is equal to the target.
		Is(tar float64) Float64Checker
		// LT checks the gotten float64 is lesser than the target.
		LT(tar float64) Float64Checker
		// LTE checks the gotten float64 is lesser or equal to the target.
		LTE(tar float64) Float64Checker
		// Not checks the gotten float64 is not equal to the target.
		Not(values ...float64) Float64Checker
		// OutRange checks the gotten float64 is not in the closed interval [lo:hi].
		OutRange(lo, hi float64) Float64Checker
	}

	// HTTPHeaderCheckerProvider provides checks on type http.Header.
	HTTPHeaderCheckerProvider interface {
		// CheckValue checks the gotten http.Header has a value for the matching key
		// that passes the given StringChecker.
		// It only checks the first result for the given key.
		CheckValue(key string, c StringChecker) HTTPHeaderChecker
		// HasKey checks the gotten http.Header has a specific key set.
		// The corresponding value is ignored, meaning an empty value
		// for that key passes the check.
		HasKey(key string) HTTPHeaderChecker
		// HasNotKey checks the gotten http.Header does not have
		// a specific key set.
		HasNotKey(key string) HTTPHeaderChecker
		// HasNotValue checks the gotten http.Header does not have a value equal to val.
		// It only compares the first result for each key.
		HasNotValue(val string) HTTPHeaderChecker
		// HasValue checks the gotten http.Header has any value equal to val.
		// It only compares the first result for each key.
		HasValue(val string) HTTPHeaderChecker
	}

	// HTTPRequestCheckerProvider provides checks on type *http.Request.
	HTTPRequestCheckerProvider interface {
		// Body checks the gotten *http.Request Body passes the input BytesChecker.
		// It should be used only once on a same *http.Request as it closes its body
		// after reading it.
		Body(c BytesChecker) HTTPRequestChecker
		// ContentLength checks the gotten *http.Request ContentLength passes
		// the input IntChecker.
		ContentLength(c IntChecker) HTTPRequestChecker
		// Context checks the gotten *http.Request Context passes
		// the input ContextChecker.
		Context(c ContextChecker) HTTPRequestChecker
		// Header checks the gotten *http.Request Header passes
		// the input HTTPHeaderChecker.
		Header(c HTTPHeaderChecker) HTTPRequestChecker
	}

	// HTTPResponseCheckerProvider provides checks on type *http.Response.
	HTTPResponseCheckerProvider interface {
		// Body checks the gotten *http.Response Body passes the input BytesChecker.
		// It should be used only once on a same *http.Response as it closes its body
		// after reading it.
		Body(c BytesChecker) HTTPResponseChecker
		// ContentLength checks the gotten *http.Response ContentLength passes
		// the input IntChecker.
		ContentLength(c IntChecker) HTTPResponseChecker
		// Header checks the gotten *http.Response Header passes
		// the input HTTPHeaderChecker.
		Header(c HTTPHeaderChecker) HTTPResponseChecker
		// Status checks the gotten *http.Response Status passes
		// the input StringChecker.
		Status(c StringChecker) HTTPResponseChecker
		// StatusCode checks the gotten *http.Response StatusCode passes
		// the input IntChecker.
		StatusCode(c IntChecker) HTTPResponseChecker
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
		Not(values ...int) IntChecker
		// OutRange checks the gotten int is not in the closed interval [lo:hi].
		OutRange(lo, hi int) IntChecker
	}

	// MapCheckerProvider provides checks on kind map.
	MapCheckerProvider interface {
		ValueCheckerProvider

		// CheckValues checks the gotten map's values corresponding to the given keys
		// pass the given checker. A key not found is considered a fail.
		// If len(keys) == 0, the check is made on all map values.
		CheckValues(c ValueChecker, keys ...interface{}) ValueChecker
		// HasKeys checks the gotten map has the given keys set.
		HasKeys(keys ...interface{}) ValueChecker
		// HasNotKeys checks the gotten map has the given keys set.
		HasNotKeys(keys ...interface{}) ValueChecker
		// HasNotValues checks the gotten map has not the given values set.
		HasNotValues(values ...interface{}) ValueChecker
		// HasValues checks the gotten map has the given values set.
		HasValues(values ...interface{}) ValueChecker
		// Len checks the gotten map passes the given IntChecker.
		Len(c IntChecker) ValueChecker
	}

	// SliceCheckerProvider provides checks on kind slice.
	SliceCheckerProvider interface {
		ValueCheckerProvider

		// Cap checks the capacity of the gotten slice passes the given IntChecker.
		Cap(c IntChecker) ValueChecker
		// CheckValues checks the values of the gotten slice pass the given ValueChecker.
		// If a filterFunc is provided, the values not passing it are ignored.
		CheckValues(c ValueChecker, filters ...func(i int, v interface{}) bool) ValueChecker
		// HasNotValues checks the gotten slice has not the given values set.
		HasNotValues(values ...interface{}) ValueChecker
		// HasValues checks the gotten slice has the given values set.
		HasValues(values ...interface{}) ValueChecker
		// Len checks the length of the gotten slice passes the given IntChecker.
		Len(c IntChecker) ValueChecker
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
		// Not checks the gotten string is not equal to the target.
		Not(values ...string) StringChecker
		// NotContains checks the gotten string do not contain the target
		// substring.
		NotContains(sub string) StringChecker
		// NotMatch checks the gotten string do not match the given regexp.
		NotMatch(rgx *regexp.Regexp) StringChecker
	}

	// StructCheckerProvider provides checks on kind struct.
	StructCheckerProvider interface {
		ValueCheckerProvider

		// CheckFields checks all given fields pass the ValueChecker.
		// It panics if the fields do not exist or are not exported,
		// or if the tested value is not a struct.
		CheckFields(c ValueChecker, fields []string) ValueChecker
		// FieldsEqual checks all given fields equal the exp value.
		// It panics if the fields do not exist or are not exported,
		// or if the tested value is not a struct.
		FieldsEqual(exp interface{}, fields []string) ValueChecker
	}

	// ValueCheckerProvider provides checks on type interface{}.
	ValueCheckerProvider interface {
		// Custom checks the gotten value passes the given ValuePassFunc.
		// The description should give information about the expected value,
		// as it outputs in format "exp <desc>" in case of failure.
		Custom(desc string, f ValuePassFunc) ValueChecker
		// Is checks the gotten value is equal to the target.
		Is(tar interface{}) ValueChecker
		// IsZero checks the gotten value is a zero value, indicating it might not
		// have been initialized.
		IsZero() ValueChecker
		// Not checks the gotten value is not equal to the target.
		Not(values ...interface{}) ValueChecker
		// NotZero checks the gotten struct contains at least 1 non-zero value,
		// meaning it has been initialized.
		NotZero() ValueChecker
		// SameJSON checks the gotten value and the target value
		// produce the same JSON, ignoring formatting and keys order.
		// It panics if any error occurs in the marshaling process.
		SameJSON(tar interface{}) ValueChecker
	}
)

var (
	// Bool implements BoolCheckerProvider.
	Bool BoolCheckerProvider = boolCheckerProvider{}
	// Bytes implements BytesCheckerProvider.
	Bytes BytesCheckerProvider = bytesCheckerProvider{}
	// Context implements ContextCheckerProvider.
	Context ContextCheckerProvider = contextCheckerProvider{}
	// Duration implements DurationCheckerProvider.
	Duration DurationCheckerProvider = durationCheckerProvider{}
	// Float64 implements Float64CheckerProvider.
	Float64 Float64CheckerProvider = float64CheckerProvider{}
	// HTTPHeader implements HTTPHeaderCheckerProvider.
	HTTPHeader HTTPHeaderCheckerProvider = httpHeaderCheckerProvider{}
	// HTTPRequest implements HTTPRequestCheckerProvider.
	HTTPRequest HTTPRequestCheckerProvider = httpRequestCheckerProvider{}
	// HTTPResponse implements HTTPResponseCheckerProvider.
	HTTPResponse HTTPResponseCheckerProvider = httpResponseCheckerProvider{}
	// Int implements IntCheckerProvider.
	Int IntCheckerProvider = intCheckerProvider{}
	// Map implements MapCheckerProvider.
	Map MapCheckerProvider = mapCheckerProvider{}
	// Slice implements SliceCheckerProvider.
	Slice SliceCheckerProvider = sliceCheckerProvider{}
	// String implements StringCheckerProvider.
	String StringCheckerProvider = stringCheckerProvider{}
	// Struct implements StructCheckerProvider.
	Struct StructCheckerProvider = structCheckerProvider{}
	// Value implements ValueCheckerProvider.
	Value ValueCheckerProvider = valueCheckerProvider{}
)
