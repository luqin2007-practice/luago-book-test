package state

import "go-luacompiler/api"

type luaState struct {
	registry *luaTable
	stack    *luaStack
}

func New() *luaState {
	registry := newLuaTable(0, 0)
	// 全局变量表
	registry.put(api.LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaState(api.LUA_MINSTACK, ls))

	// 绑定 API 函数
	for _, f := range api.ApiFunctions {
		ls.Register(f.Name, f.Function)
	}
	return ls
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
