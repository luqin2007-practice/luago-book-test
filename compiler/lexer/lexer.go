package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

/*
Lexer 词法分析器

	chunk     源代码
	chunkName 源文件名
	line      当前行号
*/
type Lexer struct {
	chunk     string
	chunkName string
	line      int
	// 下一个 Token
	nextToken string
	nextKind  int
	nextLine  int
}

func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{
		chunk:     chunk,
		chunkName: chunkName,
		line:      1,
	}
}

/*
NextToken 获取下一个 Token

	line  所在行
	kind  Token 类型
	token Token 字符串
*/
func (self *Lexer) NextToken() (line, kind int, token string) {
	// 已经看过下一个 Token，直接使用
	if len(self.nextToken) > 0 {
		line = self.nextLine
		kind = self.nextKind
		token = self.nextToken
		self.line = self.nextLine
		self.nextToken = ""
		self.nextLine = 0
		return
	}

	self.skipWhiteSpace()
	// 文档结束
	if len(self.chunk) == 0 {
		return self.line, TOKEN_EOF, "EOF"
	}

	// Token: 运算符、分隔符
	switch self.chunk[0] {
	case ';':
		return self._nextToken(TOKEN_SEP_SEMI, ";")
	case ',':
		return self._nextToken(TOKEN_SEP_COMMA, ",")
	case '(':
		return self._nextToken(TOKEN_SEP_LPAREN, "(")
	case ')':
		return self._nextToken(TOKEN_SEP_RPAREN, ")")
	case ']':
		return self._nextToken(TOKEN_SEP_RBRACK, "]")
	case '{':
		return self._nextToken(TOKEN_SEP_LCURLY, "{")
	case '}':
		return self._nextToken(TOKEN_SEP_RCURLY, "}")
	case '+':
		return self._nextToken(TOKEN_OP_ADD, "+")
	case '-':
		return self._nextToken(TOKEN_OP_MINUS, "-")
	case '*':
		return self._nextToken(TOKEN_OP_MUL, "*")
	case '^':
		return self._nextToken(TOKEN_OP_POW, "^")
	case '%':
		return self._nextToken(TOKEN_OP_MOD, "%")
	case '&':
		return self._nextToken(TOKEN_OP_BAND, "&")
	case '|':
		return self._nextToken(TOKEN_OP_BOR, "|")
	case '#':
		return self._nextToken(TOKEN_OP_LEN, "#")
	case ':':
		if self.test("::") {
			return self._nextToken(TOKEN_SEP_LABEL, "::")
		}
		return self._nextToken(TOKEN_SEP_COLON, ":")
	case '/':
		if self.test("//") {
			return self._nextToken(TOKEN_OP_IDIV, "//")
		}
		return self._nextToken(TOKEN_OP_DIV, "/")
	case '~':
		if self.test("~=") {
			return self._nextToken(TOKEN_OP_NE, "~=")
		}
		return self._nextToken(TOKEN_OP_WAVE, "~")
	case '=':
		if self.test("==") {
			return self._nextToken(TOKEN_OP_EQ, "==")
		}
		return self._nextToken(TOKEN_OP_ASSIGN, "=")
	case '<':
		if self.test("<<") {
			return self._nextToken(TOKEN_OP_SHL, "<<")
		} else if self.test("<=") {
			return self._nextToken(TOKEN_OP_LE, "<=")
		}
		return self._nextToken(TOKEN_OP_LT, "<")
	case '>':
		if self.test(">>") {
			return self._nextToken(TOKEN_OP_SHR, ">>")
		} else if self.test(">=") {
			return self._nextToken(TOKEN_OP_GE, ">=")
		}
		return self._nextToken(TOKEN_OP_GT, ">")
	case '.':
		if self.test("...") {
			return self._nextToken(TOKEN_VARARG, "...")
		} else if self.test("..") {
			return self._nextToken(TOKEN_OP_CONCAT, "..")
		} else if len(self.chunk) == 1 || !isDigit(self.chunk[1]) {
			return self._nextToken(TOKEN_SEP_DOT, ".")
		}
	case '[':
		if self.test("[[") || self.test("[=") {
			// 长字符串
			return self.line, TOKEN_STRING, self.scanLongString()
		} else {
			return self._nextToken(TOKEN_SEP_LBRACK, "[")
		}
	case '\'', '"':
		return self.line, TOKEN_STRING, self.scanShortString()
	}

	// 数字
	c := self.chunk[0]
	if c == '.' || isDigit(c) {
		return self.line, TOKEN_NUMBER, self.scanNumber()
	}

	// 关键字、标识符
	if c == '_' || isLetter(c) {
		token := self.scanIdentifier()
		if kind, found := keywords[token]; found {
			return self.line, kind, token
		} else {
			return self.line, TOKEN_IDENTIFIER, token
		}
	}

	self.error("unexpected symbol near %q", c)
	return
}

// LookAhead 查看下一个 Token
func (self *Lexer) LookAhead() (line, kind int, token string) {
	if len(self.nextToken) > 0 {
		return self.line, self.nextKind, self.nextToken
	}

	currentLine := self.line
	line, kind, token = self.NextToken()
	self.line = currentLine
	self.nextToken = token
	self.nextKind = kind
	self.nextLine = line
	return
}

