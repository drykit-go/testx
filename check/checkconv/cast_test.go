package checkconv_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

type checkerTestcase struct {
	checker interface{}
	in      interface{}
	expPass bool
	expExpl string
}

func TestCast(t *testing.T) {
	t.Run("native checker", func(t *testing.T) {
		cases := []checkerTestcase{
			{
				checker: check.Bytes.Is([]byte{42}),
				in:      []byte{42},
				expPass: true,
				expExpl: "",
			},
			{
				checker: check.Int.InRange(41, 43),
				in:      -1,
				expPass: false,
				expExpl: "expect value in range [41:43], got -1",
			},
			{
				checker: check.Value.Custom("", func(got interface{}) bool { return true }),
				in:      "",
				expPass: true,
				expExpl: "",
			},
		}

		for _, c := range cases {
			assertCasted(t, c)
		}
	})

	t.Run("custom checker", func(t *testing.T) {
		cases := []checkerTestcase{
			{
				checker: check.NewIntChecker(isEven, isEvenExpl),
				in:      0,
				expPass: true,
				expExpl: "",
			},
			{
				checker: check.NewIntChecker(isEven, isEvenExpl),
				in:      1,
				expPass: false,
				expExpl: "expect value to be even, got 1",
			},
		}

		for _, c := range cases {
			assertCasted(t, c)
		}
	})

	t.Run("unknown checker", func(t *testing.T) {
		cases := []checkerTestcase{
			{
				checker: validCheckerInt{},
				in:      0,
				expPass: true,
				expExpl: "ok",
			},
			{
				checker: validCheckerInterface{},
				in:      "anything",
				expPass: true,
				expExpl: "ok",
			},
		}

		for _, c := range cases {
			assertCasted(t, c)
		}
	})

	t.Run("invalid checker", func(t *testing.T) {
		for _, c := range badCheckers {
			assertNotCasted(t, c)
		}
	})
}

func assertCasted(t *testing.T, tc checkerTestcase) {
	c, ok := checkconv.Cast(tc.checker)
	if !ok {
		t.Errorf("failed to retrieve checker: %#v", tc.checker)
	}
	assertValidValueChecker(t, c, tc)
}

func assertNotCasted(t *testing.T, badChecker interface{}) {
	got, ok := checkconv.Cast(badChecker)
	if ok {
		t.Errorf("returned ok from bad input: %#v", badChecker)
	}
	if got != nil {
		t.Errorf("returned a non-nil checker from bad input: %#v", badChecker)
	}
}
