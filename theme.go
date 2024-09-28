package ebitencp

import (
	"image/color"

	"github.com/jakecoffman/cp/v2"
)

type Theme struct {
	Outline                         color.RGBA
	Shape, ShapeSleeping, ShapeIdle color.RGBA
	Constraint, CollisionPoint      color.RGBA
}

func toFColor(c color.RGBA) cp.FColor {
	r := float32(c.R) / 255.0
	g := float32(c.G) / 255.0
	b := float32(c.B) / 255.0
	a := float32(c.A) / 255.0
	return cp.FColor{R: r, G: g, B: b, A: a}
}

func DefaultTheme() *Theme {
	return &Theme{
		Outline:        color.RGBA{0xC8, 0xD2, 0xE6, 0xFF},
		ShapeSleeping:  color.RGBA{0x33, 0x33, 0x33, 0x80},
		ShapeIdle:      color.RGBA{0xA8, 0xA8, 0xA8, 0x80},
		Shape:          color.RGBA{0xB2, 0x4C, 0x99, 0x80},
		Constraint:     color.RGBA{0x00, 0xBF, 0x00, 255},
		CollisionPoint: color.RGBA{0xFF, 0x19, 0x33, 255},
	}
}
