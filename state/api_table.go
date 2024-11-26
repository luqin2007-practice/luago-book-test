package state

import "go-luacompiler/api"

func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}

func (self *luaState) CreateTable(nArr, nRec int) {
	table := newLuaTable(nArr, nRec)
	self.stack.push(table)
}

func (self *luaState) GetTable(index int) api.LuaType {
	table := self.stack.get(index)
	key := self.stack.pop()
	return self.getTable(table, key)
}

func (self *luaState) GetField(index int, k string) api.LuaType {
	table := self.stack.get(index)
	return self.getTable(table, k)
}

func (self *luaState) GetI(index int, i int64) api.LuaType {
	table := self.stack.get(index)
	return self.getTable(table, i)
}

func (self *luaState) GetGlobal(name string) api.LuaType {
	table := self.registry.get(api.LUA_REGISTRYINDEX)
	return self.getTable(table, name)
}

// getTable 从表中获取值，并将结果放入栈顶
func (self *luaState) getTable(table luaValue, key luaValue) api.LuaType {
	if t, ok := table.(*luaTable); ok {
		v := t.get(key)
		self.stack.push(v)
		return typeOf(v)
	}
	panic("not a table!")
}
