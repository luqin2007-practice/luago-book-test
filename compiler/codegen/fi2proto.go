package codegen

import . "go-luacompiler/binchunk"

func toProto(fi *funcInfo) *Prototype {
	return &Prototype{
		NumParams:    byte(fi.numParams),
		IsVararg:     _isVararg(fi),
		MaxStackSize: byte(fi.maxRegs),
		Code:         fi.insts,
		Constants:    _getConstants(fi),
		Upvalues:     _getUpvalues(fi),
		Protos:       _toProtos(fi.subFuncs),
		// debug
		Source:          "",
		LineDefined:     0,
		LastLineDefined: 0,
		LineInfo:        []uint32{},
		LocVars:         []LocVar{},
		UpvalueNames:    []string{},
	}
}

func _isVararg(fi *funcInfo) byte {
	if fi.isVararg {
		return 1
	}
	return 0
}

func _getConstants(fi *funcInfo) []interface{} {
	consts := make([]interface{}, len(fi.constants))
	for k, idx := range fi.constants {
		consts[idx] = k
	}
	return consts
}

func _getUpvalues(fi *funcInfo) []Upvalue {
	upvals := make([]Upvalue, len(fi.upvalues))
	for _, uv := range fi.upvalues {
		if uv.locVarSlot >= 0 {
			upvals[uv.index] = Upvalue{
				Instack: 1,
				Idx:     byte(uv.locVarSlot),
			}
		} else {
			upvals[uv.index] = Upvalue{
				Instack: 0,
				Idx:     byte(uv.upvalIndex),
			}
		}
	}
	return upvals
}

func _toProtos(subf []*funcInfo) []*Prototype {
	protos := make([]*Prototype, len(subf))
	for i, s := range subf {
		protos[i] = toProto(s)
	}
	return protos
}