// skipWhiteSpace 跳过空白字符和注释
func (self *Lexer) skipWhiteSpace() {
	for len(self.chunk) > 0 {
		switch {
		case self.test("--"):
			self.skipComment()
			continue
		case self.test("\r\n") || self.test("\n\r"):
			self.next(2)
			self.line += 1
			continue
		case isNewLine(self.chunk[0]):
			self.next(1)
			self.line += 1
			continue
		case isWhitespace(self.chunk[0]):
			self.next(1)
			continue
		}
		break
	}
}

// next 向后移动 n 个字符
func (self *Lexer) next(n int) {
	self.chunk = self.chunk[n:]
}

// test 以某字符开头
func (self *Lexer) test(s string) bool {
	return strings.HasPrefix(self.chunk, s)
}

// skipComment 跳过注释
func (self *Lexer) skipComment() {
	self.next(2) // --

	if self.test("[") && reOpeningLongBracket.FindString(self.chunk) != "" {
		// long comment
		self.scanLongString()
		return
	}

	for len(self.chunk) > 0 && !isNewLine(self.chunk[0]) {
		self.next(1)
	}
}

// scanLongString 跳过长方括号，获取长方括号内字符串
func (self *Lexer) scanLongString() string {
	// 长括号前半截
	openingLongBracket := reOpeningLongBracket.FindString(self.chunk)
	if openingLongBracket == "" {
		self.error("Invalid long string delimiter near '%s'", self.chunk[0:2])
	}

	// 长括号后半截
	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	closingLongBracketIdx := strings.Index(self.chunk, closingLongBracket)
	if closingLongBracketIdx < 0 {
		self.error("Unfinished long string or comment")
	}

	// 取括号内字符串
	str := self.chunk[len(openingLongBracket):closingLongBracketIdx]
	self.next(closingLongBracketIdx + len(closingLongBracket))

	// 替换换行符，统计行数
	str = reNewLine.ReplaceAllString(str, "\n")
	self.line += strings.Count(str, "\n")
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}

	return str
}

// scanShortString 跳过短字符串，获取字符串内容
func (self *Lexer) scanShortString() string {
	// 使用正则提取短字符串
	if str := reShortStr.FindString(self.chunk); str != "" {
		self.next(len(str))
		// 去除两端引号
		str = str[1 : len(str)-1]
		// 处理转义字符
		if strings.Index(str, `\`) >= 0 {
			self.line += len(reNewLine.FindAllString(str, -1))
			str = self._escape(str)
		}
		return str
	}

	self.error("Unfinished string!")
	return ""
}

func (self *Lexer) scanNumber() string {
	return self._scan(reNumber)
}

func (self *Lexer) scanIdentifier() string {
	return self._scan(reIdentifier)
}

func isNewLine(ch byte) bool {
	return ch == '\r' || ch == '\n'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\v' || ch == '\f'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (self *Lexer) error(format string, args ...interface{}) {
	err := fmt.Sprintf(format, args...)
	err = fmt.Sprintf("%s:%d: %s", self.chunkName, self.line, err)
	panic(err)
}

var _excapes = map[byte]byte{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'\n': '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'"':  '"',
	'\'': '\'',
	'\\': '\\',
}

func (self *Lexer) _nextToken(tKind int, sToken string) (line, kind int, token string) {
	self.next(len(sToken))
	return self.line, tKind, sToken
}

func (self *Lexer) _escape(str string) string {
	var buf bytes.Buffer
	for len(str) > 0 {
		// 非转义部分
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}
		// 长度不足
		if len(str) == 1 {
			self.error("Unfinished string!")
		}

		ch := str[1]
		if replace, found := _excapes[ch]; found {
			// 一般转义字符
			buf.WriteByte(replace)
			str = str[2:]
			continue
		} else if isDigit(ch) {
			// \ddd
			if found := reDecEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[1:], 10, 32)
				if d <= 0xFF {
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
				self.error("decimal escape too large near '%s'", found)
			}
		} else if ch == 'x' {
			// \xXXX
			if found := reHexEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[2:], 16, 32)
				buf.WriteByte(byte(d))
				str = str[len(found):]
				continue
			}
		} else if ch == 'u' {
			// \u{XXX}
			if found := reUnicodeEscapeSeq.FindString(str); found != "" {
				d, err := strconv.ParseUint(found[3:len(found)-1], 16, 32)
				if err == nil && d <= 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				self.error("UTF-8 value too large near '%s'", found)
			}
		} else if ch == 'z' {
			// \z 换行
			str = str[2:]
			for len(str) > 0 && isWhitespace(str[0]) {
				str = str[1:]
			}
			continue
		}
		self.error("invalid escape sequence near '\\%c'", ch)
	}

	return buf.String()
}

func (self *Lexer) _scan(re *regexp.Regexp) string {
	if token := re.FindString(self.chunk); token != "" {
		self.next(len(token))
		return token
	}

	panic("unreachable!")
}
