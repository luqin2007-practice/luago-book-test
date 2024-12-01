package state

import "go-luacompiler/api"

type cmp_operator struct {
	metamethod string
	fun        func(a, b luaValue) (bool, bool)
}

var cmp_operators = []cmp_operator{
	{"__eq", _eq},
	{"__lt", _lt},
	{"__le", _le},
}

func (self *luaState) Compare(index1, index2 int, op api.CompareOp) bool {
	return self._compare(index1, index2, op, false)
}

func (self *luaState) RawEqual(index1, index2 int) bool {
	return self._compare(index1, index2, api.LUA_OPEQ, true)
}

func (self *luaState) _compare(index1, index2 int, op api.CompareOp, raw bool) bool {
	a := self.stack.get(index1)
	b := self.stack.get(index2)

	// 默认方法
	if result, ok := cmp_operators[op].fun(a, b); ok {
		return result
	}

	// 自定义方法
	if !raw {
		result, ok := callMetamethod(a, b, cmp_operators[op].metamethod, self)
		if ok {
			return convertToBoolean(result)
		}
	}

	panic("comparison error!")
}

// 比较两个值是否相等
func _eq(a, b luaValue) (bool, bool) {
	switch x := a.(type) {
	case nil:
		return b == nil, true
	case bool:
		y, ok := b.(bool)
		return ok && x == y, true
	case string:
		y, ok := b.(string)
		return ok && x == y, true
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y, true
		case float64:
			return float64(x) == y, true
		default:
			return false, true
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x == float64(y), true
		case float64:
			return x == y, true
		default:
			return false, true
		}
	default:
		return a == b, false
	}
}

// 比较 a < b
func _lt(a, b luaValue) (bool, bool) {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y, true
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y, true
		case float64:
			return float64(x) < y, true
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x < float64(y), true
		case float64:
			return x < y, true
		}
	}
	return false, false
}

// 比较 a <= b
func _le(a, b luaValue) (bool, bool) {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y, true
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y, true
		case float64:
			return float64(x) <= y, true
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x <= float64(y), true
		case float64:
			return x <= y, true
		}
	}
	return false, false
}
