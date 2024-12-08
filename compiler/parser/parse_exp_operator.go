package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

// exp1  ::= exp0 {'^' exp2} 右结合性
func parseExp1(lexer *Lexer) ast.Exp {
	exp := parseExp0(lexer)
	if lexer.LookAhead() == TOKEN_OP_POW {
		exp = newBinopExp(lexer, exp, parseExp2)
	}
	return exp
}

// exp2  ::= {not | '#' | '-' | '~'} exp1
func parseExp2(lexer *Lexer) ast.Exp {
	switch lexer.LookAhead() {
	case TOKEN_OP_NOT, TOKEN_OP_LEN, TOKEN_OP_UNM, TOKEN_OP_BNOT:
		return newUnopExp(lexer, parseExp2)
	default:
		return parseExp1(lexer)
	}
}

// exp3  ::= exp2 {('*' | '/' | '//' | '%') exp2}
func parseExp3(lexer *Lexer) ast.Exp {
	exp := parseExp2(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_MUL, TOKEN_OP_DIV, TOKEN_OP_IDIV, TOKEN_OP_MOD:
			return newBinopExp(lexer, exp, parseExp2)
		default:
			return exp
		}
	}
}

// exp4  ::= exp3 {('+' | '-') exp3}
func parseExp4(lexer *Lexer) ast.Exp {
	exp := parseExp3(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_ADD, TOKEN_OP_SUB:
			exp = newBinopExp(lexer, exp, parseExp3)
		default:
			return exp
		}
	}
}

// exp5  ::= exp4 {'..' exp4} 右结合性
func parseExp5(lexer *Lexer) ast.Exp {
	exp := parseExp4(lexer)

	if lexer.LookAhead() == TOKEN_OP_CONCAT {
		exps := []ast.Exp{exp}
		line := lexer.Line()
		for lexer.LookAhead() == TOKEN_OP_CONCAT {
			lexer.NextToken()
			exps = append(exps, parseExp4(lexer))
		}
		return &ast.ConcatExp{
			Line: line,
			Exps: exps,
		}
	}

	return exp
}

// exp6  ::= exp5 {('<<' | '>>') exp5}
func parseExp6(lexer *Lexer) ast.Exp {
	exp := parseExp5(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_SHL, TOKEN_OP_SHR:
			exp = newBinopExp(lexer, exp, parseExp5)
		default:
			return exp
		}
	}
}

// exp7  ::= exp6 {'&' exp6}
func parseExp7(lexer *Lexer) ast.Exp {
	exp := parseExp6(lexer)
	for lexer.LookAhead() == TOKEN_OP_BAND {
		exp = newBinopExp(lexer, exp, parseExp6)
	}
	return exp
}

// exp8  ::= exp7 {'~' exp7}
func parseExp8(lexer *Lexer) ast.Exp {
	exp := parseExp7(lexer)
	for lexer.LookAhead() == TOKEN_OP_BXOR {
		exp = newBinopExp(lexer, exp, parseExp7)
	}
	return exp
}

// exp9  ::= exp8 {'|' exp8}
func parseExp9(lexer *Lexer) ast.Exp {
	exp := parseExp8(lexer)
	for lexer.LookAhead() == TOKEN_OP_BOR {
		exp = newBinopExp(lexer, exp, parseExp8)
	}
	return exp
}

// exp10 ::= exp9 {('<' | '>' | '<=' | '>=' | '~=' | '==') exp9}
func parseExp10(lexer *Lexer) ast.Exp {
	exp := parseExp9(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_LT, TOKEN_OP_LE, TOKEN_OP_EQ, TOKEN_OP_GT, TOKEN_OP_GE, TOKEN_OP_NE:
			exp = newBinopExp(lexer, exp, parseExp9)
		default:
			return exp
		}
	}
}

// exp11 ::= exp10 {'and' exp10}
func parseExp11(lexer *Lexer) ast.Exp {
	exp := parseExp10(lexer)
	for lexer.LookAhead() == TOKEN_OP_AND {
		exp = newBinopExp(lexer, exp, parseExp10)
	}
	return exp
}

// exp12 ::= exp11 {'or' exp11}
func parseExp12(lexer *Lexer) ast.Exp {
	exp := parseExp11(lexer)
	for lexer.LookAhead() == TOKEN_OP_OR {
		exp = newBinopExp(lexer, exp, parseExp11)
	}
	return exp
}

func newBinopExp(lexer *Lexer, exp ast.Exp, nextExp func(*Lexer) ast.Exp) ast.Exp {
	line, op, _ := lexer.NextToken()
	return &ast.BinopExp{
		Line: line,
		Op:   op,
		Exp1: exp,
		Exp2: nextExp(lexer),
	}
}

func newUnopExp(lexer *Lexer, nextExp func(*Lexer) ast.Exp) ast.Exp {
	line, op, _ := lexer.NextToken()
	exp := &ast.UnopExp{
		Line: line,
		Op:   op,
		Exp:  nextExp(lexer),
	}
	return optimizeUnaryOp(exp)
}
