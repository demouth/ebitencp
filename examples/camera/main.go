package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"
	"math/rand/v2"

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
	camera    Camera
}

type Camera struct {
	Offset cp.Vector
	Zoom   float64
	Rotate float64
}

func (g *Game) Update() error {
	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)

	// Camera.Offset is deprecated
	// g.drawer.Camera.Offset.X = g.ball1.Position().X
	// g.drawer.Camera.Offset.Y = g.ball1.Position().Y

	g.drawer.GeoM.Reset()
	g.drawer.GeoM.Translate(-g.ball1.Position().X, -g.ball1.Position().Y)
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.camera = Camera{Zoom: 1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camera.Offset.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camera.Offset.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camera.Offset.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camera.Offset.Y += 2
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.drawer.FlipYAxis = !g.drawer.FlipYAxis
	}
	g.drawer.GeoM.Scale(g.camera.Zoom, g.camera.Zoom)
	g.drawer.GeoM.Rotate(g.camera.Rotate)
	g.drawer.GeoM.Translate(g.camera.Offset.X, g.camera.Offset.Y)

	g.space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cp.DrawSpace(g.space, g.drawer.WithScreen(screen))

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
  Drag Object = Cursor
  Flip Y axis = SPACE`,
			g.camera.Offset,
			g.camera.Zoom,
			g.camera.Rotate,
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
		{X: -100, Y: -200}, {X: 100, Y: -180},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		shape := space.AddShape(cp.NewSegment(space.StaticBody, walls[i], walls[i+1], 10))
		shape.SetElasticity(0.5)
		shape.SetFriction(0.5)
	}
	addBall(space, 80, 100, 10)
	addBall(space, -100, 150, 20)
	ball1 := addBall(space, 0, 0, 25)
	addChains(space)
	addPentagon(space)

	// Initialising Ebitengine/v2
	game := &Game{}
	game.space = space
	game.drawer = ebitencp.NewDrawer(screenWidth, screenHeight)
	game.drawer.FlipYAxis = false
	game.ball1 = ball1
	game.flipYAxis = false
	game.camera = Camera{Zoom: 1}
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
func addWall(space *cp.Space, x1, y1, x2, y2, radius float64) {
	pos1 := cp.Vector{X: x1, Y: y1}
	pos2 := cp.Vector{X: x2, Y: y2}
	shape := space.AddShape(cp.NewSegment(space.StaticBody, pos1, pos2, radius))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
}
func addChains(space *cp.Space) {
	const (
		CHAIN_COUNT = 8
		LINK_COUNT  = 10
	)
	var (
		mass    = 1.0
		width   = 20.0
		height  = 30.0
		spacing = width * 0.3
	)
	BreakableJointPostStepRemove := func(space *cp.Space, joint interface{}, _ interface{}) {
		space.RemoveConstraint(joint.(*cp.Constraint))
	}
	BreakableJointPostSolve := func(joint *cp.Constraint, space *cp.Space) {
		dt := space.TimeStep()
		// Convert the impulse to a force by dividing it by the timestep.
		force := joint.Class.GetImpulse() / dt
		maxForce := joint.MaxForce()
		// If the force is almost as big as the joint's max force, break it.
		if force > 0.9*maxForce {
			space.AddPostStepCallback(BreakableJointPostStepRemove, joint, nil)
		}
	}

	var i, j float64
	for i = 0; i < CHAIN_COUNT; i++ {
		var prev *cp.Body

		for j = 0; j < LINK_COUNT; j++ {
			pos := cp.Vector{X: 40 * (i - (CHAIN_COUNT-1)/2.0), Y: 240 - (j+0.5)*height - (j+1)*spacing}

			body := space.AddBody(cp.NewBody(mass, cp.MomentForBox(mass, width, height)))
			body.SetPosition(pos)

			shape := space.AddShape(cp.NewSegment(body, cp.Vector{X: 0, Y: (height - width) / 2}, cp.Vector{X: 0, Y: (width - height) / 2}, width/2))
			shape.SetFriction(0.8)

			breakingForce := 80000.0

			var constraint *cp.Constraint
			if prev == nil {
				constraint = space.AddConstraint(cp.NewSlideJoint(body, space.StaticBody, cp.Vector{X: 0, Y: height / 2}, cp.Vector{X: pos.X, Y: 240}, 0, spacing))
			} else {
				constraint = space.AddConstraint(cp.NewSlideJoint(body, prev, cp.Vector{X: 0, Y: height / 2}, cp.Vector{X: 0, Y: -height / 2}, 0, spacing))
			}

			constraint.SetMaxForce(breakingForce)
			constraint.PostSolve = BreakableJointPostSolve
			constraint.SetCollideBodies(false)

			prev = body
		}
	}
}

func addPentagon(space *cp.Space) {
	const numVerts = 5
	var (
		pentagonMass   = 0.0
		pentagonMoment = 0.0
		body           *cp.Body
		shape          *cp.Shape
	)

	verts := []cp.Vector{}
	for i := 0; i < numVerts; i++ {
		angle := -2.0 * math.Pi * float64(i) / numVerts
		verts = append(verts, cp.Vector{X: 20 * math.Cos(angle), Y: 20 * math.Sin(angle)})
	}

	pentagonMass = 1.0
	pentagonMoment = cp.MomentForPoly(1, numVerts, verts, cp.Vector{}, 0)

	for i := 0; i < 10; i++ {
		body = space.AddBody(cp.NewBody(pentagonMass, pentagonMoment))
		x := rand.Float64()*640 - 320
		body.SetPosition(cp.Vector{X: x, Y: 0})

		shape = space.AddShape(cp.NewPolyShape(body, numVerts, verts, cp.NewTransformIdentity(), 0))
		shape.SetElasticity(0)
		shape.SetFriction(0.4)
	}
}
