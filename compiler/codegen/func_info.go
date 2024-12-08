package codegen

type funcInfo struct {
	// 常量表: map[常量]索引
	constants map[interface{}]int
	// 寄存器
	usedRegs int
	maxRegs  int
	// 局部变量表
	scopeLv  int
	locVars  []*locVarInfo
	locNames map[string]*locVarInfo
	// break
	breaks [][]int
}

type locVarInfo struct {
	prev     *locVarInfo // 所有同名变量在一个链表上
	name     string
	scopeLv  int
	slot     int
	captured bool
}

// indexOfConstant 在常量表中查找索引，若不存在则存入常量表
func (self *funcInfo) indexOfConstant(constVal interface{}) int {
	if index, ok := self.constants[constVal]; ok {
		return index
	}

	index := len(self.constants)
	self.constants[constVal] = index
	return index
}

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

// enterScope 进入一个作用域 scope
func (self *funcInfo) enterScope(breakable bool) {
	self.scopeLv++

	// 处理循环块
	if breakable {
		self.breaks = append(self.breaks, []int{})
	} else {
		self.breaks = append(self.breaks, nil)
	}
}

// addLocVar 向该作用域添加一个局部变量，返回分配的寄存器索引
func (self *funcInfo) addLocVar(name string) int {
	newVar := &locVarInfo{
		name:    name,
		prev:    self.locNames[name],
		scopeLv: self.scopeLv,
		slot:    self.allocReg(),
	}
	self.locVars = append(self.locVars, newVar)
	self.locNames[name] = newVar
	return newVar.slot
}

// slotOfLocVar 获取局部变量绑定的寄存器，未绑定返回 -1
func (self *funcInfo) slotOfLocVar(name string) int {
	if locVar, ok := self.locNames[name]; ok {
		return locVar.slot
	}
	return -1
}

// exitScope 离开作用域
func (self *funcInfo) exitScope() {
	// 处理循环块
	pendingBreakJmps := self.breaks[len(self.breaks)-1]
	self.breaks := self.breaks[:len(self.breaks)-1]
	// TODO Here

	self.scopeLv--
	for _, locVar := range self.locNames {
		if locVar.scopeLv > self.scopeLv {
			self.removeLocVar(locVar)
		}
	}
}

// removeLocVar 移除作用域变量
func (self *funcInfo) removeLocVar(locVar *locVarInfo) {
	self.freeReg()
	if locVar.prev == nil {
		// 没有变量覆盖 - 直接删除
		delete(self.locNames, locVar.name)
	} else if locVar.prev.scopeLv == locVar.scopeLv {
		// 覆盖的变量与当前变量作用域相同 - 删除覆盖的变量
		self.removeLocVar(locVar.prev)
	} else {
		// 使用被覆盖的变量
		self.locNames[locVar.name] = locVar.prev
	}
}

// addBreakJmp 将跳转指令位置写入最近的循环块
func (self *funcInfo) addBreakJmp(pc int) {
	for i := self.scopeLv; i >= 0; i-- {
		if self.breaks[i] != nil {
			self.breaks[i] = append(self.breaks[i], pc)
			return
		}
	}
}
