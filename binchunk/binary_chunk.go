package binchunk

import (
	"fmt"
)

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
	varName string
	startPC uint32
	endPC   uint32
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

func PrintList(f *Prototype) {

	// printHead
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varagFlag := ""
	if f.IsVararg > 0 {
		varagFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d %d> (%d instructions)", funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))
	fmt.Printf("\n%d%s params, %d slots, %d upvalues, %d locals, %d constants, %d functions",
		f.NumParams, varagFlag, f.MaxStackSize, len(f.Upvalues), len(f.LocVars), len(f.Constants), len(f.Protos))

	// printCode
	if len(f.Code) > 0 {
		fmt.Printf("\n\tIdx\tLine\tCode")
	}
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		fmt.Printf("\n\t%d\t[%3s]\t0x%08X", pc+1, line, c)
	}

	// printDetail
	fmt.Printf("\nconstants (%d):", len(f.Constants))
	if len(f.Constants) > 0 {
		fmt.Printf("\n\tIdx\tValue")
	}
	for i, k := range f.Constants {
		var constant string
		switch k.(type) {
		case nil:
			constant = "null"
		case bool:
			constant = fmt.Sprintf("%t", k)
		case float64:
			constant = fmt.Sprintf("%g", k)
		case int64:
			constant = fmt.Sprintf("%d", k)
		case string:
			constant = fmt.Sprintf("%q", k)
		default:
			constant = fmt.Sprintf("%v", k)
		}
		fmt.Printf("\n\t%d\t%s", i+1, constant)
	}

	fmt.Printf("\nlocals (%d):", len(f.LocVars))
	if len(f.LocVars) > 0 {
		fmt.Printf("\n\t\tName\tStartPC\tEndPC")
	}
	for i, locVar := range f.LocVars {
		fmt.Printf("\n\t%d\t%s\t%d\t%d", i, locVar.varName, locVar.startPC+1, locVar.endPC+1)
	}

	fmt.Printf("\nupvalues (%d):", len(f.Upvalues))
	if len(f.Upvalues) > 0 {
		fmt.Printf("\n\t\tName\tInstack\tIdx")
	}
	for i, upval := range f.Upvalues {
		upvalueName := "-"
		if len(f.UpvalueNames) > 0 {
			upvalueName = f.UpvalueNames[i]
		}
		fmt.Printf("\n\t%d\t%s\t%d\t%d", i, upvalueName, upval.Instack, upval.Idx)
	}

	// printPrototypes
	for _, p := range f.Protos {
		PrintList(p)
	}
}
