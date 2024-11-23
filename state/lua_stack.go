package state

import (
	"fmt"
	"go-luacompiler/binchunk"
)

type luaStack struct {
	slots   []luaValue
	top     int
	prev    *luaStack
	closure *Closure
	varargs []luaValue
	pc      int
}

func newLuaState(size int, proto *binchunk.Prototype) *luaStack {
	return &luaStack{
		slots:   make([]luaValue, size),
		top:     0,
		prev:    nil,
		closure: newLuaClosure(proto),
		varargs: []luaValue{},
		pc:      0,
	}
}

// check 检查 LuaStack 是否可以容纳 n 个值，并尝试扩容
func (self *luaStack) check(n int) {
	// 剩余空间
	free := len(self.slots) - self.top
	// 扩容
	for i := free; i < n; i += 1 {
		self.slots = append(self.slots, nil)
	}
}

func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		// TODO
		panic("stack overflow!")
	}

	self.slots[self.top] = val
	self.top++
}

func (self *luaStack) pop() (val luaValue) {
	if self.top < 1 {
		// TODO
		panic("stack underflow!")
	}
	self.top--
	val = self.slots[self.top]
	self.slots[self.top] = nil
	return
}

// absIndex 将相对索引转换成绝对索引
func (self *luaStack) absIndex(n int) int {
	if n >= 0 {
		return n
	}
	return self.top + n + 1
}

// isInvalid 判断一个索引是否是有效索引
func (self *luaStack) isValid(n int) bool {
	n = self.absIndex(n)
	return n > 0 && n <= self.top
}

// set 向 LuaStack 中写入值
func (self *luaStack) set(n int, val luaValue) {
	absIndex := self.absIndex(n)
	if !self.isValid(absIndex) {
		panic(fmt.Sprintf("invalid index %d!", n))
	}
	self.slots[absIndex-1] = val
}

// get 从 LuaStack 中读取值
func (self *luaStack) get(n int) luaValue {
	n = self.absIndex(n)
	if self.isValid(n) {
		return self.slots[n-1]
	}
	return nil
}

// reverse 交替操作
func (self *luaStack) reverse(from, to int) {
	for from < to {
		self.slots[from], self.slots[to] = self.slots[to], self.slots[from]
		from++
		to--
	}
}
