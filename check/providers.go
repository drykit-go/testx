package check

import (
	"context"
	"net/http"
	"regexp"
	"time"
)

type (
	// BoolCheckerProvider provides checks on type bool.
	BoolCheckerProvider interface {
		// Is checks the gotten bool is equal to the target.
		Is(tar bool) Checker[bool]
	}

	// BytesCheckerProvider provides checks on type []byte.
	BytesCheckerProvider interface {
		// AsMap checks the gotten []byte passes the given mapChecker
		// once json-unmarshaled to a map[string]any.
		// It fails if it is not a valid JSON.
		AsMap(mapChecker Checker[any]) Checker[[]byte]
		// AsString checks the gotten []byte passes the given Checker[string]
		// once converted to a string.
		AsString(c Checker[string]) Checker[[]byte]
		// Contains checks the gotten []byte contains a specific subslice.
		Contains(subslice []byte) Checker[[]byte]
		// Is checks the gotten []byte is equal to the target.
		Is(tar []byte) Checker[[]byte]
		// Len checks the gotten []byte's length passes the provided
		// Checker[int].
		Len(c Checker[int]) Checker[[]byte]
		// Not checks the gotten []byte is not equal to the target.
		Not(values ...[]byte) Checker[[]byte]
		// NotContains checks the gotten []byte contains a specific subslice.
		NotContains(subslice []byte) Checker[[]byte]
		// SameJSON checks the gotten []byte and the target read as the same
		// JSON value, ignoring formatting and keys order.
		SameJSON(tar []byte) Checker[[]byte]
	}

	// ContextCheckerProvider provides checks on type context.Context.
	ContextCheckerProvider interface {
		// Done checks the gotten context is done.
		Done(expectDone bool) Checker[context.Context]
		// HasKeys checks the gotten context has the given keys set.
		HasKeys(keys ...any) Checker[context.Context]
		// Value checks the gotten context's value for the given key passes
		// the given Checker[any]. It fails if value is nil.
		//
		// Examples:
		// 	Context.Value("userID", Value.Is("abcde"))
		// 	Context.Value("userID", checkconv.Assert(String.Contains("abc")))
		Value(key any, c Checker[any]) Checker[context.Context]
	}

	// DurationCheckerProvider provides checks on type time.Duration.
	DurationCheckerProvider interface {
		// InRange checks the gotten time.Duration is in range [lo:hi]
		InRange(lo, hi time.Duration) Checker[time.Duration]
		// OutRange checks the gotten time.Duration is not in range [lo:hi]
		OutRange(lo, hi time.Duration) Checker[time.Duration]
		// Over checks the gotten time.Duration is over the target duration.
		Over(tar time.Duration) Checker[time.Duration]
		// Under checks the gotten time.Duration is under the target duration.
		Under(tar time.Duration) Checker[time.Duration]
	}

	// Float64CheckerProvider provides checks on type float64.
	Float64CheckerProvider interface {
		// GT checks the gotten float64 is greater than the target.
		GT(tar float64) Checker[float64]
		// GTE checks the gotten float64 is greater or equal to the target.
		GTE(tar float64) Checker[float64]
		// InRange checks the gotten float64 is in the closed interval [lo:hi].
		InRange(lo, hi float64) Checker[float64]
		// Is checks the gotten float64 is equal to the target.
		Is(tar float64) Checker[float64]
		// LT checks the gotten float64 is lesser than the target.
		LT(tar float64) Checker[float64]
		// LTE checks the gotten float64 is lesser or equal to the target.
		LTE(tar float64) Checker[float64]
		// Not checks the gotten float64 is not equal to the target.
		Not(values ...float64) Checker[float64]
		// OutRange checks the gotten float64 is not in the closed interval [lo:hi].
		OutRange(lo, hi float64) Checker[float64]
	}

	// HTTPHeaderCheckerProvider provides checks on type http.Header.
	HTTPHeaderCheckerProvider interface {
		// CheckValue checks the gotten http.Header has a value for the matching key
		// that passes the given Checker[string].
		// It only checks the first result for the given key.
		CheckValue(key string, c Checker[string]) Checker[http.Header]
		// HasKey checks the gotten http.Header has a specific key set.
		// The corresponding value is ignored, meaning an empty value
		// for that key passes the check.
		HasKey(key string) Checker[http.Header]
		// HasNotKey checks the gotten http.Header does not have
		// a specific key set.
		HasNotKey(key string) Checker[http.Header]
		// HasNotValue checks the gotten http.Header does not have a value equal to val.
		// It only compares the first result for each key.
		HasNotValue(val string) Checker[http.Header]
		// HasValue checks the gotten http.Header has any value equal to val.
		// It only compares the first result for each key.
		HasValue(val string) Checker[http.Header]
	}

	// HTTPRequestCheckerProvider provides checks on type *http.Request.
	HTTPRequestCheckerProvider interface {
		// Body checks the gotten *http.Request Body passes the input Checker[[]byte].
		// It should be used only once on a same *http.Request as it closes its body
		// after reading it.
		Body(c Checker[[]byte]) Checker[*http.Request]
		// ContentLength checks the gotten *http.Request ContentLength passes
		// the input Checker[int].
		ContentLength(c Checker[int]) Checker[*http.Request]
		// Context checks the gotten *http.Request Context passes
		// the input Checker[context.Context].
		Context(c Checker[context.Context]) Checker[*http.Request]
		// Header checks the gotten *http.Request Header passes
		// the input Checker[http.Header].
		Header(c Checker[http.Header]) Checker[*http.Request]
	}

	// HTTPResponseCheckerProvider provides checks on type *http.Response.
	HTTPResponseCheckerProvider interface {
		// Body checks the gotten *http.Response Body passes the input Checker[[]byte].
		// It should be used only once on a same *http.Response as it closes its body
		// after reading it.
		Body(c Checker[[]byte]) Checker[*http.Response]
		// ContentLength checks the gotten *http.Response ContentLength passes
		// the input Checker[int].
		ContentLength(c Checker[int]) Checker[*http.Response]
		// Header checks the gotten *http.Response Header passes
		// the input Checker[http.Header].
		Header(c Checker[http.Header]) Checker[*http.Response]
		// Status checks the gotten *http.Response Status passes
		// the input Checker[string].
		Status(c Checker[string]) Checker[*http.Response]
		// StatusCode checks the gotten *http.Response StatusCode passes
		// the input Checker[int].
		StatusCode(c Checker[int]) Checker[*http.Response]
	}

	// IntCheckerProvider provides checks on type int.
	IntCheckerProvider interface {
		// GT checks the gotten int is greater than the target.
		GT(tar int) Checker[int]
		// GTE checks the gotten int is greater or equal to the target.
		GTE(tar int) Checker[int]
		// InRange checks the gotten int is in the closed interval [lo:hi].
		InRange(lo, hi int) Checker[int]
		// Is checks the gotten int is equal to the target.
		Is(tar int) Checker[int]
		// LT checks the gotten int is lesser than the target.
		LT(tar int) Checker[int]
		// LTE checks the gotten int is lesser or equal to the target.
		LTE(tar int) Checker[int]
		// Not checks the gotten int is not equal to the target.
		Not(values ...int) Checker[int]
		// OutRange checks the gotten int is not in the closed interval [lo:hi].
		OutRange(lo, hi int) Checker[int]
	}

	// MapCheckerProvider provides checks on kind map.
	MapCheckerProvider interface {
		ValueCheckerProvider

		// CheckValues checks the gotten map's values corresponding to the given keys
		// pass the given checker. A key not found is considered a fail.
		// If len(keys) == 0, the check is made on all map values.
		CheckValues(c Checker[any], keys ...any) Checker[any]
		// HasKeys checks the gotten map has the given keys set.
		HasKeys(keys ...any) Checker[any]
		// HasNotKeys checks the gotten map has the given keys set.
		HasNotKeys(keys ...any) Checker[any]
		// HasNotValues checks the gotten map has not the given values set.
		HasNotValues(values ...any) Checker[any]
		// HasValues checks the gotten map has the given values set.
		HasValues(values ...any) Checker[any]
		// Len checks the gotten map passes the given Checker[int].
		Len(c Checker[int]) Checker[any]
	}

	// SliceCheckerProvider provides checks on kind slice.
	SliceCheckerProvider interface {
		ValueCheckerProvider

		// Cap checks the capacity of the gotten slice passes the given Checker[int].
		Cap(c Checker[int]) Checker[any]
		// CheckValues checks the values of the gotten slice pass the given Checker[any].
		// If a filterFunc is provided, the values not passing it are ignored.
		CheckValues(c Checker[any], filters ...func(i int, v any) bool) Checker[any]
		// HasNotValues checks the gotten slice has not the given values set.
		HasNotValues(values ...any) Checker[any]
		// HasValues checks the gotten slice has the given values set.
		HasValues(values ...any) Checker[any]
		// Len checks the length of the gotten slice passes the given Checker[int].
		Len(c Checker[int]) Checker[any]
	}

	// StringCheckerProvider provides checks on type string.
	StringCheckerProvider interface {
		// Contains checks the gotten string contains the target substring.
		Contains(sub string) Checker[string]
		// Is checks the gotten string is equal to the target.
		Is(tar string) Checker[string]
		// Len checks the gotten string's length passes the given Checker[int].
		Len(c Checker[int]) Checker[string]
		// Match checks the gotten string matches the given regexp.
		Match(rgx *regexp.Regexp) Checker[string]
		// Not checks the gotten string is not equal to the target.
		Not(values ...string) Checker[string]
		// NotContains checks the gotten string do not contain the target
		// substring.
		NotContains(sub string) Checker[string]
		// NotMatch checks the gotten string do not match the given regexp.
		NotMatch(rgx *regexp.Regexp) Checker[string]
	}

	// StructCheckerProvider provides checks on kind struct.
	StructCheckerProvider interface {
		ValueCheckerProvider

		// CheckFields checks all given fields pass the Checker[any].
		// It panics if the fields do not exist or are not exported,
		// or if the tested value is not a struct.
		CheckFields(c Checker[any], fields []string) Checker[any]
		// FieldsEqual checks all given fields equal the exp value.
		// It panics if the fields do not exist or are not exported,
		// or if the tested value is not a struct.
		FieldsEqual(exp any, fields []string) Checker[any]
	}

	// ValueCheckerProvider provides checks on type any.
	ValueCheckerProvider interface {
		// Custom checks the gotten value passes the given PassFunc[any].
		// The description should give information about the expected value,
		// as it outputs in format "exp <desc>" in case of failure.
		Custom(desc string, f PassFunc[any]) Checker[any]
		// Is checks the gotten value is equal to the target.
		Is(tar any) Checker[any]
		// IsZero checks the gotten value is a zero value, indicating it might not
		// have been initialized.
		IsZero() Checker[any]
		// Not checks the gotten value is not equal to the target.
		Not(values ...any) Checker[any]
		// NotZero checks the gotten struct contains at least 1 non-zero value,
		// meaning it has been initialized.
		NotZero() Checker[any]
		// SameJSON checks the gotten value and the target value
		// produce the same JSON, ignoring formatting and keys order.
		// It panics if any error occurs in the marshaling process.
		SameJSON(tar any) Checker[any]
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
