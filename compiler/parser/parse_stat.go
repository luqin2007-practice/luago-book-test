package parser

import (
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

func parseStat(lexer *Lexer) ast.Stat {
	switch lexer.LookAhead() {
	case TOKEN_SEP_SEMI: // ;
		return parseEmptyStat(lexer)
	case TOKEN_KW_BREAK: // break
		return parseBreakStat(lexer)
	case TOKEN_SEP_LABEL: // :: ...
		return parseLabelStat(lexer)
	case TOKEN_KW_GOTO: // goto ...
		return parseGotoStat(lexer)
	case TOKEN_KW_DO: // do ...
		return parseDoStat(lexer)
	case TOKEN_KW_WHILE: // while ...
		return parseWhileStat(lexer)
	case TOKEN_KW_REPEAT: // repeat ...
		return parseRepeatStat(lexer)
	case TOKEN_KW_IF: // if ...
		return parseIfStat(lexer)
	case TOKEN_KW_FOR: // for ...
		return parseForStat(lexer)
	case TOKEN_KW_FUNCTION: // function ...
		return parseFuncDefStat(lexer)
	case TOKEN_KW_LOCAL: // local ...
		return parseLocalAssignOrFuncDefStat(lexer)
	default: // 无法通过前瞻一次获取
		return parseAssignOrFuncCallStat(lexer)
	}
}

func parseEmptyStat(lexer *Lexer) *ast.EmptyStat {
	lexer.NextToken() // ;
	return &ast.EmptyStat{}
}

func parseBreakStat(lexer *Lexer) *ast.BreakStat {
	lexer.NextToken() // break
	return &ast.BreakStat{}
}

func parseLabelStat(lexer *Lexer) *ast.LabelStat {
	lexer.NextToken()                 // ::
	_, name := lexer.NextIdentifier() // Name
	return &ast.LabelStat{Name: name}
}

func parseGotoStat(lexer *Lexer) *ast.GotoStat {
	lexer.NextToken()                 // 'goto'
	_, name := lexer.NextIdentifier() // Name
	return &ast.GotoStat{Name: name}
}

func parseDoStat(lexer *Lexer) *ast.DoStat {
	lexer.NextToken()               // 'do'
	block := parseBlock(lexer)      // block
	lexer.NextKeyword(TOKEN_KW_END) // 'end'
	return &ast.DoStat{Block: block}
}

func parseWhileStat(lexer *Lexer) *ast.WhileStat {
	lexer.NextToken()               // 'while'
	exp := parseExp(lexer)          // exp
	lexer.NextKeyword(TOKEN_KW_DO)  // 'do'
	block := parseBlock(lexer)      // block
	lexer.NextKeyword(TOKEN_KW_END) // 'end'
	return &ast.WhileStat{Exp: exp, Block: block}
}

func parseRepeatStat(lexer *Lexer) *ast.RepeatStat {
	lexer.NextToken()                 // 'repeat'
	block := parseBlock(lexer)        // block
	lexer.NextKeyword(TOKEN_KW_UNTIL) // 'until'
	exp := parseExp(lexer)            // exp
	return &ast.RepeatStat{Block: block, Exp: exp}
}

func parseIfStat(lexer *Lexer) *ast.IfStat {
	exps := make([]ast.Exp, 4)
	blocks := make([]*ast.Block, 4)

	lexer.NextToken()                          // 'if'
	exps = append(exps, parseExp(lexer))       // exp
	lexer.NextKeyword(TOKEN_KW_THEN)           // 'then'
	blocks = append(blocks, parseBlock(lexer)) // block
	for lexer.LookAhead() == TOKEN_KW_ELSEIF {
		lexer.NextToken()                          // 'elseif'
		exps = append(exps, parseExp(lexer))       // exp
		lexer.NextKeyword(TOKEN_KW_THEN)           // 'then'
		blocks = append(blocks, parseBlock(lexer)) // block
	}
	if lexer.LookAhead() == TOKEN_KW_ELSE {
		lexer.NextToken()                          // 'else'
		blocks = append(blocks, parseBlock(lexer)) // block
	}
	lexer.NextKeyword(TOKEN_KW_END) // 'end'
	return &ast.IfStat{Exps: exps, Blocks: blocks}
}

func parseForStat(lexer *Lexer) ast.Stat {
	forLine, _, _ := lexer.NextToken() // 'for'
	_, name := lexer.NextIdentifier()  // Name
	if lexer.LookAhead() == TOKEN_OP_ASSIGN {
		return parseForNumStat(lexer, forLine, name)
	} else {
		return parseForInStat(lexer, name)
	}
}

func parseFuncDefStat(lexer *Lexer) *ast.AssignStat {
	lexer.NextToken()                          // 'function'
	nameList, hasColon := parseFuncName(lexer) // funcname
	funcDefExp := parseFuncDefExp(lexer)       // funcbody

	// 解语法糖 a:b(...) -> a.b(self, ...)
	if hasColon {
		// 在 ParList 开头添加一个 self
		funcDefExp.ParList = append(funcDefExp.ParList, "")
		copy(funcDefExp.ParList[1:], funcDefExp.ParList)
		funcDefExp.ParList[0] = "self"
	}

	// 解语法糖 function a(...){...} -> a=function(...){...}
	return &ast.AssignStat{
		LastLine: lexer.Line(),
		VarList:  []ast.Exp{nameList},
		ExpList:  []ast.Exp{funcDefExp},
	}
}

func parseLocalAssignOrFuncDefStat(lexer *Lexer) ast.Stat {
	lexer.NextToken()
	if lexer.LookAhead() == TOKEN_KW_FUNCTION {
		return parseLocalFuncDefStat(lexer) // 'local' 'function' ...
	} else {
		return parseLocalVarDeclStat(lexer) // 'local' namelist ...
	}
}

/*
prefixexp ::= var | functioncall | '(' exp ')'

functioncall ::= prefixexp args | prefixexp ':' Name args

varlist '=' explist
varlist ::= var {',' var}
var ::= Name | prefixexp '[' exp ']' | prefixexp '.' Name
*/
func parseAssignOrFuncCallStat(lexer *Lexer) ast.Stat {
	prefixExp := parsePrefixExp(lexer)
	if fc, ok := prefixExp.(*ast.FuncCallExp); ok {
		// functioncall
		return fc
	} else {
		// varlist '=' explist
		return parseAssignStat(lexer, prefixExp)
	}
}
