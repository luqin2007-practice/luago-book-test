package state

import "go-luacompiler/binchunk"

type luaState struct {
	stack *luaStack
}

func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaState(stackSize, proto),
	}
}

// pushLuaStack 调用栈入栈
func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
}

// popLuaStack 调用栈出栈
func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
}
