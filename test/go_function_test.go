package test

import (
	"fmt"
	"go-luacompiler/api"
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
	ls.Register("print", _print)
	ls.Load(data, "chunk", "b")
	ls.Call(0, 0)
}

// _print 实现替代 Lua 的 print 函数
func _print(ls api.LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("bool: %t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print("str: " + ls.ToString(i))
		} else {
			fmt.Print("obj: " + ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}
