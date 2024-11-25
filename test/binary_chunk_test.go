package test

import (
	"go-luacompiler/binchunk"
	"os"
	"testing"
)

func TestPrototype(t *testing.T) {
	data, err := os.ReadFile("helloworld.out")
	if err != nil {
		panic(err)
	}
	f := binchunk.Undump(data)
	PrintPrototype(f)
}
