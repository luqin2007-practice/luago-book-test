local u, v, w

function f()
    u = v
end

--[[
main <.\upvalue.lua:0,0> (4 instructions at 0000000000bc90d0)
0+ params, 4 slots, 1 upvalue, 3 locals, 1 constant, 1 function
        1       [1]     LOADNIL         0 2
        2       [5]     CLOSURE         3 0     ; 0000000000bc93e0
        3       [3]     SETTABUP        0 -1 3  ; _ENV "f"
        4       [5]     RETURN          0 1
constants (1) for 0000000000bc90d0:
        1       "f"
locals (3) for 0000000000bc90d0:
        0       u       2       5
        1       v       2       5
        2       w       2       5
upvalues (1) for 0000000000bc90d0:
        0       _ENV    1       0

function <.\upvalue.lua:3,5> (3 instructions at 0000000000bc93e0)
0 params, 2 slots, 2 upvalues, 0 locals, 0 constants, 0 functions
        1       [4]     GETUPVAL        0 1     ; v
        2       [4]     SETUPVAL        0 0     ; u
        3       [5]     RETURN          0 1
constants (0) for 0000000000bc93e0:
locals (0) for 0000000000bc93e0:
upvalues (2) for 0000000000bc93e0:
        0       u       1       0
        1       v       1       1
]]--