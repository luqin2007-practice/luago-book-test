package binchunk

type Header struct {
	Signature       [4]byte
	Version         byte
	Format          byte
	LuacData        [6]byte
	CintSize        byte
	SizetSize       byte
	InstructionSize byte
	LuaIntegerSize  byte
	LuaNumberSize   byte
	LuacInt         int64
	LuacNum         float64
} // 头部 共 30 字节

type Prototype struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
} // 函数原型

type Upvalue struct {
	Instack byte
	Idx     byte
} // upvalue 表

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
} // 局部变量表

// Header
const (
	LUA_SIGNATURE    = "\x1bLua"
	LUA_VERSION      = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

// Constants
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

func Undump(data []byte) *Prototype {
	reader := reader{data, 0}
	reader.readHeader(true)
	reader.readByte() // skip upvalue
	body := reader.readProto("")
	return body
}
