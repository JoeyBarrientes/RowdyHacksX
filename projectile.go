package main

import (
	"main/physics"
	"main/renderer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	Body     physics.CircleBody
	Position rl.Vector2
	Speed    float32
	Color    rl.Color
}

type SpriteProjectile struct {
	Body     physics.CircleBody
	Sprite   renderer.CharacterSprite
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

// Initializes Projectile entity
func (p *Projectile) NewSpriteProjectile(Sprite rl.Texture2D, Position rl.Vector2, Velocity rl.Vector2, Speed float32, Radius float32, Angle float32, Color rl.Color, Scale float32) SpriteProjectile {
	pb := physics.NewCirclePhysicsBody(Velocity, Radius, Angle)
	nb := SpriteProjectile{
		Body: pb,
		Sprite: renderer.CharacterSprite{
			Render: renderer.SpriteRenderer{
				Sprite: Sprite,
				Color:  Color,
				Angle:  0,
				Scale:  Scale,
			},
			SourceRect:  rl.NewRectangle(0, 0, 22, 22),
			IsMoving:    false,
			SpriteFrame: 0,
		},
		Position: Position,
		Speed:    Speed,
		Color:    Color,
	}

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

func (p SpriteProjectile) Draw() {
	rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), p.Body.Radius, p.Color)
}

func (p *SpriteProjectile) VelocityTick() {
	adjustedVel := rl.Vector2Scale(p.Body.Velocity, rl.GetFrameTime())
	p.Position = rl.Vector2Add(p.Position, adjustedVel)
}

func (p *SpriteProjectile) PhysicsUpdate() {
	p.VelocityTick()
}

func (p *SpriteProjectile) DrawSprite() {
	// sourceRect := rl.NewRectangle(0, 0, float32(renderer.Sprite.Width), float32(renderer.Sprite.Height))

	// sourceRect = rl.NewRectangle(0, 0, 48, 48)
	// if character.IsMoving {
	// 	sourceRect.X = sourceRect.Width * float32(character.SpriteFrame)
	// }
	// fmt.Println(sourceRect)

	// sourceRect.Y = sourceRect.Height * float32(character.Direction.Facing)
	// destRect := rl.NewRectangle(Position.X, Position.Y, float32(renderer.Sprite.Width)*renderer.Scale, float32(renderer.Sprite.Height)*renderer.Scale)
	destRect := rl.NewRectangle(p.Position.X, p.Position.Y, 100*p.Sprite.Render.Scale, 100*p.Sprite.Render.Scale)
	origin := rl.Vector2Scale(rl.NewVector2(float32(p.Sprite.Render.Sprite.Width)/2, float32(p.Sprite.Render.Sprite.Height)/2), p.Sprite.Render.Scale)
	rl.DrawTexturePro(p.Sprite.Render.Sprite, p.Sprite.SourceRect,
		destRect,
		origin, p.Sprite.Render.Angle, p.Sprite.Render.Color)
}
