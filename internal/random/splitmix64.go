package random

// SplitMix64 is a fast, simple PRNG often used for seeding other PRNGs.
type SplitMix64 struct {
	state uint64
}

func NewSplitMix64(seed uint64) *SplitMix64 {
	if seed == 0 {
		seed = 0x106689d45497fdb5
	}
	return &SplitMix64{state: seed}
}

func (s *SplitMix64) Next() uint64 {
	s.state += 0x9e3779b97f4a7c15
	z := s.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func (s *SplitMix64) NextN(max uint64) uint64 {
	if max == 0 {
		return 0
	}
	if (max & (max - 1)) == 0 {
		return s.Next() & (max - 1)
	}
	return s.Next() % max
}

func (s *SplitMix64) NextBetween(min, max uint64) uint64 {
	rangeSize := max - min
	if max <= min || rangeSize == 0 {
		return min
	}
	if (rangeSize & (rangeSize - 1)) == 0 {
		return min + (s.Next() & (rangeSize - 1))
	}
	return min + s.NextN(rangeSize)
}
