package check_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
)

func TestDurationCheckerProvider(t *testing.T) {
	const (
		d        = 1 * time.Second
		less     = d - time.Millisecond
		more     = d + time.Millisecond
		moremore = more + time.Millisecond
	)

	ms := func(d time.Duration) int64 {
		return d.Milliseconds()
	}

	t.Run("Under pass", func(t *testing.T) {
		c := check.Duration.Under(more)
		assertPassChecker(t, "Duration.Under", c, d)
	})

	t.Run("Under fail", func(t *testing.T) {
		c := check.Duration.Under(less)
		assertFailChecker(t, "Duration.Under", c, d, makeExpl(
			fmt.Sprintf("under %dms", ms(less)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = check.Duration.Under(d)
		assertFailChecker(t, "Duration.Under", c, d, makeExpl(
			fmt.Sprintf("under %dms", ms(d)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})

	t.Run("Over pass", func(t *testing.T) {
		c := check.Duration.Over(less)
		assertPassChecker(t, "Duration.Over", c, d)
	})

	t.Run("Over fail", func(t *testing.T) {
		c := check.Duration.Over(more)
		assertFailChecker(t, "Duration.Over", c, d, makeExpl(
			fmt.Sprintf("over %dms", ms(more)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = check.Duration.Over(d)
		assertFailChecker(t, "Duration.Over", c, d, makeExpl(
			fmt.Sprintf("over %dms", ms(d)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := check.Duration.InRange(less, more)
		assertPassChecker(t, "Duration.InRange", c, d)

		c = check.Duration.InRange(d, d)
		assertPassChecker(t, "Duration.InRange", c, d)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := check.Duration.InRange(more, moremore)
		assertFailChecker(t, "Duration.InRange", c, d, makeExpl(
			fmt.Sprintf("in range [%dms:%dms]", ms(more), ms(moremore)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = check.Duration.InRange(more, less)
		assertFailChecker(t, "Duration.InRange", c, d, makeExpl(
			fmt.Sprintf("in range [%dms:%dms]", ms(more), ms(less)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := check.Duration.OutRange(more, moremore)
		assertPassChecker(t, "Duration.OutRange", c, d)

		c = check.Duration.OutRange(more, less)
		assertPassChecker(t, "Duration.OutRange", c, d)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := check.Duration.OutRange(less, more)
		assertFailChecker(t, "Duration.OutRange", c, d, makeExpl(
			fmt.Sprintf("not in range [%dms:%dms]", ms(less), ms(more)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = check.Duration.OutRange(d, d)
		assertFailChecker(t, "Duration.OutRange", c, d, makeExpl(
			fmt.Sprintf("not in range [%dms:%dms]", ms(d), ms(d)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})
}
