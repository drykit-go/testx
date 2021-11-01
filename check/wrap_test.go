package check_test

import (
	"testing"

	"github.com/drykit-go/slicex"
	"github.com/drykit-go/testx/check"
)

type checkerTestcase[T any] struct {
	checker check.Checker[T]
	in      T
	exppass bool
	expexpl string
}

func TestWrap(t *testing.T) {
	t.Run("native checkers", func(t *testing.T) {
		testcases := []checkerTestcase[any]{
			{
				checker: check.Wrap(check.Bytes.Is([]byte{42})),
				in:      []byte{42},
				exppass: true,
				expexpl: "",
			},
			{
				checker: check.Wrap(check.Int.InRange(41, 43)),
				in:      -1,
				exppass: false,
				expexpl: "value:\nexp in range [41:43]\ngot -1",
			},
			{
				checker: check.Wrap(check.Value.Custom("", func(got any) bool { return true })),
				in:      "",
				exppass: true,
				expexpl: "",
			},
		}

		for _, tc := range testcases {
			assertValidAnyChecker(t, tc)
		}
	})

	t.Run("new checkers", func(t *testing.T) {
		testcases := []checkerTestcase[any]{
			{
				checker: check.Wrap(newComplex128Checker(true, "")),
				in:      1i + 1,
				exppass: true,
				expexpl: "",
			},
			{
				checker: check.Wrap(newComplex128Checker(false, "bad")),
				in:      1i + 1,
				exppass: false,
				expexpl: "bad",
			},
		}

		for _, tc := range testcases {
			assertValidAnyChecker(t, tc)
		}
	})

	t.Run("custom checkers", func(t *testing.T) {
		testcases := []checkerTestcase[any]{
			{
				checker: check.Wrap(customComplex128Checker(true, "")),
				in:      1i + 1,
				exppass: true,
				expexpl: "",
			},
			{
				checker: check.Wrap(customComplex128Checker(false, "bad")),
				in:      1i + 1,
				exppass: false,
				expexpl: "bad",
			},
		}

		for _, tc := range testcases {
			assertValidAnyChecker(t, tc)
		}
	})
}

func TestWrapMany(t *testing.T) {
	testcases := []checkerTestcase[complex128]{
		{
			checker: newComplex128Checker(true, ""),
			in:      1i + 1,
			exppass: true,
			expexpl: "",
		},
		{
			checker: newComplex128Checker(false, "bad"),
			in:      1i + 1,
			exppass: false,
			expexpl: "bad",
		},
		{
			checker: customComplex128Checker(true, ""),
			in:      1i + 1,
			exppass: true,
			expexpl: "",
		},
		{
			checker: customComplex128Checker(false, "bad"),
			in:      1i + 1,
			exppass: false,
			expexpl: "bad",
		},
	}
	wrappedCheckers := check.WrapMany(slicex.Map(testcases,
		func(v checkerTestcase[complex128]) check.Checker[complex128] {
			return v.checker
		},
	)...)

	for i, tc := range testcases {
		assertValidAnyChecker(t, checkerTestcase[any]{
			checker: wrappedCheckers[i],
			in:      tc.in,
			exppass: tc.exppass,
			expexpl: tc.expexpl,
		})
	}
}

func assertValidAnyChecker(t *testing.T, tc checkerTestcase[any]) {
	t.Helper()
	c := tc.checker
	if pass := c.Pass(tc.in); pass != tc.exppass {
		t.Errorf(
			"unexpected Pass return value with checker %v: exp %v, got %v",
			tc.checker, tc.exppass, pass,
		)
	}
	if expl := c.Explain("value", tc.in); tc.expexpl != "" && expl != tc.expexpl {
		t.Errorf(
			"unexpected Explain return value with checker %#v:\nexp:\n%v\n\ngot:\n%v",
			tc.checker, tc.expexpl, expl,
		)
	}
}

// Helpers

type customComplex128CheckerImpl struct {
	exppass bool
	expexpl string
}

func (c customComplex128CheckerImpl) Pass(complex128) bool               { return c.exppass }
func (c customComplex128CheckerImpl) Explain(string, interface{}) string { return c.expexpl }

func newComplex128Checker(exppass bool, expexpl string) check.Checker[complex128] {
	return check.NewChecker(
		func(complex128) bool { return exppass },
		func(string, interface{}) string { return expexpl },
	)
}

func customComplex128Checker(exppass bool, expexpl string) check.Checker[complex128] {
	return customComplex128CheckerImpl{exppass: exppass, expexpl: expexpl}
}
