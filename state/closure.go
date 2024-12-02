package state

import (
	"go-luacompiler/api"
	"go-luacompiler/binchunk"
)

type closure struct {
	proto  *binchunk.Prototype
	goFunc api.GoFunction
	upvals []*upvalue
}

type upvalue struct {
	val *luaValue
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	c := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(goFunc api.GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: goFunc}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
