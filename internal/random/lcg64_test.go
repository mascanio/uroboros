//go:build randomtest

package random

import "testing"

func BenchmarkLCG64_NextN_Pow2(b *testing.B) {
	rng := NewLCG64(123)
	for b.Loop() {
		_ = rng.NextN(1024)
	}
}

func BenchmarkLCG64_NextN_NonPow2(b *testing.B) {
	rng := NewLCG64(123)
	for b.Loop() {
		_ = rng.NextN(1000)
	}
}

func BenchmarkLCG64_NextBetween_Pow2(b *testing.B) {
	rng := NewLCG64(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1124)
	}
}

func BenchmarkLCG64_NextBetween_NonPow2(b *testing.B) {
	rng := NewLCG64(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1100)
	}
}
