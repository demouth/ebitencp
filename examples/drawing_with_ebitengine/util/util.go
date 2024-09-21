package util

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage  *ebiten.Image
	runnerImage *ebiten.Image
)

func init() {
	whiteImage = ebiten.NewImage(3, 3)
	whiteImage.Fill(color.White)

	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)
}

func DrawLine(screen *ebiten.Image, x1, y1, x2, y2, width float32, c color.RGBA) {
	path := vector.Path{}
	path.MoveTo(x1, y1)
	path.LineTo(x2, y2)
	path.Close()
	sop := &vector.StrokeOptions{}
	sop.Width = width
	sop.LineJoin = vector.LineJoinRound
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, sop)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(c.R) / float32(0xff)
		vs[i].ColorG = float32(c.G) / float32(0xff)
		vs[i].ColorB = float32(c.B) / float32(0xff)
		vs[i].ColorA = float32(c.A) / float32(0xff)
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.FillRule = ebiten.FillAll
	screen.DrawTriangles(vs, is, whiteImage, op)
}
func DrawFill(screen *ebiten.Image, path vector.Path, c color.RGBA) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(c.R) / float32(0xff)
		vs[i].ColorG = float32(c.G) / float32(0xff)
		vs[i].ColorB = float32(c.B) / float32(0xff)
		vs[i].ColorA = float32(c.A) / float32(0xff)
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.FillRule = ebiten.FillAll
	screen.DrawTriangles(vs, is, whiteImage, op)
}

func DrawCircle(screen *ebiten.Image, x, y, radius float32, c color.RGBA) {
	path := vector.Path{}
	path.Arc(x, y, radius, 0, 2*math.Pi, vector.Clockwise)
	DrawFill(screen, path, c)
}
func DrawRunner(screen *ebiten.Image, x, y, radius, rotate float32) {
	const frameSize = 32
	r := radius / frameSize * 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-frameSize/2, -frameSize/2)
	op.GeoM.Translate(0, -3) // fine tuning
	op.GeoM.Scale(float64(r), float64(r))
	op.GeoM.Scale(1.3, 1.3) // fine tuning
	op.GeoM.Rotate(float64(rotate))
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(runnerImage.SubImage(image.Rect(0, 0, frameSize, frameSize)).(*ebiten.Image), op)
}
