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
		stepExp = &ast.IntegerExp{lexer.Line(), 1}
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
	lexer.NextToken()                 // function
	_, name := lexer.NextIdentifier() // Name
	exp := parseFuncDefExp(lexer)     // funcbody
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

/*
prefixexp ::= var | functioncall | '(' exp ')'

varlist '=' explist
varlist ::= var {',' var}
var ::= Name | prefixexp '[' exp ']' | prefixexp '.' Name
*/
func parseAssignStat(lexer *Lexer, firstExp ast.Exp) *ast.AssignStat {
	varList := parseVarList(lexer, firstExp)
	lexer.NextKeyword(TOKEN_OP_ASSIGN)
	expList := parseExpList(lexer)
	return &ast.AssignStat{
		LastLine: lexer.Line(),
		VarList:  varList,
		ExpList:  expList,
	}
}

// namelist ::= Name {',' Name}
func parseNameList(lexer *Lexer, firstName string) []string {
	nameList := []string{firstName}
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		_, nextName := lexer.NextIdentifier()
		nameList = append(nameList, nextName)
	}
	return nameList
}

/*
varlist ::= var {',' var}
var ::= Name | prefixexp '[' exp ']' | prefixexp '.' Name
*/
func parseVarList(lexer *Lexer, firstExp ast.Exp) []ast.Exp {
	varList := []ast.Exp{_ensureVarType(firstExp)}
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		nextExp := parsePrefixExp(lexer)
		varList = append(varList, _ensureVarType(nextExp))
	}
	return varList
}

// funcname ::= Name {'.' Name} [':' Name] '{'
func parseFuncName(lexer *Lexer) (exp ast.Exp, hasColon bool) {
	line, name := lexer.NextIdentifier() // Name
	exp = &ast.NameExp{line, name}
	hasColon = false

	for lexer.LookAhead() == TOKEN_SEP_DOT { // {*}
		lexer.NextToken()                    // .
		line, name := lexer.NextIdentifier() // Name
		exp = _dot(exp, name, line)
	}

	if lexer.LookAhead() == TOKEN_SEP_COMMA { // [*]
		hasColon = true
		lexer.NextToken()                    // :
		line, name := lexer.NextIdentifier() // Name
		exp = _dot(exp, name, line)
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
