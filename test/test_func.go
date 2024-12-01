package test

import (
	"fmt"
	"go-luacompiler/api"
)

func bindFunc(ls api.LuaState) {
	ls.Register("print", _print)
	ls.Register("getmetatable", _getmetatable)
	ls.Register("setmetatable", _setmetatable)
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

func _getmetatable(ls api.LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func _setmetatable(ls api.LuaState) int {
	ls.SetMetatable(1)
	return 1
}
