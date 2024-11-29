package test

import (
	"go-luacompiler/state"
	"os"
	"testing"
)

/*
main <closure.lua:0,0> (29 instructions at 0000000000b991b0)
0+ params, 2 slots, 1 upvalue, 0 locals, 4 constants, 1 function

	1       [7]     CLOSURE         0 0     ; 0000000000b99320
	2       [1]     SETTABUP        0 -1 0  ; _ENV "newCounter"
	3       [9]     GETTABUP        0 0 -1  ; _ENV "newCounter"
	4       [9]     CALL            0 1 2
	5       [9]     SETTABUP        0 -2 0  ; _ENV "c1"
	6       [10]    GETTABUP        0 0 -3  ; _ENV "print"
	7       [10]    GETTABUP        1 0 -2  ; _ENV "c1"
	8       [10]    CALL            1 1 0
	9       [10]    CALL            0 0 1
	10      [11]    GETTABUP        0 0 -3  ; _ENV "print"
	11      [11]    GETTABUP        1 0 -2  ; _ENV "c1"
	12      [11]    CALL            1 1 0
	13      [11]    CALL            0 0 1
	14      [13]    GETTABUP        0 0 -1  ; _ENV "newCounter"
	15      [13]    CALL            0 1 2
	16      [13]    SETTABUP        0 -4 0  ; _ENV "c2"
	17      [14]    GETTABUP        0 0 -3  ; _ENV "print"
	18      [14]    GETTABUP        1 0 -4  ; _ENV "c2"
	19      [14]    CALL            1 1 0
	20      [14]    CALL            0 0 1
	21      [15]    GETTABUP        0 0 -3  ; _ENV "print"
	22      [15]    GETTABUP        1 0 -4  ; _ENV "c2"
	23      [15]    CALL            1 1 0
	24      [15]    CALL            0 0 1
	25      [16]    GETTABUP        0 0 -3  ; _ENV "print"
	26      [16]    GETTABUP        1 0 -4  ; _ENV "c2"
	27      [16]    CALL            1 1 0
	28      [16]    CALL            0 0 1
	29      [16]    RETURN          0 1

constants (4) for 0000000000b991b0:

	1       "newCounter"
	2       "c1"
	3       "print"
	4       "c2"

locals (0) for 0000000000b991b0:
upvalues (1) for 0000000000b991b0:

	0       _ENV    1       0

function <closure.lua:1,7> (4 instructions at 0000000000b99320)
0 params, 2 slots, 0 upvalues, 1 local, 1 constant, 1 function

	1       [2]     LOADK           0 -1    ; 0
	2       [6]     CLOSURE         1 0     ; 0000000000b99540
	3       [6]     RETURN          1 2
	4       [7]     RETURN          0 1

constants (1) for 0000000000b99320:

	1       0

locals (1) for 0000000000b99320:

	0       count   2       5

upvalues (0) for 0000000000b99320:

function <closure.lua:3,6> (6 instructions at 0000000000b99540)
0 params, 2 slots, 1 upvalue, 0 locals, 1 constant, 0 functions

	1       [4]     GETUPVAL        0 0     ; count
	2       [4]     ADD             0 0 -1  ; - 1
	3       [4]     SETUPVAL        0 0     ; count
	4       [5]     GETUPVAL        0 0     ; count
	5       [5]     RETURN          0 2
	6       [6]     RETURN          0 1

constants (1) for 0000000000b99540:

	1       1

locals (0) for 0000000000b99540:
upvalues (1) for 0000000000b99540:

	0       count   1       0
*/
func TestClosure(t *testing.T) {
	data, err := os.ReadFile("closure.out")
	if err != nil {
		panic(err)
	}

	ls := state.New()

	ls.OnBeforeInstExecuted(func(i uint32) {
		PrintOpcodes(i)
	})
	ls.OnAfterInstExecuted(func(i uint32) {
		PrintStack(ls)
	})

	ls.Register("print", _print)
	ls.Load(data, "closure", "b")
	ls.Call(0, 0)
}
