package test

import (
	"fmt"
	"go-luacompiler/compiler/lexer"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	data, err := os.ReadFile("lexer.lua")
	if err != nil {
		panic(err)
	}

	lex := lexer.NewLexer(string(data), "lexer")
	for {
		line, kind, token := lex.NextToken()
		fmt.Printf("line: %d, kind: %d, token: %s\n", line, kind, token)
		if kind == lexer.TOKEN_EOF {
			break
		}
	}
}
