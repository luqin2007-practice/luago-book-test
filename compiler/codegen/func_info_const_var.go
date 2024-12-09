package codegen

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
