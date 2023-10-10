package main

import (
	"testing"
)

func BenchmarkRun(b *testing.B) {
	b.ResetTimer()
	b.Run("2 active threads", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Run(2, "pleasure", "explain", "itself", "mistaken", "pain")
		}
	})
	b.Run("32 active threads", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Run(32, "pleasure", "explain", "itself", "mistaken", "pain")
		}
	})
}
