package state

func (self *luaState) Len(index int) {
	val := self.stack.get(index)

	// 自定义运算
	if result, ok := callMetamethod(val, val, "__len", self); ok {
		self.stack.push(result)
		return
	}

	// 默认运算符
	if s, ok := val.(string); ok {
		self.stack.push(int64(len(s)))
	} else if t, ok := val.(*luaTable); ok {
		self.stack.push(t.len())
	} else {
		panic("length error!")
	}
}

func (self *luaState) Concat(n int) {
	if n == 0 {
		self.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			// 默认连接
			if self.IsString(-1) && self.IsString(-2) {
				s2 := self.ToString(-1)
				s1 := self.ToString(-2)
				self.stack.pop()
				self.stack.pop()
				self.stack.push(s1 + s2)
				continue
			}

			// 自定义运算符
			b := self.stack.pop()
			a := self.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", self); ok {
				self.stack.push(result)
				continue
			}
			panic("concatenation error!")
		}
	}
	// n == 1: do nothing
}

func (self *luaState) RawLen(index int) uint {
	val := self.stack.get(index)
	switch v := val.(type) {
	case *luaTable:
		return uint(v.len())
	case string:
		return uint(len(v))
	default:
		return 0
	}
}

func (self *luaState) Next(index int) bool {
	table := self.stack.get(index)
	t, ok := table.(*luaTable)
	if !ok {
		panic("table expected!")
	}

	key := self.stack.pop()
	if ne := t.nextKey(key); ne != nil {
		self.stack.push(ne)
		self.stack.push(t.get(ne))
		return true
	}
	return false
}
