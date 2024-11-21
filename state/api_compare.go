package state

import "go-luacompiler/api"

func (self *luaState) Compare(index1, index2 int, op api.CompareOp) bool {
	a := self.stack.get(index1)
	b := self.stack.get(index2)
	switch op {
	case api.LUA_OPEQ:
		return _eq(a, b)
	case api.LUA_OPLT:
		return _lt(a, b)
	case api.LUA_OPLE:
		return _le(a, b)
	default:
		panic("invalid compare op")
	}
}

// 比较两个值是否相等
func _eq(a, b luaValue) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x == float64(y)
		case float64:
			return x == y
		default:
			return false
		}
	default:
		return a == b
	}
}

// 比较 a < b
func _lt(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x < float64(y)
		case float64:
			return x < y
		}
	}
	panic("comparison error")
}

// 比较 a <= b
func _le(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x <= float64(y)
		case float64:
			return x <= y
		}
	}
	panic("comparison error")
}
