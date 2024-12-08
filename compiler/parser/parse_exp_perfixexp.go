package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

/*
	prefixexp ::= (Name | '(' exp ')') {
	  '[' exp ']' |
	  '.' Name    |
	  [':' Name] args
	}
*/
func parsePrefixExp(lexer *Lexer) ast.Exp {
	// Name | '(' exp ')'
	var exp ast.Exp
	if lexer.LookAhead() == TOKEN_IDENTIFIER {
		line, name := lexer.NextIdentifier() // Name
		exp = &ast.NameExp{
			Line: line,
			Name: name,
		}
	} else {
		exp = parseParensExp(lexer)
	}

	switch lexer.LookAhead() {
	case TOKEN_SEP_LBRACK:
		// '[' exp ']' 表访问
		line, _, _ := lexer.NextToken()     // '['
		key := parseExp(lexer)              // exp
		lexer.NextKeyword(TOKEN_SEP_RBRACK) // ']'
		return &ast.TableAccessExp{
			LastLine:  line,
			PrefixExp: exp,
			KeyExp:    key,
		}
	case TOKEN_SEP_DOT:
		// '.' Name 表访问
		line, _, _ := lexer.NextToken()       // '.'
		idLine, key := lexer.NextIdentifier() // Name
		return &ast.TableAccessExp{
			LastLine:  line,
			PrefixExp: exp,
			KeyExp:    &ast.NameExp{Line: idLine, Name: key},
		}
	case TOKEN_SEP_COLON, TOKEN_SEP_LPAREN, TOKEN_SEP_LCURLY, TOKEN_STRING:
		// [':' Name] args 函数调用
		return parseFuncCallExp(lexer, exp)
	default:
		return exp
	}
}

func parseParensExp(lexer *Lexer) ast.Exp {
	lexer.NextKeyword(TOKEN_SEP_LPAREN) // '('
	exp := parseExp(lexer)              // exp
	lexer.NextKeyword(TOKEN_SEP_RPAREN) // ')'
	switch exp.(type) {
	case *ast.VarargExp, *ast.FuncCallExp, *ast.NameExp, *ast.TableAccessExp:
		return &ast.ParensExp{Exp: exp}
	}
	return exp
}

func parseFuncCallExp(lexer *Lexer, prefixExp ast.Exp) *ast.FuncCallExp {
	nameExp := parseFuncCallName(lexer)
	line := lexer.Line()
	argExp := parseFuncCallArgs(lexer)
	lastLine := lexer.Line()
	return &ast.FuncCallExp{
		Line:      line,
		LastLine:  lastLine,
		PrefixExp: prefixExp,
		NameExp:   nameExp,
		Args:      argExp,
	}
}

func parseFuncCallName(lexer *Lexer) *ast.StringExp {
	if lexer.LookAhead() == TOKEN_SEP_COLON {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		return &ast.StringExp{
			Line: line,
			Val:  name,
		}
	}

	return nil
}

func parseFuncCallArgs(lexer *Lexer) (args []ast.Exp) {
	switch lexer.LookAhead() {
	case TOKEN_SEP_LPAREN: // '(' 正常函数列表
		lexer.NextToken()
		if lexer.LookAhead() != TOKEN_SEP_RPAREN {
			args = parseExpList(lexer)
		}
		lexer.NextKeyword(TOKEN_SEP_RPAREN)
	case TOKEN_SEP_LCURLY: // '{' 表构造，省略括号
		args = []ast.Exp{parseTableConstructorExp(lexer)}
	default: // 字符串字面量
		line, val := lexer.NextTokenOfKind(TOKEN_STRING)
		args = []ast.Exp{&ast.StringExp{
			Line: line,
			Val:  val,
		}}
	}
	return
}
