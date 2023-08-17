package main

import (
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	adjRate                   = .005
	screenWidth, screenHeight = 640, 360
	boidsCount                = 500
	viewRadius                = 13
	screenScale               = 2
	boidSpeedTimeout          = 10
)

var (
	green   = color.RGBA{10, 255, 50, 255}
	boids   [boidsCount]*Boid
	boidMap [screenWidth + 1][screenHeight + 1]int
	lock    = sync.RWMutex{}
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	for _, boid := range boids {
		boid.draw(screen)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for i := 0; i < boidsCount; i++ {
		createBoid(i)
	}
	ebiten.SetWindowSize(screenWidth*screenScale, screenHeight*screenScale)
	ebiten.SetWindowTitle("Boids in a box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
