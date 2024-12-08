package parser

import (
	"fmt"
	"go-luacompiler/compiler/ast"
	. "go-luacompiler/compiler/lexer"
)

func parseForNumStat(lexer *Lexer, forLine int, varName string) *ast.ForNumStat {
	lexer.NextToken()                  // =
	initExp := parseExp(lexer)         // exp
	lexer.NextKeyword(TOKEN_SEP_COMMA) // ,
	limitExp := parseExp(lexer)        // exp
	var stepExp ast.Exp
	if lexer.LookAhead() == TOKEN_SEP_COMMA { // [?]
		lexer.NextToken()         // ,
		stepExp = parseExp(lexer) // exp
	} else { // [? 1]
		stepExp = &ast.IntegerExp{Line: lexer.Line(), Val: 1}
	}
	doLine := lexer.NextKeyword(TOKEN_KW_DO) // 'do'
	block := parseBlock(lexer)               // block
	lexer.NextKeyword(TOKEN_KW_END)          // 'end'
	return &ast.ForNumStat{
		LineOfFor: forLine,
		LineOfDo:  doLine,
		VarName:   varName,
		InitExp:   initExp,
		LimitExp:  limitExp,
		StepExp:   stepExp,
		Block:     block,
	}
}

func parseForInStat(lexer *Lexer, firstName string) *ast.ForInStat {
	nameList := parseNameList(lexer, firstName) // namelist
	lexer.NextKeyword(TOKEN_KW_IN)              // 'in'
	explist := parseExpList(lexer)              // explist
	doLine := lexer.NextKeyword(TOKEN_KW_DO)    // 'do'
	block := parseBlock(lexer)                  // block
	lexer.NextKeyword(TOKEN_KW_END)             // 'end'
	return &ast.ForInStat{
		LineOfDo: doLine,
		NameList: nameList,
		ExpList:  explist,
		Block:    block,
	}
}

func parseLocalFuncDefStat(lexer *Lexer) ast.LocalFuncDefStat {
	lexer.NextToken()                    // function
	line, name := lexer.NextIdentifier() // Name
	exp := parseFuncDefExp(lexer, line)  // funcbody
	return ast.LocalFuncDefStat{Name: name, Exp: exp}
}

func parseLocalVarDeclStat(lexer *Lexer) *ast.LocalVarDeclStat {
	_, firstName := lexer.NextIdentifier()
	nameList := parseNameList(lexer, firstName) // namelist
	var expList []ast.Exp
	if lexer.LookAhead() == TOKEN_OP_ASSIGN { // [*]
		lexer.NextToken()             // =
		expList = parseExpList(lexer) // explist
	} else {
		expList = make([]ast.Exp, 0)
	}
	return &ast.LocalVarDeclStat{
		LastLine: lexer.Line(),
		NameList: nameList,
		ExpList:  expList,
	}
}

func parseAssignStat(lexer *Lexer, firstExp ast.Exp) *ast.AssignStat {
	varList := parseVarList(lexer, firstExp) // var {',' var}
	lexer.NextKeyword(TOKEN_OP_ASSIGN)       // =
	expList := parseExpList(lexer)           // explist
	return &ast.AssignStat{
		LastLine: lexer.Line(),
		VarList:  varList,
		ExpList:  expList,
	}
}

func parseNameList(lexer *Lexer, firstName string) []string {
	nameList := []string{firstName}            // Name
	for lexer.LookAhead() == TOKEN_SEP_COMMA { // {*}
		lexer.NextToken()                     // ,
		_, nextName := lexer.NextIdentifier() // Name
		nameList = append(nameList, nextName)
	}
	return nameList
}

func parseVarList(lexer *Lexer, firstExp ast.Exp) []ast.Exp {
	varList := []ast.Exp{_ensureVarType(firstExp)} // var
	for lexer.LookAhead() == TOKEN_SEP_COMMA {     // {*}
		lexer.NextToken()                // ,
		nextExp := parsePrefixExp(lexer) // var
		varList = append(varList, _ensureVarType(nextExp))
	}
	return varList
}

func parseFuncName(lexer *Lexer) (exp ast.Exp, hasColon bool, line int) {
	_line, name := lexer.NextIdentifier() // Name
	line = _line
	exp = &ast.NameExp{Line: line, Name: name}
	hasColon = false

	for lexer.LookAhead() == TOKEN_SEP_DOT { // {*}
		lexer.NextToken()                    // .
		line, name := lexer.NextIdentifier() // Name
		exp = _dot(exp, name, line)
	}

	if lexer.LookAhead() == TOKEN_SEP_COMMA { // [*]
		lexer.NextToken()                    // :
		line, name := lexer.NextIdentifier() // Name
		exp = _dot(exp, name, line)
		hasColon = true
	}

	return
}

// 保证表达式为 NameExp 和 TableAccessExp
func _ensureVarType(exp ast.Exp) ast.Exp {
	switch exp := exp.(type) {
	case *ast.NameExp, *ast.TableAccessExp:
		return exp
	}
	panic(fmt.Sprintf("Unknown exp type %v!", exp))
}

// 生成 parent[key] 的形式
func _dot(parent ast.Exp, key string, line int) *ast.TableAccessExp {
	return &ast.TableAccessExp{
		LastLine:  line,
		PrefixExp: parent,
		KeyExp: ast.StringExp{
			Line: line,
			Val:  key,
		},
	}
}
