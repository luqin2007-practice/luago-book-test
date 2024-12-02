package state

import (
	"fmt"
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
	case int64, float64:
		return api.LUA_TNUMBER
	case string:
		return api.LUA_TSTRING
	case *luaTable:
		return api.LUA_TTABLE
	case *closure:
		return api.LUA_TFUNCTION
	default:
		panic(fmt.Sprintf("Unknown type from value %v", v)) // TODO 其他类型暂未实现
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

func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

// getMetatable 获取数据的元表
func getMetatable(val luaValue, vm *luaState) *luaTable {
	// 表元数据
	if t, ok := val.(*luaTable); ok {
		return t.metatable
	}

	// 其他值元数据
	if mt := vm.registry.get(metatabkeKey(val)); mt != nil {
		return mt.(*luaTable)
	}

	return nil
}

// setMetatable 设置数据元表
func setMetatable(val luaValue, mt *luaTable, ls *luaState) {
	if t, ok := val.(*luaTable); ok {
		// 表元数据
		t.metatable = mt
	} else {
		// 其他值元数据
		ls.registry.put(metatabkeKey(val), mt)
	}
}

// callMetamethod 调用元表中的函数，输入两个参数，输出一个参数
func callMetamethod(arg1, arg2 luaValue, name string, ls *luaState) (luaValue, bool) {
	var mf luaValue
	if mf = getMetafield(ls, arg1, name); mf == nil {
		if mf = getMetafield(ls, arg2, name); mf == nil {
			return nil, false
		}
	}

	// 执行自定义函数
	ls.stack.check(4)
	ls.stack.push(mf)
	ls.stack.push(arg1)
	ls.stack.push(arg2)
	ls.Call(2, 1)
	return ls.stack.pop(), true
}

// getMetafield 获取元表内容
func getMetafield(ls *luaState, val luaValue, fieldName string) luaValue {
	if mt := getMetatable(val, ls); mt != nil {
		return mt.get(fieldName)
	}
	return nil
}

// metatabkeKey 获取元表在注册表中的键
func metatabkeKey(val luaValue) interface{} {
	return fmt.Sprintf("_MT%d", typeOf(val))
}
