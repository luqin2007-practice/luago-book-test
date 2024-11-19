package main

import (
	"fmt"
	"go-luacompiler/binchunk"
	"os"
)

func main() {
	fmt.Println("Hello Lua")
	data, err := os.ReadFile("C:\\Dev\\projects\\go-luacompiler\\lua-test\\helloworld53.lunc")
	if err != nil {
		panic(err)
	}
	f := binchunk.Undump(data)
	PrintList(f)
}

func PrintList(f *binchunk.Prototype) {

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

	// printPrototypes
	for _, p := range f.Protos {
		PrintList(p)
	}
}
