package renderer

import (
	"main/physics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SpriteRenderer struct {
	Sprite rl.Texture2D
	Color  rl.Color
	Scale  float32
}

func (renderer *SpriteRenderer) Draw(Position rl.Vector2, angle float32) {
	sourceRect := rl.NewRectangle(0, 0, float32(renderer.Sprite.Width), float32(renderer.Sprite.Height))
	destRect := rl.NewRectangle(Position.X, Position.Y, float32(renderer.Sprite.Width)*renderer.Scale, float32(renderer.Sprite.Height)*renderer.Scale)
	origin := rl.Vector2Scale(rl.NewVector2(float32(renderer.Sprite.Width)/2, float32(renderer.Sprite.Height)/2), renderer.Scale)
	rl.DrawTexturePro(renderer.Sprite, sourceRect,
		destRect,
		origin, angle, renderer.Color)
}

// Display entity with rectangle hit box
func DrawRectEntity(body *physics.RectangleBody, renderer *SpriteRenderer, Width, Height, angle float32) {
	renderer.Draw(body.Position, angle)
}

// Display entity with circle hit box
func DrawCircleEntity(body *physics.CircleBody, renderer *SpriteRenderer, angle float32) {
	renderer.Draw(body.Position, angle)
	// renderer.DrawCircleBoundary(body.Position, body.Radius)
}

func (renderer *SpriteRenderer) DrawCircleBoundary(Position rl.Vector2, Radius float32) {
	rl.DrawCircleLines(int32(Position.X), int32(Position.Y), Radius, rl.SkyBlue)
}
