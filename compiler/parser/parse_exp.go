package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
	"go-luacompiler/number"
)

func parseExp(lexer *Lexer) ast.Exp {
	return parseExp12(lexer)
}

// nil | false | true | Numeral | LiteralString | '...' | functiondef | prefixexp | tableconstructor
func parseExp0(lexer *Lexer) ast.Exp {
	switch lexer.LookAhead() {
	case TOKEN_VARARG:
		line, _, _ := lexer.NextToken()
		return ast.VarargExp{Line: line}
	case TOKEN_KW_NIL:
		line, _, _ := lexer.NextToken()
		return &ast.NilExp{Line: line}
	case TOKEN_KW_FALSE:
		line, _, _ := lexer.NextToken()
		return &ast.FalseExp{Line: line}
	case TOKEN_KW_TRUE:
		line, _, _ := lexer.NextToken()
		return &ast.TrueExp{Line: line}
	case TOKEN_STRING:
		line, _, val := lexer.NextToken()
		return &ast.StringExp{
			Line: line,
			Val:  val,
		}
	case TOKEN_NUMBER:
		return parseNumberExp(lexer)
	case TOKEN_SEP_LCURLY:
		return parseTableConstructorExp(lexer)
	case TOKEN_KW_FUNCTION:
		line, _, _ := lexer.NextToken()
		return parseFuncDefExp(lexer, line)
	default:
		return parsePrefixExp(lexer)
	}
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

func parseNumberExp(lexer *Lexer) ast.Exp {
	line, _, val := lexer.NextToken()
	if f, ok := number.ParseFloat(val); ok {
		return &ast.FloatExp{
			Line: line,
			Val:  f,
		}
	} else if i, ok := number.ParseInteger(val); ok {
		return &ast.IntegerExp{
			Line: line,
			Val:  i,
		}
	}
	panic("Not a number: " + val)
}
