package codegen

import (
	. "go-luacompiler/vm"
)
import . "go-luacompiler/compiler/lexer"

func (self *funcInfo) pc() int {
	return len(self.insts) - 1
}

func (self *funcInfo) emitABC(opcode, a, b, c int) {
	self.insts = append(self.insts, uint32(b<<23|c<<14|a<<6|opcode))
}

func (self *funcInfo) emitABx(opcode, a, bx int) {
	self.insts = append(self.insts, uint32(bx<<14|a<<6|opcode))
}

func (self *funcInfo) emitAsBx(opcode, a, sBx int) {
	self.insts = append(self.insts, uint32((sBx+MAXARG_sBx)<<14|a<<6|opcode))
}

func (self *funcInfo) emitAx(opcode, ax int) {
	self.insts = append(self.insts, uint32(ax<<6|opcode))
}

// r[a] = r[b]
func (self *funcInfo) emitMove(a, b int) {
	self.emitABC(OP_MOVE, a, b, 0)
}

// r[a], r[a+1], ..., r[a+b] = nil
func (self *funcInfo) emitLoadNil(a, n int) {
	self.emitABC(OP_LOADNIL, a, n-1, 0)
}

// r[a] = (bool)b; if (c) pc++
func (self *funcInfo) emitLoadBool(a, b, c int) {
	self.emitABC(OP_LOADBOOL, a, b, c)
}

// r[a] = kst[bx]
func (self *funcInfo) emitLoadK(a int, k interface{}) {
	idx := self.indexOfConstant(k)
	if idx < (1 << 18) {
		self.emitABx(OP_LOADK, a, idx)
	} else {
		self.emitABx(OP_LOADKX, a, 0)
		self.emitAx(OP_EXTRAARG, idx)
	}
}

// r[a], r[a+1], ..., r[a+b-2] = vararg
func (self *funcInfo) emitVararg(a, n int) {
	self.emitABC(OP_VARARG, a, n+1, 0)
}

// r[a] = emitClosure(proto[bx])
func (self *funcInfo) emitClosure(a, bx int) {
	self.emitABx(OP_CLOSURE, a, bx)
}

// r[a] = {}
func (self *funcInfo) emitNewTable(a, nArr, nRec int) {
	self.emitABC(OP_NEWTABLE,
		a, Int2fb(nArr), Int2fb(nRec))
}

// r[a][(c-1)*FPF+i] := r[a+i], 1 <= i <= b
func (self *funcInfo) emitSetList(a, b, c int) {
	self.emitABC(OP_SETLIST, a, b, c)
}

// r[a] := r[b][rk(c)]
func (self *funcInfo) emitGetTable(a, b, c int) {
	self.emitABC(OP_GETTABLE, a, b, c)
}

// r[a][rk(b)] = rk(c)
func (self *funcInfo) emitSetTable(a, b, c int) {
	self.emitABC(OP_SETTABLE, a, b, c)
}

// r[a] = upval[b]
func (self *funcInfo) emitGetUpval(a, b int) {
	self.emitABC(OP_GETUPVAL, a, b, 0)
}

// upval[b] = r[a]
func (self *funcInfo) emitSetUpval(a, b int) {
	self.emitABC(OP_SETUPVAL, a, b, 0)
}

// r[a] = upval[b][rk(c)]
func (self *funcInfo) emitGetTabUp(a, b, c int) {
	self.emitABC(OP_GETTABUP, a, b, c)
}

// upval[a][rk(b)] = rk(c)
func (self *funcInfo) emitSetTabUp(a, b, c int) {
	self.emitABC(OP_SETTABUP, a, b, c)
}

// r[a], ..., r[a+c-2] = r[a](r[a+1], ..., r[a+b-1])
func (self *funcInfo) emitCall(a, nArgs, nRet int) {
	self.emitABC(OP_CALL, a, nArgs+1, nRet+1)
}

// return r[a](r[a+1], ... ,r[a+b-1])
func (self *funcInfo) emitTailCall(a, nArgs int) {
	self.emitABC(OP_TAILCALL, a, nArgs+1, 0)
}

