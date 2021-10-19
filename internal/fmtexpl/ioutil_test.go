package fmtexpl_test

import (
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/internal/fmtexpl"
)

func TestDefault(t *testing.T) {
	exp := "int value:\nexp 42\ngot -1"
	got := fmtexpl.Default("int value", 42, -1)
	if got != exp {
		t.Errorf("\nexp %s\ngot %s", exp, got)
	}
}

func TestPretty(t *testing.T) {
	exp := "\u274c int value:\nexp 42\ngot -1"
	got := fmtexpl.Pretty("int value", 42, -1)
	if got != exp {
		t.Errorf("\nexp %s\ngot %s", exp, got)
	}
}

func TestChecker(t *testing.T) {
	exp := "int value:\nexp to pass IntChecker\ngot explanation: Oh hi Mark!"
	got := fmtexpl.Checker("int value", "to pass IntChecker", "Oh hi Mark!")
	if got != exp {
		t.Errorf("\nexp %s\ngot %s", exp, got)
	}
}

func TestTableCaseLabel(t *testing.T) {
	t.Run("with label input", func(t *testing.T) {
		exp := `Table.Cases[3] "division by 0" divide(42, 0)`
		got := fmtexpl.TableCaseLabel(
			"divide",
			3,
			"division by 0",
			testx.Args{42, 0},
		)
		if got != exp {
			t.Errorf("\nexp %s\ngot %s", exp, got)
		}
	})

	t.Run("no label input", func(t *testing.T) {
		exp := `Table.Cases[3] divide(42, 0)`
		got := fmtexpl.TableCaseLabel("divide", 3, "", testx.Args{42, 0})
		if got != exp {
			t.Errorf("\nexp %s\ngot %s", exp, got)
		}
	})
}
