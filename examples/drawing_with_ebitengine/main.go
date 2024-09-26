package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"

	"github.com/demouth/ebitencp"
	"github.com/demouth/ebitencp/examples/drawing_with_ebitengine/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp/v2"
)

const (
	screenWidth  = 800
	screenHeight = 800
)

var (
	space  *cp.Space
	drawer *ebitencp.Drawer

	drawingWithEbitengine = true
)

type Game struct {
	camera Camera
}
type Camera struct {
	Offset cp.Vector
	Zoom   float64
	Rotate float64
}

func (g *Game) Update() error {
	space.Step(1 / 60.0)
	drawer.HandleMouseEvent(space)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		drawingWithEbitengine = !drawingWithEbitengine
	}
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.camera = Camera{
			Zoom: 1,
			// Set the camera offset to the center of the screen
			Offset: cp.Vector{X: -screenWidth / 2, Y: -screenHeight / 2},
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camera.Offset.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camera.Offset.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camera.Offset.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camera.Offset.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.Rotate += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.Rotate -= 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.camera.Zoom += 0.05
	}

	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.camera.Zoom -= 0.05
		if g.camera.Zoom < 0.05 {
			g.camera.Zoom = 0.05
		}
	}
	drawer.GeoM.Reset()
	drawer.GeoM.Translate(g.camera.Offset.X, g.camera.Offset.Y)
	drawer.GeoM.Scale(g.camera.Zoom, g.camera.Zoom)
	drawer.GeoM.Rotate(g.camera.Rotate)
	drawer.GeoM.Translate(screenWidth/2, screenHeight/2)
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	if drawingWithEbitengine {
		space.EachShape(func(s *cp.Shape) {
			switch s.Class.(type) {
			case *cp.Circle:
				circle := s.Class.(*cp.Circle)
				body := circle.Body()
				util.DrawRunner(
					screen,
					*drawer.GeoM,
					float32(body.Position().X),
					float32(body.Position().Y),
					float32(circle.Radius()),
					float32(circle.Body().Angle()),
				)
			case *cp.PolyShape:
				poly := s.Class.(*cp.PolyShape)
				body := poly.Body()
				r := (poly.TransformVert(0).Distance(poly.TransformVert(1))) * 0.5
				util.DrawRunner(
					screen,
					*drawer.GeoM,
					float32(body.Position().X),
					float32(body.Position().Y),
					float32(r),
					float32(poly.Body().Angle()),
				)
			case *cp.Segment:
				segment := s.Class.(*cp.Segment)
				ta := segment.TransformA()
				tb := segment.TransformB()
				util.DrawLine(
					screen,
					*drawer.GeoM,
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
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			`Offset: %v
Zoom: %v
Rotation: %v
FlipYAxis: %v
Usage:
  Camera Position = WASD
  Camera Rotation = Q / E
  Camera Zoom = Z / X
  Reset Camera = Backspace
  Space = Toggle drawing process
  Drag Object = Cursor`,
			g.camera.Offset,
			g.camera.Zoom,
			g.camera.Rotate,
			drawer.FlipYAxis,
		),
	)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func main() {
	space = cp.NewSpace()
	// Gravity is set to a positive value to match the Ebitengine coordinate system
	space.SetGravity(cp.Vector{X: 0, Y: 200})
	for i := 0; i < 100; i++ {
		size := rand.Float64()*30 + 20
		addBox(space, size, size*2, screenWidth*rand.Float64(), screenHeight*rand.Float64())
	}
	addWall(space, 0, screenHeight, 0, 0, 5)
	addWall(space, screenWidth, screenHeight, screenWidth, 0, 5)
	addWall(space, 0, 0, screenWidth, 0, 5)
	addWall(space, 0, screenHeight, screenWidth, screenHeight, 5)

	game := &Game{}
	drawer = ebitencp.NewDrawer(0, 0)
	drawer.FlipYAxis = true
	game.camera = Camera{
		Zoom: 1,
		// Set the camera offset to the center of the screen
		Offset: cp.Vector{X: -screenWidth / 2, Y: -screenHeight / 2},
	}
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
func addBox(space *cp.Space, w, h float64, x, y float64) *cp.Body {
	mass := w * h / 400.0
	body := space.AddBody(cp.NewBody(mass, cp.MomentForBox(mass, w, h)))
	body.SetPosition(cp.Vector{X: x, Y: y})

	shape := space.AddShape(cp.NewBox(body, w, h, 0))
	shape.SetElasticity(0.9)
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
