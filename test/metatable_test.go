package test

import (
	"go-luacompiler/state"
	"os"
	"testing"
)

func TestMetatable(t *testing.T) {
	data, err := os.ReadFile("metatable.out")
	if err != nil {
		panic(err)
	}

	ls := state.New()
	//addMonitor(ls)
	ls.Load(data, "metatable", "b")
	ls.Call(0, 0)
}
