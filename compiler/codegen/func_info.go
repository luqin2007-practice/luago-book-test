package codegen

import (
	"fmt"
	"go-luacompiler/compiler/ast"
	"go-luacompiler/vm"
)

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
	// break 位置表
	breaks [][]int
	// Upvalue 表
	parent   *funcInfo
	upvalues map[string]upvalInfo
	// 字节码
	insts []uint32
	// 子函数
	subFuncs  []*funcInfo
	numParams int
	isVararg  bool
}

func newFuncInfo(parent *funcInfo, fd *ast.FuncDefExp) *funcInfo {
	return &funcInfo{
		constants: map[interface{}]int{},
		locVars:   make([]*locVarInfo, 0, 8),
		locNames:  map[string]*locVarInfo{},
		breaks:    make([][]int, 1),
		parent:    parent,
		upvalues:  map[string]upvalInfo{},
		insts:     make([]uint32, 0, 8),
		subFuncs:  []*funcInfo{},
		numParams: len(fd.ParList),
		isVararg:  fd.IsVararg,
	}
}

// enterScope 进入一个作用域 scope
func (self *funcInfo) enterScope(breakable bool) {
	self.scopeLv++

	// 记录最近循环块中的 break 位置
	if breakable {
		self.breaks = append(self.breaks, []int{})
	} else {
		self.breaks = append(self.breaks, nil)
	}
}

// exitScope 离开作用域
func (self *funcInfo) exitScope() {
	// 修复跳转指令
	pendingBreakJmps := self.breaks[len(self.breaks)-1]
	self.breaks = self.breaks[:len(self.breaks)-1]
	a := self.getJmpArgA()
	for _, pc := range pendingBreakJmps {
		sBx := self.pc() - pc
		i := (sBx+vm.MAXARG_sBx)<<14 | a<<6 | vm.OP_JMP
		self.insts[pc] = uint32(i)
	}

	// 清除局部变量
	self.scopeLv--
	for _, locVar := range self.locNames {
		if locVar.scopeLv > self.scopeLv {
			self.removeLocVar(locVar)
		}
	}
}

// closeOpenUpvals 闭合 Upvalue，实际是跳转到第一个被捕获变量
func (self *funcInfo) closeOpenUpvals() {
	a := self.getJmpArgA()
	if a > 0 {
		self.emitJmp(a, 0)
	}
}

func (self *funcInfo) getJmpArgA() int {
	hasCapturedLocVars := false
	minSlotOfLocVars := self.maxRegs
	for _, locVar := range self.locNames {
		if locVar.scopeLv == self.scopeLv {
			for v := locVar; v != nil && v.scopeLv == self.scopeLv; v = v.prev {
				if v.captured {
					hasCapturedLocVars = true
				}
				if v.slot < minSlotOfLocVars && v.name[0] != '(' {
					minSlotOfLocVars = v.slot
				}
			}
		}
	}
	if hasCapturedLocVars {
		return minSlotOfLocVars + 1
	} else {
		return 0
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
	panic(fmt.Sprintf("<break> at line %d not inside a loop!", pc))
}

// fixSbx 将正确的 break 位置写入指令
func (self *funcInfo) fixSbx(pc, sBx int) {
	i := self.insts[pc]
	i = i << 18 >> 18                     // clear sBx
	i = i | uint32(sBx+vm.MAXARG_sBx)<<14 // reset sBx
	self.insts[pc] = i
}
