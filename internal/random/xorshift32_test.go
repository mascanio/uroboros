//go:build randomtest

package random

import "testing"

func BenchmarkNextN32_Pow2(b *testing.B) {
	rng := NewXorShift32(123)
	for b.Loop() {
		_ = rng.NextN(1024)
	}
}

func BenchmarkNextN32_NonPow2(b *testing.B) {
	rng := NewXorShift32(123)
	for b.Loop() {
		_ = rng.NextN(1000)
	}
}

func BenchmarkNextBetween32_Pow2(b *testing.B) {
	rng := NewXorShift32(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1124)
	}
}

func BenchmarkNextBetween32_NonPow2(b *testing.B) {
	rng := NewXorShift32(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1100)
	}
}
