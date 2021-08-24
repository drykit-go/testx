package checkconv_test

import (
	"fmt"
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

func TestRetrieve(t *testing.T) {
	t.Run("native checker", func(t *testing.T) {
		cases := []checkerTestcase{
			{
				checker: check.Bytes.Equal([]byte{42}),
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
				checker: check.Untyped.Custom("", func(got interface{}) bool { return true }),
				in:      "",
				expPass: true,
				expExpl: "",
			},
		}

		for _, c := range cases {
			assertRetrieved(t, c)
		}
	})

	t.Run("custom checker", func(t *testing.T) {
		cases := []checkerTestcase{
			{
				checker: check.NewIntCheck(isEven, isEvenExpl),
				in:      0,
				expPass: true,
				expExpl: "",
			},
			{
				checker: check.NewIntCheck(isEven, isEvenExpl),
				in:      1,
				expPass: false,
				expExpl: "expect value to be even, got 1",
			},
		}

		for _, c := range cases {
			assertRetrieved(t, c)
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
			assertRetrieved(t, c)
		}
	})

	t.Run("invalid checker", func(t *testing.T) {
		for _, c := range badCheckers {
			assertNotRetrieved(t, c)
		}
	})
}

func assertRetrieved(t *testing.T, tc checkerTestcase) {
	c, ok := checkconv.Retrieve(tc.checker)
	if !ok {
		t.Errorf("failed to retrieve checker: %#v", tc.checker)
	}
	if pass := c.Pass(tc.in); pass != tc.expPass {
		t.Errorf(
			"unexpected Pass return value with checker %#v: exp %v, got %v",
			tc.checker, tc.expPass, pass,
		)
	}
	if expl := c.Explain("value", tc.in); tc.expExpl != "" && expl != tc.expExpl {
		t.Errorf(
			"unexpected Explain return value with checker %#v: exp %v, got %v",
			tc.checker, tc.expPass, expl,
		)
	}
}

func assertNotRetrieved(t *testing.T, badChecker interface{}) {
	got, ok := checkconv.Retrieve(badChecker)
	if ok {
		t.Errorf("returned ok from bad input: %#v", badChecker)
	}
	if got != nil {
		t.Errorf("returned a non-nil checker from bad input: %#v", badChecker)
	}
}

// isEven is a dummy func for custom checkers
func isEven(n int) bool { return n&1 == 0 }

func isEvenExpl(_ string, got interface{}) string {
	return fmt.Sprintf("expect value to be even, got %v", got)
}
