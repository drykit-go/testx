package providers_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/internal/providers"
)

func TestBoolCheckerProvider(t *testing.T) {
	checkBool := providers.BoolCheckerProvider{}
	const b = true

	t.Run("Is pass", func(t *testing.T) {
		c := checkBool.Is(b)
		fmt.Print(c)
		assertPassChecker(t, "Bool.Is", c, b)
		c = checkBool.Is(!b)
		assertPassChecker(t, "Bool.Is", c, !b)
	})

	t.Run("Is fail", func(t *testing.T) {
		c := checkBool.Is(!b)
		assertFailChecker(t, "Bool.Is", c, b, makeExpl("false", "true"))
		c = checkBool.Is(b)
		assertFailChecker(t, "Bool.Is", c, !b, makeExpl("true", "false"))
	})
}
