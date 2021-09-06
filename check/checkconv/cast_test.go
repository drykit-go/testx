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
				expExpl: "value:\nexp in range [41:43]\ngot -1",
			},
			{
				checker: check.Value.Custom("", func(got interface{}) bool { return true }),
				in:      "",
				expPass: true,
				expExpl: "",
			},
		}

		for _, c := range cases {
			assertCastable(t, c)
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
			assertCastable(t, c)
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
			assertCastable(t, c)
		}
	})

	t.Run("invalid checker", func(t *testing.T) {
		for _, c := range badCheckers {
			assertNotCastable(t, c)
		}
	})
}

func TestCastMany(t *testing.T) {
	t.Run("valid checkers", func(t *testing.T) {
		res, ok := checkconv.CastMany(goodCheckers...)
		if !ok {
			t.Error("returned !ok from good checkers")
		}
		if len(res) != len(goodCheckers) {
			t.Error("failed to cast valid checkers")
		}
	})

	t.Run("invalid checkers", func(t *testing.T) {
		res, ok := checkconv.CastMany(badCheckers...)
		if ok {
			t.Error("returned ok from invalid checkers")
		}
		if len(res) != 0 {
			t.Error("casted invalid checkers")
		}
	})
}

func assertCastable(t *testing.T, tc checkerTestcase) {
	t.Helper()
	c, ok := checkconv.Cast(tc.checker)
	if !ok {
		t.Errorf("failed to cast checker: %#v", tc.checker)
	}
	assertValidValueChecker(t, c, tc)
}

func assertNotCastable(t *testing.T, badChecker interface{}) {
	t.Helper()
	got, ok := checkconv.Cast(badChecker)
	if ok {
		t.Errorf("returned ok from bad input: %#v", badChecker)
	}
	if got != nil {
		t.Errorf("returned a non-nil checker from bad input: %#v", badChecker)
	}
}
