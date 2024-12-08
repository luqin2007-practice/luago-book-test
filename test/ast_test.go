package test

import (
	"encoding/json"
	"fmt"
	"go-luacompiler/compiler/parser"
	"os"
	"testing"
)

func TestAST(t *testing.T) {
	data, err := os.ReadFile("helloworld.lua")
	if err != nil {
		panic(err)
	}

	ast := parser.Parse(string(data), "helloworld")
	b, err := json.Marshal(ast)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
