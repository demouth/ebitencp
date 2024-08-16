package main

// This is based on "jakecoffman/cp-examples/plink".

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/demouth/ebitencp"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	pentagonMass   = 0.0
	pentagonMoment = 0.0
)

const numVerts = 5

type Game struct {
	space  *cp.Space
	drawer *ebitencp.Drawer
}

func (g *Game) Update() error {
	g.space.EachBody(func(body *cp.Body) {
		pos := body.Position()
		if pos.Y < -260 || math.Abs(pos.X) > 340 {
			x := rand.Float64()*640 - 320
			body.SetPosition(cp.Vector{X: x, Y: 260})
		}
	})
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

func main() {
	// Initialising Chipmunk

	space := cp.NewSpace()
	space.Iterations = 5
	space.SetGravity(cp.Vector{X: 0, Y: -100})

	var body *cp.Body
	var shape *cp.Shape

	tris := []cp.Vector{
		{X: -15, Y: -15},
		{X: 0, Y: 10},
		{X: 15, Y: -15},
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 6; j++ {
			stagger := (j % 2) * 40
			offset := cp.Vector{X: float64(i*80 - 320 + stagger), Y: float64(j*70 - 240)}
			shape = space.AddShape(cp.NewPolyShape(space.StaticBody, 3, tris, cp.NewTransformTranslate(offset), 0))
			shape.SetElasticity(1)
			shape.SetFriction(1)
		}
	}

	verts := []cp.Vector{}
	for i := 0; i < numVerts; i++ {
		angle := -2.0 * math.Pi * float64(i) / numVerts
		verts = append(verts, cp.Vector{X: 10 * math.Cos(angle), Y: 10 * math.Sin(angle)})
	}

	pentagonMass = 1.0
	pentagonMoment = cp.MomentForPoly(1, numVerts, verts, cp.Vector{}, 0)

	for i := 0; i < 300; i++ {
		body = space.AddBody(cp.NewBody(pentagonMass, pentagonMoment))
		x := rand.Float64()*640 - 320
		body.SetPosition(cp.Vector{X: x, Y: 350})

		shape = space.AddShape(cp.NewPolyShape(body, numVerts, verts, cp.NewTransformIdentity(), 0))
		shape.SetElasticity(0)
		shape.SetFriction(0.4)
	}

	// Initialising Ebitengine/v2

	game := &Game{}
	game.space = space
	game.drawer = ebitencp.NewDrawer(screenWidth, screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - plink")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
