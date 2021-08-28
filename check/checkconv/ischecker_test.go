package checkconv_test

import (
	"testing"

	"github.com/drykit-go/testx/check/checkconv"
)

func TestIsChecker(t *testing.T) {
	t.Run("invalid checkers", func(t *testing.T) {
		values := append([]interface{}{
			"a string",
			42,
			func(int) bool { return true },
		}, badCheckers...)

		for _, v := range values {
			if checkconv.IsChecker(v) {
				t.Errorf("value %v was wrongly considered a checker", v)
			}
		}
	})

	t.Run("valid checkers", func(t *testing.T) {
		for _, v := range goodCheckers {
			if !checkconv.IsChecker(v) {
				t.Errorf("checker %v was wrongly considered not a checker", v)
			}
		}
	})
}