// return r[a], ... ,r[a+b-2]
func (self *funcInfo) emitReturn(a, n int) {
	self.emitABC(OP_RETURN, a, n+1, 0)
}

// r[a+1] := r[b]; r[a] := r[b][rk(c)]
func (self *funcInfo) emitSelf(a, b, c int) {
	self.emitABC(OP_SELF, a, b, c)
}

// pc+=sBx; if (a) close all upvalues >= r[a - 1]
func (self *funcInfo) emitJmp(a, sBx int) int {
	self.emitAsBx(OP_JMP, a, sBx)
	return len(self.insts) - 1
}

// if not (r[a] <=> c) then pc++
func (self *funcInfo) emitTest(a, c int) {
	self.emitABC(OP_TEST, a, 0, c)
}

// if (r[b] <=> c) then r[a] := r[b] else pc++
func (self *funcInfo) emitTestSet(a, b, c int) {
	self.emitABC(OP_TESTSET, a, b, c)
}

func (self *funcInfo) emitForPrep(a, sBx int) int {
	self.emitAsBx(OP_FORPREP, a, sBx)
	return len(self.insts) - 1
}

func (self *funcInfo) emitForLoop(a, sBx int) int {
	self.emitAsBx(OP_FORLOOP, a, sBx)
	return len(self.insts) - 1
}

func (self *funcInfo) emitTForCall(a, c int) {
	self.emitABC(OP_TFORCALL, a, 0, c)
}

func (self *funcInfo) emitTForLoop(a, sBx int) {
	self.emitAsBx(OP_TFORLOOP, a, sBx)
}

// r[a] = op r[b]
func (self *funcInfo) emitUnaryOp(op, a, b int) {
	switch op {
	case TOKEN_OP_NOT:
		self.emitABC(OP_NOT, a, b, 0)
	case TOKEN_OP_BNOT:
		self.emitABC(OP_BNOT, a, b, 0)
	case TOKEN_OP_LEN:
		self.emitABC(OP_LEN, a, b, 0)
	case TOKEN_OP_UNM:
		self.emitABC(OP_UNM, a, b, 0)
	}
}

type logicBinop struct {
	opcode int  // 操作码
	result int  // 比较结果
	swap   bool // 是否交换形参顺序
}

var _arithAndBitwiseBinops = map[int]int{
	TOKEN_OP_ADD:  OP_ADD,
	TOKEN_OP_SUB:  OP_SUB,
	TOKEN_OP_MUL:  OP_MUL,
	TOKEN_OP_MOD:  OP_MOD,
	TOKEN_OP_POW:  OP_POW,
	TOKEN_OP_DIV:  OP_DIV,
	TOKEN_OP_IDIV: OP_IDIV,
	TOKEN_OP_BAND: OP_BAND,
	TOKEN_OP_BOR:  OP_BOR,
	TOKEN_OP_BXOR: OP_BXOR,
	TOKEN_OP_SHL:  OP_SHL,
	TOKEN_OP_SHR:  OP_SHR,
}

var _logicBinops = map[int]logicBinop{
	TOKEN_OP_EQ: {OP_EQ, 1, false},
	TOKEN_OP_NE: {OP_EQ, 0, false},
	TOKEN_OP_LT: {OP_LT, 1, false},
	TOKEN_OP_GT: {OP_LT, 1, true},
	TOKEN_OP_LE: {OP_LE, 1, false},
	TOKEN_OP_GE: {OP_LE, 1, true},
}

// r[a] = rk[b] op rk[c]
// arith & bitwise & relational
func (self *funcInfo) emitBinaryOp(op, a, b, c int) {
	if opcode, found := _arithAndBitwiseBinops[op]; found {
		// 算数运算符
		self.emitABC(opcode, a, b, c)
	} else if opcode, found := _logicBinops[op]; found {
		// 逻辑运算符
		if opcode.swap {
			self.emitABC(opcode.opcode, opcode.result, b, a)
		} else {
			self.emitABC(opcode.opcode, opcode.result, a, b)
		}
		self.emitJmp(0, 1)
		self.emitLoadBool(a, 0, 1)
		self.emitLoadBool(a, 1, 0)
	}
}
