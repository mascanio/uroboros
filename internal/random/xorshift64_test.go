//go:build randomtest

package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestN64_Next(t *testing.T) {
	rng := NewXorShift64(42)
	zero := false
	one := false
	for range 1_000_000 {
		n := rng.NextBetween(0, 2)
		if n == 0 {
			zero = true
		}
		if n == 1 {
			one = true
		}
		assert.GreaterOrEqual(t, n, uint64(0))
		assert.Less(t, n, uint64(2))
	}
	assert.True(t, zero)
	assert.True(t, one)
}

func FuzzN64_NextN(f *testing.F) {
	f.Add(uint64(42), uint64(3))
	f.Fuzz(func(t *testing.T, seed uint64, max uint64) {
		rng := NewXorShift64(seed)
		for range 1_000_000 {
			n := rng.NextN(max)
			assert.GreaterOrEqual(t, n, uint64(0))
			if max != 0 {
				assert.Less(t, n, uint64(max))
			}
		}
	})
}

func FuzzN64_NextBetween(f *testing.F) {
	f.Add(uint64(42), uint64(3), uint64(10))
	f.Fuzz(func(t *testing.T, seed uint64, min, max uint64) {
		rng := NewXorShift64(seed)
		for range 1_000_000 {
			n := rng.NextBetween(min, max)
			fmt.Println(min, max, n)
			if min >= max {
				require.Equal(t, uint64(min), n)
				return
			}
			require.GreaterOrEqual(t, n, uint64(min))
			require.Less(t, n, uint64(max))
		}
	})
}

func BenchmarkNextN64_Pow2(b *testing.B) {
	rng := NewXorShift64(123)
	for b.Loop() {
		_ = rng.NextN(1024)
	}
}

func BenchmarkNextN64_NonPow2(b *testing.B) {
	rng := NewXorShift64(123)
	for b.Loop() {
		_ = rng.NextN(1000)
	}
}

func BenchmarkNextBetween64_Pow2(b *testing.B) {
	rng := NewXorShift64(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1124)
	}
}

func BenchmarkNextBetween64_NonPow2(b *testing.B) {
	rng := NewXorShift64(123)
	for b.Loop() {
		_ = rng.NextBetween(100, 1100)
	}
}
