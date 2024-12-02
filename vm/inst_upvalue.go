package vm

import "go-luacompiler/api"

/*
指令：GETUPVAL (iABC)

伪代码：

	R(A) := Upvalue[B]
*/
func getUpval(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(vm.LuaUpvalueIndex(b), a)
}

/*
指令：SETUPVAL (iABC)

伪代码：

	Upvalue[B] := R(A)
*/
func setUpval(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, vm.LuaUpvalueIndex(b))
}

/*
指令：GETTABUP (iABC)

伪代码：

	key := Kst(C)
	table := R(B)
	R(A) = table[key]
*/
func getTabUp(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(vm.LuaUpvalueIndex(b))
	vm.Replace(a)
}

/*
指令：SETTABUP (iABC)

伪代码：

	Upvalue[A][Kst(B)] = Kst(C)
*/
func setTabUp(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(vm.LuaUpvalueIndex(a))
}
