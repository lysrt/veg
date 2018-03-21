package main

import (
	"github.com/fogleman/gg"
)

type svg struct {
	shapes        []painter
	width, height int
}

type shape struct {
	x, y        float64
	strokeColor string
	strokeWidth float64
	fillColor   string
}

type painter interface {
	paint(dc *gg.Context)
}

type circle struct {
	shape
	radius float64
}

func (c *circle) paint(dc *gg.Context) {
	x := c.shape.x
	y := c.shape.y
	r := c.radius
	dc.DrawCircle(x, y, r)
	dc.SetHexColor(c.shape.fillColor)
	dc.Fill()

	dc.DrawCircle(x, y, r)
	dc.SetLineWidth(c.shape.strokeWidth)
	dc.SetHexColor(c.shape.strokeColor)
	dc.Stroke()
}
