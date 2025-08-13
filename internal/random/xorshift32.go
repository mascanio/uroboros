package random

// XorShift32 is a fast, simple PRNG for non-crypto use.
type XorShift32 struct {
	state uint32
}

func NewXorShift32(seed uint32) *XorShift32 {
	if seed == 0 {
		seed = 2463534242
	}
	return &XorShift32{state: seed}
}

func (x *XorShift32) Next() uint32 {
	x.state ^= x.state << 13
	x.state ^= x.state >> 17
	x.state ^= x.state << 5
	return x.state
}

func (x *XorShift32) NextN(max uint32) uint32 {
	if max == 0 {
		return 0
	}
	if (max & (max - 1)) == 0 {
		return x.Next() & (max - 1)
	}
	return x.Next() % max
}

func (x *XorShift32) NextBetween(min, max uint32) uint32 {
	rangeSize := max - min
	if max <= min || rangeSize == 0 {
		return min
	}
	if (rangeSize & (rangeSize - 1)) == 0 {
		return min + (x.Next() & (rangeSize - 1))
	}
	return min + x.NextN(rangeSize)
}
