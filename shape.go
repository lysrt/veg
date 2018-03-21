package main

import (
	"github.com/fogleman/gg"
)

type svg struct {
	drawables     []drawable
	width, height int
}

type drawable interface {
	draw(dc *gg.Context)
}

type shape struct {
	x, y        float64
	fillColor   string
	strokeColor string
	strokeWidth float64
}

type circle struct {
	shape
	radius float64
}

func (c *circle) draw(dc *gg.Context) {
	x := c.x
	y := c.y
	r := c.radius

	dc.DrawCircle(x, y, r)
	dc.SetHexColor(c.fillColor)
	dc.Fill()

	dc.DrawCircle(x, y, r)
	dc.SetLineWidth(c.strokeWidth)
	dc.SetHexColor(c.strokeColor)
	dc.Stroke()
}
