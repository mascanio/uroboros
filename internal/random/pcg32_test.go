//go:build randomtest

package random

import "testing"

func BenchmarkPCG32_NextN_Pow2(b *testing.B) {
	rng := NewPCG32(123, 456)
	for b.Loop() {
		_ = rng.NextN(1024)
	}
}

func BenchmarkPCG32_NextN_NonPow2(b *testing.B) {
	rng := NewPCG32(123, 456)
	for b.Loop() {
		_ = rng.NextN(1000)
	}
}

func BenchmarkPCG32_NextBetween_Pow2(b *testing.B) {
	rng := NewPCG32(123, 456)
	for b.Loop() {
		_ = rng.NextBetween(100, 1124)
	}
}

func BenchmarkPCG32_NextBetween_NonPow2(b *testing.B) {
	rng := NewPCG32(123, 456)
	for b.Loop() {
		_ = rng.NextBetween(100, 1100)
	}
}
