package codegen

import . "go-luacompiler/compiler/ast"

func cgStat(fi *funcInfo, stat Stat) {
	switch stat := stat.(type) {
	case *FuncCallStat:
		cgFuncCallStat(fi, stat)
	case *BreakStat:
		cgBreakStat(fi, stat)
	case *DoStat:
		cgDoStat(fi, stat)
	case *WhileStat:
		cgWhileStat(fi, stat)
	case *RepeatStat:
		cgRepeatStat(fi, stat)
	case *IfStat:
		cgIfStat(fi, stat)
	case *ForNumStat:
		cgForNumStat(fi, stat)
	case *ForInStat:
		cgForInStat(fi, stat)
	case *AssignStat:
		cgAssignStat(fi, stat)
	case *LocalVarDeclStat:
		cgLocalVarDeclStat(fi, stat)
	case *LocalFuncDefStat:
		cgLocalFuncDefStat(fi, stat)
	case *LabelStat, *GotoStat:
		panic("label and goto statements are not supported!")
	}
}

func cgFuncCallStat(fi *funcInfo, stat *FuncCallStat) {
	r := fi.allocReg()
	cgFuncCallExp(fi, stat, r, 0)
	fi.freeReg()
}

func cgBreakStat(fi *funcInfo, stat *BreakStat) {
	pc := fi.emitJmp(0, 0)
	fi.addBreakJmp(pc)
}

func cgDoStat(fi *funcInfo, stat *DoStat) {
	fi.enterScope(true)
	cgBlock(fi, stat.Block)
	fi.closeOpenUpvals()
	fi.exitScope()
}

func cgWhileStat(fi *funcInfo, stat *WhileStat) {
	pcBeforeWhile := fi.pc()

	// 条件判断
	r := pushExp(fi, stat.Exp)
	fi.freeReg()
	fi.emitTest(r, 0)
	pcJmpToEnd := fi.emitJmp(0, 0)

	// 进入循环体
	fi.enterScope(true)
	cgBlock(fi, stat.Block)
	fi.closeOpenUpvals()
	fi.emitJmp(0, pcBeforeWhile-fi.pc()-1)
	fi.exitScope()

	fi.fixSbx(pcJmpToEnd, fi.pc()-pcJmpToEnd)
}

func cgRepeatStat(fi *funcInfo, stat *RepeatStat) {
	// 循环体
	fi.enterScope(true)
	pcBeforeRepeat := fi.pc()
	cgBlock(fi, stat.Block)

	// 条件判断
	r := pushExp(fi, stat.Exp)
	fi.freeReg()
	fi.emitTest(r, 0)
	fi.emitJmp(fi.getJmpArgA(), pcBeforeRepeat-fi.pc()-1)
	fi.closeOpenUpvals()
	fi.exitScope()
}

func cgIfStat(fi *funcInfo, stat *IfStat) {
	pcJmpToEnds := make([]int, len(stat.Exps))
	pcJmpToNext := -1

	for i, exp := range stat.Exps {
		if pcJmpToNext >= 0 {
			fi.fixSbx(pcJmpToNext, fi.pc()-pcJmpToNext)
		}

		// if (...)
		r := pushExp(fi, exp)
		fi.emitTest(r, 0)
		pcJmpToNext = fi.emitJmp(0, 0)

		// [body]
		fi.enterScope(false)
		cgBlock(fi, stat.Blocks[i])
		fi.closeOpenUpvals()
		fi.exitScope()

		// -> else
		if i < len(stat.Exps)-1 {
			pcJmpToEnds[i] = fi.emitJmp(0, 0)
		} else {
			pcJmpToEnds[i] = fi.emitJmp(0, pcJmpToNext)
		}
	}

	for _, pc := range pcJmpToEnds {
		fi.fixSbx(pc, fi.pc()-pc)
	}
}

func cgForNumStat(fi *funcInfo, stat *ForNumStat) {
	fi.enterScope(true)

	// 生成循环变量
	cgLocalVarDeclStat(fi, &LocalVarDeclStat{
		NameList: []string{"(for index)", "(for limit)", "(for step)"},
		ExpList:  []Exp{stat.InitExp, stat.LimitExp, stat.StepExp},
	})
	fi.addLocVar(stat.VarName)

	// 循环体
	a := fi.usedRegs - 4
	pcForPrep := fi.emitForPrep(a, 0)
	cgBlock(fi, stat.Block)
	fi.closeOpenUpvals()
	pcForLoop := fi.emitForLoop(a, 0)

	// 循环跳转指令
	fi.fixSbx(pcForPrep, pcForLoop-pcForPrep-1)
	fi.fixSbx(pcForLoop, pcForPrep-pcForLoop)
	fi.exitScope()
}

