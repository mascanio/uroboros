package payloadgenerator

import (
	"testing"

	sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type luaMockSequenceGenerator struct {
	calls int
}

func (m *luaMockSequenceGenerator) Next() (int64, sequencegenerator.IsDone) {
	if m.calls == 0 {
		m.calls++
		return 42, sequencegenerator.NotDone
	}
	if m.calls == 1 {
		m.calls++
		return 43, sequencegenerator.NotDone
	}
	return 0, sequencegenerator.Done
}

func TestNewLuaMessageGenerator_Success(t *testing.T) {
	luaCode := `
function generate(i)
  return "msg-" .. tostring(i)
end
`
	seq := &luaMockSequenceGenerator{}
	gen, err := NewLuaMessageGenerator(luaCode, seq)
	require.NoError(t, err, "expected success, got error")
	require.NotNil(t, gen, "expected generator, got nil")
	gen.Finalize()
}

func TestNewLuaMessageGenerator_SyntaxError(t *testing.T) {
	luaCode := `function generate(i) return "msg-" .. tostring(i) -- missing end`
	seq := &luaMockSequenceGenerator{}
	gen, err := NewLuaMessageGenerator(luaCode, seq)
	assert.Error(t, err, "expected error for syntax error")
	assert.Nil(t, gen, "expected nil generator on error")
}

func TestNewLuaMessageGenerator_MissingFunction(t *testing.T) {
	luaCode := `function notgenerate(i) return "msg-" .. tostring(i) end`
	seq := &luaMockSequenceGenerator{}
	gen, err := NewLuaMessageGenerator(luaCode, seq)
	assert.Error(t, err, "expected error for missing generate function")
	assert.Nil(t, gen, "expected nil generator on error")
}

func TestNewLuaMessageGenerator_NonStringReturn(t *testing.T) {
	luaCode := `function generate(i) return 123 end`
	seq := &luaMockSequenceGenerator{}
	gen, err := NewLuaMessageGenerator(luaCode, seq)
	assert.Error(t, err, "expected error for non-string return")
	assert.Nil(t, gen, "expected nil generator on error")
}

func TestLuaMessageGenerator_GenerateMessage(t *testing.T) {
	luaCode := `function generate(i) return "msg-" .. tostring(i) end`
	seq := &luaMockSequenceGenerator{}
	gen, err := NewLuaMessageGenerator(luaCode, seq)
	require.NoError(t, err, "unexpected error")
	msg := gen.GenerateMessage()
	assert.Equal(t, "msg-42", string(msg), "expected 'msg-42'")
	msg2 := gen.GenerateMessage()
	assert.Equal(t, "msg-43", string(msg2), "expected 'msg-43'")
	msg3 := gen.GenerateMessage()
	assert.Nil(t, msg3, "expected nil when sequence is done")
	gen.Finalize()
}

func TestLuaMessageGenerator_GenerateMessage_PanicOnLuaError(t *testing.T) {
	luaCode := `
call_count = 0
function generate(i)
  call_count = call_count + 1
  if call_count == 2 then error("fail") end
  return "ok"
end
`
	seq := &luaMockSequenceGenerator{}
	gen, err := NewLuaMessageGenerator(luaCode, seq)
	require.NoError(t, err, "unexpected error")
	defer gen.Finalize()
	_ = gen.GenerateMessage() // first call, should succeed
	assert.Panics(t, func() { _ = gen.GenerateMessage() })
}

func TestLuaMessageGenerator_Finalize(t *testing.T) {
	luaCode := `function generate(i) return "msg-" .. tostring(i) end`
	seq := &luaMockSequenceGenerator{}
	gen, err := NewLuaMessageGenerator(luaCode, seq)
	require.NoError(t, err, "unexpected error")
	assert.NotPanics(t, func() { gen.Finalize() }, "Finalize should not panic")
}
