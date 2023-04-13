package utils

import lua "github.com/yuin/gopher-lua"

type luaUtil struct {
}

var Lua = luaUtil{}

func (l *luaUtil) ExecLua(
	state *lua.LState,
	script string,
	funcName string,
	args ...lua.LValue) (results lua.LValue, err error) {
	err = state.DoString(script)
	if err != nil {
		return nil, err
	}
	fn := state.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		return nil, nil
	}
	state.Push(fn)
	for _, arg := range args {
		state.Push(arg)
	}
	state.Call(len(args), 1)
	return state.Get(state.GetTop()), nil
}
