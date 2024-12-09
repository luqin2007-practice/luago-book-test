package codegen

import (
	. "go-luacompiler/compiler/ast"
	"go-luacompiler/compiler/lexer"
	"go-luacompiler/vm"
)

var _unaryOps = map[int]int{
	lexer.TOKEN_OP_NOT:  vm.OP_NOT,
	lexer.TOKEN_OP_BNOT: vm.OP_BNOT,
	lexer.TOKEN_OP_LEN:  vm.OP_LEN,
	lexer.TOKEN_OP_UNM:  vm.OP_UNM,
}

// cgExp 编译表达式
func cgExp(fi *funcInfo, exp Exp, a, n int) {
	switch exp := exp.(type) {
	case *NilExp:
		fi.emitLoadNil(a, n)
	case *FalseExp:
		fi.emitLoadBool(a, 0, 0)
	case *TrueExp:
		fi.emitLoadBool(a, 1, 0)
	case *IntegerExp:
		fi.emitLoadK(a, exp.Val)
	case *FloatExp:
		fi.emitLoadK(a, exp.Val)
	case *StringExp:
		fi.emitLoadK(a, exp.Val)
	case *ParensExp:
		cgExp(fi, exp.Exp, a, 1)
	case *VarargExp:
		cgVarargExp(fi, exp, a, n)
	case *FuncDefExp:
		cgFuncDefExp(fi, exp, a)
	case *TableConstructorExp:
		cgTableConstructorExp(fi, exp, a)
	case *UnopExp:
		cgUnopExp(fi, exp, a)
	case *BinopExp:
		cgBinopExp(fi, exp, a)
	case *ConcatExp:
		cgConcatExp(fi, exp, a)
	case *NameExp:
		cgNameExp(fi, exp, a)
	case *TableAccessExp:
		cgTableAccessExp(fi, exp, a)
	case *FuncCallExp:
		cgFuncCallExp(fi, exp, a, n)
	}
}

// pushExp 申请一个寄存器，将一个表达式结果存入寄存器中
func pushExp(fi *funcInfo, exp Exp) int {
	a := fi.allocReg()
	cgExp(fi, exp, a, 1)
	return a
}

func cgVarargExp(fi *funcInfo, exp *VarargExp, a int, n int) {
	if !fi.isVararg {
		panic("cannot use ... outside a vararg function")
	}
	fi.emitVararg(a, n)
}

func cgFuncDefExp(fi *funcInfo, exp *FuncDefExp, a int) {
	subf := newFuncInfo(fi, exp)
	fi.subFuncs = append(fi.subFuncs, subf)

	for _, param := range exp.ParList {
		subf.addLocVar(param)
	}
	cgBlock(subf, exp.Block)
	subf.exitScope()
	// 为子函数添加一条 RETURN
	subf.emitReturn(0, 0)

	// 生成 CLOSURE 指令
	bx := len(fi.subFuncs) - 1
	fi.emitClosure(a, bx)
}

func cgTableConstructorExp(fi *funcInfo, exp *TableConstructorExp, a int) {
	nArr := 0

	// 预处理 计算数组长度
	for _, key := range exp.KeyExps {
		if key == nil {
			nArr++
		}
	}

	nExps := len(exp.KeyExps)
	multRet := nExps > 0 && isVarargOrFuncCall(exp.ValExps[nExps-1])

	// 创建表，赋值
	fi.emitNewTable(a, nArr, nExps-nArr)
	arrIdx := 0
	for i, key := range exp.KeyExps {
		val := exp.ValExps[i]
		if key == nil {
			// 数组
			arrIdx++
			tmp := fi.allocReg()
			if i == nExps-1 && multRet {
				cgExp(fi, val, tmp, -1)
			} else {
				cgExp(fi, val, tmp, 1)
			}

			if arrIdx%50 == 0 || arrIdx == nArr {
				n := arrIdx % 50
				if n == 0 {
					n = 50
				}

				c := (arrIdx-1)/50 + 1
				if i == nExps-1 && multRet {
					fi.emitSetList(a, 0, c)
				} else {
					fi.emitSetList(a, n, c)
				}
			}
		} else {
			// 表
			b := pushExp(fi, key)
			c := pushExp(fi, val)
			fi.freeRegs(2)
			fi.emitSetTable(a, b, c)
		}
	}
}

