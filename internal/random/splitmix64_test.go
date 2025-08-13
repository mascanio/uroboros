//go:build randomtest

package random

import "testing"

func BenchmarkSplitMix64_NextN_Pow2(b *testing.B) {
	rng := NewSplitMix64(123)
	for b.Loop() {
		_ = rng.NextN(1024)
	}
}

func BenchmarkSplitMix64_NextN_NonPow2(b *testing.B) {
	rng := NewSplitMix64(123)
	for b.Loop() {
		_ = rng.NextN(1000)
	}
}

func BenchmarkSplitMix64_NextBetween_Pow2(b *testing.B) {
	rng := NewSplitMix64(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1124)
	}
}

func BenchmarkSplitMix64_NextBetween_NonPow2(b *testing.B) {
	rng := NewSplitMix64(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1100)
	}
}
