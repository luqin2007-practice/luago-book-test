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
	return self.getTable(table, key, false)
}

func (self *luaState) GetField(index int, k string) api.LuaType {
	table := self.stack.get(index)
	return self.getTable(table, k, false)
}

func (self *luaState) GetI(index int, i int64) api.LuaType {
	table := self.stack.get(index)
	return self.getTable(table, i, false)
}

func (self *luaState) GetGlobal(name string) api.LuaType {
	table := self.registry.get(api.LUA_RIDX_GLOBALS)
	return self.getTable(table, name, false)
}

func (self *luaState) SetTable(index int) {
	table := self.stack.get(index)
	val := self.stack.pop()
	key := self.stack.pop()
	self.setTable(table, key, val, false)
}

func (self *luaState) SetField(index int, k string) {
	table := self.stack.get(index)
	val := self.stack.pop()
	self.setTable(table, k, val, false)
}

func (self *luaState) SetI(index int, i int64) {
	table := self.stack.get(index)
	val := self.stack.pop()
	self.setTable(table, i, val, false)
}

func (self *luaState) SetGlobal(name string) {
	table := self.registry.get(api.LUA_RIDX_GLOBALS)
	value := self.stack.pop()
	self.setTable(table, name, value, false)
}

func (self *luaState) GetMetatable(index int) bool {
	val := self.stack.get(index)
	mt := getMetatable(val, self)
	if mt != nil {
		self.stack.push(mt)
		return true
	}
	return false
}

func (self *luaState) SetMetatable(index int) {
	mt := self.stack.pop()
	val := self.stack.get(index)
	if mt == nil {
		setMetatable(val, nil, self)
	} else if mt, ok := mt.(*luaTable); ok {
		setMetatable(val, mt, self)
	} else {
		panic("not a table")
	}
}

func (self *luaState) RawGet(index int) api.LuaType {
	table := self.stack.get(index)
	key := self.stack.pop()
	return self.getTable(table, key, true)
}

func (self *luaState) RawSet(index int) {
	table := self.stack.get(index)
	val := self.stack.pop()
	key := self.stack.pop()
	self.setTable(table, key, val, true)
}

func (self *luaState) RawGetI(index int, i int64) api.LuaType {
	table := self.stack.get(index)
	return self.getTable(table, i, true)
}

func (self *luaState) RawSetI(index int, i int64) {
	table := self.stack.get(index)
	val := self.stack.pop()
	self.setTable(table, i, val, true)
}

// getTable 从表中获取值，并将结果放入栈顶
//
//	raw: 是否忽略元方法
func (self *luaState) getTable(table, key luaValue, raw bool) api.LuaType {
	if t, ok := table.(*luaTable); ok {
		v := t.get(key)
		if raw || v != nil || !t.hasMetafield("__index") {
			self.stack.push(v)
			return typeOf(v)
		}
	}

	// 处理元方法
	if !raw {
		if mf := getMetafield(self, table, "__index"); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				// __index 为表时，将行为转发给表
				return self.getTable(x, key, false)
			case *closure:
				// __index 为方法时，调用方法
				self.stack.push(mf)
				self.stack.push(table)
				self.stack.push(key)
				self.Call(2, 1)
				return typeOf(self.stack.get(-1))
			}
		}
	}
	panic("not a table or metatable!")
}

func (self *luaState) setTable(table, key, val luaValue, raw bool) {
	if t, ok := table.(*luaTable); ok {
		if raw || t.get(key) != nil || !t.hasMetafield("__newindex") {
			t.put(key, val)
			return
		}
	}

	if !raw {
		if mf := getMetafield(self, table, "__newindex"); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				// __index 为表时，将行为转发给表
				self.setTable(x, key, val, false)
				return
			case *closure:
				// __index 为方法时，调用方法
				self.stack.push(mf)
				self.stack.push(table)
				self.stack.push(key)
				self.stack.push(val)
				self.Call(3, 0)
				return
			}
		}
	}
	panic("not a table or metatable!")
}
