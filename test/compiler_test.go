package test

import (
	"go-luacompiler/state"
	"os"
	"testing"
)

func TestLuaCompiler(t *testing.T) {
	data, err := os.ReadFile("helloworld.lua")
	if err != nil {
		panic(err)
	}

	ls := state.New()
	ls.Load(data, "chunk", "f")
	ls.Call(0, 0)
}
