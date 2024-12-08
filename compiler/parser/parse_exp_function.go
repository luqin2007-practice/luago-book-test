package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

// parseFuncDefExp 解析函数体构造表达式，从形参列表开始
func parseFuncDefExp(lexer *Lexer, line int) *ast.FuncDefExp {
	lexer.NextKeyword(TOKEN_SEP_LPAREN)         // '('
	parList, isVararg := parseParList(lexer)    // [parlist]
	lexer.NextKeyword(TOKEN_SEP_RPAREN)         // ')'
	block := parseBlock(lexer)                  // coding...
	lastLine := lexer.NextKeyword(TOKEN_KW_END) // end
	return &ast.FuncDefExp{
		Line:     line,
		LastLine: lastLine,
		ParList:  parList,
		IsVararg: isVararg,
		Block:    block,
	}
}

// parlist ::= namelist [',' '...'] | '...'
func parseParList(lexer *Lexer) (parList []string, isVararg bool) {
	parList = []string{}
	isVararg = false
	for lexer.LookAhead() != TOKEN_SEP_RPAREN { // [*]
		switch lexer.LookAhead() {
		case TOKEN_IDENTIFIER: // Name
			_, _, val := lexer.NextToken()
			parList = append(parList, val)
		case TOKEN_VARARG: // ...
			lexer.NextToken()
			isVararg = true
		default:
			lexer.NextToken()
		}
	}
	return
}