func cgForInStat(fi *funcInfo, stat *ForInStat) {
	fi.enterScope(true)

	// 循环变量
	cgLocalVarDeclStat(fi, &LocalVarDeclStat{
		NameList: []string{"(for generator)", "(for stat)", "(for control)"},
		ExpList:  stat.ExpList,
	})
	for _, name := range stat.NameList {
		fi.addLocVar(name)
	}

	// 循环体
	pcJmpToTFC := fi.emitJmp(0, 0)
	cgBlock(fi, stat.Block)
	fi.fixSbx(pcJmpToTFC, fi.pc()-pcJmpToTFC)

	// 循环跳转指令
	rGenerator := fi.slotOfLocVar("(for generator)")
	fi.emitTForCall(rGenerator, len(stat.NameList))
	fi.emitTForLoop(rGenerator+2, pcJmpToTFC-fi.pc()-1)

	fi.exitScope()
}

func cgAssignStat(fi *funcInfo, stat *AssignStat) {
	nExps := len(stat.ExpList)
	nVars := len(stat.VarList)
	oldRegs := fi.usedRegs
	tRegs := make([]int, nVars) // 表临时变量
	kRegs := make([]int, nVars) // 键临时变量
	vRegs := make([]int, nVars) // 值临时变量

	// 处理左侧 t[k] 类型键，对右侧表达式求值
	for i, exp := range stat.ExpList {
		if taExp, ok := exp.(*TableAccessExp); ok {
			tRegs[i] = pushExp(fi, taExp.PrefixExp)
			kRegs[i] = pushExp(fi, taExp.KeyExp)
		}
	}

	// 分配寄存器
	for i := 0; i < nVars; i++ {
		vRegs[i] = fi.usedRegs + i
	}

	// 生成最终赋值代码
	if nExps >= nVars {
		for i, exp := range stat.ExpList {
			a := fi.allocReg()
			if i >= nVars && i == nExps-1 && isVarargOrFuncCall(exp) {
				cgExp(fi, exp, a, 0)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
	} else {
		multRet := false
		for i, exp := range stat.ExpList {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				// - 特殊处理最后一个值为函数调用或 ...
				multRet = true
				n := nVars - nExps + 1
				cgExp(fi, exp, a, n)
				fi.freeRegs(n - 1)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
		if !multRet {
			// - 还不足够则插入 nil
			n := nVars - nExps
			a := fi.allocRegs(n)
			fi.emitLoadNil(a, n)
		}
	}

	for i, exp := range stat.VarList {
		if nameExp, ok := exp.(*NameExp); ok {
			// 直接赋值
			if a := fi.slotOfLocVar(nameExp.Name); a >= 0 {
				// 函数内局部变量
				fi.emitMove(a, vRegs[i])
			} else if u := fi.indexOfUpval(nameExp.Name); u >= 0 {
				// 外层函数局部变量 Upvalue
				fi.emitSetUpval(u, vRegs[i])
			} else {
				// 全局变量
				env := fi.indexOfUpval("_ENV")
				idx := 0x100 + fi.indexOfConstant(nameExp.Name)
				fi.emitSetTabUp(env, idx, vRegs[i])
			}
		} else {
			// 表索引
			fi.emitSetTable(tRegs[i], kRegs[i], vRegs[i])
		}
	}
	fi.usedRegs = oldRegs
}

func cgLocalVarDeclStat(fi *funcInfo, stat *LocalVarDeclStat) {
	nExps := len(stat.ExpList)
	nNames := len(stat.NameList)

	oldRegs := fi.usedRegs
	if nExps == nNames {
		// 表达式数量与变量名数量相等
		for _, exp := range stat.ExpList {
			pushExp(fi, exp)
		}
	} else if nExps > nNames {
		// 表达式数量大于变量名
		for i, exp := range stat.ExpList {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				// 特殊处理最后一个值为函数调用或 ...
				cgExp(fi, exp, a, 0)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
	} else {
		// 表达式数量小于变量名
		multRet := false
		for i, exp := range stat.ExpList {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				// - 特殊处理最后一个值为函数调用或 ...
				multRet = true
				n := nNames - nExps + 1
				cgExp(fi, exp, a, n)
				fi.freeRegs(n - 1)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
		if !multRet {
			// - 还不足够则插入 nil
			n := nNames - nExps
			a := fi.allocRegs(n)
			fi.emitLoadNil(a, n)
		}
	}

	// 释放临时变量，声明局部变量
	fi.usedRegs = oldRegs
	for _, name := range stat.NameList {
		fi.addLocVar(name)
	}
}

func cgLocalFuncDefStat(fi *funcInfo, stat *LocalFuncDefStat) {

}
