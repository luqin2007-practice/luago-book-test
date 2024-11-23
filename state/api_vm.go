package state

func (self *luaState) PC() int {
	return self.stack.pc
}

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

func (self *luaState) Fetch() uint32 {
	i := self.stack.closure.proto.Code[self.stack.pc]
	self.stack.pc++
	return i
}

func (self *luaState) GetConst(index int) {
	val := self.stack.closure.proto.Constants[index]
	self.stack.push(val)
}

func (self *luaState) GetRK(rk int) {
	if rk > 0xFF {
		// 常量池
		self.GetConst(rk & 0xFF)
	} else {
		// 寄存器索引
		self.PushValue(rk + 1)
	}
}
