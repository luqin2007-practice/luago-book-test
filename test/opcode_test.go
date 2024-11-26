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
0+ params, 6 slots, 1 upvalue, 5 locals, 4 constants, 0 functions

	1       [1]     LOADK           0 -1    ; 0
	2       [2]     LOADK           1 -2    ; 1
	3       [2]     LOADK           2 -3    ; 100
	4       [2]     LOADK           3 -2    ; 1
	5       [2]     FORPREP         1 4     ; to 10
	6       [3]     MOD             5 4 -4  ; - 2
	7       [3]     EQ              0 5 -1  ; - 0
	8       [3]     JMP             0 1     ; to 10
	9       [4]     ADD             0 0 4
	10      [2]     FORLOOP         1 -5    ; to 6
	11      [6]     RETURN          0 1

constants (4) for 0000000000c89050:

	1       0
	2       1
	3       100
	4       2

locals (5) for 0000000000c89050:

	0       sum     2       12
	1       (for index)     5       11
	2       (for limit)     5       11
	3       (for step)      5       11
	4       i       6       10

upvalues (1) for 0000000000c89050:

	0       _ENV    1       0
*/
func TestOpcode(t *testing.T) {
	// 2550
	data, err := os.ReadFile("sum.out")
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
	assert.Equal(t, "[2500][101][100][1][100][0]", PrintStack(ls))
}
