package main

import (
	"fmt"
	"main/physics"
	"main/renderer"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vehicle struct {
	Body           physics.RectangleBody
	Sprite         renderer.CharacterSprite
	Position       rl.Vector2
	Bullets        []Projectile
	BulletVelocity rl.Vector2
	SlowingDown    bool
	Speed          float32
	Acceleration   float32
	SpeedBar
}

type SpeedBar struct {
	SpeedTracker float32
	TextColor    rl.Color
	ProgressBar
}

// func NewVehicle(Position rl.Vector2, Sprite rl.Texture2D, Speed float32) {

// }

// Initializes vehicle entity
func NewVehicle(Sprite rl.Texture2D, Color rl.Color, Position rl.Vector2, Velocity rl.Vector2, Width, Height, Scale float32) Vehicle {
	Vehicle := Vehicle{
		Body: physics.RectangleBody{
			Velocity: Velocity,
			Width:    Width,
			Height:   Height,
			Angle:    0,
		},
		Sprite: renderer.CharacterSprite{
			Render: renderer.SpriteRenderer{
				Sprite: Sprite,
				Color:  Color,
				Scale:  Scale,
			},
			SourceRect:  rl.NewRectangle(0, 0, 64, 64),
			IsMoving:    false,
			SpriteFrame: 0,
		}, SpeedBar: SpeedBar{
			SpeedTracker: 0,
			TextColor:    rl.White,
			ProgressBar:  NewProgressBar(int32(screenSize.X), 25, 4, 600, 70, &theme),
		},
		Speed:          1,
		Acceleration:   1,
		Position:       Position,
		Bullets:        []Projectile{},
		BulletVelocity: rl.NewVector2(100, 100),
		SlowingDown:    false,
	}
	return Vehicle
}

func (vehicle *Vehicle) DrawCharacter() {
	// sourceRect := rl.NewRectangle(0, 0, float32(renderer.Sprite.Width), float32(renderer.Sprite.Height))

	// sourceRect = rl.NewRectangle(0, 0, 48, 48)
	// if character.IsMoving {
	// 	sourceRect.X = sourceRect.Width * float32(character.SpriteFrame)
	// }
	// fmt.Println(sourceRect)

	// sourceRect.Y = sourceRect.Height * float32(character.Direction.Facing)
	// destRect := rl.NewRectangle(Position.X, Position.Y, float32(renderer.Sprite.Width)*renderer.Scale, float32(renderer.Sprite.Height)*renderer.Scale)
	destRect := rl.NewRectangle(vehicle.Position.X, vehicle.Position.Y, 100*vehicle.Sprite.Render.Scale, 100*vehicle.Sprite.Render.Scale)
	origin := rl.Vector2Scale(rl.NewVector2(float32(vehicle.Sprite.Render.Sprite.Width)/2, float32(vehicle.Sprite.Render.Sprite.Height)/2), vehicle.Sprite.Render.Scale)
	rl.DrawTexturePro(vehicle.Sprite.Render.Sprite, vehicle.Sprite.SourceRect,
		destRect,
		origin, vehicle.Sprite.Render.Angle, vehicle.Sprite.Render.Color)
}

// func DrawRectEntity(Position rl.Vector2, body *physics.RectangleBody, renderer *SpriteRenderer, Width, Height, angle float32) {
// 	vehicle.DrawCharacter(Position, angle)
// }

func (vehicle *Vehicle) VelocityTick() {
	adjustedVel := rl.Vector2Scale(vehicle.Body.Velocity, rl.GetFrameTime())
	vehicle.Position = rl.Vector2Add(vehicle.Position, adjustedVel)
}

func (vehicle *Vehicle) PhysicsUpdate() {
	vehicle.VelocityTick()
}

// Creates and shoots a projectile in the direction of the ship
func (vehicle *Vehicle) shoot() {
	offset := rl.NewVector2(vehicle.Position.X-300, vehicle.Position.Y+40)
	mouse := rl.GetMousePosition()
	direction := rl.Vector2Subtract(mouse, offset)
	direction = normalizeVector(direction)
	newBullet := Projectile{
		Body:     physics.NewCirclePhysicsBody(vehicle.BulletVelocity, 15, 0),
		Position: offset,
		Speed:    700,
		Color:    rl.White,
	}
	// bullet := vehicle.Bullets.NewProjectile(vehicle.Position, vehicle.Body.Velocity, vehicle.Projectile.Speed, 10, 0, rl.White)
	newBullet.Body.Velocity.X = direction.X * newBullet.Speed
	newBullet.Body.Velocity.Y = direction.Y * newBullet.Speed
	vehicle.Bullets = append(vehicle.Bullets, newBullet)
	// ve = append(*projectiles, pr)
}

// Draws each mine to screen
func (vehicle *Vehicle) drawBullets() {
	for i := len(vehicle.Bullets) - 1; i >= 0; i-- {
		bullet := &(vehicle.Bullets)[i]
		rl.DrawCircle(int32(bullet.Position.X), int32(bullet.Position.Y), bullet.Body.Radius, bullet.Color)
	}
}

// Draws each mine to screen
func (vehicle *Vehicle) updateBullets() {
	for i := len(vehicle.Bullets) - 1; i >= 0; i-- {
		bullet := &(vehicle.Bullets)[i]
		bullet.PhysicsUpdate()
	}
}

// Removes projectile from slice when it goes offscreen
func (vehicle *Vehicle) despawnBullets() {
	for i := len(vehicle.Bullets) - 1; i >= 0; i-- {
		bullet := &(vehicle.Bullets)[i]

		if bullet.Position.X < 0 ||
			bullet.Position.X > screenSize.X ||
			bullet.Position.Y < 0 ||
			bullet.Position.Y > screenSize.Y {

			// Remove the projectile at index i
			vehicle.Bullets = append(vehicle.Bullets[:i], vehicle.Bullets[i+1:]...)
		}
	}
}

func (vehicle *Vehicle) increaseSpeed() {
	if !vehicle.SlowingDown {
		// vehicle.Speed += rl.GetFrameTime() * 2.5
		vehicle.Acceleration = rl.Clamp(vehicle.Acceleration+0.5*rl.GetFrameTime(), 0, 12)
		fmt.Println(vehicle.Acceleration)
		vehicle.Speed = rl.Clamp(vehicle.Speed+rl.GetFrameTime()*vehicle.Acceleration, 0, 200)
	}
}

func (vehicle *Vehicle) decreaseSpeed() {
	if rl.IsKeyPressed(rl.KeySpace) {
		vehicle.SlowingDown = true
	}
	if vehicle.SlowingDown {
		vehicle.Speed -= 14
		vehicle.Acceleration = rl.Clamp(vehicle.Acceleration/2, 2, 8)
		vehicle.SlowingDown = false
	}
}

func (vehicle *Vehicle) drawSpeed() {
	vehicle.ProgressBar.DrawBar()
	speedText := strconv.Itoa(int(vehicle.Speed)) + " MPH"

	if vehicle.SpeedTracker > 0.341 {
		vehicle.TextColor = rl.Black
	} else {
		vehicle.TextColor = rl.White
	}

	rl.DrawText(speedText, vehicle.SpeedBar.X+25, 40, 50, vehicle.TextColor)

}

func (vehicle *Vehicle) updateSpeed() {
	vehicle.increaseSpeed()
	vehicle.drawSpeed()
	vehicle.updateSpeedBar()
}

func (vehicle *Vehicle) updateSpeedBar() {
	vehicle.SpeedBar.X = int32(screenSize.X) - vehicle.Width - 25
	vehicle.SpeedTracker = (vehicle.Speed / 88)
	vehicle.SpeedBar.SetProgress(vehicle.SpeedTracker)
}
