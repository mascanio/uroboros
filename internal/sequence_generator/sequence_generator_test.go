package sequencegenerator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSequenceGenerator_NormalRange(t *testing.T) {
	gen := NewSequenceGenerator(10, 13)
	v, done := gen.Next()
	require.Equal(t, int64(10), v)
	require.Equal(t, NotDone, done)

	v, done = gen.Next()
	require.Equal(t, int64(11), v)
	require.Equal(t, NotDone, done)

	v, done = gen.Next()
	require.Equal(t, int64(12), v)
	require.Equal(t, NotDone, done)

	v, done = gen.Next()
	require.Equal(t, int64(0), v)
	require.Equal(t, Done, done)
}

func TestSequenceGenerator_EmptyRange(t *testing.T) {
	gen := NewSequenceGenerator(5, 5)
	v, done := gen.Next()
	require.Equal(t, int64(0), v)
	require.Equal(t, Done, done)
}

func TestSequenceGenerator_SingleValue(t *testing.T) {
	gen := NewSequenceGenerator(7, 8)
	v, done := gen.Next()
	require.Equal(t, int64(7), v)
	require.Equal(t, NotDone, done)

	v, done = gen.Next()
	require.Equal(t, int64(0), v)
	require.Equal(t, Done, done)
}