func cgUnopExp(fi *funcInfo, exp *UnopExp, a int) {
	b := pushExp(fi, exp.Exp)
	fi.emitABC(_unaryOps[exp.Op], a, b, 0)
	fi.freeReg()
}

func cgBinopExp(fi *funcInfo, exp *BinopExp, a int) {
	switch exp.Op {
	case lexer.TOKEN_OP_AND, lexer.TOKEN_OP_OR:
		// 判断 AND OR
		b := pushExp(fi, exp.Exp1)
		fi.freeReg()
		if exp.Op == lexer.TOKEN_OP_AND {
			fi.emitTestSet(a, b, 0)
		} else {
			fi.emitTestSet(a, b, 1)
		}
		pcJmp := fi.emitJmp(0, 0)

		b = pushExp(fi, exp.Exp2)
		fi.freeReg()
		fi.emitMove(a, b)
		fi.fixSbx(pcJmp, fi.pc()-pcJmp)
	default:
		b := pushExp(fi, exp.Exp1)
		c := pushExp(fi, exp.Exp2)
		fi.emitBinaryOp(exp.Op, a, b, c)
		fi.freeRegs(2)
	}
}

func cgConcatExp(fi *funcInfo, exp *ConcatExp, a int) {
	for _, part := range exp.Exps {
		pushExp(fi, part)
	}

	c := fi.usedRegs - 1
	b := c - len(exp.Exps) + 1
	fi.freeRegs(c - b + 1)
	fi.emitABC(vm.OP_CONCAT, a, b, c)
}

func cgNameExp(fi *funcInfo, exp *NameExp, a int) {
	if r := fi.slotOfLocVar(exp.Name); r >= 0 {
		// 局部变量
		fi.emitMove(a, r)
	} else if u := fi.indexOfUpval(exp.Name); u >= 0 {
		// Upvalue
		fi.emitGetUpval(a, u)
	} else {
		// 表
		ta := &TableAccessExp{
			PrefixExp: &NameExp{
				Line: 0,
				Name: "_ENV",
			},
			KeyExp: &StringExp{
				Line: 0,
				Val:  exp.Name,
			},
		}
		cgTableAccessExp(fi, ta, a)
	}
}

func cgTableAccessExp(fi *funcInfo, exp *TableAccessExp, a int) {
	b := pushExp(fi, exp.PrefixExp)
	c := pushExp(fi, exp.KeyExp)
	fi.emitGetTable(a, b, c)
	fi.freeRegs(2)
}

func cgFuncCallExp(fi *funcInfo, exp *FuncCallExp, a int, n int) {
	nArgs := _prepFuncCall(fi, exp, a)
	fi.emitCall(a, nArgs, n)
}

func _prepFuncCall(fi *funcInfo, exp *FuncCallExp, a int) int {
	nArgs := len(exp.Args)
	lastArgIsVerargOrFuncCall := false

	// SELF 尾调用优化
	if exp.NameExp == nil {
		c := 0x100 + fi.indexOfConstant(exp.NameExp.Val)
		fi.emitSelf(a, a, c)
	}

	for i, arg := range exp.Args {
		tmp := fi.allocReg()
		if i == nArgs-1 && isVarargOrFuncCall(arg) {
			lastArgIsVerargOrFuncCall = true
			cgExp(fi, arg, tmp, -1)
		} else {
			cgExp(fi, arg, tmp, 1)
		}
	}
	fi.freeRegs(nArgs)

	if exp.NameExp != nil {
		nArgs++
	}

	if lastArgIsVerargOrFuncCall {
		nArgs = -1
	}

	return nArgs
}
