package testx_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

var deq = reflect.DeepEqual

//nolint: structcheck // false-positive
type valueRunnerResultCheck struct {
	passed bool
	reason string
}

type valueRunnerResults struct {
	baseResults
	checks []valueRunnerResultCheck
}

func TestValueRunner(t *testing.T) {
	t.Run("should pass", func(t *testing.T) {
		res := testx.Value(42).
			MustBe(42).
			MustNotBe(3, "hello").
			MustPass(check.Int.InRange(41, 43)).
			DryRun()

		exp := valueRunnerResults{
			baseResults: baseResults{
				Passed:  true,
				Failed:  false,
				NPassed: 4,
				NFailed: 0,
				NChecks: 4,
			},
			checks: []valueRunnerResultCheck{
				{passed: true, reason: ""},
				{passed: true, reason: ""},
				{passed: true, reason: ""},
				{passed: true, reason: ""},
			},
		}

		checkValueRunnerResults(t, res, exp)
	})

	t.Run("should fail", func(t *testing.T) {
		res := testx.Value(42).
			MustBe(99).
			MustNotBe(99, 42).
			MustPass(check.Int.LesserThan(10)).
			DryRun()

		exp := valueRunnerResults{
			baseResults: baseResults{
				Passed:  false,
				Failed:  true,
				NPassed: 1,
				NFailed: 3,
				NChecks: 4,
			},
			checks: []valueRunnerResultCheck{
				{passed: false, reason: "value: expect 99, got 42"},
				{passed: true, reason: ""},
				{passed: false, reason: "value: expect not 42, got 42"},
				{passed: false, reason: "expect value < 10, got 42"},
			},
		}

		checkValueRunnerResults(t, res, exp)
	})
}

func checkValueRunnerResults(t *testing.T, res testx.ValueRunnerResults, exp valueRunnerResults) {
	if got := toValueRunnerResults(res); !deq(got, exp) {
		t.Errorf("bad results\nexp %#v\ngot %#v", exp, got)
	}
}

func toValueRunnerResults(res testx.ValueRunnerResults) valueRunnerResults {
	return *(*valueRunnerResults)(unsafe.Pointer(&res))
}
