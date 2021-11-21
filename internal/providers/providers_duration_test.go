package providers_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/drykit-go/testx/internal/providers"
)

func TestDurationCheckerProvider(t *testing.T) {
	checkDuration := providers.DurationCheckerProvider{}

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
		c := checkDuration.Under(more)
		assertPassChecker(t, "Duration.Under", c, d)
	})

	t.Run("Under fail", func(t *testing.T) {
		c := checkDuration.Under(less)
		assertFailChecker(t, "Duration.Under", c, d, makeExpl(
			fmt.Sprintf("under %dms", ms(less)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = checkDuration.Under(d)
		assertFailChecker(t, "Duration.Under", c, d, makeExpl(
			fmt.Sprintf("under %dms", ms(d)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})

	t.Run("Over pass", func(t *testing.T) {
		c := checkDuration.Over(less)
		assertPassChecker(t, "Duration.Over", c, d)
	})

	t.Run("Over fail", func(t *testing.T) {
		c := checkDuration.Over(more)
		assertFailChecker(t, "Duration.Over", c, d, makeExpl(
			fmt.Sprintf("over %dms", ms(more)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = checkDuration.Over(d)
		assertFailChecker(t, "Duration.Over", c, d, makeExpl(
			fmt.Sprintf("over %dms", ms(d)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})

	t.Run("InRange pass", func(t *testing.T) {
		c := checkDuration.InRange(less, more)
		assertPassChecker(t, "Duration.InRange", c, d)

		c = checkDuration.InRange(d, d)
		assertPassChecker(t, "Duration.InRange", c, d)
	})

	t.Run("InRange fail", func(t *testing.T) {
		c := checkDuration.InRange(more, moremore)
		assertFailChecker(t, "Duration.InRange", c, d, makeExpl(
			fmt.Sprintf("in range [%dms:%dms]", ms(more), ms(moremore)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = checkDuration.InRange(more, less)
		assertFailChecker(t, "Duration.InRange", c, d, makeExpl(
			fmt.Sprintf("in range [%dms:%dms]", ms(more), ms(less)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})

	t.Run("OutRange pass", func(t *testing.T) {
		c := checkDuration.OutRange(more, moremore)
		assertPassChecker(t, "Duration.OutRange", c, d)

		c = checkDuration.OutRange(more, less)
		assertPassChecker(t, "Duration.OutRange", c, d)
	})

	t.Run("OutRange fail", func(t *testing.T) {
		c := checkDuration.OutRange(less, more)
		assertFailChecker(t, "Duration.OutRange", c, d, makeExpl(
			fmt.Sprintf("not in range [%dms:%dms]", ms(less), ms(more)),
			fmt.Sprintf("%dms", ms(d)),
		))

		c = checkDuration.OutRange(d, d)
		assertFailChecker(t, "Duration.OutRange", c, d, makeExpl(
			fmt.Sprintf("not in range [%dms:%dms]", ms(d), ms(d)),
			fmt.Sprintf("%dms", ms(d)),
		))
	})
}
