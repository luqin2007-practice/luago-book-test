package state

import "go-luacompiler/binchunk"

type luaState struct {
	stack *luaStack
	proto *binchunk.Prototype
	pc    int
}

func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaState(stackSize),
		proto: proto,
		pc:    0,
	}
}
