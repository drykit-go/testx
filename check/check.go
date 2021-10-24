// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 24 Oct 21 15:20 UTC

// Package check provides types to perform checks on values
// in a testing context.
package check

import (
	"context"
	"net/http"
	"time"
)

type (
	// BoolPassFunc is the required method to implement BoolPasser.
	// It returns a boolean that indicates whether the gotten bool value
	// passes the current check.
	BoolPassFunc func(got bool) bool
	// BytesPassFunc is the required method to implement BytesPasser.
	// It returns a boolean that indicates whether the gotten []byte value
	// passes the current check.
	BytesPassFunc func(got []byte) bool
	// StringPassFunc is the required method to implement StringPasser.
	// It returns a boolean that indicates whether the gotten string value
	// passes the current check.
	StringPassFunc func(got string) bool
	// IntPassFunc is the required method to implement IntPasser.
	// It returns a boolean that indicates whether the gotten int value
	// passes the current check.
	IntPassFunc func(got int) bool
	// Float64PassFunc is the required method to implement Float64Passer.
	// It returns a boolean that indicates whether the gotten float64 value
	// passes the current check.
	Float64PassFunc func(got float64) bool
	// DurationPassFunc is the required method to implement DurationPasser.
	// It returns a boolean that indicates whether the gotten time.Duration value
	// passes the current check.
	DurationPassFunc func(got time.Duration) bool
	// ContextPassFunc is the required method to implement ContextPasser.
	// It returns a boolean that indicates whether the gotten context.Context value
	// passes the current check.
	ContextPassFunc func(got context.Context) bool
	// HTTPHeaderPassFunc is the required method to implement HTTPHeaderPasser.
	// It returns a boolean that indicates whether the gotten http.Header value
	// passes the current check.
	HTTPHeaderPassFunc func(got http.Header) bool
	// HTTPRequestPassFunc is the required method to implement HTTPRequestPasser.
	// It returns a boolean that indicates whether the gotten *http.Request value
	// passes the current check.
	HTTPRequestPassFunc func(got *http.Request) bool
	// HTTPResponsePassFunc is the required method to implement HTTPResponsePasser.
	// It returns a boolean that indicates whether the gotten *http.Response value
	// passes the current check.
	HTTPResponsePassFunc func(got *http.Response) bool
	// ValuePassFunc is the required method to implement ValuePasser.
	// It returns a boolean that indicates whether the gotten interface{} value
	// passes the current check.
	ValuePassFunc func(got interface{}) bool

	// ExplainFunc is the required method to implement Explainer.
	// It returns a string explaining why the gotten value failed the check.
	// The label provides some context, such as "response code".
	ExplainFunc func(label string, got interface{}) string
)

type (
	// BoolPasser provides a method Pass that returns a bool that indicates
	// whether the gotten bool value passes the current check.
	BoolPasser interface{ Pass(got bool) bool }
	// BytesPasser provides a method Pass that returns a bool that indicates
	// whether the gotten []byte value passes the current check.
	BytesPasser interface{ Pass(got []byte) bool }
	// StringPasser provides a method Pass that returns a bool that indicates
	// whether the gotten string value passes the current check.
	StringPasser interface{ Pass(got string) bool }
	// IntPasser provides a method Pass that returns a bool that indicates
	// whether the gotten int value passes the current check.
	IntPasser interface{ Pass(got int) bool }
	// Float64Passer provides a method Pass that returns a bool that indicates
	// whether the gotten float64 value passes the current check.
	Float64Passer interface{ Pass(got float64) bool }
	// DurationPasser provides a method Pass that returns a bool that indicates
	// whether the gotten time.Duration value passes the current check.
	DurationPasser interface{ Pass(got time.Duration) bool }
	// ContextPasser provides a method Pass that returns a bool that indicates
	// whether the gotten context.Context value passes the current check.
	ContextPasser interface {
		Pass(got context.Context) bool
	}
	// HTTPHeaderPasser provides a method Pass that returns a bool that indicates
	// whether the gotten http.Header value passes the current check.
	HTTPHeaderPasser interface{ Pass(got http.Header) bool }
	// HTTPRequestPasser provides a method Pass that returns a bool that indicates
	// whether the gotten *http.Request value passes the current check.
	HTTPRequestPasser interface{ Pass(got *http.Request) bool }
	// HTTPResponsePasser provides a method Pass that returns a bool that indicates
	// whether the gotten *http.Response value passes the current check.
	HTTPResponsePasser interface{ Pass(got *http.Response) bool }
	// ValuePasser provides a method Pass that returns a bool that indicates
	// whether the gotten interface{} value passes the current check.
	ValuePasser interface{ Pass(got interface{}) bool }

	// Explainer provides a method Explain describing the reason of a failed check.
	Explainer interface {
		Explain(label string, got interface{}) string
	}
)

type (
	// BoolChecker implements both BoolPasser and Explainer interfaces.
	BoolChecker interface {
		BoolPasser
		Explainer
	}
	// BytesChecker implements both BytesPasser and Explainer interfaces.
	BytesChecker interface {
		BytesPasser
		Explainer
	}
	// StringChecker implements both StringPasser and Explainer interfaces.
	StringChecker interface {
		StringPasser
		Explainer
	}
	// IntChecker implements both IntPasser and Explainer interfaces.
	IntChecker interface {
		IntPasser
		Explainer
	}
	// Float64Checker implements both Float64Passer and Explainer interfaces.
	Float64Checker interface {
		Float64Passer
		Explainer
	}
	// DurationChecker implements both DurationPasser and Explainer interfaces.
	DurationChecker interface {
		DurationPasser
		Explainer
	}
	// ContextChecker implements both ContextPasser and Explainer interfaces.
	ContextChecker interface {
		ContextPasser
		Explainer
	}
	// HTTPHeaderChecker implements both HTTPHeaderPasser and Explainer interfaces.
	HTTPHeaderChecker interface {
		HTTPHeaderPasser
		Explainer
	}
	// HTTPRequestChecker implements both HTTPRequestPasser and Explainer interfaces.
	HTTPRequestChecker interface {
		HTTPRequestPasser
		Explainer
	}
	// HTTPResponseChecker implements both HTTPResponsePasser and Explainer interfaces.
	HTTPResponseChecker interface {
		HTTPResponsePasser
		Explainer
	}
	// ValueChecker implements both ValuePasser and Explainer interfaces.
	ValueChecker interface {
		ValuePasser
		Explainer
	}
)
