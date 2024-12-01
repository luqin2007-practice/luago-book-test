local mt = {}

function vector(x, y)
    local v = { x = x, y = y }
    setmetatable(v, mt)
    return v
end

mt.__add = function(v1, v2)
    return vector(v1.x + v2.x, v1.y + v2.y)
end

mt.__sub = function(v1, v2)
    return vector(v1.x - v2.x, v1.y - v2.y)
end

mt.__mul = function(v, n)
    return vector(v.x * n, v.y * n)
end

mt.__div = function(v, n)
    return vector(v.x / n, v.y / n)
end

mt.__len = function(v)
    return (v.x * v.x + v.y * v.y) ^ 0.5
end

mt.__eq = function(v1, v2)
    return v1.x == v2.x and v1.y == v2.y
end

mt.__index = function(v, k)
    if k == "print" then
        return function()
            print("[" .. v.x .. ", " .. v.y .. "]")
        end
    end
end

mt.__call = function(v)
    print("[" .. v.x .. ", " .. v.y .. "]")
end

v1 = vector(1, 2); v1:print()
v2 = vector(3, 4); v2:print()
v3 = v1 * 2;       v3:print()
v4 = v1 + v3;      v4:print()
print(#v2)
print(v1 == v2)
print(v2 == vector(3, 4))
v4()

--[[
[1, 2]
[3, 4]
[2, 4]
[3, 6]
5.0
false
true
[3, 6]
]]--

--[[
main <.\metatable.lua:0,0> (74 instructions at 00000000007d9220)
0+ params, 6 slots, 1 upvalue, 1 local, 18 constants, 9 functions
        1       [1]     NEWTABLE        0 0 0
        2       [7]     CLOSURE         1 0     ; 00000000007d9460
        3       [3]     SETTABUP        0 -1 1  ; _ENV "vector"
        4       [11]    CLOSURE         1 1     ; 00000000007d9670
        5       [11]    SETTABLE        0 -2 1  ; "__add" -
        6       [15]    CLOSURE         1 2     ; 00000000007d9a40
        7       [15]    SETTABLE        0 -3 1  ; "__sub" -
        8       [19]    CLOSURE         1 3     ; 00000000007da060
        9       [19]    SETTABLE        0 -4 1  ; "__mul" -
        10      [23]    CLOSURE         1 4     ; 00000000007d9780
        11      [23]    SETTABLE        0 -5 1  ; "__div" -
        12      [27]    CLOSURE         1 5     ; 00000000007da680
        13      [27]    SETTABLE        0 -6 1  ; "__len" -
        14      [31]    CLOSURE         1 6     ; 00000000007da840
        15      [31]    SETTABLE        0 -7 1  ; "__eq" -
        16      [39]    CLOSURE         1 7     ; 00000000007daa00
        17      [39]    SETTABLE        0 -8 1  ; "__index" -
        18      [43]    CLOSURE         1 8     ; 00000000007da1d0
        19      [43]    SETTABLE        0 -9 1  ; "__call" -
        20      [45]    GETTABUP        1 0 -1  ; _ENV "vector"
        21      [45]    LOADK           2 -11   ; 1
        22      [45]    LOADK           3 -12   ; 2
        23      [45]    CALL            1 3 2
        24      [45]    SETTABUP        0 -10 1 ; _ENV "v1"
        25      [45]    GETTABUP        1 0 -10 ; _ENV "v1"
        26      [45]    SELF            1 1 -13 ; "print"
        27      [45]    CALL            1 2 1
        28      [46]    GETTABUP        1 0 -1  ; _ENV "vector"
        29      [46]    LOADK           2 -15   ; 3
        30      [46]    LOADK           3 -16   ; 4
        31      [46]    CALL            1 3 2
        32      [46]    SETTABUP        0 -14 1 ; _ENV "v2"
        33      [46]    GETTABUP        1 0 -14 ; _ENV "v2"
        34      [46]    SELF            1 1 -13 ; "print"
        35      [46]    CALL            1 2 1
        36      [47]    GETTABUP        1 0 -10 ; _ENV "v1"
        37      [47]    MUL             1 1 -12 ; - 2
        38      [47]    SETTABUP        0 -17 1 ; _ENV "v3"
        39      [47]    GETTABUP        1 0 -17 ; _ENV "v3"
        40      [47]    SELF            1 1 -13 ; "print"
        41      [47]    CALL            1 2 1
        42      [48]    GETTABUP        1 0 -10 ; _ENV "v1"
        43      [48]    GETTABUP        2 0 -17 ; _ENV "v3"
        44      [48]    ADD             1 1 2
        45      [48]    SETTABUP        0 -18 1 ; _ENV "v4"
        46      [48]    GETTABUP        1 0 -18 ; _ENV "v4"
        47      [48]    SELF            1 1 -13 ; "print"
        48      [48]    CALL            1 2 1
        49      [49]    GETTABUP        1 0 -13 ; _ENV "print"
        50      [49]    GETTABUP        2 0 -14 ; _ENV "v2"
        51      [49]    LEN             2 2
        52      [49]    CALL            1 2 1
        53      [50]    GETTABUP        1 0 -13 ; _ENV "print"
        54      [50]    GETTABUP        2 0 -10 ; _ENV "v1"
        55      [50]    GETTABUP        3 0 -14 ; _ENV "v2"
        56      [50]    EQ              1 2 3
        57      [50]    JMP             0 1     ; to 59
        58      [50]    LOADBOOL        2 0 1
        59      [50]    LOADBOOL        2 1 0
        60      [50]    CALL            1 2 1
        61      [51]    GETTABUP        1 0 -13 ; _ENV "print"
        62      [51]    GETTABUP        2 0 -14 ; _ENV "v2"
        63      [51]    GETTABUP        3 0 -1  ; _ENV "vector"
        64      [51]    LOADK           4 -15   ; 3
        65      [51]    LOADK           5 -16   ; 4
        66      [51]    CALL            3 3 2
        67      [51]    EQ              1 2 3
        68      [51]    JMP             0 1     ; to 70
        69      [51]    LOADBOOL        2 0 1
        70      [51]    LOADBOOL        2 1 0
        71      [51]    CALL            1 2 1
        72      [52]    GETTABUP        1 0 -18 ; _ENV "v4"
        73      [52]    CALL            1 1 1
        74      [52]    RETURN          0 1
constants (18) for 00000000007d9220:
        1       "vector"
        2       "__add"
        3       "__sub"
        4       "__mul"
        5       "__div"
        6       "__len"
        7       "__eq"
        8       "__index"
        9       "__call"
        10      "v1"
        11      1
        12      2
        13      "print"
        14      "v2"
        15      3
        16      4
        17      "v3"
        18      "v4"
locals (1) for 00000000007d9220:
        0       mt      2       75
upvalues (1) for 00000000007d9220:
        0       _ENV    1       0

function <.\metatable.lua:3,7> (9 instructions at 00000000007d9460)
2 params, 6 slots, 2 upvalues, 3 locals, 3 constants, 0 functions
        1       [4]     NEWTABLE        2 0 2
        2       [4]     SETTABLE        2 -1 0  ; "x" -
        3       [4]     SETTABLE        2 -2 1  ; "y" -
        4       [5]     GETTABUP        3 0 -3  ; _ENV "setmetatable"
        5       [5]     MOVE            4 2
        6       [5]     GETUPVAL        5 1     ; mt
        7       [5]     CALL            3 3 1
        8       [6]     RETURN          2 2
        9       [7]     RETURN          0 1
constants (3) for 00000000007d9460:
        1       "x"
        2       "y"
        3       "setmetatable"
locals (3) for 00000000007d9460:
        0       x       1       10
        1       y       1       10
        2       v       4       10
upvalues (2) for 00000000007d9460:
        0       _ENV    0       0
        1       mt      1       0

function <.\metatable.lua:9,11> (10 instructions at 00000000007d9670)
2 params, 6 slots, 1 upvalue, 2 locals, 3 constants, 0 functions
        1       [10]    GETTABUP        2 0 -1  ; _ENV "vector"
        2       [10]    GETTABLE        3 0 -2  ; "x"
        3       [10]    GETTABLE        4 1 -2  ; "x"
        4       [10]    ADD             3 3 4
        5       [10]    GETTABLE        4 0 -3  ; "y"
        6       [10]    GETTABLE        5 1 -3  ; "y"
        7       [10]    ADD             4 4 5
        8       [10]    TAILCALL        2 3 0
        9       [10]    RETURN          2 0
        10      [11]    RETURN          0 1
constants (3) for 00000000007d9670:
        1       "vector"
        2       "x"
        3       "y"
locals (2) for 00000000007d9670:
        0       v1      1       11
        1       v2      1       11
upvalues (1) for 00000000007d9670:
        0       _ENV    0       0

function <.\metatable.lua:13,15> (10 instructions at 00000000007d9a40)
2 params, 6 slots, 1 upvalue, 2 locals, 3 constants, 0 functions
        1       [14]    GETTABUP        2 0 -1  ; _ENV "vector"
        2       [14]    GETTABLE        3 0 -2  ; "x"
        3       [14]    GETTABLE        4 1 -2  ; "x"
        4       [14]    SUB             3 3 4
        5       [14]    GETTABLE        4 0 -3  ; "y"
        6       [14]    GETTABLE        5 1 -3  ; "y"
        7       [14]    SUB             4 4 5
        8       [14]    TAILCALL        2 3 0
        9       [14]    RETURN          2 0
        10      [15]    RETURN          0 1
constants (3) for 00000000007d9a40:
        1       "vector"
        2       "x"
        3       "y"
locals (2) for 00000000007d9a40:
        0       v1      1       11
        1       v2      1       11
upvalues (1) for 00000000007d9a40:
        0       _ENV    0       0

function <.\metatable.lua:17,19> (8 instructions at 00000000007da060)
2 params, 5 slots, 1 upvalue, 2 locals, 3 constants, 0 functions
        1       [18]    GETTABUP        2 0 -1  ; _ENV "vector"
        2       [18]    GETTABLE        3 0 -2  ; "x"
        3       [18]    MUL             3 3 1
        4       [18]    GETTABLE        4 0 -3  ; "y"
        5       [18]    MUL             4 4 1
        6       [18]    TAILCALL        2 3 0
        7       [18]    RETURN          2 0
        8       [19]    RETURN          0 1
constants (3) for 00000000007da060:
        1       "vector"
        2       "x"
        3       "y"
locals (2) for 00000000007da060:
        0       v       1       9
        1       n       1       9
upvalues (1) for 00000000007da060:
        0       _ENV    0       0

function <.\metatable.lua:21,23> (8 instructions at 00000000007d9780)
2 params, 5 slots, 1 upvalue, 2 locals, 3 constants, 0 functions
        1       [22]    GETTABUP        2 0 -1  ; _ENV "vector"
        2       [22]    GETTABLE        3 0 -2  ; "x"
        3       [22]    DIV             3 3 1
        4       [22]    GETTABLE        4 0 -3  ; "y"
        5       [22]    DIV             4 4 1
        6       [22]    TAILCALL        2 3 0
        7       [22]    RETURN          2 0
        8       [23]    RETURN          0 1
constants (3) for 00000000007d9780:
        1       "vector"
        2       "x"
        3       "y"
locals (2) for 00000000007d9780:
        0       v       1       9
        1       n       1       9
upvalues (1) for 00000000007d9780:
        0       _ENV    0       0

function <.\metatable.lua:25,27> (10 instructions at 00000000007da680)
1 param, 4 slots, 0 upvalues, 1 local, 3 constants, 0 functions
        1       [26]    GETTABLE        1 0 -1  ; "x"
        2       [26]    GETTABLE        2 0 -1  ; "x"
        3       [26]    MUL             1 1 2
        4       [26]    GETTABLE        2 0 -2  ; "y"
        5       [26]    GETTABLE        3 0 -2  ; "y"
        6       [26]    MUL             2 2 3
        7       [26]    ADD             1 1 2
        8       [26]    POW             1 1 -3  ; - 0.5
        9       [26]    RETURN          1 2
        10      [27]    RETURN          0 1
constants (3) for 00000000007da680:
        1       "x"
        2       "y"
        3       0.5
locals (1) for 00000000007da680:
        0       v       1       11
upvalues (0) for 00000000007da680:

function <.\metatable.lua:29,31> (12 instructions at 00000000007da840)
2 params, 4 slots, 0 upvalues, 2 locals, 2 constants, 0 functions
        1       [30]    GETTABLE        2 0 -1  ; "x"
        2       [30]    GETTABLE        3 1 -1  ; "x"
        3       [30]    EQ              0 2 3
        4       [30]    JMP             0 4     ; to 9
        5       [30]    GETTABLE        2 0 -2  ; "y"
        6       [30]    GETTABLE        3 1 -2  ; "y"
        7       [30]    EQ              1 2 3
        8       [30]    JMP             0 1     ; to 10
        9       [30]    LOADBOOL        2 0 1
        10      [30]    LOADBOOL        2 1 0
        11      [30]    RETURN          2 2
        12      [31]    RETURN          0 1
constants (2) for 00000000007da840:
        1       "x"
        2       "y"
locals (2) for 00000000007da840:
        0       v1      1       13
        1       v2      1       13
upvalues (0) for 00000000007da840:

function <.\metatable.lua:33,39> (5 instructions at 00000000007daa00)
2 params, 3 slots, 1 upvalue, 2 locals, 1 constant, 1 function
        1       [34]    EQ              0 1 -1  ; - "print"
        2       [34]    JMP             0 2     ; to 5
        3       [37]    CLOSURE         2 0     ; 00000000007dab20
        4       [37]    RETURN          2 2
        5       [39]    RETURN          0 1
constants (1) for 00000000007daa00:
        1       "print"
locals (2) for 00000000007daa00:
        0       v       1       6
        1       k       1       6
upvalues (1) for 00000000007daa00:
        0       _ENV    0       0

function <.\metatable.lua:35,37> (9 instructions at 00000000007dab20)
0 params, 6 slots, 2 upvalues, 0 locals, 6 constants, 0 functions
        1       [36]    GETTABUP        0 0 -1  ; _ENV "print"
        2       [36]    LOADK           1 -2    ; "["
        3       [36]    GETTABUP        2 1 -3  ; v "x"
        4       [36]    LOADK           3 -4    ; ", "
        5       [36]    GETTABUP        4 1 -5  ; v "y"
        6       [36]    LOADK           5 -6    ; "]"
        7       [36]    CONCAT          1 1 5
        8       [36]    CALL            0 2 1
        9       [37]    RETURN          0 1
constants (6) for 00000000007dab20:
        1       "print"
        2       "["
        3       "x"
        4       ", "
        5       "y"
        6       "]"
locals (0) for 00000000007dab20:
upvalues (2) for 00000000007dab20:
        0       _ENV    0       0
        1       v       1       0

function <.\metatable.lua:41,43> (9 instructions at 00000000007da1d0)
1 param, 7 slots, 1 upvalue, 1 local, 6 constants, 0 functions
        1       [42]    GETTABUP        1 0 -1  ; _ENV "print"
        2       [42]    LOADK           2 -2    ; "["
        3       [42]    GETTABLE        3 0 -3  ; "x"
        4       [42]    LOADK           4 -4    ; ", "
        5       [42]    GETTABLE        5 0 -5  ; "y"
        6       [42]    LOADK           6 -6    ; "]"
        7       [42]    CONCAT          2 2 6
        8       [42]    CALL            1 2 1
        9       [43]    RETURN          0 1
constants (6) for 00000000007da1d0:
        1       "print"
        2       "["
        3       "x"
        4       ", "
        5       "y"
        6       "]"
locals (1) for 00000000007da1d0:
        0       v       1       10
upvalues (1) for 00000000007da1d0:
        0       _ENV    0       0
]]--