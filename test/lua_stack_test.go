package test

import (
	"github.com/stretchr/testify/assert"
	"go-luacompiler/state"
	"testing"
)

func TestLuaStack(t *testing.T) {
	ls := state.New()

	/*
	   true
	*/
	ls.PushBoolean(true)
	assert.Equal(t, "[true]", PrintStack(ls))

	/*
	   10
	   true
	*/
	ls.PushInteger(10)
	assert.Equal(t, "[true][10]", PrintStack(ls))

	/*
	   nil
	   10
	   true
	*/
	ls.PushNil()
	assert.Equal(t, "[true][10][nil]", PrintStack(ls))

	/*
	   hello
	   nil
	   10
	   true
	*/
	ls.PushString("hello")
	assert.Equal(t, "[true][10][nil][\"hello\"]", PrintStack(ls))

	/*
	         | true
	   hello | hello
	   nil   | nil
	   10    | 10
	   true  | true
	*/
	ls.PushValue(-4)
	assert.Equal(t, "[true][10][nil][\"hello\"][true]", PrintStack(ls))

	/*
	   true  |
	   hello | hello
	   nil   | true
	   10    | 10
	   true  | true
	*/
	ls.Replace(3)
	assert.Equal(t, "[true][10][true][\"hello\"]", PrintStack(ls))

	/*
	         | nil
	         | nil
	   hello | hello
	   nil   | true
	   10    | 10
	   true  | true
	*/
	ls.SetTop(6)
	assert.Equal(t, "[true][10][true][\"hello\"][nil][nil]", PrintStack(ls))

	/*
	   nil   |
	   nil   | nil
	   hello | nil
	   true  | true
	   10    | 10
	   true  | true
	*/
	ls.Remove(-3)
	assert.Equal(t, "[true][10][true][nil][nil]", PrintStack(ls))

	/*
	   nil  |
	   nil  |
	   true |
	   10   |
	   true | true
	*/
	ls.SetTop(-5)
	assert.Equal(t, "[true]", PrintStack(ls))
}
