package main

type Lexer struct {
	input           string
	currentPosition int
	nextPosition    int
	char            byte
}

func NewLexer(input string) *Lexer {
	lexer := Lexer{input: input}
	lexer.readChar()
	return &lexer
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) readChar() {
	l.char = l.peekChar()
	l.currentPosition = l.nextPosition
	l.nextPosition++
}

type TokenType string

const (
	EOF             TokenType = "EOF"
	ILLEGAL         TokenType = "ILLEGAL"
	TokenLt         TokenType = "LT"
	TokenGt         TokenType = "GT"
	TokenSlash      TokenType = "SLASH"
	TokenQuote      TokenType = "QUOTE"
	TokenEqual      TokenType = "EQUAL"
	TokenIdentifier TokenType = "IDENT"
)

type Token struct {
	Type    TokenType
	Literal string
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()
	switch l.char {
	case '<':
		tok = Token{TokenLt, string(l.char)}
	case '>':
		tok = Token{TokenGt, string(l.char)}
	case '/':
		tok = Token{TokenSlash, string(l.char)}
	case '"':
		tok = Token{TokenQuote, string(l.char)}
	case '=':
		tok = Token{TokenEqual, string(l.char)}
	case 0:
		tok.Type = EOF
	default:
		if isAlphaNum(l.char) {
			tok.Literal = l.readAlphaNum()
			tok.Type = TokenIdentifier
			return tok
		} else {
			tok = Token{ILLEGAL, string(l.char)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func isAlphaNum(char byte) bool {
	return isLetter(char) || isDigit(char)
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char == '-'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) readAlphaNum() string {
	startPos := l.currentPosition
	for isAlphaNum(l.char) {
		l.readChar()
	}
	return l.input[startPos:l.currentPosition]
}
