package payloadgenerator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Reuse mockSequenceGenerator from seq_message_test.go

func TestStaticMessageGenerator_GenerateMessage(t *testing.T) {
	seq := &mockSequenceGenerator{values: []int64{1, 2}}
	gen := NewStaticMessageGenerator("static", seq)

	msg1 := gen.GenerateMessage()
	require.Equal(t, "static", string(msg1), "first message")

	msg2 := gen.GenerateMessage()
	require.Equal(t, "static", string(msg2), "second message")

	msg3 := gen.GenerateMessage()
	require.Nil(t, msg3, "should be nil after sequence done")
}

func TestStaticMessageGenerator_EmptySequence(t *testing.T) {
	seq := &mockSequenceGenerator{values: []int64{}}
	gen := NewStaticMessageGenerator("foo", seq)
	msg := gen.GenerateMessage()
	require.Nil(t, msg, "should be nil for empty sequence")
}
