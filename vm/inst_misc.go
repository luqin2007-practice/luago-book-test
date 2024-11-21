package vm

import "go-luacompiler/api"

/*
指令：MOVE (iABC)

伪代码：

	R(A) := R(B)
*/
func move(instruction Instruction, vm api.LuaVM) {
	a, b, _ := instruction.ABC()
	a += 1
	b += 1

	vm.Copy(b, a)
}

/*
指令：JMP (iAsBx)

伪代码：

	PC += sBx
	TODO A 与 Upvalue 有关
*/
func jmp(instruction Instruction, vm api.LuaVM) {
	a, sbx := instruction.AsBx()
	vm.AddPC(sbx)
	if a != 0 {
		// ...
	}
}
