local u, v, w

local function f()
    -- 捕获 u, v 供 g() 捕获
    local function g()
        u = v -- 捕获 u, v
    end
end

--[[
main <.\flatclosure.lua:0,0> (3 instructions at 00000000006d9110)
0+ params, 4 slots, 1 upvalue, 4 locals, 0 constants, 1 function
        1       [1]     LOADNIL         0 2
        2       [8]     CLOSURE         3 0     ; 00000000006d9280
        3       [8]     RETURN          0 1
constants (0) for 00000000006d9110:
locals (4) for 00000000006d9110:
        0       u       2       4
        1       v       2       4
        2       w       2       4
        3       f       3       4
upvalues (1) for 00000000006d9110:
        0       _ENV    1       0

function <.\flatclosure.lua:3,8> (2 instructions at 00000000006d9280)
0 params, 2 slots, 2 upvalues, 1 local, 0 constants, 1 function
        1       [7]     CLOSURE         0 0     ; 00000000006d9470
        2       [8]     RETURN          0 1
constants (0) for 00000000006d9280:
locals (1) for 00000000006d9280:
        0       g       2       3
upvalues (2) for 00000000006d9280:
        0       u       1       0
        1       v       1       1

function <.\flatclosure.lua:5,7> (3 instructions at 00000000006d9470)
0 params, 2 slots, 2 upvalues, 0 locals, 0 constants, 0 functions
        1       [6]     GETUPVAL        0 1     ; v
        2       [6]     SETUPVAL        0 0     ; u
        3       [7]     RETURN          0 1
constants (0) for 00000000006d9470:
locals (0) for 00000000006d9470:
upvalues (2) for 00000000006d9470:
        0       u       0       0
        1       v       0       1
]]--