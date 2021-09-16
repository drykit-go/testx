package testx

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/drykit-go/cond"

	"github.com/drykit-go/testx/check"
)

/*
	Runners interfaces
*/

// Runner provides a method Run to perform various tests.
type Runner interface {
	// Run runs a test and fails it if a check does not pass.
	Run(t *testing.T)
}

// ValueRunner provides methods to perform tests on a single value.
type ValueRunner interface {
	Runner
	// DryRun returns a Resulter to access test results
	// without running it.
	DryRun() Resulter
	// Exp adds an equality check on the tested value.
	Exp(value interface{}) ValueRunner
	// Not adds inequality checks on the tested value.
	Not(values ...interface{}) ValueRunner
	// Pass adds checkers on the tested value.
	Pass(checkers ...check.ValueChecker) ValueRunner
}

// TableRunner provides methods to perform tests on a given func
// using a slice of Case.
type TableRunner interface {
	Runner
	// DryRun returns a TableResulter to access test results
	// without running it.
	DryRun() TableResulter
	// Cases adds test cases to be run on the tested func.
	Cases(cases []Case) TableRunner
}

// HTTPHandlerRunner provides methods to perform tests on a http.Handler
// or http.HandlerFunc.
type HTTPHandlerRunner interface {
	Runner
	// DryRun returns a HandlerResulter to access test results
	// without running it.
	DryRun() HandlerResulter
	// ResponseHeader adds checkers on the response header.
	ResponseHeader(...check.HTTPHeaderChecker) HTTPHandlerRunner
	// ResponseHeader adds checkers on the response status.
	ResponseStatus(...check.StringChecker) HTTPHandlerRunner
	// ResponseHeader adds checkers on the response code.
	ResponseCode(...check.IntChecker) HTTPHandlerRunner
	// ResponseHeader adds checkers on the response body.
	ResponseBody(...check.BytesChecker) HTTPHandlerRunner
	// ResponseHeader adds checkers on the handling duration.
	Duration(...check.DurationChecker) HTTPHandlerRunner
}

/*
	Results interfaces
*/

// Resulter provides methods to read test results after a dry run.
type Resulter interface {
	// Check returns a slice of CheckResults listing the runned checks
	Checks() []CheckResult
	// Passed returns true if all checks passed.
	Passed() bool
	// Failed returns true if one check or more failed.
	Failed() bool
	// NChecks returns the number of checks.
	NChecks() int
	// NPassed returns the number of checks that passed.
	NPassed() int
	// NPassed returns the number of checks that failed.
	NFailed() int
}

// HandlerResulter provides methods to read HandlerRunner results
// after a dry run.
type HandlerResulter interface {
	Resulter
	// ResponseHeader returns the gotten response header.
	ResponseHeader() http.Header
	// ResponseStatus returns the gotten response status.
	ResponseStatus() string
	// ResponseCode returns the gotten response code.
	ResponseCode() int
	// ResponseBody returns the gotten response body.
	ResponseBody() []byte
	// ResponseDuration returns the handler's execution time.
	ResponseDuration() time.Duration
}

// TableResulter provides methods to read TableRunner results
// after a dry run.
type TableResulter interface {
	Resulter
	// PassedAt returns true if the ith test case passed.
	PassedAt(index int) bool
	// PassedAt returns true if the ith test case failed.
	FailedAt(index int) bool
	// PassedAt returns true if the test case with matching label passed.
	PassedLabel(label string) bool
	// PassedAt returns true if the test case with matching label failed.
	FailedLabel(label string) bool
}

// CheckResult is a single check result after a dry run.
type CheckResult struct {
	// Passed is true if the current check passed
	Passed bool
	// Reason is the string output of a failed test as returned by a
	// check.Explainer, typically in format "expect X, got Y".
	Reason string
	label  string
}

func (cr CheckResult) String() string {
	return fmt.Sprintf("{%s%s}", cond.String("passed", "failed ", cr.Passed), cr.Reason)
}

/*
	Runners
*/

// Value returns a ValueRunner to run tests on a single value.
func Value(v interface{}) ValueRunner {
	return newValueRunner(v)
}

// HTTPHandler returns a HandlerRunner to run tests on a http.HTTPHandler
// response to given request.
func HTTPHandler(h http.Handler, r *http.Request) HTTPHandlerRunner {
	return newHandlerRunner(h.ServeHTTP, r)
}

// HTTPHandlerFunc returns a HandlerRunner to run tests on a http.HTTPHandlerFunc
// response to a given request.
func HTTPHandlerFunc(hf http.HandlerFunc, r *http.Request) HTTPHandlerRunner {
	return newHandlerRunner(hf, r)
}

// Table returns a TableRunner to run test cases on a func. By default,
// it works with funcs having a single input and output value. However,
// with an appropriate config it is compatible with any func signature.
func Table(testedFunc interface{}, cfg *TableConfig) TableRunner {
	return newTableRunner(testedFunc, cfg)
}
