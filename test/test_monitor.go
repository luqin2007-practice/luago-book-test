package test

import "go-luacompiler/api"

func addMonitor(ls api.LuaVM) {
	ls.OnBeforeInstExecuted(func(i uint32) {
		PrintOpcodes(i)
	})
	ls.OnAfterInstExecuted(func(i uint32) {
		PrintStack(ls)
	})
}
