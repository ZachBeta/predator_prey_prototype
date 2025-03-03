package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	WindowWidth            = 800
	WindowHeight           = 600
	CircleRadius           = 10
	GridCellSize           = 20
	NumCircles             = 50
	MaxAge                 = 1000
	ReproductionDistance   = GridCellSize * 3
	ReproductionChance     = 0.01
	ReproductionCooldown   = 300
	InitialEnergy          = 100
	EnergyLossPerTick      = 0.1
	FoodEnergy             = 50
	NumFoodParticles       = 20
	ReproductionEnergyCost = 0.33
)

type Brain struct {
	targetMate *Circle
	rng        *rand.Rand
}

func NewBrain(rng *rand.Rand) *Brain {
	return &Brain{
		rng: rng,
	}
}

func (b *Brain) think(self *Circle, circles []*Circle) {
	if b.rng.Float64() < 0.01 {
		self.dx = b.rng.Float64()*2 - 1
		self.dy = b.rng.Float64()*2 - 1
		return
	}

	if self.reproductionTime > 0 || !self.alive {
		b.targetMate = nil
		return
	}

	if b.targetMate != nil && (!b.targetMate.alive || b.targetMate.reproductionTime > 0) {
		b.targetMate = nil
	}

	if b.targetMate == nil {
		b.findNewMate(self, circles)
	}

	if b.targetMate != nil {
		b.moveTowardsMate(self)
	}
}

func (b *Brain) findNewMate(self *Circle, circles []*Circle) {
	var closestMate *Circle
	closestDist := math.MaxFloat64

	for _, other := range circles {
		if other == self || !other.alive || other.reproductionTime > 0 {
			continue
		}

		dx := other.x - self.x
		dy := other.y - self.y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < closestDist {
			closestDist = dist
			closestMate = other
		}
	}

	b.targetMate = closestMate
}

func (b *Brain) moveTowardsMate(self *Circle) {
	dx := b.targetMate.x - self.x
	dy := b.targetMate.y - self.y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0 {
		self.dx = (dx / dist) * 2
		self.dy = (dy / dist) * 2
	}
}

type Food struct {
	x, y float64
}

type Circle struct {
	x                float64
	y                float64
	dx               float64
	dy               float64
	alive            bool
	color            color.RGBA
	age              int
	reproductionTime int
	brain            *Brain
	energy           float64
	rng              *rand.Rand
}

func NewCircle(x, y float64, rng *rand.Rand) *Circle {
	return &Circle{
		x:                x,
		y:                y,
		dx:               rng.Float64()*2 - 1,
		dy:               rng.Float64()*2 - 1,
		alive:            true,
		color:            color.RGBA{R: uint8(rng.Intn(256)), G: uint8(rng.Intn(256)), B: uint8(rng.Intn(256)), A: 255},
		age:              0,
		reproductionTime: 0,
		brain:            NewBrain(rng),
		energy:           InitialEnergy,
		rng:              rng,
	}
}

func (c *Circle) update(circles []*Circle) {
	c.brain.think(c, circles)

	// Calculate potential new position
	newX := c.x + c.dx
	newY := c.y + c.dy

	// Check collisions with other circles
	for _, other := range circles {
		if other != c && other.alive {
			dx := newX - other.x
			dy := newY - other.y
			distSquared := dx*dx + dy*dy
			minDist := float64(CircleRadius * 2)

			if distSquared < minDist*minDist {
				// Calculate collision response
				dist := math.Sqrt(distSquared)
				if dist == 0 {
					// If exactly overlapping, nudge slightly
					newX += c.rng.Float64() - 0.5
					newY += c.rng.Float64() - 0.5
					continue
				}

				// Normal vector of collision
				nx := dx / dist
				ny := dy / dist

				// Separate the circles
				overlap := minDist - dist
				newX += nx * overlap * 0.5
				newY += ny * overlap * 0.5

				// Bounce velocity along normal
				dotProduct := c.dx*nx + c.dy*ny
				c.dx -= 2 * dotProduct * nx
				c.dy -= 2 * dotProduct * ny

				// Add small random variation to prevent getting stuck
				c.dx += (c.rng.Float64()*0.2 - 0.1)
				c.dy += (c.rng.Float64()*0.2 - 0.1)
			}
		}
	}

	// Boundary checking
	if newX < CircleRadius {
		newX = CircleRadius
		c.dx *= -1
	} else if newX > WindowWidth-CircleRadius {
		newX = WindowWidth - CircleRadius
		c.dx *= -1
	}

	if newY < CircleRadius {
		newY = CircleRadius
		c.dy *= -1
	} else if newY > WindowHeight-CircleRadius {
		newY = WindowHeight - CircleRadius
		c.dy *= -1
	}

	c.x = newX
	c.y = newY
}

