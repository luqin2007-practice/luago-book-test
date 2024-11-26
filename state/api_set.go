package state

import "go-luacompiler/api"

func (self *luaState) SetTable(index int) {
	table := self.stack.get(index)
	key := self.stack.pop()
	val := self.stack.pop()
	setTable(table, key, val)
}

func (self *luaState) SetField(index int, k string) {
	table := self.stack.get(index)
	val := self.stack.pop()
	setTable(table, k, val)
}

func (self *luaState) SetI(index int, i int64) {
	table := self.stack.get(index)
	val := self.stack.pop()
	setTable(table, i, val)
}

func (self *luaState) SetGlobal(name string) {
	table := self.registry.get(api.LUA_RIDX_GLOBALS)
	value := self.stack.pop()
	setTable(table, name, value)
}

func (self *luaState) Register(name string, f api.GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}

func setTable(table luaValue, key luaValue, val luaValue) {
	if t, ok := table.(*luaTable); ok {
		t.put(key, val)
	} else {
		panic("not a table!")
	}
}
