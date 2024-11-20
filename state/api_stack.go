package state

import "fmt"

func (self *luaState) GetTop() int {
	return self.stack.top
}

func (self *luaState) AbsIndex(index int) int {
	return self.stack.absIndex(index)
}

func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true // never fails
}

func (self *luaState) Pop(n int) {
	self.SetTop(-n - 1)
}

func (self *luaState) Copy(from, to int) {
	val := self.stack.get(from)
	self.stack.set(to, val)
}

func (self *luaState) PushValue(index int) {
	val := self.stack.get(index)
	self.stack.push(val)
}

func (self *luaState) Replace(index int) {
	val := self.stack.pop()
	self.stack.set(index, val)
}

func (self *luaState) Insert(index int) {
	self.Rotate(index, 1)
}

func (self *luaState) Remove(index int) {
	self.Rotate(index, -1)
	self.Pop(1)
}

func (self *luaState) Rotate(index, n int) {
	t := self.stack.top - 1
	p := self.stack.absIndex(index) - 1
	var m int

	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}

	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)
}

func (self *luaState) SetTop(index int) {
	absIndex := self.stack.absIndex(index)
	if absIndex < 0 {
		panic(fmt.Sprintf("index %d(%d) out of range", index, absIndex))
	}

	n := self.stack.top - absIndex
	if n > 0 {
		// 给定 index 小于 top, 相当于出栈
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else {
		// 给定 index 大于 top, 相当于向栈插入 n 个 Nil
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}
