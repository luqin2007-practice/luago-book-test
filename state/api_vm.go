package state

func (self *luaState) PC() int {
	return self.stack.pc
}

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

func (self *luaState) Fetch() uint32 {
	self.stack.pc++
	return self.stack.closure.proto.Code[self.stack.pc-1]
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

func (self *luaState) RegisterCount() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

func (self *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(self.stack.varargs)
	}

	self.stack.check(n)
	self.stack.pushN(self.stack.varargs, n)
}

func (self *luaState) LoadProto(n int) {
	stack := self.stack
	proto := stack.closure.proto.Protos[n]
	closure := newLuaClosure(proto)
	stack.push(closure)

	for i, uvInfo := range proto.Upvalues {
		idx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 {
			// 是当前函数上下文的局部变量：访问局部变量
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}
			if openuv, ok := stack.openuvs[idx]; ok {
				// Open：栈中直接引用
				closure.upvals[i] = openuv
			} else {
				// Closed：存于其他位置
				closure.upvals[i] = &upvalue{&stack.slots[idx]}
				stack.openuvs[idx] = closure.upvals[i]
			}
		} else {
			// 非局部变量：该参数已被外层函数捕获，直接从外部函数获取
			closure.upvals[i] = stack.closure.upvals[idx]
		}
	}
}

func (self *luaState) CloseUpvalues(a int) {
	for i, openuv := range self.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(self.stack.openuvs, i)
		}
	}
}
