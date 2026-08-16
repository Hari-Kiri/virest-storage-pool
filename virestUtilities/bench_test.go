package virestUtilities

import (
	"errors"
	"testing"
)

func BenchmarkWrap(b *testing.B) {
	err := errors.New("x")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Wrap("op", err)
	}
}
