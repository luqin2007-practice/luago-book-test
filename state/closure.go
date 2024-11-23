package state

import "go-luacompiler/binchunk"

type Closure struct {
	proto *binchunk.Prototype
}

func newLuaClosure(proto *binchunk.Prototype) *Closure {
	return &Closure{proto: proto}
}
