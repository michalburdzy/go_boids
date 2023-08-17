package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Boid struct {
	id       int
	position Vector2D
	velocity Vector2D
}

func (b *Boid) MoveOne() {
	accel := b.calcAcceleration()
	lock.Lock()
	b.velocity = b.velocity.add(accel).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	next := b.position.add(b.velocity)

	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
	lock.Unlock()
}

func (b *Boid) start() {
	for {
		b.MoveOne()
		time.Sleep(boidSpeedTimeout * time.Millisecond)
	}
}

func (b *Boid) draw(screen *ebiten.Image) {
	screen.Set(int(b.position.x+1), int(b.position.y), green)
	screen.Set(int(b.position.x-1), int(b.position.y), green)
	screen.Set(int(b.position.x), int(b.position.y-1), green)
	screen.Set(int(b.position.x), int(b.position.y+1), green)

}

func createBoid(id int) {
	b := Boid{
		id:       id,
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
	}

	boids[id] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}

func (b *Boid) calcAcceleration() Vector2D {
	accel := Vector2D{b.borderBounce(b.position.x, screenWidth), b.borderBounce(b.position.y, screenHeight)}
	avgVelocity, avgPosition, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0

	minX, maxX, minY, maxY :=
		int(math.Max(0, math.Floor(b.position.x)-viewRadius)),
		int(math.Min(screenWidth, math.Floor(b.position.x)+viewRadius)),
		int(math.Max(math.Floor(b.position.y)-viewRadius, 0)),
		int(math.Min(math.Floor(b.position.y)+viewRadius, float64(screenHeight)))

	lock.RLock()
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			if boidId := boidMap[x][y]; boidId != -1 && boidId != b.id {
				if dist := boids[boidId].position.Distance(b.position); dist < viewRadius {
					avgVelocity = avgVelocity.add(boids[boidId].velocity)
					avgPosition = avgPosition.add(boids[boidId].position)
					separation = separation.add(b.position.subtract(boids[boidId].position).divideV(dist))
					count++
				}
			}
		}
	}
	lock.RUnlock()

	if count > 0 {
		avgVelocity = avgVelocity.divideV(float64(count))
		avgPosition = avgPosition.divideV(float64(count))
		accelAlignment := avgVelocity.subtract(b.velocity).multiplyV(adjRate)
		accelCohesion := avgPosition.subtract(b.position).multiplyV(adjRate)
		accelSeparation := separation.multiplyV(adjRate)
		accel = accel.add(accelAlignment).add(accelCohesion).add(accelSeparation)
	}

	return accel
}

func (b *Boid) borderBounce(position, limit float64) float64 {
	if position < viewRadius {
		return 1 / position
	} else if position > limit-viewRadius {
		return 1 / (position - limit)
	}

	return 0
}
