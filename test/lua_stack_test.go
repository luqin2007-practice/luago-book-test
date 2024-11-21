package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-luacompiler/api"
	"go-luacompiler/state"
	"strings"
	"testing"
)

func TestLuaStack(t *testing.T) {
	ls := state.New(20, nil)

	/*
	   true
	*/
	ls.PushBoolean(true)
	assert.Equal(t, "[true]", printStack(ls))

	/*
	   10
	   true
	*/
	ls.PushInteger(10)
	assert.Equal(t, "[true][10]", printStack(ls))

	/*
	   nil
	   10
	   true
	*/
	ls.PushNil()
	assert.Equal(t, "[true][10][nil]", printStack(ls))

	/*
	   hello
	   nil
	   10
	   true
	*/
	ls.PushString("hello")
	assert.Equal(t, "[true][10][nil][\"hello\"]", printStack(ls))

	/*
	         | true
	   hello | hello
	   nil   | nil
	   10    | 10
	   true  | true
	*/
	ls.PushValue(-4)
	assert.Equal(t, "[true][10][nil][\"hello\"][true]", printStack(ls))

	/*
	   true  |
	   hello | hello
	   nil   | true
	   10    | 10
	   true  | true
	*/
	ls.Replace(3)
	assert.Equal(t, "[true][10][true][\"hello\"]", printStack(ls))

	/*
	         | nil
	         | nil
	   hello | hello
	   nil   | true
	   10    | 10
	   true  | true
	*/
	ls.SetTop(6)
	assert.Equal(t, "[true][10][true][\"hello\"][nil][nil]", printStack(ls))

	/*
	   nil   |
	   nil   | nil
	   hello | nil
	   true  | true
	   10    | 10
	   true  | true
	*/
	ls.Remove(-3)
	assert.Equal(t, "[true][10][true][nil][nil]", printStack(ls))

	/*
	   nil  |
	   nil  |
	   true |
	   10   |
	   true | true
	*/
	ls.SetTop(-5)
	assert.Equal(t, "[true]", printStack(ls))
}

func printStack(ls api.LuaState) string {
	top := ls.GetTop()
	builder := strings.Builder{}
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			builder.WriteString(fmt.Sprintf("[%t]", ls.ToBoolean(i)))
		case api.LUA_TNUMBER:
			if ls.IsInteger(i) {
				builder.WriteString(fmt.Sprintf("[%d]", ls.ToInteger(i)))
			} else {
				builder.WriteString(fmt.Sprintf("[%f]", ls.ToNumber(i)))
			}
		case api.LUA_TSTRING:
			builder.WriteString(fmt.Sprintf("[%q]", ls.ToString(i)))
		default:
			builder.WriteString(fmt.Sprintf("[%s]", ls.TypeName(t)))
		}
	}
	str := builder.String()
	fmt.Println(str)
	return str
}
