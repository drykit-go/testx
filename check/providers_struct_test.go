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
	s := structTest{A: vAB, B: vAB, X: vXY, Y: vXY}
	// FIXME: remove forced conversion
	itf := func(v structTest) any {
		return v
	}

	t.Run("FieldsEqual pass", func(t *testing.T) {
		c := check.Struct.FieldsEqual(vAB, []string{"A", "B"})
		assertPassChecker(t, "Struct.FieldsEqual", c, itf(s))
	})

	t.Run("FieldsEqual fail", func(t *testing.T) {
		c := check.Struct.FieldsEqual(vAB, []string{"A", "B", "X", "Y"})
		assertFailChecker(t, "Struct.FieldsEqual", c, itf(s), makeExpl(
			fmt.Sprintf("fields [.A, .B, .X, .Y] to equal %v", vAB),
			fmt.Sprintf(".X=%v, .Y=%v", vXY, vXY),
		))
	})

	t.Run("CheckFields pass", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.LT(vAB+1)),
			[]string{"A", "B"},
		)
		assertPassChecker(t, "Struct.CheckFields", c, itf(s))
	})

	t.Run("CheckFields fail", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.LT(vAB+1)),
			[]string{"A", "B", "X", "Y"},
		)
		assertFailChecker(t, "Struct.CheckFields", c, itf(s), makeExpl(
			"fields [.A, .B, .X, .Y] to pass Checker[any]",
			"explanation: fields:\n"+makeExpl("< 11", ".X=20, .Y=20"),
		))
	})
}
