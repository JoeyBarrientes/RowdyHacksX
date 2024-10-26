package main

import (
	"main/physics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	Body  physics.CircleBody
	Speed float32
	Color rl.Color
}

// Initializes Projectile entity
func (p *Projectile) NewProjectile(Position rl.Vector2, Velocity rl.Vector2, Speed float32, Radius float32, Angle float32, color rl.Color) Projectile {
	pb := physics.NewCirclePhysicsBody(Position, Velocity, Radius, Angle)
	nb := Projectile{Body: pb, Speed: Speed, Color: color}

	return nb
}

func (p Projectile) Draw() {
	rl.DrawCircle(int32(p.Body.Position.X), int32(p.Body.Position.Y), p.Body.Radius, p.Color)
}
