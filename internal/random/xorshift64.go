package random

// XorShift64 is a fast, simple PRNG for non-crypto use.
type XorShift64 struct {
	state uint64
}

func NewXorShift64(seed uint64) *XorShift64 {
	if seed == 0 {
		seed = 88172645463325252
	}
	return &XorShift64{state: seed}
}

func (x *XorShift64) Next() uint64 {
	x.state ^= x.state << 13
	x.state ^= x.state >> 7
	x.state ^= x.state << 17
	return x.state
}

func (x *XorShift64) NextN(max uint64) uint64 {
	if max == 0 {
		return 0
	}
	if (max & (max - 1)) == 0 {
		return x.Next() & (max - 1)
	}
	return x.Next() % max
}

func (x *XorShift64) NextBetween(min, max uint64) uint64 {
	rangeSize := max - min
	if max <= min || rangeSize == 0 {
		return min
	}
	if (rangeSize & (rangeSize - 1)) == 0 {
		return min + (x.Next() & (rangeSize - 1))
	}
	return min + x.NextN(rangeSize)
}
