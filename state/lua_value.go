package state

import (
	"go-luacompiler/api"
	"go-luacompiler/number"
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

func convertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case int64:
		return float64(x), true
	case float64:
		return x, true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		// 尝试直接转换为整型
		if i, ok := number.ParseInteger(x); ok {
			return i, true
		}
		// 尝试通过浮点数转换为整型
		if f, ok := number.ParseFloat(x); ok {
			return number.FloatToInteger(f)
		}

		return 0, false
	default:
		return 0, false
	}
}
