package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SpriteRenderer struct {
	Sprite rl.Texture2D
	Color  rl.Color
	Angle  float32
	Scale  float32
}

type StillSprite struct {
	Position rl.Vector2
	Sprite   rl.Texture2D
	Color    rl.Color
	Angle    float32
	Scale    float32
}

type CharacterSprite struct {
	Render      SpriteRenderer
	SourceRect  rl.Rectangle
	IsMoving    bool
	SpriteFrame int
}

func NewStillSprite(Position rl.Vector2, Sprite rl.Texture2D, Color rl.Color, Angle, Scale float32) StillSprite {
	Character := StillSprite{
		Position: Position,
		Sprite:   Sprite,
		Color:    Color,
		Angle:    Angle,
		Scale:    Scale,
	}
	return Character
}

func (renderer *SpriteRenderer) Draw(Position rl.Vector2, angle float32) {
	sourceRect := rl.NewRectangle(0, 0, float32(renderer.Sprite.Width), float32(renderer.Sprite.Height))
	destRect := rl.NewRectangle(Position.X, Position.Y, float32(renderer.Sprite.Width)*renderer.Scale, float32(renderer.Sprite.Height)*renderer.Scale)
	origin := rl.Vector2Scale(rl.NewVector2(float32(renderer.Sprite.Width)/2, float32(renderer.Sprite.Height)/2), renderer.Scale)
	rl.DrawTexturePro(renderer.Sprite, sourceRect,
		destRect,
		origin, angle, renderer.Color)
}

func (renderer *StillSprite) Draw(angle float32) {
	sourceRect := rl.NewRectangle(0, 0, float32(renderer.Sprite.Width), float32(renderer.Sprite.Height))
	destRect := rl.NewRectangle(renderer.Position.X, renderer.Position.Y, float32(renderer.Sprite.Width)*renderer.Scale, float32(renderer.Sprite.Height)*renderer.Scale)
	origin := rl.Vector2Scale(rl.NewVector2(float32(renderer.Sprite.Width)/2, float32(renderer.Sprite.Height)/2), renderer.Scale)
	rl.DrawTexturePro(renderer.Sprite, sourceRect,
		destRect,
		origin, angle, renderer.Color)
}

// Display entity with rectangle hit box
func DrawRectEntity(position rl.Vector2, renderer *SpriteRenderer, Width, Height, angle float32) {
	renderer.Draw(position, angle)
}

func (renderer *SpriteRenderer) DrawCircleBoundary(Position rl.Vector2, Radius float32) {
	rl.DrawCircleLines(int32(Position.X), int32(Position.Y), Radius, rl.SkyBlue)
}
