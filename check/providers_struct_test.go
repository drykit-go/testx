package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

type structTest struct {
	A, B, X, Y int
}

func TestStructCheckerProvider(t *testing.T) {
	const (
		vAB = 10
		vXY = 20
	)
	var (
		vABs = fmt.Sprint(vAB)
		vXYs = fmt.Sprint(vXY)
	)
	s := structTest{A: vAB, B: vAB, X: vXY, Y: vXY}

	t.Run("FieldsEqual pass", func(t *testing.T) {
		c := check.Struct.FieldsEqual(vAB, []string{"A", "B"})
		assertPassStructChecker(t, "FieldsEqual", c, s)
	})

	t.Run("FieldsEqual fail", func(t *testing.T) {
		c := check.Struct.FieldsEqual(vAB, []string{"A", "B", "X", "Y"})
		assertFailStructChecker(t, "FieldsEqual", c, s, makeExpl(
			"fields [.A, .B, .X, .Y] to equal "+vABs,
			".X="+vXYs+", .Y="+vXYs,
		))
	})

	t.Run("CheckFields pass", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.LT(vAB+1)),
			[]string{"A", "B"},
		)
		assertPassStructChecker(t, "CheckFields", c, s)
	})

	t.Run("CheckFields fail", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.LT(vAB+1)),
			[]string{"A", "B", "X", "Y"},
		)
		assertFailStructChecker(t, "CheckFields", c, s, makeExpl(
			"fields [.A, .B, .X, .Y] to pass ValueChecker",
			"explanation: fields:\nexp < 11\ngot .X=20, .Y=20",
		))
	})
}

// Helpers

func assertPassStructChecker(t *testing.T, method string, c check.ValueChecker, s structTest) {
	t.Helper()
	if !c.Pass(s) {
		failStructCheckerTest(t, true, method, s, c.Explain)
	}
}

func assertFailStructChecker(t *testing.T, method string, c check.ValueChecker, s structTest, expexpl string) {
	t.Helper()
	if c.Pass(s) {
		failStructCheckerTest(t, false, method, s, c.Explain)
	}
	if expexpl != "" {
		assertGoodExplain(t, c, s, expexpl)
	}
}

func failStructCheckerTest(t *testing.T, expPass bool, method string, s structTest, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Struct."+method, explain("struct value", s))
}
