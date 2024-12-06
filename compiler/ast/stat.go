package ast

/*
Stat

	stat ::= ';'
	    | varlist '=' explist
	    | 'break'
	    | '::' Name
	    | 'goto' Name
	    | 'do' block 'end'
	    | 'while' exp 'do' block 'end'
	    | 'repeat' block 'until' exp
	    | 'if' exp 'then' block {'elseif' exp 'then' block} ['else' block] 'end'
	    | 'for' Name '=' exp ',' exp [',' exp] 'do' block 'end'
	    | 'for' namelist 'in' explist 'do' block 'end'
	    | 'function' funcname funcbody
	    | 'local' function Name funcbody
	    | 'local' name ['=' explist]
	    | functioncall
*/
type Stat interface{}

type EmptyStat struct{}
type BreakStat struct{ Line int }
type LabelStat struct{ Name string }
type GotoStat struct{ Name string }
type DoStat struct{ Block *Block }
type FuncCallStat = FuncCallExp

type WhileStat struct {
	Exp   Exp
	Block *Block
}

type RepeatStat struct {
	Block *Block
	Exp   Exp
}

type IfStat struct {
	Exps   []Exp
	Blocks []*Block
}

type ForNumStat struct {
	LineOfFor int
	LineOfDo  int

	VarName  string
	InitExp  Exp
	LimitExp Exp
	StepExp  Exp
	Block    *Block
}

type ForInStat struct {
	LineOfDo int

	NameList []string
	ExpList  []Exp
	Block    *Block
}

type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}

type AssignStat struct {
	LastLine int
	VarList  []Exp
	ExpList  []Exp
}

type LocalFuncDefStat struct {
	Name string
	Exp  *FuncDefExp
}