type Game struct {
	circles   []*Circle
	food      []*Food
	rng       *rand.Rand
	ticks     int
	startTime time.Time
}

func (g *Game) respawnFood() {
	for i, f := range g.food {
		if f == nil && g.rng.Float64() < 0.1 {
			g.food[i] = &Food{
				x: g.rng.Float64() * WindowWidth,
				y: g.rng.Float64() * WindowHeight,
			}
		}
	}
}

func (g *Game) Update() error {
	g.ticks++
	// Update positions and age
	for _, circle := range g.circles {
		if circle.alive {
			circle.update(g.circles)
			circle.age++
			circle.energy -= EnergyLossPerTick

			// Die if no energy
			if circle.energy <= 0 {
				circle.alive = false
			}

			// Check for food collision
			for i, f := range g.food {
				if f != nil {
					dx := circle.x - f.x
					dy := circle.y - f.y
					if dx*dx+dy*dy < CircleRadius*CircleRadius {
						circle.energy += FoodEnergy
						g.food[i] = nil
					}
				}
			}

			if circle.reproductionTime > 0 {
				circle.reproductionTime--
			}
			if circle.age > MaxAge {
				circle.alive = false
			}
		}
	}

	// Respawn food
	g.respawnFood()

	// Handle reproduction
	newCircles := []*Circle{}
	for i, c1 := range g.circles {
		if !c1.alive || c1.reproductionTime > 0 || c1.energy < InitialEnergy*ReproductionEnergyCost {
			continue
		}

		for j, c2 := range g.circles {
			if i != j && c2.alive && c2.reproductionTime == 0 && c2.energy >= InitialEnergy*ReproductionEnergyCost {
				dx := c1.x - c2.x
				dy := c1.y - c2.y
				dist := dx*dx + dy*dy

				if dist < ReproductionDistance*ReproductionDistance && g.rng.Float64() < ReproductionChance {
					newX := (c1.x + c2.x) / 2
					newY := (c1.y + c2.y) / 2
					newCircle := NewCircle(newX, newY, g.rng)

					// Transfer energy from parents to child
					energyContribution := InitialEnergy * ReproductionEnergyCost
					c1.energy -= energyContribution
					c2.energy -= energyContribution
					newCircle.energy = energyContribution * 2

					newCircles = append(newCircles, newCircle)
					c1.reproductionTime = ReproductionCooldown
					c2.reproductionTime = ReproductionCooldown
				}
			}
		}
	}

	g.circles = append(g.circles, newCircles...)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// Draw food
	for _, f := range g.food {
		if f != nil {
			drawCircle(screen, int(f.x), int(f.y), 3, color.RGBA{0, 255, 0, 255})
		}
	}

	for _, circle := range g.circles {
		if circle.alive {
			drawCircle(screen, int(circle.x), int(circle.y), CircleRadius, circle.color)
		}
	}

	// Draw stats overlay
	aliveCount := 0
	for _, c := range g.circles {
		if c.alive {
			aliveCount++
		}
	}

	foodCount := 0
	for _, f := range g.food {
		if f != nil {
			foodCount++
		}
	}

	elapsed := time.Since(g.startTime).Milliseconds()
	stats := []string{
		fmt.Sprintf("Ticks: %d", g.ticks),
		fmt.Sprintf("Time: %dms", elapsed),
		fmt.Sprintf("Alive: %d", aliveCount),
		fmt.Sprintf("Food: %d", foodCount),
	}

	for i, text := range stats {
		ebitenutil.DebugPrintAt(screen, text, 10, 20*i+10)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

func drawCircle(screen *ebiten.Image, x, y, radius int, c color.Color) {
	for py := -radius; py <= radius; py++ {
		for px := -radius; px <= radius; px++ {
			if px*px+py*py <= radius*radius {
				screen.Set(x+px, y+py, c)
			}
		}
	}
}

func main() {
	rng := rand.New(rand.NewSource(1234))
	game := &Game{
		circles:   make([]*Circle, NumCircles),
		food:      make([]*Food, NumFoodParticles),
		rng:       rng,
		startTime: time.Now(),
	}

	// Initialize circles
	for i := range game.circles {
		x := rng.Float64() * WindowWidth
		y := rng.Float64() * WindowHeight
		game.circles[i] = NewCircle(x, y, rng)
	}

	// Initialize food
	for i := range game.food {
		game.food[i] = &Food{
			x: rng.Float64() * WindowWidth,
			y: rng.Float64() * WindowHeight,
		}
	}

	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Moving Circles Life Simulation")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
