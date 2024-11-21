package vm

import "go-luacompiler/api"

const LFIELDS_PER_FLUSH = 50

/*
指令：NEWTABLE (iABC)

伪代码：

	B = fb2int(B)
	C = fb2int(C)
	R(A) := CreateTable(B, C)
*/
func newTable(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

/*
指令：GETTABLE (iABC)

伪代码：

	key := Kst(C)
	table := R(B)
	R(A) := table[key]
*/
func getTable(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

/*
指令：SETTABLE (iABC)

伪代码：

	key := Kst(B)
	value := Kst(C)
	R(A)[key] = [value]
*/
func setTable(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(c)
	vm.GetRK(b)
	vm.SetTable(a)
}

/*
指令：

	SETLIST  (iABC)
	EXTRAARG (iAx) (仅 C > 25600)

伪代码：

	table := R(A)
	i = C 或 iAx(EXTRAARG).Ax

	table[R(A+1)] = R(i+1)
	table[R(A+2)] = R(i+2)
	...
	table[R(A+b)] = R(i+b)
*/
func setList(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	// 下标
	if c > 0 {
		c = c - 1
	} else {
		// 读 EXTRAARG
		c = Instruction(vm.Fetch()).Ax()
	}
	idx := int64(c * LFIELDS_PER_FLUSH)

	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j)
		vm.SetI(a, idx)
	}
}
