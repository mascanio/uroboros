package random

// PCG32 is a modern, high-quality, fast PRNG.
type PCG32 struct {
	state uint64
	inc   uint64
}

func NewPCG32(seed, seq uint64) *PCG32 {
	if seed == 0 {
		seed = 0x853c49e6748fea9b
	}
	if seq == 0 {
		seq = 0xda3e39cb94b95bdb
	}
	return &PCG32{state: seed, inc: (seq << 1) | 1}
}

func (p *PCG32) Next() uint32 {
	oldstate := p.state
	p.state = oldstate*6364136223846793005 + p.inc
	xorshifted := uint32(((oldstate >> 18) ^ oldstate) >> 27)
	rot := uint32(oldstate >> 59)
	return (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
}

func (p *PCG32) NextN(max uint32) uint32 {
	if max == 0 {
		return 0
	}
	if (max & (max - 1)) == 0 {
		return p.Next() & (max - 1)
	}
	return p.Next() % max
}

func (p *PCG32) NextBetween(min, max uint32) uint32 {
	rangeSize := max - min
	if max <= min || rangeSize == 0 {
		return min
	}
	if (rangeSize & (rangeSize - 1)) == 0 {
		return min + (p.Next() & (rangeSize - 1))
	}
	return min + p.NextN(rangeSize)
}
