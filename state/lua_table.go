package state

import (
	"go-luacompiler/number"
	"math"
)

type luaTable struct {
	arr       []luaValue
	_map      map[luaValue]luaValue
	metatable *luaTable
}

// newLuaTable 创建 Table
func newLuaTable(nArr, nRec int) *luaTable {
	table := &luaTable{}

	if nArr > 0 {
		table.arr = make([]luaValue, 0, nArr)
	}

	if nRec > 0 {
		table._map = make(map[luaValue]luaValue, nRec)
	}

	return table
}

// get 取数据
func (self *luaTable) get(key luaValue) luaValue {
	key = _keyToInteger(key)
	// 注意 arr 和 _map 可能为 nil
	if i, ok := key.(int64); ok && self.arr != nil && i > 0 && i <= self.len() {
		return self.arr[i-1]
	}
	if self._map != nil {
		if v, ok := self._map[key]; ok {
			return v
		}
	}
	return nil
}

// put 存数据 / 删数据
func (self *luaTable) put(key luaValue, value luaValue) {
	// 校验键类型 nil NaN
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}
	// 尝试转换为整型
	key = _keyToInteger(key)
	if i, ok := key.(int64); ok && i > 0 {
		// 正整数下标，操作列表部分
		if value == nil {
			if i <= self.len() {
				self.arr[i-1] = nil
				// 检查洞 nil 是否在末尾，并压缩数组
				self._shrinkArray()
			}
			// nil 且在列表长度之外，不需要记录
		} else {
			// 扩充数组长度
			self._expandArray(i)
			self.arr[i-1] = value
		}
	} else {
		// 其他类型下标，操作记录部分
		if value == nil && self._map != nil {
			delete(self._map, key)
		} else {
			if self._map == nil {
				self._map = make(map[luaValue]luaValue)
			}
			self._map[key] = value
		}
	}
}

// len 获取表的列表部分长度
func (self *luaTable) len() int64 {
	self._shrinkArray()
	return int64(len(self.arr))
}

func (self *luaTable) hasMetafield(fieldName string) bool {
	return self.metatable != nil && self.metatable.get(fieldName) != nil
}

// _keyToInteger 若 key 为 float64，尝试转换为 int64
func _keyToInteger(key luaValue) luaValue {
	if i, ok := key.(float64); ok {
		if j, ok := number.FloatToInteger(i); ok {
			return j
		}
	}
	return key
}

// _expandArray 扩充数组，直到长度为 n
func (self *luaTable) _expandArray(n int64) {
	if self.arr == nil {
		self.arr = make([]luaValue, n)
	} else {
		count := n - self.len()
		if count > 0 {
			self.arr = append(self.arr, make([]luaValue, n)...)
		}
	}
}

// _shrinkArray 缩减数组，去除末尾的洞
func (self *luaTable) _shrinkArray() {
	if self.arr == nil {
		return
	}

	// i: 最后一个非 nil 位置 + 1
	i := len(self.arr)
	for ; i > 0 && self.arr[i-1] == nil; i-- {
	}
	if i != len(self.arr) {
		self.arr = self.arr[:i]
	}
}
