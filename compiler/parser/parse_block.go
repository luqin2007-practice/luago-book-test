package parser

import (
	"go-luacompiler/compiler/ast"
)

import . "go-luacompiler/compiler/lexer"

// block ::= {stat} [retstat]
func parseBlock(lexer *Lexer) *ast.Block {
	return &ast.Block{
		Stats:    parseStats(lexer),
		RetExp:   parseRetExp(lexer),
		LastLine: lexer.Line(),
	}
}

// {stat}
func parseStats(lexer *Lexer) []ast.Stat {
	stats := make([]ast.Stat, 0, 8)
	for !IsReturnOrBlockEndToken(lexer.LookAhead()) {
		stat := parseStat(lexer)
		if _, ok := stat.(*ast.EmptyStat); !ok {
			stats = append(stats, stat)
		}
	}
	return stats
}

// retstat ::= return [explist] [';']
func parseRetExp(lexer *Lexer) []ast.Exp {
	// return
	if lexer.LookAhead() != TOKEN_KW_RETURN {
		return nil
	}
	lexer.NextToken()

	switch lexer.LookAhead() {
	case TOKEN_EOF, TOKEN_KW_END, TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		// 块结束 无返回值
		return []ast.Exp{}
	case TOKEN_SEP_SEMI:
		// ; 无返回值
		lexer.NextToken()
		return []ast.Exp{}
	default:
		// 有返回值
		exps := parseExpList(lexer)
		if lexer.LookAhead() == TOKEN_SEP_SEMI {
			lexer.NextToken()
		}
		return exps
	}
}
