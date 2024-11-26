package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-luacompiler/binchunk"
	"go-luacompiler/state"
	"go-luacompiler/vm"
	"os"
	"testing"
)

/*
main <.\table.lua:0,0> (14 instructions at 0000000000be9160)
0+ params, 6 slots, 1 upvalue, 2 locals, 9 constants, 0 functions

	1       [1]     NEWTABLE        0 3 0
	2       [1]     LOADK           1 -1    ; "a"
	3       [1]     LOADK           2 -2    ; "b"
	4       [1]     LOADK           3 -3    ; "c"
	5       [1]     SETLIST         0 3 1   ; 1
	6       [2]     SETTABLE        0 -4 -5 ; 2 "B"
	7       [3]     SETTABLE        0 -6 -7 ; "foo" "Bar"
	8       [5]     GETTABLE        1 0 -8  ; 3
	9       [5]     GETTABLE        2 0 -4  ; 2
	10      [5]     GETTABLE        3 0 -9  ; 1
	11      [5]     GETTABLE        4 0 -6  ; "foo"
	12      [5]     LEN             5 0
	13      [5]     CONCAT          1 1 5
	14      [5]     RETURN          0 1

constants (9) for 0000000000be9160:

	1       "a"
	2       "b"
	3       "c"
	4       2
	5       "B"
	6       "foo"
	7       "Bar"
	8       3
	9       1

locals (2) for 0000000000be9160:

	0       t       6       15
	1       s       14      15

upvalues (1) for 0000000000be9160:

	0       _ENV    1       0
*/
func TestTable(t *testing.T) {
	data, err := os.ReadFile("table.out")
	if err != nil {
		panic(err)
	}

	proto := binchunk.Undump(data)

	// 初始化 LuaVM
	regs := int(proto.MaxStackSize)
	ls := state.New()
	ls.SetTop(regs)

	// 运行
	fmt.Printf("PC\t\t%-10s\t栈空间\n", "指令")
	for {
		// 读指令
		pc := ls.PC()
		inst := vm.Instruction(ls.Fetch())

		if inst.Opcode() == vm.OP_RETURN {
			// 遇到 RETURN 指令退出
			break
		}

		// 执行指令
		inst.Execute(ls)
		fmt.Printf("[%02d]\t%-10s\t", pc+1, inst.OpName())
		PrintStack(ls)
	}

	assert.Equal(t, "[table][\"cBaBar3\"][\"B\"][\"a\"][\"Bar\"][3]", PrintStack(ls))
}
