package main

type Lexer struct {
	input           string
	currentPosition int
	nextPosition    int
	ch              byte
}

type TokenType string

const (
	EOF             TokenType = "EOF"
	ILLEGAL         TokenType = "ILLEGAL"
	TokenOpen       TokenType = "<"
	TokenClose      TokenType = ">"
	TokenSlash      TokenType = "/"
	TokenQuote      TokenType = "\""
	TokenEqual      TokenType = "="
	TokenIdentifier TokenType = "IDENT"
)

type Token struct {
	Type    TokenType
	Literal string
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
	l.ch = l.peekChar()
	l.currentPosition = l.nextPosition
	l.nextPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.readWhitespace()

	switch l.ch {
	case '<':
		tok = Token{TokenOpen, string(l.ch)}
	case '>':
		tok = Token{TokenClose, string(l.ch)}
	case '/':
		tok = Token{TokenSlash, string(l.ch)}
	case '"':
		tok = Token{TokenQuote, string(l.ch)}
	case '=':
		tok = Token{TokenEqual, string(l.ch)}
	case 0:
		tok.Type = EOF
	default:
		if isAlphaNum(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = TokenIdentifier
			// Don't call readChar(), as the chars were consummed by readIdentifier()
			return tok
		} else {
			tok = Token{ILLEGAL, string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	startPos := l.currentPosition
	for isAlphaNum(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.currentPosition]
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}
func isAlphaNum(char byte) bool {
	return isLetter(char) || isDigit(char)
}
func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '-'
}
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
