package state

import (
	"go-luacompiler/api"
)

type luaValue interface{}

func typeOf(v interface{}) api.LuaType {
	switch v.(type) {
	case nil:
		return api.LUA_TNIL
	case bool:
		return api.LUA_TBOOLEAN
	case int64:
		return api.LUA_TNUMBER
	case float64:
		return api.LUA_TNUMBER
	case string:
		return api.LUA_TSTRING
	default:
		panic("TODO") // TODO 其他类型暂未实现
	}
}
