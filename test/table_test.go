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

func TestTable(t *testing.T) {
	data, err := os.ReadFile("table.out")
	if err != nil {
		panic(err)
	}

	proto := binchunk.Undump(data)

	// 初始化 LuaVM
	regs := int(proto.MaxStackSize)
	ls := state.New(regs+8, proto)
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
		printStack(ls)
	}

	assert.Equal(t, "[table][\"cBaBar3\"][\"B\"][\"a\"][\"Bar\"][3]", printStack(ls))
}
