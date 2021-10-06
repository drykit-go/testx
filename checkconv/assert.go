// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 15 Sep 21 22:10 UTC

// Package checkconv provides functions to convert typed checks
// into generic ones.
package checkconv

import (
	"log"
	"net/http"
	"time"

	"github.com/drykit-go/testx/check"
)

// FromBool returns a new check.ValueChecker from the given Bool checker.
// It can be used to facilitate checkers usage by test runners.
func FromBool(c check.BoolChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(bool)) },
		c.Explain,
	)
}

// FromBytes returns a new check.ValueChecker from the given Bytes checker.
// It can be used to facilitate checkers usage by test runners.
func FromBytes(c check.BytesChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.([]byte)) },
		c.Explain,
	)
}

// FromString returns a new check.ValueChecker from the given String checker.
// It can be used to facilitate checkers usage by test runners.
func FromString(c check.StringChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(string)) },
		c.Explain,
	)
}

// FromInt returns a new check.ValueChecker from the given Int checker.
// It can be used to facilitate checkers usage by test runners.
func FromInt(c check.IntChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(int)) },
		c.Explain,
	)
}

// FromFloat64 returns a new check.ValueChecker from the given Float64 checker.
// It can be used to facilitate checkers usage by test runners.
func FromFloat64(c check.Float64Checker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(float64)) },
		c.Explain,
	)
}

// FromDuration returns a new check.ValueChecker from the given Duration checker.
// It can be used to facilitate checkers usage by test runners.
func FromDuration(c check.DurationChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(time.Duration)) },
		c.Explain,
	)
}

// FromHTTPHeader returns a new check.ValueChecker from the given HTTPHeader checker.
// It can be used to facilitate checkers usage by test runners.
func FromHTTPHeader(c check.HTTPHeaderChecker) check.ValueChecker {
	return check.NewValueChecker(
		func(got interface{}) bool { return c.Pass(got.(http.Header)) },
		c.Explain,
	)
}

// Assert takes a known typed checker (such as check.IntChecker)
// and returns it as a check.ValueChecker.
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
	case check.HTTPHeaderChecker:
		return FromHTTPHeader(c)
	case check.ValueChecker:
		return c
	default:
		log.Panic("assert from unknown checker type")
		return nil
	}
}

// AssertMany converts the given known checkers as described by Assert,
// and returns them as a slice of check.ValueChecker.
func AssertMany(knownCheckers ...interface{}) []check.ValueChecker {
	valueCheckers := []check.ValueChecker{}
	for _, c := range knownCheckers {
		valueCheckers = append(valueCheckers, Assert(c))
	}
	return valueCheckers
}
