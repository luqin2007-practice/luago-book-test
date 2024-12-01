package state

import "go-luacompiler/api"

type operator struct {
	metamethod  string
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var operators = []operator{
	{"__add", iadd, fadd},
	{"__sub", isub, fsub},
	{"__mul", imul, fmul},
	{"__mod", imod, fmod},
	{"__pow", nil, pow},
	{"__div", nil, div},
	{"__idiv", iidiv, fidiv},
	{"__band", band, nil},
	{"__bor", bor, nil},
	{"__bxor", bxor, nil},
	{"__shl", shl, nil},
	{"__shr", shr, nil},
	{"__unm", iunm, funm},
	{"__bnot", bnot, nil},
}

func (self *luaState) Arith(op api.ArithOp) {
	var a, b luaValue
	b = self.stack.pop()
	if op != api.LUA_OPUNM && op != api.LUA_OPBNOT {
		a = self.stack.pop()
	} else {
		// 一元运算符
		a = b
	}
	operator := operators[op]

	// 普通运算符
	if result := _arith(a, b, operator); result != nil {
		self.stack.push(result)
		return
	}

	// 自定义运算符
	if result, ok := callMetamethod(a, b, operator.metamethod, self); ok {
		self.stack.push(result)
		return
	}

	panic("arithmetic error!")
}

func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil {
		// 仅整型运算（位运算）
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else {
		if op.integerFunc != nil {
			// 先尝试整型运算
			if x, ok := convertToInteger(a); ok {
				if y, ok := convertToInteger(b); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		// 再尝试浮点运算
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	return nil
}
