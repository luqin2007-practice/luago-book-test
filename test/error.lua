function div0(a, b)
    if b == 0 then
        error("DIV BY ZERO!")
    else
        return a / b
    end
end

function div1(a, b) return div0(a, b) end
function div2(a, b) return div1(a, b) end

ok, result = pcall(div2, 4, 2); print(ok, result)
ok, err = pcall(div2, 5, 0);    print(ok, err)
ok, err = pcall(div2, {}, {});  print(ok, err)

--[[
main <.\error.lua:0,0> (40 instructions at 0000000000c791b0)
0+ params, 4 slots, 1 upvalue, 0 locals, 12 constants, 3 functions
        1       [7]     CLOSURE         0 0     ; 0000000000c79320
        2       [1]     SETTABUP        0 -1 0  ; _ENV "div0"
        3       [9]     CLOSURE         0 1     ; 0000000000c797a0
        4       [9]     SETTABUP        0 -2 0  ; _ENV "div1"
        5       [10]    CLOSURE         0 2     ; 0000000000c798f0
        6       [10]    SETTABUP        0 -3 0  ; _ENV "div2"
        7       [12]    GETTABUP        0 0 -6  ; _ENV "pcall"
        8       [12]    GETTABUP        1 0 -3  ; _ENV "div2"
        9       [12]    LOADK           2 -7    ; 4
        10      [12]    LOADK           3 -8    ; 2
        11      [12]    CALL            0 4 3
        12      [12]    SETTABUP        0 -5 1  ; _ENV "result"
        13      [12]    SETTABUP        0 -4 0  ; _ENV "ok"
        14      [12]    GETTABUP        0 0 -9  ; _ENV "print"
        15      [12]    GETTABUP        1 0 -4  ; _ENV "ok"
        16      [12]    GETTABUP        2 0 -5  ; _ENV "result"
        17      [12]    CALL            0 3 1
        18      [13]    GETTABUP        0 0 -6  ; _ENV "pcall"
        19      [13]    GETTABUP        1 0 -3  ; _ENV "div2"
        20      [13]    LOADK           2 -11   ; 5
        21      [13]    LOADK           3 -12   ; 0
        22      [13]    CALL            0 4 3
        23      [13]    SETTABUP        0 -10 1 ; _ENV "err"
        24      [13]    SETTABUP        0 -4 0  ; _ENV "ok"
        25      [13]    GETTABUP        0 0 -9  ; _ENV "print"
        26      [13]    GETTABUP        1 0 -4  ; _ENV "ok"
        27      [13]    GETTABUP        2 0 -10 ; _ENV "err"
        28      [13]    CALL            0 3 1
        29      [14]    GETTABUP        0 0 -6  ; _ENV "pcall"
        30      [14]    GETTABUP        1 0 -3  ; _ENV "div2"
        31      [14]    NEWTABLE        2 0 0
        32      [14]    NEWTABLE        3 0 0
        33      [14]    CALL            0 4 3
        34      [14]    SETTABUP        0 -10 1 ; _ENV "err"
        35      [14]    SETTABUP        0 -4 0  ; _ENV "ok"
        36      [14]    GETTABUP        0 0 -9  ; _ENV "print"
        37      [14]    GETTABUP        1 0 -4  ; _ENV "ok"
        38      [14]    GETTABUP        2 0 -10 ; _ENV "err"
        39      [14]    CALL            0 3 1
        40      [14]    RETURN          0 1
constants (12) for 0000000000c791b0:
        1       "div0"
        2       "div1"
        3       "div2"
        4       "ok"
        5       "result"
        6       "pcall"
        7       4
        8       2
        9       "print"
        10      "err"
        11      5
        12      0
locals (0) for 0000000000c791b0:
upvalues (1) for 0000000000c791b0:
        0       _ENV    1       0

function <.\error.lua:1,7> (9 instructions at 0000000000c79320)
2 params, 4 slots, 1 upvalue, 2 locals, 3 constants, 0 functions
        1       [2]     EQ              0 1 -1  ; - 0
        2       [2]     JMP             0 4     ; to 7
        3       [3]     GETTABUP        2 0 -2  ; _ENV "error"
        4       [3]     LOADK           3 -3    ; "DIV BY ZERO!"
        5       [3]     CALL            2 2 1
        6       [3]     JMP             0 2     ; to 9
        7       [5]     DIV             2 0 1
        8       [5]     RETURN          2 2
        9       [7]     RETURN          0 1
constants (3) for 0000000000c79320:
        1       0
        2       "error"
        3       "DIV BY ZERO!"
locals (2) for 0000000000c79320:
        0       a       1       10
        1       b       1       10
upvalues (1) for 0000000000c79320:
        0       _ENV    0       0

function <.\error.lua:9,9> (6 instructions at 0000000000c797a0)
2 params, 5 slots, 1 upvalue, 2 locals, 1 constant, 0 functions
        1       [9]     GETTABUP        2 0 -1  ; _ENV "div0"
        2       [9]     MOVE            3 0
        3       [9]     MOVE            4 1
        4       [9]     TAILCALL        2 3 0
        5       [9]     RETURN          2 0
        6       [9]     RETURN          0 1
constants (1) for 0000000000c797a0:
        1       "div0"
locals (2) for 0000000000c797a0:
        0       a       1       7
        1       b       1       7
upvalues (1) for 0000000000c797a0:
        0       _ENV    0       0

function <.\error.lua:10,10> (6 instructions at 0000000000c798f0)
2 params, 5 slots, 1 upvalue, 2 locals, 1 constant, 0 functions
        1       [10]    GETTABUP        2 0 -1  ; _ENV "div1"
        2       [10]    MOVE            3 0
        3       [10]    MOVE            4 1
        4       [10]    TAILCALL        2 3 0
        5       [10]    RETURN          2 0
        6       [10]    RETURN          0 1
constants (1) for 0000000000c798f0:
        1       "div1"
locals (2) for 0000000000c798f0:
        0       a       1       7
        1       b       1       7
upvalues (1) for 0000000000c798f0:
        0       _ENV    0       0
]]--