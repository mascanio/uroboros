package payloadgenerator

import (
	"bytes"
	"testing"

	sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRNG implements the RNG interface for deterministic tests
// It cycles through a fixed sequence or always returns 0 for simplicity
// (can be extended for more complex patterns if needed)
type mockRNG struct {
	val      uint64
	inc      uint64
	sequence []uint64 // If set, use these values in order for NextBetween (for length)
	idx      int      // Current index in sequence
}

func (m *mockRNG) NextN(max uint64) uint64 {
	return m.NextBetween(0, max)
}

func (m *mockRNG) NextBetween(min, max uint64) uint64 {
	if len(m.sequence) > 0 && m.idx < len(m.sequence) {
		v := min + m.sequence[m.idx]
		m.idx++
		return v
	}
	if max <= min {
		return min
	}
	v := min + (m.val % (max - min))
	m.val += m.inc
	return v
}

type dummySeq struct {
	cur int64
	max int64
}

func (d *dummySeq) Next() (int64, sequencegenerator.IsDone) {
	if d.cur >= d.max {
		return 0, sequencegenerator.Done
	}
	v := d.cur
	d.cur++
	return v, sequencegenerator.NotDone
}

func TestRandomVariableLengthMessageGenerator(t *testing.T) {
	minLen, maxLen := 5, 10
	nLen := maxLen - minLen + 1
	rngSeq := make([]uint64, 0)
	for i := range nLen {
		length := minLen + i
		rngSeq = append(rngSeq, uint64(i))
		for range length {
			rngSeq = append(rngSeq, 0) // always pick 'a' for content
		}
	}
	seq := &dummySeq{max: int64(nLen)}
	gen := NewRandomMessageGenerator(
		WithMinMaxLength(minLen, maxLen),
		WithRNG(&mockRNG{sequence: rngSeq}),
		WithSequenceGenerator(seq),
		WithIncludeSeqInMsg(false),
	)
	for i := range nLen {
		msg := gen.GenerateMessage()
		require.NotNil(t, msg)
		assert.Equal(t, bytes.Repeat([]byte("a"), minLen+i), msg)
	}
	// Next call should return nil
	msgNil := gen.GenerateMessage()
	require.Nil(t, msgNil)
}

func TestRandomFixedLengthMessageGenerator(t *testing.T) {
	length := 8
	n := 5
	rngSeq := make([]uint64, 0)
	// For fixed length, no length RNG is consumed, only chars
	for range n {
		for range length {
			rngSeq = append(rngSeq, 0) // always pick 'a' for content
		}
	}
	seq := &dummySeq{max: int64(n)}
	gen := NewRandomMessageGenerator(
		WithFixedLength(length),
		WithRNG(&mockRNG{sequence: rngSeq}),
		WithSequenceGenerator(seq),
		WithIncludeSeqInMsg(false),
	)
	for range n {
		msg := gen.GenerateMessage()
		require.NotNil(t, msg)
		assert.Equal(t, bytes.Repeat([]byte("a"), length), msg)
	}
}

func TestRandomMessageGenerator_AllRepeatedCharacters(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-"
	for idx := 0; idx < len(alphabet); idx++ {
		seq := &dummySeq{max: 1}
		gen := NewRandomMessageGenerator(
			WithFixedLength(3),
			WithSequenceGenerator(seq),
			WithRNG(&mockRNG{sequence: []uint64{uint64(idx), uint64(idx), uint64(idx)}}),
			WithIncludeSeqInMsg(false),
		)
		msg := gen.GenerateMessage()
		require.NotNil(t, msg)
		expected := bytes.Repeat([]byte{alphabet[idx]}, 3)
		assert.Equal(t, expected, msg, "idx=%d char=%c", idx, alphabet[idx])
	}
}

func TestRandomMessageGenerator_AllAlphabetCharacters(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-"
	seq := &dummySeq{max: 1}
	indices := make([]uint64, len(alphabet))
	for i := range indices {
		indices[i] = uint64(i)
	}
	gen := NewRandomMessageGenerator(
		WithFixedLength(len(alphabet)),
		WithSequenceGenerator(seq),
		WithRNG(&mockRNG{sequence: indices}),
		WithIncludeSeqInMsg(false),
	)
	msg := gen.GenerateMessage()
	require.NotNil(t, msg)
	assert.Equal(t, []byte(alphabet), msg)
}

func TestRandomMessageGenerator_Options(t *testing.T) {
	tests := []struct {
		name        string
		options     []RandomMessageGeneratorOption
		seqMax      int64
		expectLen   int
		includeID   bool
		shouldPanic bool
	}{
		{
			name: "min/max length",
			options: []RandomMessageGeneratorOption{
				WithMinMaxLength(4, 7),
				WithRNG(&mockRNG{sequence: []uint64{2, 0, 0, 0, 0, 0}}), // length=6, all 'a'
				WithIncludeSeqInMsg(false),
			},
			seqMax:    1,
			expectLen: 6,
			includeID: false,
		},
		{
			name: "fixed length",
			options: []RandomMessageGeneratorOption{
				WithFixedLength(5),
				WithRNG(&mockRNG{sequence: []uint64{0, 0, 0, 0, 0}}),
				WithIncludeSeqInMsg(false),
			},
			seqMax:    1,
			expectLen: 5,
			includeID: false,
		},
		{
			name: "include sequence ID",
			options: []RandomMessageGeneratorOption{
				WithFixedLength(3),
				WithRNG(&mockRNG{sequence: []uint64{0, 0, 0}}),
				WithIncludeSeqInMsg(true),
			},
			seqMax:    1,
			expectLen: 5, // "0 " + 3 chars
			includeID: true,
		},
		{
			name: "default RNG",
			options: []RandomMessageGeneratorOption{
				WithFixedLength(4),
				WithIncludeSeqInMsg(false),
			},
			seqMax:    1,
			expectLen: 4,
			includeID: false,
		},
		{
			name: "minLen > maxLen",
			options: []RandomMessageGeneratorOption{
				WithMinMaxLength(7, 3),
				WithRNG(&mockRNG{sequence: []uint64{0, 0, 0, 0, 0, 0, 0}}),
				WithIncludeSeqInMsg(false),
			},
			seqMax:    1,
			expectLen: 7,
			includeID: false,
		},
		{
			name: "missing sequence generator",
			options: []RandomMessageGeneratorOption{
				WithFixedLength(3),
				WithIncludeSeqInMsg(false),
			},
			seqMax:      0,
			expectLen:   0,
			includeID:   false,
			shouldPanic: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var seq sequencegenerator.SequenceGenerator
			if tc.seqMax > 0 {
				seq = &dummySeq{max: tc.seqMax}
			}
			opts := tc.options
			if seq != nil {
				opts = append(opts, WithSequenceGenerator(seq))
			}
			if tc.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic for missing sequence generator")
					}
				}()
			}
			gen := NewRandomMessageGenerator(opts...)
			if tc.shouldPanic {
				return
			}
			msg := gen.GenerateMessage()
			require.NotNil(t, msg)
			assert.Len(t, msg, tc.expectLen)
			if tc.includeID {
				assert.Contains(t, string(msg), "0 ")
			}
		})
	}
}
