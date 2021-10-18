package checkconv_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

func TestAssert(t *testing.T) {
	t.Run("known checker type", func(t *testing.T) {
		cases := []checkerTestcase{
			{
				checker: check.String.Contains("a"),
				in:      "aaa",
				expPass: true,
				expExpl: "",
			},
			{
				checker: check.NewIntChecker(isEven, isEvenExpl),
				in:      -1,
				expPass: false,
				expExpl: "expect value to be even, got -1",
			},
			{
				checker: validCheckerInt{},
				in:      0,
				expPass: true,
				expExpl: "ok",
			},
		}

		for _, tc := range cases {
			c := checkconv.Assert(tc.checker)
			assertValidValueChecker(t, c, tc)
		}
	})

	t.Run("unknown checker type", func(t *testing.T) {
		defer assertPanic(t, "assert from unknown checker type")
		checkconv.Assert(validCheckerFloat32{})
	})

	t.Run("invalid checkers", func(t *testing.T) {
		for _, badChecker := range badCheckers {
			func() {
				defer assertPanic(t, "assert from unknown checker type")
				checkconv.Assert(badChecker)
			}()
		}
	})
}

func TestAssertMany(t *testing.T) {
	t.Run("all provided checkers", func(t *testing.T) {
		testcases := []struct {
			checker interface{}
			in      interface{}
		}{
			{checker: check.Bool.Is(true), in: true},
			{checker: check.Bytes.Is([]byte{'a'}), in: []byte{'a'}},
			{checker: check.String.Is("a"), in: "a"},
			{checker: check.Int.Is(1), in: 1},
			{checker: check.Float64.Is(1), in: 1.},
			{checker: check.Duration.Over(time.Millisecond), in: time.Second},
			{checker: check.Context.Done(true), in: ctxDone()},
			{checker: check.HTTPHeader.HasKey("a"), in: http.Header{"a": []string{"b"}}},
			{checker: check.HTTPRequest.ContentLength(check.Int.GT(1)), in: &http.Request{ContentLength: 2}},
			{checker: check.HTTPResponse.StatusCode(check.Int.GT(1)), in: &http.Response{StatusCode: 2}},
			{checker: check.Value.Is(1), in: 1},
			{checker: check.Map.HasKeys("a"), in: map[string]int{"a": 1}},
			{checker: check.Slice.HasValues("a"), in: []string{"a"}},
			{checker: check.Struct.NotZero(), in: struct{ n int }{1}},
		}

		providedCheckers := func() (checkers []interface{}) {
			for _, tc := range testcases {
				checkers = append(checkers, tc.checker)
			}
			return
		}()

		res := checkconv.AssertMany(providedCheckers...)

		if gotLen, expLen := len(res), len(providedCheckers); gotLen != expLen {
			t.Errorf(
				"failed to assert provided checkers: exp len %d, got %d",
				expLen, gotLen,
			)
		}

		for i, tc := range testcases {
			if !res[i].Pass(tc.in) {
				t.Errorf(
					"failed to assert provided checkers: exp test case %d to pass",
					i,
				)
			}
		}
	})

	t.Run("custom checkers known type", func(t *testing.T) {
		knownCheckers := []interface{}{
			check.Value.Custom("", func(_ interface{}) bool { return true }),
			check.NewIntChecker(isEven, validExplainFunc),
			validCheckerInt{},
			validCheckerInterface{},
		}
		res := checkconv.AssertMany(knownCheckers...)
		if len(res) != len(knownCheckers) {
			t.Error("failed to assert many known checkers")
		}
	})

	t.Run("custom checkers unknown type", func(t *testing.T) {
		defer assertPanic(t, "assert from unknown checker type")
		unknownCheckers := []interface{}{validCheckerFloat32{}}
		checkconv.AssertMany(unknownCheckers...)
	})
}

func assertPanic(t *testing.T, expMessage string) {
	t.Helper()
	r := recover()
	if r == nil {
		t.Errorf("expected to panic but did not")
	} else if r != expMessage {
		t.Errorf("bad panic message:\nexp %s\ngot %s", expMessage, r)
	}
}

func ctxDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
