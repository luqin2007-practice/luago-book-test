package parser

import (
	"go-luacompiler/compiler/ast"
	"go-luacompiler/compiler/lexer"
)

func optimizeUnaryOp(exp *ast.UnopExp) ast.Exp {
	switch exp.Op {
	case lexer.TOKEN_OP_UNM:
		return optimizeUnm(exp)
	case lexer.TOKEN_OP_NOT:
		return optimizeNot(exp)
	case lexer.TOKEN_OP_BNOT:
		return optimizeBnot(exp)
	default:
		return exp
	}
}

func optimizeUnm(exp *ast.UnopExp) ast.Exp {
	switch x := exp.Exp.(type) {
	case *ast.IntegerExp:
		x.Val = -x.Val
		return x
	case *ast.FloatExp:
		x.Val = -x.Val
		return x
	default:
		return exp
	}
}

func optimizeNot(exp *ast.UnopExp) ast.Exp {
	switch x := exp.Exp.(type) {
	case *ast.TrueExp:
		return &ast.FalseExp{Line: x.Line}
	case *ast.FalseExp:
		return &ast.TrueExp{Line: x.Line}
	default:
		return exp
	}
}

func optimizeBnot(exp *ast.UnopExp) ast.Exp {
	switch x := exp.Exp.(type) {
	case *ast.IntegerExp:
		return &ast.IntegerExp{Line: x.Line, Val: ^x.Val}
	default:
		return exp
	}
}
