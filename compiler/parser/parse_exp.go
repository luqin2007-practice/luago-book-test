package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

func parseExp(lexer *Lexer) ast.Exp {

}

// exp {',' exp}
func parseExpList(lexer *Lexer) []ast.Exp {
	firstExp := parseExp(lexer)
	exps := []ast.Exp{firstExp}

	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		exps = append(exps, parseExp(lexer))
	}
	return exps
}

func parsePrefixExp(lexer *Lexer) ast.Exp {}

// funcbody ::= '(' [parlist] ')' block 'end'
func parseFuncDefExp(lexer *Lexer) *ast.FuncDefExp {}
