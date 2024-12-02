t = {a = 1, b = 2, c = 3}
for k, v in pairs(t) do
  print(k, v)
end

t = {"a", "b", "c"}
for k, v in ipairs(t) do
  print(k, v)
end

--[[
main <./iterator.lua:0,0> (32 instructions at 0000000000ad91e0)
0+ params, 8 slots, 1 upvalue, 10 locals, 10 constants, 0 functions
        1       [1]     NEWTABLE        0 0 3
        2       [1]     SETTABLE        0 -2 -3 ; "a" 1
        3       [1]     SETTABLE        0 -4 -5 ; "b" 2
        4       [1]     SETTABLE        0 -6 -7 ; "c" 3
        5       [1]     SETTABUP        0 -1 0  ; _ENV "t"
        6       [2]     GETTABUP        0 0 -8  ; _ENV "pairs"
        7       [2]     GETTABUP        1 0 -1  ; _ENV "t"
        8       [2]     CALL            0 2 4
        9       [2]     JMP             0 4     ; to 14
        10      [3]     GETTABUP        5 0 -9  ; _ENV "print"
        11      [3]     MOVE            6 3
        12      [3]     MOVE            7 4
        13      [3]     CALL            5 3 1
        14      [2]     TFORCALL        0 2
        15      [2]     TFORLOOP        2 -6    ; to 10
        16      [6]     NEWTABLE        0 3 0
        17      [6]     LOADK           1 -2    ; "a"
        18      [6]     LOADK           2 -4    ; "b"
        19      [6]     LOADK           3 -6    ; "c"
        20      [6]     SETLIST         0 3 1   ; 1
        21      [6]     SETTABUP        0 -1 0  ; _ENV "t"
        22      [7]     GETTABUP        0 0 -10 ; _ENV "ipairs"
        23      [7]     GETTABUP        1 0 -1  ; _ENV "t"
        24      [7]     CALL            0 2 4
        25      [7]     JMP             0 4     ; to 30
        26      [8]     GETTABUP        5 0 -9  ; _ENV "print"
        27      [8]     MOVE            6 3
        28      [8]     MOVE            7 4
        29      [8]     CALL            5 3 1
        30      [7]     TFORCALL        0 2
        31      [7]     TFORLOOP        2 -6    ; to 26
        32      [9]     RETURN          0 1
constants (10) for 0000000000ad91e0:
        1       "t"
        2       "a"
        3       1
        4       "b"
        5       2
        6       "c"
        7       3
        8       "pairs"
        9       "print"
        10      "ipairs"
locals (10) for 0000000000ad91e0:
        0       (for generator) 9       16
        1       (for state)     9       16
        2       (for control)   9       16
        3       k       10      14
        4       v       10      14
        5       (for generator) 25      32
        6       (for state)     25      32
        7       (for control)   25      32
        8       k       26      30
        9       v       26      30
upvalues (1) for 0000000000ad91e0:
        0       _ENV    1       0
]]--