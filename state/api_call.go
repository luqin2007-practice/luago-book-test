package state

import (
	"fmt"
	"go-luacompiler/binchunk"
	"go-luacompiler/vm"
)

func (self *luaState) Load(chunk []byte, chunkName string, mode string) int {
	if "b" == mode {
		prop := binchunk.Undump(chunk)
		closure := newLuaClosure(prop)
		self.stack.push(closure)
	} else {
		// TODO 暂时先只实现加载二进制数据
		panic(fmt.Sprintf("Chunk mode %s not supported!", mode))
	}
	return 0
}

func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-nArgs - 1)
	c, ok := val.(*Closure)
	if !ok {
		panic("not a function or closure!")
	}
	fmt.Printf("call %s<%d,%d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)

	if c.proto != nil {
		self.callLuaClosure(c, nArgs, nResults)
	} else {
		self.callGoClosure(c, nArgs, nResults)
	}
}

func (self *luaState) LoadProto(n int) {
	proto := self.stack.closure.proto.Protos[n]
	closure := newLuaClosure(proto)
	self.stack.push(closure)
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

func (self *luaState) callLuaClosure(c *Closure, nArgs, nResults int) {
	nRegs := int(c.proto.MaxStackSize) // 函数所需寄存器大小
	nParams := int(c.proto.NumParams)  // 函数声明参数数量
	isVararg := c.proto.IsVararg == 1  // 函数是否包含变长参数

	// 创建闭包（函数）调用栈
	newStack := newLuaState(nRegs + 20)
	newStack.closure = c

	// 从当前栈中提取参数和闭包（函数），并将参数存入被调函数栈
	funcAndArgs := self.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// 将被调函数插入主调函数栈帧，执行后出栈
	self.pushLuaStack(newStack)
	self.runLuaClosure()
	self.popLuaStack()

	// 提取函数执行结果，存入主调函数栈帧
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

// 执行闭包
func (self *luaState) runLuaClosure() {
	for {
		// 执行指令
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)
		// RETURN 指令退出
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func (self *luaState) callGoClosure(c *Closure, nArgs, nResults int) {
	newStack := newLuaState(nArgs + 20)
	newStack.closure = c

	args := self.stack.popN(nArgs)
	newStack.pushN(args, nArgs)
	self.stack.pop()

	self.pushLuaStack(newStack)
	r := c.goFunc(self)
	self.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(r)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}
