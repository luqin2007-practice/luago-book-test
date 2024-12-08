package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

func Parse(chunk, chunkName string) *ast.Block {
	lexer := NewLexer(chunk, chunkName)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_EOF)
	return block
}
