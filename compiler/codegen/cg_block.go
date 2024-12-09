package codegen

import "go-luacompiler/compiler/ast"

// cgBlock 编译函数/代码块
func cgBlock(fi *funcInfo, node *ast.Block) {
	// 编译指令
	for _, stat := range node.Stats {
		cgStat(fi, stat)
	}

	// 处理返回值
	for node.RetExp != nil {
		cgRetStat(fi, node.RetExp)
	}
}

func cgRetStat(fi *funcInfo, exps []ast.Exp) {
	nExp := len(exps)

	// 无返回值
	if len(exps) == 0 {
		fi.emitReturn(0, 0)
		return
	}

	// 检查返回值最后一项是否为函数调用或 ...
	multRet := isVarargOrFuncCall(exps[nExp-1])

	// 申请寄存器存储值
	for i, exp := range exps {
		r := fi.allocReg()
		if i == nExp-1 && multRet {
			cgExp(fi, exp, r, -1)
		} else {
			cgExp(fi, exp, r, 1)
		}
	}
	fi.freeRegs(nExp)

	// 返回
	if multRet {
		fi.emitReturn(fi.usedRegs, -1)
	} else {
		fi.emitReturn(fi.usedRegs, nExp)
	}
}

func isVarargOrFuncCall(exp ast.Exp) bool {
	switch exp.(type) {
	case *ast.VarargExp, *ast.FuncCallExp:
		return true
	default:
		return false
	}
}
