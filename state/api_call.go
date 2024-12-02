package state

import (
	"fmt"
	"go-luacompiler/api"
	"go-luacompiler/binchunk"
	"go-luacompiler/vm"
)

func (self *luaState) Load(chunk []byte, chunkName string, mode string) int {
	var proto *binchunk.Prototype
	var closure *closure
	if "b" == mode {
		proto = binchunk.Undump(chunk)
		closure = newLuaClosure(proto)
	} else {
		// TODO 暂时先只实现加载二进制数据
		panic(fmt.Sprintf("Chunk mode %s not supported!", mode))
	}
	self.stack.push(closure)

	// 设置 _ENV
	if len(proto.Protos) > 0 {
		env := self.registry.get(api.LUA_RIDX_GLOBALS)
		closure.upvals[0] = &upvalue{&env}
	}
	return 0
}

func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-nArgs - 1)
	c, ok := val.(*closure)
	//fmt.Printf("call %s<%d,%d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)

	// 查找元表
	if !ok {
		if mf := getMetafield(self, val, "__call"); mf != nil {
			if c, ok = mf.(*closure); ok {
				self.stack.push(val)
				self.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			self.callLuaClosure(c, nArgs, nResults)
		} else {
			self.callGoClosure(c, nArgs, nResults)
		}
	} else {
		panic("not a function or closure!")
	}
}

func (self *luaState) callLuaClosure(c *closure, nArgs, nResults int) {
	nRegs := int(c.proto.MaxStackSize) // 函数所需寄存器大小
	nParams := int(c.proto.NumParams)  // 函数声明参数数量
	isVararg := c.proto.IsVararg == 1  // 函数是否包含变长参数

	// 创建闭包（函数）调用栈
	newStack := newLuaState(nRegs+api.LUA_MINSTACK, self)
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

func (self *luaState) callGoClosure(c *closure, nArgs, nResults int) {
	newStack := newLuaState(nArgs+api.LUA_MINSTACK, self)
	newStack.closure = c

	if nArgs > 0 {
		args := self.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
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
