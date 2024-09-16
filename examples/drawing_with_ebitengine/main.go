package main

import (
	"image/color"
	_ "image/png"

	"github.com/demouth/ebitencp"
	"github.com/demouth/ebitencp/examples/drawing_with_ebitengine/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp/v2"
)

const (
	screenWidth  = 600
	screenHeight = 600
)

var (
	space  *cp.Space
	drawer *ebitencp.Drawer

	drawingWithEbitengine = false
)

type Game struct{}

func (g *Game) Update() error {
	space.Step(1 / 60.0)
	drawer.HandleMouseEvent(space)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		drawingWithEbitengine = !drawingWithEbitengine
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	if drawingWithEbitengine {
		space.EachShape(func(s *cp.Shape) {
			switch s.Class.(type) {
			case *cp.Circle:
				circle := s.Class.(*cp.Circle)
				body := circle.Body()
				util.DrawCircle(
					screen,
					float32(body.Position().X),
					float32(body.Position().Y),
					float32(circle.Radius()),
					color.RGBA{0xff, 0xff, 0xff, 0xff},
				)
			case *cp.Segment:
				segment := s.Class.(*cp.Segment)
				ta := segment.TransformA()
				tb := segment.TransformB()
				util.DrawLine(
					screen,
					float32(ta.X), float32(ta.Y),
					float32(tb.X), float32(tb.Y),
					float32(segment.Radius()*2),
					color.RGBA{0xff, 0xff, 0xff, 0xff},
				)
			}
		})
	} else {
		cp.DrawSpace(space, drawer.WithScreen(screen))
	}
	ebitenutil.DebugPrint(screen, "\n Press 'Space' to toggle drawing with ebitengine")
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func main() {
	space = cp.NewSpace()
	// Gravity is set to a positive value to match the Ebitengine coordinate system
	space.SetGravity(cp.Vector{X: 0, Y: 200})
	addBall(space, screenWidth/2, screenHeight/2, 50)
	addWall(space, 0, screenHeight, 0, 0, 5)
	addWall(space, screenWidth, screenHeight, screenWidth, 0, 5)
	addWall(space, 0, 0, screenWidth, 0, 5)
	addWall(space, 0, screenHeight, screenWidth, screenHeight, 5)
	addWall(space, 200, 500, 400, 510, 5)
	addWall(space, 400, 110, 200, 100, 5)

	game := &Game{}
	drawer = ebitencp.NewDrawer(screenWidth, screenHeight)
	drawer.FlipYAxis = true

	// Set the camera offset to the center of the screen
	drawer.Camera.Offset = cp.Vector{X: screenWidth / 2, Y: screenHeight / 2}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.RunGame(game)
}
func addBall(space *cp.Space, x, y, radius float64) *cp.Body {
	mass := radius * radius / 100.0
	body := space.AddBody(
		cp.NewBody(
			mass,
			cp.MomentForCircle(mass, 0, radius, cp.Vector{}),
		),
	)
	body.SetPosition(cp.Vector{X: x, Y: y})
	shape := space.AddShape(
		cp.NewCircle(
			body,
			radius,
			cp.Vector{},
		),
	)
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
	return body
}
func addWall(space *cp.Space, x1, y1, x2, y2, radius float64) {
	pos1 := cp.Vector{X: x1, Y: y1}
	pos2 := cp.Vector{X: x2, Y: y2}
	shape := space.AddShape(cp.NewSegment(space.StaticBody, pos1, pos2, radius))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
}
