package codegen

import . "go-luacompiler/compiler/ast"

func cgStat(fi *funcInfo, stat Stat) {
	switch stat := node.(type) {
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
	r := fi.allocReg()
	cgExp(fi, stat.Exp)
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

}

func cgIfStat(fi *funcInfo, stat *IfStat) {

}

func cgForNumStat(fi *funcInfo, stat *ForNumStat) {

}

func cgForInStat(fi *funcInfo, stat *ForInStat) {

}

func cgAssignStat(fi *funcInfo, stat *AssignStat) {

}

func cgLocalVarDeclStat(fi *funcInfo, stat *LocalVarDeclStat) {

}

func cgLocalFuncDefStat(fi *funcInfo, stat *LocalFuncDefStat) {

}
