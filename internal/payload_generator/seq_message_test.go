package payloadgenerator

import (
	sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"
	"github.com/stretchr/testify/require"
	"testing"
)

// mockSequenceGenerator implements SequenceGenerator for testing
// It returns a fixed sequence of numbers, then signals done.
type mockSequenceGenerator struct {
	values []int64
	idx    int
}

func (m *mockSequenceGenerator) Next() (int64, sequencegenerator.IsDone) {
	if m.idx >= len(m.values) {
		return 0, sequencegenerator.Done
	}
	v := m.values[m.idx]
	m.idx++
	return v, sequencegenerator.NotDone
}

func TestSeqMessageGenerator_GenerateMessage(t *testing.T) {
	seq := &mockSequenceGenerator{values: []int64{1, 2}}
	gen := NewIDMessageGenerator("hello", seq)

	msg1 := gen.GenerateMessage()
	require.Equal(t, "id: 1 hello", string(msg1), "first message")

	msg2 := gen.GenerateMessage()
	require.Equal(t, "id: 2 hello", string(msg2), "second message")

	msg3 := gen.GenerateMessage()
	require.Nil(t, msg3, "should be nil after sequence done")
}

func TestSeqMessageGenerator_EmptySequence(t *testing.T) {
	seq := &mockSequenceGenerator{values: []int64{}}
	gen := NewIDMessageGenerator("foo", seq)
	msg := gen.GenerateMessage()
	require.Nil(t, msg, "should be nil for empty sequence")
}
