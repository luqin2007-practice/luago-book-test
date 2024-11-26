package state

import (
	"go-luacompiler/api"
	"go-luacompiler/binchunk"
)

type Closure struct {
	proto  *binchunk.Prototype
	goFunc api.GoFunction
}

func newLuaClosure(proto *binchunk.Prototype) *Closure {
	return &Closure{proto: proto}
}

func newGoClosure(goFunc api.GoFunction) *Closure {
	return &Closure{goFunc: goFunc}
}
