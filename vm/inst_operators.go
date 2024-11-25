package vm

import "go-luacompiler/api"

// ADD +
func add(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPADD)
}

// SUB -
func sub(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPSUB)
}

// MUL *
func mul(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPMUL)
}

// MOD %
func mod(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPMOD)
}

// POW ^
func pow(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPPOW)
}

// DIV /
func div(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPDIV)
}

// IDIV //
func idiv(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPIDIV)
}

// BAND &
func band(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPBAND)
}

// BOR |
func bor(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPBOR)
}

// BXOR ~
func bxor(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPBXOR)
}

// SHL <<
func shl(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPSHL)
}

// SHR >>
func shr(i Instruction, vm api.LuaVM) {
	_binaryArith(i, vm, api.LUA_OPSHR)
}

// UNM -
func unm(i Instruction, vm api.LuaVM) {
	_unaryArith(i, vm, api.LUA_OPUNM)
}

// BNOT ~
func bnot(i Instruction, vm api.LuaVM) {
	_unaryArith(i, vm, api.LUA_OPBNOT)
}

/*
二元运算 (iABC)

伪代码：

	R(A) := Kst(B) op Kst(C)
*/
func _binaryArith(i Instruction, vm api.LuaVM, op api.ArithOp) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a)
}

/*
一元运算 (iABC)

伪代码：

	R(A) := op Kst(B)
*/
func _unaryArith(i Instruction, vm api.LuaVM, op api.ArithOp) {
	a, b, _ := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.Arith(op)
	vm.Replace(a)
}

/*
指令：LEN (iABC)

伪代码：

	R(A) := length of R(B)
*/
func length(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Len(b)
	vm.Replace(a)
}

/*
指令：CONCAT (iABC)

伪代码：

	R(A) := R(B) .. .. .. R(C)
*/
func concat(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	c += 1

	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}

// EQ ==
func eq(i Instruction, vm api.LuaVM) {
	_compare(i, vm, api.LUA_OPEQ)
}

// LT <
func lt(i Instruction, vm api.LuaVM) {
	_compare(i, vm, api.LUA_OPLT)
}

// LE <=
func le(i Instruction, vm api.LuaVM) {
	_compare(i, vm, api.LUA_OPLE)
}

/*
比较指令

伪代码：

	if (Kst(B) op Kst(C)) != bool(A)
	    PC++
*/
func _compare(i Instruction, vm api.LuaVM, op api.CompareOp) {
	a, b, c := i.ABC()

	vm.GetRK(b)
	vm.GetRK(c)
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

/*
指令：NOT (iABC)

伪代码：

	R(A) := not R(B)
*/
func not(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

/*
指令：TESTSET (iABC)

伪代码：

	if (bool(R(B)) == C)
	    R(A) := R(B)
	else
	    PC++
*/
func testSet(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

/*
指令：TEST (iABC)

伪代码：

	if (bool(R(a)) != C)
	    PC++
*/
func test(i Instruction, vm api.LuaVM) {
	a, _, c := i.ABC()
	a += 1

	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}
