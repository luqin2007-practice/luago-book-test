package vm

import "go-luacompiler/api"

/*
指令：LOADNIL (iABC)

伪代码：

	R(A) := Nil
	R(A + 1) := Nil
	R(A + 2) := Nil
	...
	R(A + B) := Nil
*/
func loadNil(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

/*
指令：LOADBOOL (iABC)

伪代码：

	R(A) := (bool) B
	if (C) PC++
*/
func loadBool(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.PushBoolean(b != 0)
	vm.Replace(a)

	if c != 0 {
		vm.AddPC(1)
	}
}

/*
指令：LOADK (iABx)

伪代码：

	R(A) := Kst(Bx)
*/
func loadK(i Instruction, vm api.LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.GetConst(bx)
	vm.Replace(a)
}

/*
指令：

	LOADK    (iABx)
	EXTRAARG (iAx)

伪代码：

	A, _ := iABx(LOADK)
	Ax := iAx(EXTRAARG)
	R(A) := Kst(Ax)
*/
func loadKx(i Instruction, vm api.LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}
