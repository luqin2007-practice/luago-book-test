package codegen

import (
	"go-luacompiler/binchunk"
	"go-luacompiler/compiler/ast"
)

func GenProto(chunk *ast.Block) *binchunk.Prototype {
	fd := &ast.FuncDefExp{
		IsVararg: true,
		Block:    chunk,
	}
	fi := newFuncInfo(nil, fd)
	fi.addLocVar("_ENV")
	cgFuncDefExp(fi, fd, 0)
	return toProto(fi)
}
