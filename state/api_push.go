package state

import (
	fmt2 "fmt"
	"go-luacompiler/api"
)

func (self *luaState) PushNil() {
	self.stack.push(nil)
}

func (self *luaState) PushBoolean(b bool) {
	self.stack.push(b)
}

func (self *luaState) PushInteger(n int64) {
	self.stack.push(n)
}

func (self *luaState) PushNumber(n float64) {
	self.stack.push(n)
}

func (self *luaState) PushString(s string) {
	self.stack.push(s)
}

func (self *luaState) PushFString(fmt string, a ...interface{}) {
	self.PushString(fmt2.Sprintf(fmt, a...))
}

func (self *luaState) PushGoFunction(f api.GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := self.stack.pop()
		closure.upvals[i-1] = &upvalue{&val}
	}
	self.stack.push(closure)
}

func (self *luaState) PushGlobalTable() {
	self.stack.push(self.registry.get(api.LUA_RIDX_GLOBALS))
}

func (self *luaState) Register(name string, f api.GoFunction) {
	self.PushGoFunction(f, 0)
	self.SetGlobal(name)
}
