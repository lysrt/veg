package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Give one argument")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	input, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	lexer := NewLexer(string(input))
	parser := NewParser(lexer)
	svg := parser.ParseSvg()

	dc := gg.NewContext(svg.width, svg.height)
	// dc.SetColor(color.White)
	// dc.Clear()

	for _, s := range svg.drawables {
		s.draw(dc)
	}

	dc.SavePNG("out.png")
}
