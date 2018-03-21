package main

import (
	"log"
	"strconv"
)

type Parser struct {
	lex          *Lexer
	currentToken Token
	nextToken    Token
}

func NewParser(lexer *Lexer) *Parser {
	parser := Parser{lex: lexer}
	parser.readToken()
	parser.readToken()
	return &parser
}

func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lex.NextToken()
}

func (p *Parser) ParseSvg() *svg {
	s := svg{}
	s.shapes = []painter{}

	p.Must(TokenLt, "expect <")
	tagName := p.parseIdentifier()
	if tagName != "svg" {
		panic("expect svg")
	}

	attributes := p.parseAttributes()
	s.width = lookupInt(attributes, "width")
	s.height = lookupInt(attributes, "height")

	p.Must(TokenGt, "expect svg >")

	// Look for "</"
	for p.currentToken.Type != TokenLt || p.nextToken.Type != TokenSlash {
		pp := p.parseShape()
		s.shapes = append(s.shapes, pp)
	}

	// Consume remaining tokens, and ignore them
	for p.currentToken.Type != EOF {
		p.readToken()
	}

	return &s
}

func (p *Parser) parseAttributes() map[string]string {
	attributes := make(map[string]string)
	for p.currentToken.Type != TokenSlash && p.currentToken.Type != TokenGt {
		k, v := p.parseAttribute()
		attributes[k] = v
	}
	return attributes
}

func (p *Parser) parseAttribute() (key, value string) {
	if p.currentToken.Type != TokenIdentifier {
		panic("expect attribute identifier")
	}
	key = p.currentToken.Literal
	p.readToken()

	if p.currentToken.Type != TokenEqual {
		panic("expect attribute equal")
	}
	p.readToken()

	if p.currentToken.Type != TokenQuote {
		panic("expect attribute opening quote")
	}
	p.readToken()

	if p.currentToken.Type != TokenIdentifier {
		panic("expect attribute quoted identifier")
	}
	value = p.currentToken.Literal
	p.readToken()

	if p.currentToken.Type != TokenQuote {
		panic("expect attribute closing quote")
	}
	p.readToken()
	return
}

func (p *Parser) Must(tokenType TokenType, msg string) {
	if p.currentToken.Type != tokenType {
		panic(msg)
	}
	p.readToken()
}

func (p *Parser) parseIdentifier() string {
	if p.currentToken.Type != TokenIdentifier {
		panic("expect identifier")
	}
	id := p.currentToken.Literal
	p.readToken()
	return id
}

func (p *Parser) parseShape() painter {
	p.Must(TokenLt, "expect <")
	shapeName := p.parseIdentifier()
	attributes := p.parseAttributes()

	var pp painter

	switch shapeName {
	case "circle":
		cx := lookupFloat(attributes, "cx")
		cy := lookupFloat(attributes, "cy")
		r := lookupFloat(attributes, "r")
		strokeWidth := lookupFloat(attributes, "stroke-width")
		stroke := lookup(attributes, "stroke")
		fill := lookup(attributes, "fill")

		c := &circle{
			shape: shape{
				x:           cx,
				y:           cy,
				strokeWidth: strokeWidth,
				fillColor:   fill,
				strokeColor: stroke,
			},
			radius: r,
		}
		pp = painter(c)
	default:
		log.Printf("unknown shape %s\n", shapeName)
	}

	p.Must(TokenSlash, "expect /")
	p.Must(TokenGt, "expect >")

	return pp
}

func lookupFloat(attributes map[string]string, key string) float64 {
	if v, ok := attributes[key]; ok {
		fl, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic("cannot parse float: " + v)
		}
		return fl
	}
	panic("key not found" + key)
}

func lookupInt(attributes map[string]string, key string) int {
	if v, ok := attributes[key]; ok {
		integer, err := strconv.Atoi(v)
		if err != nil {
			panic("cannot parse integer: " + v)
		}
		return integer
	}
	panic("key not found" + key)
}

func lookup(attributes map[string]string, key string) string {
	if v, ok := attributes[key]; ok {
		return v
	}
	panic("key not found" + key)
}
