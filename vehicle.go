package main

import (
	"main/physics"
	"main/renderer"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vehicle struct {
	Body           physics.RectangleBody
	Sprite         renderer.CharacterSprite
	Position       rl.Vector2
	BulletSprite   rl.Texture2D
	Bullets        []Projectile
	BulletVelocity rl.Vector2
	SlowingDown    bool
	Speed          float32
	Acceleration   float32
	Hyperjump      bool
	SpeedBar
	Lane
}

type SpeedBar struct {
	SpeedTracker float32
	TextColor    rl.Color
	ProgressBar
}

// func NewVehicle(Position rl.Vector2, Sprite rl.Texture2D, Speed float32) {

// }

// Initializes vehicle entity
func NewVehicle(Sprite rl.Texture2D, BulletSprite rl.Texture2D, Color rl.Color, Position rl.Vector2, Velocity rl.Vector2, Width, Height, Scale float32) Vehicle {
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
		BulletSprite:   BulletSprite,
		BulletVelocity: rl.NewVector2(100, 100),
		SlowingDown:    false,
		Lane:           0,
		Hyperjump:      false,
	}
	return Vehicle
}

func (vehicle *Vehicle) DrawCharacter() {
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
	// offsetCircle := rl.NewCi

	// projOffset := rl.NewVector2(offset.X, offset.Y)
	mouse := rl.GetMousePosition()
	direction := rl.Vector2Subtract(mouse, offset)
	direction = normalizeVector(direction)
	angleInRad := math.Atan2(float64(direction.Y), float64(direction.X))
	angleInDegrees := (angleInRad * (180.0 / math.Pi))
	if angleInDegrees < 0 {
		offset.X -= 30
	} else if angleInDegrees >= 0 {
		offset.X += 50
	}
	newBullet := Projectile{
		Body:     physics.NewCirclePhysicsBody(vehicle.BulletVelocity, 20, 0),
		Position: offset,
		Speed:    700,
		Color:    rl.White,
		Lane:     vehicle.Lane,
	}
	if newBullet.Lane == TOP {
		newBullet.Color = rl.SkyBlue
	} else {
		newBullet.Color = rl.LightGray
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
		if bullet.Lane == TOP {
			bullet.Color = rl.SkyBlue
		} else {
			bullet.Color = rl.LightGray
		}
		// rl.DrawCircleLines(int32(bullet.Position.X), int32(bullet.Position.Y), bullet.Body.Radius, rl.White)
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
		// offset := rl.NewVector2(vehicle.Position.X-300, vehicle.Position.Y+30)
		// rl.DrawCircle(int32(offset.X), int32(offset.Y), 15, rl.White)
		// vehicle.Speed += rl.GetFrameTime() * 2.5
		vehicle.Acceleration = rl.Clamp(vehicle.Acceleration+0.5*rl.GetFrameTime(), 0, 12)
		// fmt.Println(vehicle.Acceleration)
		vehicle.Speed = rl.Clamp(vehicle.Speed+rl.GetFrameTime()*vehicle.Acceleration, 0, 88)
	}
	if vehicle.Speed >= 88 {
		vehicle.Hyperjump = true
	}
}

func (vehicle *Vehicle) decreaseSpeed() {
	// if rl.IsKeyPressed(rl.KeySpace) {
	// 	vehicle.SlowingDown = true
	// }
	if vehicle.SlowingDown {
		vehicle.Speed -= 15
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

func (vehicle *Vehicle) updateFrame() {
	vehicle.Sprite.SourceRect.X = vehicle.Sprite.SourceRect.Width * float32(vehicle.Sprite.SpriteFrame)
	frameCount++
	if vehicle.Sprite.SpriteFrame > 2 {
		vehicle.Sprite.SpriteFrame = 0
		vehicle.Sprite.SourceRect = rl.NewRectangle(0, 0, 64, 64)
	}
	if frameCount%10 == 1 {
		vehicle.Sprite.SpriteFrame++

	}
}

// Return and scale player rectangle hitbox
func (vehicle *Vehicle) getRectHitbox() rl.Rectangle {
	Scale := vehicle.Sprite.Render.Scale
	X := vehicle.Position.X
	Y := vehicle.Position.Y
	Width := vehicle.Body.Width
	Height := vehicle.Body.Height

	scaledX := (X - (Width*Scale)/2) + 5*Scale
	scaledY := Y - (Height*Scale)/2 + 50*Scale
	scaledWidth := (Width * Scale) - 130*Scale
	scaledHeight := Height*Scale - 20*Scale

	hitBoxRect := rl.NewRectangle(scaledX, scaledY, scaledWidth, scaledHeight)

	// rl.DrawRectangle(int32(hitBoxRect.X), int32(hitBoxRect.Y), int32(hitBoxRect.Width), int32(hitBoxRect.Height), rl.White)
	return hitBoxRect
}

func (vehicle *Vehicle) move() {
	if currentLane != BOTTOM && (rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressed(rl.KeyDown)) {
		vehicle.Position.Y = rl.Lerp(vehicle.Position.Y, vehicle.Position.Y+100, 1)
		currentLane = BOTTOM
		vehicle.Lane = BOTTOM
		// fmt.Println("down")
	}
	if currentLane != TOP && (rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp)) {
		vehicle.Position.Y = rl.Lerp(vehicle.Position.Y, vehicle.Position.Y-100, 1)
		currentLane = TOP
		vehicle.Lane = TOP
		// fmt.Println("up")
	}
}
