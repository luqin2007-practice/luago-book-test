package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

func parseTableConstructorExp(lexer *Lexer) ast.Exp {
	line := lexer.NextKeyword(TOKEN_SEP_LCURLY)     // '{'
	keys, values := parseFieldList(lexer)           // [fieldlist]
	lastLine := lexer.NextKeyword(TOKEN_SEP_RCURLY) // '}'
	return &ast.TableConstructorExp{
		Line:     line,
		LastLine: lastLine,
		KeyExps:  keys,
		ValExps:  values,
	}
}

func parseFieldList(lexer *Lexer) (keys, values []ast.Exp) {
	keys = make([]ast.Exp, 0)
	values = make([]ast.Exp, 0)
	for lexer.LookAhead() != TOKEN_SEP_RCURLY { // {*}
		key, value := parseField(lexer) // field
		switch lexer.LookAhead() {      // [,|;]
		case TOKEN_SEP_COMMA, TOKEN_SEP_SEMI:
			lexer.NextToken()
		}
		keys = append(keys, key)
		values = append(values, value)
	}
	return
}

// field ::= '[' exp ']' '=' exp | Name '=' exp | exp
func parseField(lexer *Lexer) (key, value ast.Exp) {
	if lexer.LookAhead() == TOKEN_SEP_LBRACK {
		lexer.NextToken()                   // '['
		key = parseExp(lexer)               // key
		lexer.NextKeyword(TOKEN_SEP_RBRACK) // ']'
		lexer.NextKeyword(TOKEN_OP_ASSIGN)  // '='
		value = parseExp(lexer)             // value
		return
	}

	exp := parseExp(lexer)
	if name, ok := exp.(*ast.NameExp); ok {
		key = &ast.StringExp{
			Line: name.Line,
			Val:  name.Name,
		}
		lexer.NextKeyword(TOKEN_OP_ASSIGN)
		value = parseExp(lexer)
	}

	return nil, exp
}
