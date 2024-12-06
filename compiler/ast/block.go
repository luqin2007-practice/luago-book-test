package ast

/*
Block chunk 与 block 结构

EBNF 表示：

	chunk ::= block
	block ::= {stat} [retstat]
	retstat ::= return [explist] [';']
	explist ::= exp {',' exp}
*/
type Block struct {
	LastLine int
	Stats    []Stat
	RetExp   []Exp
}
