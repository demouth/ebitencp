package main

// This is based on "jakecoffman/cp-examples/bench".

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/demouth/ebitencp"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/jakecoffman/cp/v2"
)

const (
	screenWidth   = 640
	screenHeight  = 480
	hScreenWidth  = screenWidth / 2
	hScreenHeight = screenHeight / 2
)

type Game struct {
	space  *cp.Space
	drawer *ebitencp.Drawer
}

func (g *Game) Update() error {
	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)
	g.space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cp.DrawSpace(g.space, g.drawer.WithScreen(screen))

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"FPS: %0.2f",
		ebiten.ActualFPS(),
	))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	ball *cp.Body
)

func main() {
	// Initialising Chipmunk
	space := cp.NewSpace()
	space.SleepTimeThreshold = 0.5
	space.SetGravity(cp.Vector{X: 0, Y: -100})
	simpleTerrain(space)
	var r float64 = 6.0
	for i := 0; i < 100; i++ {
		addBall(space, float64(i%10)*r*2, float64(i/10)*r*2, r)
	}

	// Initialising Ebitengine/v2
	game := &Game{}
	game.space = space
	game.drawer = ebitencp.NewDrawer(screenWidth, screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - bench")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func addBall(space *cp.Space, x, y, radius float64) {
	mass := radius * radius / 100.0
	body := space.AddBody(cp.NewBody(mass, cp.MomentForCircle(mass, 0, radius, cp.Vector{})))
	ball = body
	body.SetPosition(cp.Vector{X: x, Y: y})
	shape := space.AddShape(cp.NewCircle(body, radius, cp.Vector{}))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
}

func simpleTerrain(space *cp.Space) *cp.Space {
	var simpleTerrainVerts = []cp.Vector{
		{X: 350.00, Y: 425.07}, {X: 336.00, Y: 436.55}, {X: 272.00, Y: 435.39}, {X: 258.00, Y: 427.63}, {X: 225.28, Y: 420.00}, {X: 202.82, Y: 396.00},
		{X: 191.81, Y: 388.00}, {X: 189.00, Y: 381.89}, {X: 173.00, Y: 380.39}, {X: 162.59, Y: 368.00}, {X: 150.47, Y: 319.00}, {X: 128.00, Y: 311.55},
		{X: 119.14, Y: 286.00}, {X: 126.84, Y: 263.00}, {X: 120.56, Y: 227.00}, {X: 141.14, Y: 178.00}, {X: 137.52, Y: 162.00}, {X: 146.51, Y: 142.00},
		{X: 156.23, Y: 136.00}, {X: 158.00, Y: 118.27}, {X: 170.00, Y: 100.77}, {X: 208.43, Y: 84.00}, {X: 224.00, Y: 69.65}, {X: 249.30, Y: 68.00},
		{X: 257.00, Y: 54.77}, {X: 363.00, Y: 45.94}, {X: 374.15, Y: 54.00}, {X: 386.00, Y: 69.60}, {X: 413.00, Y: 70.73}, {X: 456.00, Y: 84.89},
		{X: 468.09, Y: 99.00}, {X: 467.09, Y: 123.00}, {X: 464.92, Y: 135.00}, {X: 469.00, Y: 141.03}, {X: 497.00, Y: 148.67}, {X: 513.85, Y: 180.00},
		{X: 509.56, Y: 223.00}, {X: 523.51, Y: 247.00}, {X: 523.00, Y: 277.00}, {X: 497.79, Y: 311.00}, {X: 478.67, Y: 348.00}, {X: 467.90, Y: 360.00},
		{X: 456.76, Y: 382.00}, {X: 432.95, Y: 389.00}, {X: 417.00, Y: 411.32}, {X: 373.00, Y: 433.19}, {X: 361.00, Y: 430.02}, {X: 350.00, Y: 425.07},
	}
	offset := cp.Vector{X: -320, Y: -240}
	for i := 0; i < len(simpleTerrainVerts)-1; i++ {
		a := simpleTerrainVerts[i]
		b := simpleTerrainVerts[i+1]
		space.AddShape(cp.NewSegment(space.StaticBody, a.Add(offset), b.Add(offset), 0))
	}
	return space
}
