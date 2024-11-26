package vm

import "go-luacompiler/api"

/*
指令：GETTABUP (iABC)

伪代码：
*/
func getTabUp(i Instruction, vm api.LuaVM) {
	a, _, c := i.ABC()
	a += 1

	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2)
	vm.Replace(a)
	vm.Pop(1)
}
