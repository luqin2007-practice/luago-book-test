package test

import (
	"go-luacompiler/state"
	"os"
	"testing"
)

func TestIterator(t *testing.T) {
	data, err := os.ReadFile("iterator.out")
	if err != nil {
		panic(err)
	}

	ls := state.New()
	//addMonitor(ls)
	ls.Load(data, "iterator", "b")
	ls.Call(0, 0)
}
