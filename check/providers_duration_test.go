package check_test

import (
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
)

func TestDurationCheckerProvider(t *testing.T) {
	const (
		d    = 1 * time.Second
		less = d - time.Millisecond
		more = d + time.Millisecond
	)

	t.Run("Under pass", func(t *testing.T) {
		c := check.Duration.Under(more)
		assertPassDurationChecker(t, "Under", c, d)
	})

	t.Run("Under fail", func(t *testing.T) {
		c := check.Duration.Under(less)
		assertFailDurationChecker(t, "Under", c, d)
		c = check.Duration.Under(d)
		assertFailDurationChecker(t, "Under", c, d)
	})

	t.Run("Over pass", func(t *testing.T) {
		c := check.Duration.Over(less)
		assertPassDurationChecker(t, "Over", c, d)
	})

	t.Run("Over fail", func(t *testing.T) {
		c := check.Duration.Over(more)
		assertFailDurationChecker(t, "Over", c, d)
		c = check.Duration.Over(d)
		assertFailDurationChecker(t, "Over", c, d)
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := check.Duration.InRange(less, more)
		assertPassDurationChecker(t, "InRange", c, d)

		c = check.Duration.InRange(d, d)
		assertPassDurationChecker(t, "InRange", c, d)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := check.Duration.InRange(more, more+time.Millisecond)
		assertFailDurationChecker(t, "InRange", c, d)

		c = check.Duration.InRange(more, less)
		assertFailDurationChecker(t, "InRange", c, d)
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := check.Duration.OutRange(more, more+time.Millisecond)
		assertPassDurationChecker(t, "OutRange", c, d)

		c = check.Duration.OutRange(more, less)
		assertPassDurationChecker(t, "OutRange", c, d)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := check.Duration.OutRange(less, more)
		assertFailDurationChecker(t, "OutRange", c, d)

		c = check.Duration.OutRange(d, d)
		assertFailDurationChecker(t, "OutRange", c, d)
	})
}

// Helpers

func assertPassDurationChecker(t *testing.T, method string, c check.DurationChecker, d time.Duration) {
	t.Helper()
	if !c.Pass(d) {
		failDurationCheckerTest(t, true, method, d, c.Explain)
	}
}

func assertFailDurationChecker(t *testing.T, method string, c check.DurationChecker, d time.Duration) {
	t.Helper()
	if c.Pass(d) {
		failDurationCheckerTest(t, false, method, d, c.Explain)
	}
}

func failDurationCheckerTest(t *testing.T, expPass bool, method string, d time.Duration, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Duration."+method, explain("Duration value", d))
}
