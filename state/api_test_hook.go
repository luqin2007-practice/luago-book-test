package state

var beforeExec func(i uint32)
var afterExec func(i uint32)

func (_ *luaState) OnBeforeInstExecuted(f func(i uint32)) {
	beforeExec = f
}

func (_ *luaState) BeforeInstExecuted(i uint32) {
	if beforeExec != nil {
		beforeExec(i)
	}
}

func (_ *luaState) OnAfterInstExecuted(f func(i uint32)) {
	afterExec = f
}

func (_ *luaState) AfterInstExecuted(i uint32) {
	if afterExec != nil {
		afterExec(i)
	}
}
