package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/demouth/ebitencp"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/jakecoffman/cp/v2"
)

const (
	screenWidth   = 640
	screenHeight  = 480
	hScreenWidth  = screenWidth / 2
	hScreenHeight = screenHeight / 2
)

type Game struct {
	space     *cp.Space
	drawer    *ebitencp.Drawer
	ball1     *cp.Body
	ball2     *cp.Body
	flipYAxis bool
}

func (g *Game) Update() error {
	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)
	g.drawer.Camera.Offset.X = g.ball1.Position().X
	g.drawer.Camera.Offset.Y = g.ball1.Position().Y
	g.space.Step(1 / 60.0)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.drawer.FlipYAxis = !g.drawer.FlipYAxis
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cp.DrawSpace(g.space, g.drawer.WithScreen(screen))

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("Camera Offset: %v\nBall Position: %v\nFlipYAxis: %v\nPress SPACE to flip Y axis",
			g.drawer.Camera.Offset,
			g.ball1.Position(),
			g.drawer.FlipYAxis,
		),
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Initialising Chipmunk
	space := cp.NewSpace()
	space.SleepTimeThreshold = 0.5
	space.SetGravity(cp.Vector{X: 0, Y: -100})
	walls := []cp.Vector{
		{X: -hScreenWidth, Y: -hScreenHeight}, {X: -hScreenWidth, Y: hScreenHeight},
		{X: hScreenWidth, Y: -hScreenHeight}, {X: hScreenWidth, Y: hScreenHeight},
		{X: -hScreenWidth, Y: -hScreenHeight}, {X: hScreenWidth, Y: -hScreenHeight},
		{X: -hScreenWidth, Y: hScreenHeight}, {X: hScreenWidth, Y: hScreenHeight},
		{X: -100, Y: -100}, {X: 100, Y: -80},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		shape := space.AddShape(cp.NewSegment(space.StaticBody, walls[i], walls[i+1], 10))
		shape.SetElasticity(0.5)
		shape.SetFriction(0.5)
	}
	ball1 := addBall(space, 0, 0, 50)
	addBall(space, 80, 100, 20)
	addBall(space, -100, 150, 40)

	// Initialising Ebitengine/v2
	game := &Game{}
	game.space = space
	game.drawer = ebitencp.NewDrawer(screenWidth, screenHeight)
	game.drawer.FlipYAxis = false
	game.ball1 = ball1
	game.flipYAxis = false
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - camera")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func addBall(space *cp.Space, x, y, radius float64) *cp.Body {
	mass := radius * radius / 100.0
	body := space.AddBody(cp.NewBody(mass, cp.MomentForCircle(mass, 0, radius, cp.Vector{})))
	body.SetPosition(cp.Vector{X: x, Y: y})
	shape := space.AddShape(cp.NewCircle(body, radius, cp.Vector{}))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
	return body
}
