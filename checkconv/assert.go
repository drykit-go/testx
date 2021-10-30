// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 30 Oct 21 12:10 UTC

// Package checkconv provides functions to convert typed checkers
// into generic ones.
package checkconv

import (
	"context"
	"net/http"
	"time"

	"github.com/drykit-go/testx/check"
)

// FromBool returns a check.ValueChecker that wraps the given
// check.BoolChecker, so it can be used as a generic checker.
func FromBool(c check.BoolChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(bool)) },
		c.Explain,
	)
}

// FromBytes returns a check.ValueChecker that wraps the given
// check.BytesChecker, so it can be used as a generic checker.
func FromBytes(c check.BytesChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.([]byte)) },
		c.Explain,
	)
}

// FromString returns a check.ValueChecker that wraps the given
// check.StringChecker, so it can be used as a generic checker.
func FromString(c check.StringChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(string)) },
		c.Explain,
	)
}

// FromInt returns a check.ValueChecker that wraps the given
// check.IntChecker, so it can be used as a generic checker.
func FromInt(c check.IntChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(int)) },
		c.Explain,
	)
}

// FromFloat64 returns a check.ValueChecker that wraps the given
// check.Float64Checker, so it can be used as a generic checker.
func FromFloat64(c check.Float64Checker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(float64)) },
		c.Explain,
	)
}

// FromDuration returns a check.ValueChecker that wraps the given
// check.DurationChecker, so it can be used as a generic checker.
func FromDuration(c check.DurationChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(time.Duration)) },
		c.Explain,
	)
}

// FromContext returns a check.ValueChecker that wraps the given
// check.ContextChecker, so it can be used as a generic checker.
func FromContext(c check.ContextChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(context.Context)) },
		c.Explain,
	)
}

// FromHTTPHeader returns a check.ValueChecker that wraps the given
// check.HTTPHeaderChecker, so it can be used as a generic checker.
func FromHTTPHeader(c check.HTTPHeaderChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(http.Header)) },
		c.Explain,
	)
}

// FromHTTPRequest returns a check.ValueChecker that wraps the given
// check.HTTPRequestChecker, so it can be used as a generic checker.
func FromHTTPRequest(c check.HTTPRequestChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(*http.Request)) },
		c.Explain,
	)
}

// FromHTTPResponse returns a check.ValueChecker that wraps the given
// check.HTTPResponseChecker, so it can be used as a generic checker.
func FromHTTPResponse(c check.HTTPResponseChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(*http.Response)) },
		c.Explain,
	)
}

// Assert returns a check.ValueChecker that wraps the given
// check.<Type>Checker (such as check.IntChecker).
//
// It panics if checker is not a known checker type. For instance,
// a custom checker that implements check.IntChecker will be successfully
// converted, while a valid implementation of an unknown interface,
// such as Complex128Checker, will panic.
// For that matter, Cast should be used instead.
func Assert(knownChecker interface{}) check.ValueChecker {
	switch c := knownChecker.(type) {
	case check.BoolChecker:
		return FromBool(c)
	case check.BytesChecker:
		return FromBytes(c)
	case check.StringChecker:
		return FromString(c)
	case check.IntChecker:
		return FromInt(c)
	case check.Float64Checker:
		return FromFloat64(c)
	case check.DurationChecker:
		return FromDuration(c)
	case check.ContextChecker:
		return FromContext(c)
	case check.HTTPHeaderChecker:
		return FromHTTPHeader(c)
	case check.HTTPRequestChecker:
		return FromHTTPRequest(c)
	case check.HTTPResponseChecker:
		return FromHTTPResponse(c)
	case check.ValueChecker:
		return c
	default:
		panic("assert from unknown checker type")
	}
}

// AssertMany returns a slice of check.ValueChecker that wrap the given
// check.<Type>Checkers (such as check.IntChecker).
//
// It panics if any checker is not a known checker type. See Assert
// for further documentation.
func AssertMany(knownCheckers ...interface{}) []check.ValueChecker {
	valueCheckers := []check.ValueChecker{}
	for _, c := range knownCheckers {
		valueCheckers = append(valueCheckers, Assert(c))
	}
	return valueCheckers
}
