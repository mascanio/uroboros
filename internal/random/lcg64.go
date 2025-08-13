package random

// LCG64 is a simple linear congruential generator (not cryptographically secure).
type LCG64 struct {
	state uint64
}

func NewLCG64(seed uint64) *LCG64 {
	if seed == 0 {
		seed = 4101842887655102017
	}
	return &LCG64{state: seed}
}

func (l *LCG64) Next() uint64 {
	// Parameters from Numerical Recipes
	l.state = 6364136223846793005*l.state + 1
	return l.state
}

func (l *LCG64) NextN(max uint64) uint64 {
	if max == 0 {
		return 0
	}
	if (max & (max - 1)) == 0 {
		return l.Next() & (max - 1)
	}
	return l.Next() % max
}

func (l *LCG64) NextBetween(min, max uint64) uint64 {
	rangeSize := max - min
	if max <= min || rangeSize == 0 {
		return min
	}
	if (rangeSize & (rangeSize - 1)) == 0 {
		return min + (l.Next() & (rangeSize - 1))
	}
	return min + l.NextN(rangeSize)
}
