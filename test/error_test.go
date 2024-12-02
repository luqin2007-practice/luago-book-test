package test

import (
	"go-luacompiler/state"
	"os"
	"testing"
)

func TestError(t *testing.T) {
	data, err := os.ReadFile("error.out")
	if err != nil {
		panic(err)
	}

	ls := state.New()
	ls.Load(data, "error", "b")
	ls.Call(0, 0)
}
