package checkconv_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

func BenchmarkAll(b *testing.B) {
	b.Run("IsChecker", BenchmarkIsChecker)
	b.Run("Cast", BenchmarkCast)
	b.Run("Assert", BenchmarkAssert)
	b.Run("From", BenchmarkFrom)
}

func BenchmarkIsChecker(b *testing.B) {
	c := validCheckerInt{}
	for i := 0; i < b.N; i++ {
		checkconv.IsChecker(c)
	}
}

func BenchmarkAssert(b *testing.B) {
	b.Run("check.BytesChecker first_case", func(b *testing.B) {
		c := check.Bytes.Is([]byte{0})
		for i := 0; i < b.N; i++ {
			checkconv.Assert(c)
		}
	})
	b.Run("check.IntChecker midway_case", func(b *testing.B) {
		c := check.Int.Is(0)
		for i := 0; i < b.N; i++ {
			checkconv.Assert(c)
		}
	})
	b.Run("check.HTTPHeaderChecker last_case", func(b *testing.B) {
		c := check.HTTPHeader.KeySet("")
		for i := 0; i < b.N; i++ {
			checkconv.Assert(c)
		}
	})
}

func BenchmarkCast(b *testing.B) {
	c := validCheckerInt{}
	for i := 0; i < b.N; i++ {
		checkconv.Cast(c)
	}
}

func BenchmarkFrom(b *testing.B) {
	c := validCheckerInt{}
	for i := 0; i < b.N; i++ {
		checkconv.FromInt(c)
	}
}
