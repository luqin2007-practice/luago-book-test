package vm

import "go-luacompiler/api"

/*
指令：CLOSURE (iABx)

伪代码：

	R(A) := closure(KPROTO[Bx])
*/
func closure(i Instruction, vm api.LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

/*
指令：CALL (iABC)

伪代码：

	func := R(A)
	R(A), R(A+1), ..., R(A+C-2) := func(R(A+1), R(A+2), ..., R(A+B-1))
*/
func call(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

/*
指令：RETURN (iABC)

伪代码：

	return R(A), R(A+1), ..., R(A+B-2)
*/
func lua_return(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 {
		// 无返回值
	} else if b > 1 {
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		_fixStack(a, vm)
	}
}

/*
指令：VARARG (iABC)

伪代码：

	R(A), R(A+1), ..., R(A+B-2) = vararg
*/
func vararg(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b != 1 {
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

/*
指令：TAILCALL (iABC)

伪代码（不优化）：

	func := R(A)
	return func(R(A+1), R(A+2), ..., R(A+B-1))
*/
func tailcall(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, -1)
	_popResults(a, 0, vm)
}

/*
指令：SELF (iABC)

伪代码：（不优化）

	R(A+1) := R(B)
	R(A) := R(B)[Kst(C)]
*/
func self(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// _pushFuncAndArgs 将函数和参数推入栈顶
//
//	函数：R(a)
//	参数：从 R(a+1) 开始，共 b-1 个
func _pushFuncAndArgs(a, b int, vm api.LuaVM) int {
	if b >= 1 {
		// 函数声明有参数
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		// 函数声明没有参数
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

// _popResults 函数调用完成后，将结果存入指定位置
//
//	数量：c-1
//	位置：R(a) 到 R(a+c-2)
func _popResults(a, c int, vm api.LuaVM) {
	if c == 1 {
		// 无返回
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// 返回所有值
		// 暂时将值放在栈中，并存入目标位置
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

func _fixStack(a int, vm api.LuaVM) {
	// 取出目标位置
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}
