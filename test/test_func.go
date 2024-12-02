package test

import (
	"fmt"
	"go-luacompiler/api"
)

func bindFunc(ls api.LuaState) {
	ls.Register("print", _print)
	ls.Register("getmetatable", _getmetatable)
	ls.Register("setmetatable", _setmetatable)
	ls.Register("next", _next)
	ls.Register("ipairs", _ipairs)
	ls.Register("pairs", _pairs)
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

func _next(ls api.LuaState) int {
	ls.SetTop(2) // 两个参数
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

/*
function pairs(t)

	return next, t, nil

end
*/
func _pairs(ls api.LuaState) int {
	ls.PushGoFunction(_next, 0)
	ls.PushValue(1)
	ls.PushNil()
	return 3
}

/*
function ipairs(t)

	return _iPairsAux, t, 0

end
*/
func _ipairs(ls api.LuaState) int {
	ls.PushGoFunction(_iPairsAux, 0)
	ls.PushValue(1)
	ls.PushInteger(0)
	return 3
}

/*
function _iPairsAux(t, i)

	local nextIndex = i + 1
	local nextVal = t[nextIndex]
	if nextVal == nil then
	    return nil
	else
	    return nextIndex, nextVal

end
*/
func _iPairsAux(ls api.LuaState) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == api.LUA_TNIL {
		return 1
	} else {
		return 2
	}
}
