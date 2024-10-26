package physics

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RectangleBody struct {
	Velocity rl.Vector2
	Width    float32
	Height   float32
	Angle    float32
}

// Defines circular body and hitbox
type CircleBody struct {
	Velocity rl.Vector2
	Radius   float32
	Angle    float32
}

// Initializes rectangle body
func NewRectanglePhysicsBody(velocity rl.Vector2, width, height float32, angle float32) RectangleBody {
	pb := RectangleBody{
		Velocity: velocity,
		Width:    width,
		Height:   height,
		Angle:    angle}
	return pb
}

// Initializes circle body
func NewCirclePhysicsBody(velocity rl.Vector2, radius float32, angle float32) CircleBody {
	pb := CircleBody{
		Velocity: velocity,
		Radius:   radius,
		Angle:    angle}
	return pb
}

// func (rb *RectangleBody) VelocityTick() {
// 	adjustedVel := rl.Vector2Scale(rb.Velocity, rl.GetFrameTime())
// 	rb.Position = rl.Vector2Add(rb.Position, adjustedVel)
// }

// func (rb *RectangleBody) PhysicsUpdate() {
// 	rb.VelocityTick()
// }

// func (cb *CircleBody) VelocityTick() {
// 	adjustedVel := rl.Vector2Scale(cb.Velocity, rl.GetFrameTime())
// 	cb.Position = rl.Vector2Add(cb.Position, adjustedVel)
// }

// func (cb *CircleBody) PhysicsUpdate() {
// 	cb.VelocityTick()
// }
