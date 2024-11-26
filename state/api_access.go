package state

import (
	"fmt"
	"go-luacompiler/api"
)

func (self *luaState) TypeName(tp api.LuaType) string {
	switch tp {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRING:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (self *luaState) Type(index int) api.LuaType {
	if self.stack.isValid(index) {
		val := self.stack.get(index)
		return typeOf(val)
	}
	return api.LUA_TNONE
}

func (self *luaState) IsNone(index int) bool {
	return self.Type(index) == api.LUA_TNONE
}

func (self *luaState) IsNil(index int) bool {
	return self.Type(index) == api.LUA_TNIL
}

func (self *luaState) IsNoneOrNil(index int) bool {
	return self.IsNone(index) || self.IsNil(index)
}

func (self *luaState) IsBoolean(index int) bool {
	return self.Type(index) == api.LUA_TBOOLEAN
}

func (self *luaState) IsInteger(index int) bool {
	_, ok := self.stack.get(index).(int64)
	return ok
}

func (self *luaState) IsNumber(index int) bool {
	_, ok := self.stack.get(index).(float64)
	return ok
}

func (self *luaState) IsString(index int) bool {
	return self.Type(index) == api.LUA_TSTRING || self.Type(index) == api.LUA_TNUMBER
}

func (self *luaState) IsGoFunction(index int) bool {
	c, ok := self.stack.get(index).(*Closure)
	return ok && c.goFunc != nil
}

func (self *luaState) ToBoolean(index int) bool {
	val := self.stack.get(index)
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func (self *luaState) ToInteger(index int) int64 {
	n, _ := self.ToIntegerX(index)
	return n
}

func (self *luaState) ToIntegerX(index int) (int64, bool) {
	val := self.stack.get(index)
	return convertToInteger(val)
}

func (self *luaState) ToNumber(index int) float64 {
	n, _ := self.ToNumberX(index)
	return n
}

func (self *luaState) ToNumberX(index int) (float64, bool) {
	val := self.stack.get(index)
	return convertToFloat(val)
}

func (self *luaState) ToString(index int) string {
	n, _ := self.ToStringX(index)
	return n
}

func (self *luaState) ToStringX(index int) (string, bool) {
	switch x := self.stack.get(index).(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		// !!! 这里会更新 stack 值
		self.stack.set(index, s)
		return s, true
	default:
		return "", false
	}
}

func (self *luaState) ToGoFunction(index int) api.GoFunction {
	c, _ := self.stack.get(index).(*Closure)
	return c.goFunc
}
