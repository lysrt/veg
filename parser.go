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
	s.drawables = []drawable{}

	p.expect(TokenOpen)
	if tagName := p.parseIdentifier(); tagName != "svg" {
		log.Fatalf("expected tag 'svg', got '%s'\n", tagName)
	}

	attributes := p.parseAttributes()
	s.width = lookupInt(attributes, "width")
	s.height = lookupInt(attributes, "height")

	p.expect(TokenClose)

	// Stop parsing shapes when finding "</"
	for p.currentToken.Type != TokenOpen || p.nextToken.Type != TokenSlash {
		d := p.parseShape()
		s.drawables = append(s.drawables, d)
	}

	// Consume and ignore remaining tokens
	for p.currentToken.Type != EOF {
		p.readToken()
	}

	return &s
}

func (p *Parser) parseAttributes() map[string]string {
	attributes := make(map[string]string)
	for p.currentToken.Type != TokenSlash && p.currentToken.Type != TokenClose {
		k, v := p.parseAttribute()
		attributes[k] = v
	}
	return attributes
}

func (p *Parser) parseAttribute() (key, value string) {
	key = p.parseIdentifier()
	p.expect(TokenEqual)
	p.expect(TokenQuote)
	value = p.parseIdentifier()
	p.expect(TokenQuote)
	return
}

func (p *Parser) expect(tokenType TokenType) {
	if p.currentToken.Type != tokenType {
		log.Fatalf("expect %q, got %q (%q)", tokenType, p.currentToken.Type, p.currentToken.Literal)
	}
	p.readToken()
}

func (p *Parser) parseIdentifier() string {
	if p.currentToken.Type != TokenIdentifier {
		log.Printf("expect identifier, got %q (%q)", p.currentToken.Type, p.currentToken.Literal)
	}
	id := p.currentToken.Literal
	p.readToken()
	return id
}

func (p *Parser) parseShape() drawable {
	p.expect(TokenOpen)
	shapeName := p.parseIdentifier()
	attributes := p.parseAttributes()

	var d drawable
	switch shapeName {
	case "circle":
		d = parseCircle(attributes)
	default:
		log.Printf("unknown shape %s\n", shapeName)
	}

	p.expect(TokenSlash)
	p.expect(TokenClose)

	return d
}

func parseCircle(attributes map[string]string) *circle {
	cx := lookupFloat(attributes, "cx")
	cy := lookupFloat(attributes, "cy")
	r := lookupFloat(attributes, "r")
	strokeWidth := lookupFloat(attributes, "stroke-width")
	stroke := attributes["stroke"]
	fill := attributes["fill"]

	return &circle{
		shape: shape{
			x:           cx,
			y:           cy,
			strokeWidth: strokeWidth,
			fillColor:   fill,
			strokeColor: stroke,
		},
		radius: r,
	}
}

func lookupFloat(m map[string]string, key string) float64 {
	fl, err := strconv.ParseFloat(m[key], 64)
	if err != nil {
		panic("cannot parse float: " + m[key])
	}
	return fl
}
func lookupInt(m map[string]string, key string) int {
	integer, err := strconv.Atoi(m[key])
	if err != nil {
		panic("cannot parse integer: " + m[key])
	}
	return integer
}
