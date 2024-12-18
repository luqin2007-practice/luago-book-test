package test

import (
	"go-luacompiler/state"
	"os"
	"testing"
)

func TestGoFunction(t *testing.T) {
	data, err := os.ReadFile("helloworld.out")
	if err != nil {
		panic(err)
	}

	ls := state.New()
	ls.Load(data, "chunk", "b")
	ls.Call(0, 0)
}
