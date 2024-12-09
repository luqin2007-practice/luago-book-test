package codegen

import . "go-luacompiler/compiler/ast"

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

func cgVarargExp(fi *funcInfo, exp *VarargExp, a int, n int) {
	if !fi.isVararg {
		panic("cannot use ... outside a vararg function")
	}
	fi.emitVararg(a, n)
}

func cgFuncDefExp(fi *funcInfo, exp *FuncDefExp, a int) {

}

func cgTableConstructorExp(fi *funcInfo, exp *TableConstructorExp, a int) {

}

func cgUnopExp(fi *funcInfo, exp *UnopExp, a int) {

}

func cgBinopExp(fi *funcInfo, exp *BinopExp, a int) {

}

func cgConcatExp(fi *funcInfo, exp *ConcatExp, a int) {

}

func cgNameExp(fi *funcInfo, exp *NameExp, a int) {

}

func cgTableAccessExp(fi *funcInfo, exp *TableAccessExp, a int) {

}

func cgFuncCallExp(fi *funcInfo, exp *FuncCallExp, a int, n int) {

}
