package main

import (
	"main/physics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	Body     physics.CircleBody
	Position rl.Vector2
	Speed    float32
	Color    rl.Color
}

// Initializes Projectile entity
func (p *Projectile) NewProjectile(Position rl.Vector2, Velocity rl.Vector2, Speed float32, Radius float32, Angle float32, color rl.Color) Projectile {
	pb := physics.NewCirclePhysicsBody(Velocity, Radius, Angle)
	nb := Projectile{Body: pb, Position: Position, Speed: Speed, Color: color}

	return nb
}

func (p Projectile) Draw() {
	rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), p.Body.Radius, p.Color)
}

func (p *Projectile) VelocityTick() {
	adjustedVel := rl.Vector2Scale(p.Body.Velocity, rl.GetFrameTime())
	p.Position = rl.Vector2Add(p.Position, adjustedVel)
}

func (p *Projectile) PhysicsUpdate() {
	p.VelocityTick()
}
