package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data    []byte
	pointer uint64
}

func (self *reader) readBytes(n uint64) []byte {
	b := self.data[self.pointer : self.pointer+n]
	self.pointer += n
	return b
}

func (self *reader) readByte() byte {
	self.pointer++
	return self.data[self.pointer-1]
}

func (self *reader) readUint32() uint32 {
	return binary.LittleEndian.Uint32(self.readBytes(4))
}

func (self *reader) readUint64() uint64 {
	return binary.LittleEndian.Uint64(self.readBytes(8))
}

func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	stype := self.readByte()
	switch stype {
	case 0x00:
		return ""
	case 0xFF: // 长字符串
		size := self.readUint64()
		str := self.readBytes(size)
		return string(str)
	default: // 短字符串
		size := uint64(stype) - 1
		str := self.readBytes(size)
		return string(str)
	}
}

func (self *reader) readUint32List() []uint32 {
	size := self.readUint32()
	array := make([]uint32, size)
	for i := 0; i < int(size); i++ {
		array[i] = self.readUint32()
	}
	return array
}

func (self *reader) readStringList() []string {
	size := self.readUint32()
	array := make([]string, size)
	for i := 0; i < int(size); i++ {
		array[i] = self.readString()
	}
	return array
}

func readList[T any](self *reader, f func() T) []T {
	size := self.readUint32()
	array := make([]T, size)
	for i := 0; i < int(size); i++ {
		array[i] = f()
	}
	return array
}

// 校验 chunk 头
func (self *reader) readHeader(check bool) *Header {
	signature := self.readBytes(4)
	if check && string(signature) != LUA_SIGNATURE {
		panic("invalid lua signature")
	}

	version := self.readByte()
	if check && version != LUA_VERSION {
		panic("invalid lua version")
	}

	format := self.readByte()
	if check && format != LUAC_FORMAT {
		panic("invalid lua format")
	}

	luacData := self.readBytes(6)
	if check && string(luacData) != LUAC_DATA {
		panic("invalid lua data")
	}

	cintSize := self.readByte()
	if check && cintSize != CINT_SIZE {
		panic("invalid lua cint size")
	}

	sizetSize := self.readByte()
	if check && sizetSize != CSIZET_SIZE {
		panic("invalid lua size_t size")
	}

	instructionSize := self.readByte()
	if check && instructionSize != INSTRUCTION_SIZE {
		panic("invalid lua instruction size")
	}

	luaIntegerSize := self.readByte()
	if check && luaIntegerSize != LUA_INTEGER_SIZE {
		panic("invalid lua lua integer size")
	}

	luaNumberSize := self.readByte()
	if check && luaNumberSize != LUA_NUMBER_SIZE {
		panic("invalid lua number size")
	}

	luacInt := self.readLuaInteger()
	if check && luacInt != LUAC_INT {
		panic("invalid endianness")
	}

	luacNum := self.readLuaNumber()
	if check && luacNum != LUAC_NUM {
		panic("invalid float format")
	}

	return &Header{
		Signature:       [4]byte(signature),
		Version:         version,
		Format:          format,
		LuacData:        [6]byte(luacData),
		CintSize:        cintSize,
		SizetSize:       sizetSize,
		InstructionSize: instructionSize,
		LuaIntegerSize:  luaIntegerSize,
		LuaNumberSize:   luaNumberSize,
		LuacInt:         luacInt,
		LuacNum:         luacNum,
	}
}

// 读取函数原型
func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if source == "" {
		source = parentSource
	}

	return &Prototype{
		Source:          source,
		LineDefined:     self.readUint32(),
		LastLineDefined: self.readUint32(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readUint32List(),
		Constants:       readList(self, self.readConstant),
		Upvalues:        readList(self, self.readUpvalue),
		Protos:          readList(self, func() *Prototype { return self.readProto(parentSource) }),
		LineInfo:        self.readUint32List(),
		LocVars:         readList(self, self.readLocVar),
		UpvalueNames:    self.readStringList(),
	}
}

func (self *reader) readConstant() interface{} {
	tag := self.readByte()
	switch tag {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return self.readByte() == 1
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_INTEGER:
		return self.readLuaInteger()
	case TAG_SHORT_STR:
		return self.readString()
	case TAG_LONG_STR:
		return self.readString()
	default:
		panic("invalid tag")
	}
}

func (self *reader) readUpvalue() Upvalue {
	return Upvalue{
		Instack: self.readByte(),
		Idx:     self.readByte(),
	}
}

func (self *reader) readLocVar() LocVar {
	return LocVar{
		varName: self.readString(),
		startPC: self.readUint32(),
		endPC:   self.readUint32(),
	}
}
