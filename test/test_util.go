package test

import (
	"fmt"
	"go-luacompiler/api"
	"go-luacompiler/binchunk"
	"go-luacompiler/vm"
	"strings"
)

func PrintPrototype(f *binchunk.Prototype) {
	PrintHead(f)
	PrintCode(f)
	PrintDetail(f)
	PrintPrototypes(f)
}

// PrintHead 打印函数头
func PrintHead(f *binchunk.Prototype) {
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
}

// PrintCode 打印函数体
func PrintCode(f *binchunk.Prototype) {
	if len(f.Code) > 0 {
		fmt.Printf("\n\tIdx\tLine\tName\t\tA\tB\tC")
	}
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		i := vm.Instruction(c)
		fmt.Printf("\n\t%d\t[%3s]\t%8s", pc+1, line, i.OpName())
		PrintInstruction(i)
	}
}

// PrintInstruction 打印函数指令
func PrintInstruction(i vm.Instruction) {
	fmt.Printf("\t")
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != vm.OpArgN {
			if b > 0xFF {
				// 常量表索引时以负数形式输出
				fmt.Printf("\t%d", -1-b&0xFF)
			} else {
				fmt.Printf("\t%d", b)
			}
		}
		if i.CMode() != vm.OpArgN {
			if c > 0xFF {
				// 常量表索引时以负数形式输出
				fmt.Printf("\t%d", -1-c&0xFF)
			} else {
				fmt.Printf("\t%d", c)
			}
		}
	case vm.IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.BMode() == vm.OpArgK {
			// 常量表索引时以负数形式输出
			fmt.Printf("\t%d", -1-bx&0xFF)
		} else if i.BMode() == vm.OpArgU {
			fmt.Printf("\t%d", bx)
		}
	case vm.IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d\t%d", a, sbx)
	case vm.IAx:
		fmt.Printf("%d", -1-i.Ax())
	}
}

// PrintDetail 打印详细信息 本地变量+upvalue
func PrintDetail(f *binchunk.Prototype) {
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
		fmt.Printf("\n\t%d\t%s\t%d\t%d", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
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
}

// PrintPrototypes 打印子函数
func PrintPrototypes(f *binchunk.Prototype) {
	for _, p := range f.Protos {
		PrintPrototype(p)
	}
}

func PrintStack(ls api.LuaState) string {
	top := ls.GetTop()
	builder := strings.Builder{}
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			builder.WriteString(fmt.Sprintf("[%t]", ls.ToBoolean(i)))
		case api.LUA_TNUMBER:
			if ls.IsInteger(i) {
				builder.WriteString(fmt.Sprintf("[%d]", ls.ToInteger(i)))
			} else {
				builder.WriteString(fmt.Sprintf("[%f]", ls.ToNumber(i)))
			}
		case api.LUA_TSTRING:
			builder.WriteString(fmt.Sprintf("[%q]", ls.ToString(i)))
		default:
			builder.WriteString(fmt.Sprintf("[%s]", ls.TypeName(t)))
		}
	}
	str := builder.String()
	fmt.Println(str)
	return str
}
