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

	// fmt.Println(string(input))

	lexer := NewLexer(string(input))
	// for {
	// 	tok := lexer.NextToken()
	// 	if tok.Type == EOF {
	// 		break
	// 	}
	// 	fmt.Println(tok.Type, tok.Literal)
	// }

	parser := NewParser(lexer)
	svg := parser.ParseSvg()

	// fmt.Println(svg)

	dc := gg.NewContext(svg.width, svg.height)
	// White background
	dc.SetRGB(1.0, 1.0, 1.0)
	dc.Clear()

	for _, s := range svg.shapes {
		s.paint(dc)
	}

	dc.SavePNG("out.png")
}
