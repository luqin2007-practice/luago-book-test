package codegen

// allocReg 分配一个寄存器，返回寄存器索引
func (self *funcInfo) allocReg() int {
	return self.allocRegs(1)
}

// allocReg 分配多个寄存器，返回第一个寄存器的索引
func (self *funcInfo) allocRegs(n int) int {
	self.usedRegs += n
	if self.usedRegs >= 255 {
		panic("function or expression needs too many registers")
	}
	if self.usedRegs > self.maxRegs {
		self.maxRegs = self.usedRegs
	}
	return self.usedRegs - n
}

// freeReg 回收上一个分配的寄存器
func (self *funcInfo) freeReg() {
	self.usedRegs--
}

// freeReg 回收最近分配的多个寄存器
func (self *funcInfo) freeRegs(n int) {
	self.usedRegs -= n
}
