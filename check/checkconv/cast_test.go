package checkconv_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

func TestUntypedChecker(t *testing.T) {
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
			c := checkconv.Cast(tc.checker)
			assertValidUntypedChecker(t, c, tc)
		}
	})

	t.Run("unknown checker type", func(t *testing.T) {
		defer assertPanic(t, "attempt to convert unknown checker type")
		checkconv.Cast(validCheckerFloat64{})
	})

	t.Run("invalid checkers", func(t *testing.T) {
		for _, badChecker := range badCheckers {
			func() {
				defer assertPanic(t, "attempt to convert unknown checker type")
				checkconv.Cast(badChecker)
			}()
		}
	})
}

func assertPanic(t *testing.T, expMessage string) {
	r := recover()
	if r == nil {
		t.Errorf("expected to panic but did not")
	} else if r != expMessage {
		t.Errorf("bad panic message:\nexp %s\ngot %s", expMessage, r)
	}
}
