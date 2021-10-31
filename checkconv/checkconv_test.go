package checkconv_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

// Common definitions used across test files.

type onlyPasser struct{}

func (onlyPasser) Pass(int) bool { return true }

type onlyExplainer struct{}

func (onlyExplainer) PassX(int) bool             { return true }
func (onlyExplainer) Explain(string, any) string { return "" }

type badPasser struct{}

func (badPasser) Pass(int) int               { return 0 }
func (badPasser) Explain(string, any) string { return "" }

type badExplainerIn struct{}

func (badExplainerIn) Pass(int) bool           { return true }
func (badExplainerIn) Explain(any, any) string { return "" }

type badExplainerOut struct{}

func (badExplainerOut) Pass(int) bool           { return true }
func (badExplainerOut) Explain(string, any) any { return "" }

type checkerAsFields struct {
	Pass    func(int) bool
	Explain func(string, any) string
}

type validCheckerInt struct{}

func (validCheckerInt) Pass(int) bool              { return true }
func (validCheckerInt) Explain(string, any) string { return "ok" }

type validCheckerFloat32 struct{}

func (validCheckerFloat32) Pass(float32) bool          { return true }
func (validCheckerFloat32) Explain(string, any) string { return "ok" }

type validCheckerInterface struct{}

func (validCheckerInterface) Pass(any) bool              { return true }
func (validCheckerInterface) Explain(string, any) string { return "ok" }

var badCheckers = []any{
	-1,
	"hi",
	errors.New(""),
	onlyPasser{},
	onlyExplainer{},
	badPasser{},
	badExplainerIn{},
	badExplainerOut{},
	checkerAsFields{
		Pass:    func(int) bool { return true },
		Explain: func(string, any) string { return "" },
	},
}

var goodCheckers = []any{
	validCheckerInt{},
	validCheckerFloat32{},
	validCheckerInterface{},
}

func validExplainFunc(_ string, _ any) string {
	return "ok"
}

// isEven is a dummy passFunc for custom checkers
func isEven(n int) bool { return n&1 == 0 }

// isEvenExpl is a dummy explainFunc for custom checkers
func isEvenExpl(_ string, got any) string {
	return fmt.Sprintf("expect value to be even, got %v", got)
}

// Test helpers

func assertValidValueChecker(t *testing.T, c check.ValueChecker, tc checkerTestcase) {
	t.Helper()
	if pass := c.Pass(tc.in); pass != tc.expPass {
		t.Errorf(
			"unexpected Pass return value with checker %v: exp %v, got %v",
			tc.checker, tc.expPass, pass,
		)
	}
	if expl := c.Explain("value", tc.in); tc.expExpl != "" && expl != tc.expExpl {
		t.Errorf(
			"unexpected Explain return value with checker %#v:\nexp:\n%v\n\ngot:\n%v",
			tc.checker, tc.expExpl, expl,
		)
	}
}
