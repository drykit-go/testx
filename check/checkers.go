// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 23 Oct 21 19:36 UTC

package check

import (
	"context"
	"net/http"
	"time"
)

// baseChecker is a base checker for all checkers. It implements Explainer.
type baseChecker struct {
	explFunc ExplainFunc
}

// Explain returns a string explaining the reason of a failed check
// for the gotten value.
func (c baseChecker) Explain(label string, got interface{}) string {
	return c.explFunc(label, got)
}

func newBaseChecker(explFunc ExplainFunc) baseChecker {
	return baseChecker{explFunc: explFunc}
}

// boolChecker is an implementation of BoolChecker interface
type boolChecker struct {
	baseChecker
	passFunc BoolPassFunc
}

// Pass returns a boolean that indicates whether the gotten bool value
// passes the current check.
func (c boolChecker) Pass(got bool) bool { return c.passFunc(got) }

// NewBoolChecker returns a custom BoolChecker with the provided
// BoolPassFunc and ExplainFunc.
func NewBoolChecker(passFunc BoolPassFunc, explainFunc ExplainFunc) BoolChecker {
	return boolChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// bytesChecker is an implementation of BytesChecker interface
type bytesChecker struct {
	baseChecker
	passFunc BytesPassFunc
}

// Pass returns a boolean that indicates whether the gotten []byte value
// passes the current check.
func (c bytesChecker) Pass(got []byte) bool { return c.passFunc(got) }

// NewBytesChecker returns a custom BytesChecker with the provided
// BytesPassFunc and ExplainFunc.
func NewBytesChecker(passFunc BytesPassFunc, explainFunc ExplainFunc) BytesChecker {
	return bytesChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// stringChecker is an implementation of StringChecker interface
type stringChecker struct {
	baseChecker
	passFunc StringPassFunc
}

// Pass returns a boolean that indicates whether the gotten string value
// passes the current check.
func (c stringChecker) Pass(got string) bool { return c.passFunc(got) }

// NewStringChecker returns a custom StringChecker with the provided
// StringPassFunc and ExplainFunc.
func NewStringChecker(passFunc StringPassFunc, explainFunc ExplainFunc) StringChecker {
	return stringChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// intChecker is an implementation of IntChecker interface
type intChecker struct {
	baseChecker
	passFunc IntPassFunc
}

// Pass returns a boolean that indicates whether the gotten int value
// passes the current check.
func (c intChecker) Pass(got int) bool { return c.passFunc(got) }

// NewIntChecker returns a custom IntChecker with the provided
// IntPassFunc and ExplainFunc.
func NewIntChecker(passFunc IntPassFunc, explainFunc ExplainFunc) IntChecker {
	return intChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// float64Checker is an implementation of Float64Checker interface
type float64Checker struct {
	baseChecker
	passFunc Float64PassFunc
}

// Pass returns a boolean that indicates whether the gotten float64 value
// passes the current check.
func (c float64Checker) Pass(got float64) bool { return c.passFunc(got) }

// NewFloat64Checker returns a custom Float64Checker with the provided
// Float64PassFunc and ExplainFunc.
func NewFloat64Checker(passFunc Float64PassFunc, explainFunc ExplainFunc) Float64Checker {
	return float64Checker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// durationChecker is an implementation of DurationChecker interface
type durationChecker struct {
	baseChecker
	passFunc DurationPassFunc
}

// Pass returns a boolean that indicates whether the gotten time.Duration value
// passes the current check.
func (c durationChecker) Pass(got time.Duration) bool { return c.passFunc(got) }

// NewDurationChecker returns a custom DurationChecker with the provided
// DurationPassFunc and ExplainFunc.
func NewDurationChecker(passFunc DurationPassFunc, explainFunc ExplainFunc) DurationChecker {
	return durationChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// contextChecker is an implementation of ContextChecker interface
type contextChecker struct {
	baseChecker
	passFunc ContextPassFunc
}

// Pass returns a boolean that indicates whether the gotten context.Context value
// passes the current check.
func (c contextChecker) Pass(got context.Context) bool { return c.passFunc(got) }

// NewContextChecker returns a custom ContextChecker with the provided
// ContextPassFunc and ExplainFunc.
func NewContextChecker(passFunc ContextPassFunc, explainFunc ExplainFunc) ContextChecker {
	return contextChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// httpHeaderChecker is an implementation of HTTPHeaderChecker interface
type httpHeaderChecker struct {
	baseChecker
	passFunc HTTPHeaderPassFunc
}

// Pass returns a boolean that indicates whether the gotten http.Header value
// passes the current check.
func (c httpHeaderChecker) Pass(got http.Header) bool { return c.passFunc(got) }

// NewHTTPHeaderChecker returns a custom HTTPHeaderChecker with the provided
// HTTPHeaderPassFunc and ExplainFunc.
func NewHTTPHeaderChecker(passFunc HTTPHeaderPassFunc, explainFunc ExplainFunc) HTTPHeaderChecker {
	return httpHeaderChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// httpRequestChecker is an implementation of HTTPRequestChecker interface
type httpRequestChecker struct {
	baseChecker
	passFunc HTTPRequestPassFunc
}

// Pass returns a boolean that indicates whether the gotten *http.Request value
// passes the current check.
func (c httpRequestChecker) Pass(got *http.Request) bool { return c.passFunc(got) }

// NewHTTPRequestChecker returns a custom HTTPRequestChecker with the provided
// HTTPRequestPassFunc and ExplainFunc.
func NewHTTPRequestChecker(passFunc HTTPRequestPassFunc, explainFunc ExplainFunc) HTTPRequestChecker {
	return httpRequestChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// httpResponseChecker is an implementation of HTTPResponseChecker interface
type httpResponseChecker struct {
	baseChecker
	passFunc HTTPResponsePassFunc
}

// Pass returns a boolean that indicates whether the gotten *http.Response value
// passes the current check.
func (c httpResponseChecker) Pass(got *http.Response) bool { return c.passFunc(got) }

// NewHTTPResponseChecker returns a custom HTTPResponseChecker with the provided
// HTTPResponsePassFunc and ExplainFunc.
func NewHTTPResponseChecker(passFunc HTTPResponsePassFunc, explainFunc ExplainFunc) HTTPResponseChecker {
	return httpResponseChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}

// valueChecker is an implementation of ValueChecker interface
type valueChecker struct {
	baseChecker
	passFunc ValuePassFunc
}

// Pass returns a boolean that indicates whether the gotten interface{} value
// passes the current check.
func (c valueChecker) Pass(got interface{}) bool { return c.passFunc(got) }

// NewValueChecker returns a custom ValueChecker with the provided
// ValuePassFunc and ExplainFunc.
func NewValueChecker(passFunc ValuePassFunc, explainFunc ExplainFunc) ValueChecker {
	return valueChecker{baseChecker: newBaseChecker(explainFunc), passFunc: passFunc}
}
