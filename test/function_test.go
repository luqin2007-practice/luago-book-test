package test

import (
	"go-luacompiler/state"
	"os"
	"testing"
)

/*
main <.\function.lua:0,0> (67 instructions at 0000000000cc91e0)
0+ params, 14 slots, 1 upvalue, 8 locals, 8 constants, 2 functions

	1       [10]    CLOSURE         0 0     ; 0000000000cc9390
	2       [14]    CLOSURE         1 1     ; 0000000000cc9680
	3       [16]    MOVE            2 0
	4       [16]    LOADK           3 -1    ; 3
	5       [16]    LOADK           4 -2    ; 9
	6       [16]    LOADK           5 -3    ; 7
	7       [16]    LOADK           6 -4    ; 128
	8       [16]    LOADK           7 -5    ; 35
	9       [16]    CALL            2 6 2
	10      [17]    MOVE            3 1
	11      [17]    EQ              1 2 -4  ; - 128
	12      [17]    JMP             0 1     ; to 14
	13      [17]    LOADBOOL        4 0 1
	14      [17]    LOADBOOL        4 1 0
	15      [17]    CALL            3 2 1
	16      [18]    MOVE            3 0
	17      [18]    LOADK           4 -1    ; 3
	18      [18]    LOADK           5 -2    ; 9
	19      [18]    LOADK           6 -3    ; 7
	20      [18]    LOADK           7 -4    ; 128
	21      [18]    LOADK           8 -5    ; 35
	22      [18]    CALL            3 6 3
	23      [19]    MOVE            5 1
	24      [19]    EQ              0 3 -4  ; - 128
	25      [19]    JMP             0 2     ; to 28
	26      [19]    EQ              1 4 -6  ; - 4
	27      [19]    JMP             0 1     ; to 29
	28      [19]    LOADBOOL        6 0 1
	29      [19]    LOADBOOL        6 1 0
	30      [19]    CALL            5 2 1
	31      [20]    MOVE            5 0
	32      [20]    MOVE            6 0
	33      [20]    LOADK           7 -1    ; 3
	34      [20]    LOADK           8 -2    ; 9
	35      [20]    LOADK           9 -3    ; 7
	36      [20]    LOADK           10 -4   ; 128
	37      [20]    LOADK           11 -5   ; 35
	38      [20]    CALL            6 6 0
	39      [20]    CALL            5 0 3
	40      [21]    MOVE            7 1
	41      [21]    EQ              0 5 -4  ; - 128
	42      [21]    JMP             0 2     ; to 45
	43      [21]    EQ              1 6 -7  ; - 1
	44      [21]    JMP             0 1     ; to 46
	45      [21]    LOADBOOL        8 0 1
	46      [21]    LOADBOOL        8 1 0
	47      [21]    CALL            7 2 1
	48      [22]    NEWTABLE        7 0 0
	49      [22]    MOVE            8 0
	50      [22]    LOADK           9 -1    ; 3
	51      [22]    LOADK           10 -2   ; 9
	52      [22]    LOADK           11 -3   ; 7
	53      [22]    LOADK           12 -4   ; 128
	54      [22]    LOADK           13 -5   ; 35
	55      [22]    CALL            8 6 0
	56      [22]    SETLIST         7 0 1   ; 1
	57      [23]    MOVE            8 1
	58      [23]    GETTABLE        9 7 -7  ; 1
	59      [23]    EQ              0 9 -4  ; - 128
	60      [23]    JMP             0 3     ; to 64
	61      [23]    GETTABLE        9 7 -8  ; 2
	62      [23]    EQ              1 9 -6  ; - 4
	63      [23]    JMP             0 1     ; to 65
	64      [23]    LOADBOOL        9 0 1
	65      [23]    LOADBOOL        9 1 0
	66      [23]    CALL            8 2 1
	67      [23]    RETURN          0 1

constants (8) for 0000000000cc91e0:

	1       3
	2       9
	3       7
	4       128
	5       35
	6       4
	7       1
	8       2

locals (8) for 0000000000cc91e0:

	0       max     2       68
	1       assert  3       68
	2       v1      10      68
	3       v2      23      68
	4       i2      23      68
	5       v3      40      68
	6       i3      40      68
	7       t       57      68

upvalues (1) for 0000000000cc91e0:

	0       _ENV    1       0

function <.\function.lua:1,10> (21 instructions at 0000000000cc9390)
0+ params, 8 slots, 0 upvalues, 7 locals, 2 constants, 0 functions

	1       [2]     NEWTABLE        0 0 0
	2       [2]     VARARG          1 0
	3       [2]     SETLIST         0 0 1   ; 1
	4       [3]     LOADNIL         1 1
	5       [4]     LOADK           3 -1    ; 1
	6       [4]     LEN             4 0
	7       [4]     LOADK           5 -1    ; 1
	8       [4]     FORPREP         3 8     ; to 17
	9       [5]     EQ              1 1 -2  ; - nil
	10      [5]     JMP             0 3     ; to 14
	11      [5]     GETTABLE        7 0 6
	12      [5]     LT              0 1 7
	13      [5]     JMP             0 3     ; to 17
	14      [6]     GETTABLE        7 0 6
	15      [6]     MOVE            2 6
	16      [6]     MOVE            1 7
	17      [4]     FORLOOP         3 -9    ; to 9
	18      [9]     MOVE            3 1
	19      [9]     MOVE            4 2
	20      [9]     RETURN          3 3
	21      [10]    RETURN          0 1

constants (2) for 0000000000cc9390:

	1       1
	2       nil

locals (7) for 0000000000cc9390:

	0       args    4       22
	1       val     5       22
	2       idx     5       22
	3       (for index)     8       18
	4       (for limit)     8       18
	5       (for step)      8       18
	6       i       9       17

upvalues (0) for 0000000000cc9390:

function <.\function.lua:12,14> (5 instructions at 0000000000cc9680)
1 param, 2 slots, 1 upvalue, 1 local, 1 constant, 0 functions

	1       [13]    TEST            0 1
	2       [13]    JMP             0 2     ; to 5
	3       [13]    GETTABUP        1 0 -1  ; _ENV "fail"
	4       [13]    CALL            1 1 1
	5       [14]    RETURN          0 1

constants (1) for 0000000000cc9680:

	1       "fail"

locals (1) for 0000000000cc9680:

	0       v       1       6

upvalues (1) for 0000000000cc9680:

	0       _ENV    0       0
*/
func TestFunction(t *testing.T) {
	data, err := os.ReadFile("function.out")
	if err != nil {
		panic(err)
	}

	ls := state.New()
	ls.Load(data, "function", "b")
	ls.Call(0, 0)
}
