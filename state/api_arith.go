package state

import "go-luacompiler/api"

type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var operators = []operator{
	{iadd, fadd},
	{isub, fsub},
	{imul, fmul},
	{imod, fmod},
	{nil, pow},
	{nil, div},
	{iidiv, fidiv},
	{band, nil},
	{bor, nil},
	{bxor, nil},
	{shl, nil},
	{shr, nil},
	{iunm, funm},
	{bnot, nil},
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
	if result := arith(a, b, operator); result != nil {
		self.stack.push(result)
	} else {
		panic("arithmetic error!")
	}
}

func arith(a, b luaValue, op operator) luaValue {
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
