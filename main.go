package main

import (
	"fmt"
	"go-luacompiler/binchunk"
	"os"
)

func main() {
	fmt.Println("Hello Lua")
	data, err := os.ReadFile("C:\\Dev\\projects\\go-luacompiler\\lua-test\\helloworld53.lunc")
	if err != nil {
		panic(err)
	}
	f := binchunk.Undump(data)
	binchunk.PrintList(f)
}
