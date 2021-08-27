package checkconv_test

import (
	"testing"

	"github.com/drykit-go/testx/check/checkconv"
)

func BenchmarkAll(b *testing.B) {
	b.Run("IsChecker", BenchmarkIsChecker)
	b.Run("Cast", BenchmarkCast)
	b.Run("Assert", BenchmarkIsChecker)
	b.Run("FromInt", BenchmarkFromInt)
}

func BenchmarkIsChecker(b *testing.B) {
	c := validCheckerInt{}
	for i := 0; i < b.N; i++ {
		checkconv.IsChecker(c)
	}
}

func BenchmarkCast(b *testing.B) {
	c := validCheckerInt{}
	for i := 0; i < b.N; i++ {
		checkconv.Cast(c)
	}
}

func BenchmarkAssert(b *testing.B) {
	c := validCheckerInt{}
	for i := 0; i < b.N; i++ {
		checkconv.IsChecker(c)
	}
}

func BenchmarkFromInt(b *testing.B) {
	c := validCheckerInt{}
	for i := 0; i < b.N; i++ {
		checkconv.FromInt(c)
	}
}
