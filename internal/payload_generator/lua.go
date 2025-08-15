package payloadgenerator

import (
	"errors"
	"fmt"

	sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"
	lua "github.com/yuin/gopher-lua"
)

type LuaMessageGenerator struct {
	luaState *lua.LState
	seq      sequencegenerator.SequenceGenerator
}

func NewLuaMessageGenerator(
	luaCode string,
	seq sequencegenerator.SequenceGenerator,
) (*LuaMessageGenerator, error) {
	rv := &LuaMessageGenerator{seq: seq}

	if err := checkFuncCorrect(luaCode); err != nil {
		return nil, err
	}

	l := lua.NewState()
	err := l.DoString(luaCode)
	if err != nil {
		l.Close()
		return nil, err
	}

	rv.luaState = l
	return rv, nil
}

func checkFuncCorrect(luaCode string) error {
	l := lua.NewState()
	defer l.Close()

	err := l.DoString(luaCode)
	if err != nil {
		return fmt.Errorf("error parsing lua code: %w", err)
	}

	// Check function
	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("generate"), // name of Lua function
		NRet:    1,                       // number of returned values
		Protect: true,                    // return err or panic
	}, lua.LNumber(0)); err != nil {
		return fmt.Errorf("error calling lua function: %w", err)
	}
	defer l.Pop(1)
	_, ok := l.Get(-1).(lua.LString)
	if !ok {
		return errors.New("lua code generate func does not return an string")
	}
	return nil
}

func (g *LuaMessageGenerator) GenerateMessage() []byte {
	i, done := g.seq.Next()
	if done == sequencegenerator.Done {
		return nil
	}
	if err := g.luaState.CallByParam(lua.P{
		Fn:      g.luaState.GetGlobal("generate"), // name of Lua function
		NRet:    1,                                // number of returned values
		Protect: true,                             // return err or panic
	}, lua.LNumber(i)); err != nil {
		panic(err)
	}
	defer g.luaState.Pop(1)
	str := g.luaState.Get(-1).(lua.LString)
	return []byte(str)
}

func (g *LuaMessageGenerator) Finalize() {
	g.luaState.Close()
}
