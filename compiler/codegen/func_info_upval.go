package codegen

type upvalInfo struct {
	locVarSlot int
	upvalIndex int
	index      int
}

// indexOfUpval 获取 Upvalue 在函数中出现的顺序；若不在当前函数中则向外层函数查询
func (self *funcInfo) indexOfUpval(name string) int {

	// 查询当前函数的 Upvalue 变量
	if upval, ok := self.upvalues[name]; ok {
		return upval.index
	}

	// 查询外层函数
	if self.parent != nil {

		// 变量是否是外部函数的局部变量
		if locVar, ok := self.parent.locNames[name]; ok {
			index := len(self.upvalues)
			locVar.captured = true
			self.upvalues[name] = upvalInfo{
				locVarSlot: locVar.slot,
				upvalIndex: -1,
				index:      index,
			}
			return index
		}

		// 变量是否是外部函数的 Upvalue 变量
		if uvIndex := self.parent.indexOfUpval(name); uvIndex >= 0 {
			index := len(self.upvalues)
			self.upvalues[name] = upvalInfo{
				locVarSlot: -1,
				upvalIndex: uvIndex,
				index:      index,
			}
			return index
		}
	}

	return -1
}
