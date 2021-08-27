package testx

import (
	"net/http"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
)

type (

	// Runner provides a method Run to perform various tests.
	Runner interface {
		// Run runs a test and fails it if a check does not pass.
		Run(t *testing.T)
	}

	// ValueRunner provides methods to perform tests on a single value.
	ValueRunner interface {
		Runner
		// DryRun returns a Resulter to access test results
		// without running it.
		DryRun() Resulter
		// Exp adds an equality check on the tested value.
		Exp(value interface{}) ValueRunner
		// ExpNot adds inequality checks on the tested value.
		ExpNot(values ...interface{}) ValueRunner
		// Pass adds checkers on the tested value.
		Pass(checkers ...interface{}) ValueRunner
	}

	// TableRunner provides methods to perform tests on a given func
	// using a slice of Case.
	TableRunner interface {
		Runner
		// DryRun returns a TableResulter to access test results
		// without running it.
		DryRun() TableResulter
		// Cases adds test cases to be run on the tested func.
		Cases(cases []Case) TableRunner
	}

	// HandlerRunner provides methods to perform tests on a http.Handler
	// or http.HandlerFunc.
	HandlerRunner interface {
		Runner
		// DryRun returns a HandlerResulter to access test results
		// without running it.
		DryRun() HandlerResulter
		// ResponseHeader adds checkers on the response header.
		ResponseHeader(...check.HTTPHeaderChecker) HandlerRunner
		// ResponseHeader adds checkers on the response status.
		ResponseStatus(...check.StringChecker) HandlerRunner
		// ResponseHeader adds checkers on the response code.
		ResponseCode(...check.IntChecker) HandlerRunner
		// ResponseHeader adds checkers on the response body.
		ResponseBody(...check.BytesChecker) HandlerRunner
		// ResponseHeader adds checkers on the handling duration.
		Duration(...check.DurationChecker) HandlerRunner
	}
)

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
	// ExecTime is not yet implemented and may be deleted in the future.
	// Do not rely on it.
	// ExecTime returns the execution time of the whole test.
	ExecTime() time.Duration
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
