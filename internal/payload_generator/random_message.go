package payloadgenerator

import (
	"strconv"

	random "github.com/mascanio/uroboros/internal/random"
	sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"
)

const (
	maxSeqDigits = 20 // max decimal digits for uint64
	alphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-"
)

// Unified random message generator

type RNG interface {
	NextN(max uint64) uint64
	NextBetween(min, max uint64) uint64
}

type RandomMessageGenerator struct {
	minLen       int
	maxLen       int
	seq          sequencegenerator.SequenceGenerator
	buf          []byte // preallocated for maxSeqDigits + 1 + maxLen/length
	rng          RNG
	includeSeqID bool
}

type RandomMessageGeneratorOption func(*RandomMessageGenerator)

func WithFixedLength(length int) RandomMessageGeneratorOption {
	return WithMinMaxLength(length, length+1)
}

func WithMinMaxLength(min, max int) RandomMessageGeneratorOption {
	return func(g *RandomMessageGenerator) {
		g.minLen = min
		g.maxLen = max
	}
}

func WithIncludeSeqInMsg(include bool) RandomMessageGeneratorOption {
	return func(g *RandomMessageGenerator) {
		g.includeSeqID = include
	}
}

func WithSequenceGenerator(seq sequencegenerator.SequenceGenerator) RandomMessageGeneratorOption {
	return func(g *RandomMessageGenerator) {
		g.seq = seq
	}
}

func WithRNG(rng RNG) RandomMessageGeneratorOption {
	return func(g *RandomMessageGenerator) {
		g.rng = rng
	}
}

func NewRandomMessageGenerator(opts ...RandomMessageGeneratorOption) *RandomMessageGenerator {
	g := &RandomMessageGenerator{
		minLen:       8,
		maxLen:       32,
		seq:          nil,
		buf:          nil,
		rng:          nil,
		includeSeqID: true,
	}
	for _, opt := range opts {
		opt(g)
	}
	if g.seq == nil {
		panic("RandomMessageGenerator: sequence generator is required. Use WithSequenceGenerator.")
	}
	// Sanity checks
	if g.maxLen < g.minLen {
		g.maxLen = g.minLen
	}
	// Buffer allocation
	bodyLen := g.maxLen
	if g.includeSeqID {
		g.buf = make([]byte, 0, maxSeqDigits+1+bodyLen)
	} else {
		g.buf = make([]byte, 0, bodyLen)
	}
	if g.rng == nil {
		g.rng = random.NewXorShift64(0)
	}
	return g
}

func (g *RandomMessageGenerator) GenerateMessage() []byte {
	i, done := g.seq.Next()
	if done == sequencegenerator.Done {
		return nil
	}
	length := g.minLen
	if g.maxLen-1 > g.minLen {
		length = int(g.rng.NextBetween(uint64(g.minLen), uint64(g.maxLen)))
	}
	g.buf = g.buf[:0]
	if g.includeSeqID {
		g.buf = strconv.AppendInt(g.buf, i, 10)
		g.buf = append(g.buf, ' ')
	}
	for range length {
		g.buf = append(g.buf, alphabet[g.rng.NextN(64)])
	}
	return g.buf
}
