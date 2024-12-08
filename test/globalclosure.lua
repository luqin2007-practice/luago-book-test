local function f()
    local function g()
        x = y
    end
end

--[[
main <.\globalclosure.lua:0,0> (2 instructions at 0000000000799110)
0+ params, 2 slots, 1 upvalue, 1 local, 0 constants, 1 function
        1       [5]     CLOSURE         0 0     ; 00000000007992c0
        2       [5]     RETURN          0 1
constants (0) for 0000000000799110:
locals (1) for 0000000000799110:
        0       f       2       3
upvalues (1) for 0000000000799110:
        0       _ENV    1       0

function <.\globalclosure.lua:1,5> (2 instructions at 00000000007992c0)
0 params, 2 slots, 1 upvalue, 1 local, 0 constants, 1 function
        1       [4]     CLOSURE         0 0     ; 0000000000799390
        2       [5]     RETURN          0 1
constants (0) for 00000000007992c0:
locals (1) for 00000000007992c0:
        0       g       2       3
upvalues (1) for 00000000007992c0:
        0       _ENV    0       0

function <.\globalclosure.lua:2,4> (3 instructions at 0000000000799390)
0 params, 2 slots, 1 upvalue, 0 locals, 2 constants, 0 functions
        1       [3]     GETTABUP        0 0 -2  ; _ENV "y"
        2       [3]     SETTABUP        0 -1 0  ; _ENV "x"
        3       [4]     RETURN          0 1
constants (2) for 0000000000799390:
        1       "x"
        2       "y"
locals (0) for 0000000000799390:
upvalues (1) for 0000000000799390:
        0       _ENV    0       0
]]--