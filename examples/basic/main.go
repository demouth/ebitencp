package main

import (
	_ "image/png"

	"github.com/demouth/ebitencp"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jakecoffman/cp/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	space  *cp.Space
	drawer *ebitencp.Drawer
)

type Game struct{}

func (g *Game) Update() error {
	// Handling dragging
	drawer.HandleMouseEvent(space)
	space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cp.DrawSpace(space, drawer.WithScreen(screen))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Initialising Chipmunk
	space = cp.NewSpace()
	space.SetGravity(cp.Vector{X: 0, Y: -100})
	addWall(space, cp.Vector{X: -200, Y: -100}, cp.Vector{X: -10, Y: -150}, 5)
	addWall(space, cp.Vector{X: 200, Y: -100}, cp.Vector{X: 10, Y: -150}, 5)
	addBall(space, -50, 0, 50)
	addBall(space, 50, 200, 20)

	// Initialising Ebitengine/v2
	game := &Game{}
	drawer = ebitencp.NewDrawer(screenWidth, screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.RunGame(game)
}

func addWall(space *cp.Space, pos1 cp.Vector, pos2 cp.Vector, radius float64) {
	shape := space.AddShape(cp.NewSegment(space.StaticBody, pos1, pos2, radius))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
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
