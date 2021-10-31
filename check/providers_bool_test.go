package check_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

func TestBoolCheckerProvider(t *testing.T) {
	const b = true

	t.Run("Is pass", func(t *testing.T) {
		c := check.Bool.Is(b)
		fmt.Print(c)
		assertPassChecker(t, "Bool.Is", c, b)
		c = check.Bool.Is(!b)
		assertPassChecker(t, "Bool.Is", c, !b)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := check.Bool.Is(!b)
		assertFailChecker(t, "Bool.Is", c, b, makeExpl("false", "true"))
		c = check.Bool.Is(b)
		assertFailChecker(t, "Bool.Is", c, !b, makeExpl("true", "false"))
	})
}
