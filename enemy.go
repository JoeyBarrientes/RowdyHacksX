package main

import (
	"main/physics"
	"main/renderer"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Defines game Enemy struct containing enemies
// of type Bat and Zombie
type Enemies struct {
	Shooting []Shooting
	Normal   []Normal
}

// Base Enemy properties
type Enemy struct {
	Body     physics.CircleBody
	Sprite   renderer.CharacterSprite
	Position rl.Vector2
	// Damage       float32
	Health       float32
	IsDeflecting bool
	DeflectTime  float32
	// DropChance   float32
	XOffset int32
	YOffset int32
	Lane
}

type Normal struct {
	Enemy
}

type Shooting struct {
	Enemy
	Projectile    []Projectile
	shootTimer    float32
	shootInterval float32
}

// func NewEnemy(Sprite rl.Texture2D, Color rl.Color, Position rl.Vector2, Velocity rl.Vector2, Radius float32, Scale float32) Enemy {
// 	enemy := Enemy{
// 		Body: physics.CircleBody{
// 			Velocity: Velocity,
// 			Radius:   Radius,
// 			Angle:    0,
// 		},
// 		Sprite: renderer.CharacterSprite{
// 			Render: renderer.SpriteRenderer{
// 				Sprite: Sprite,
// 				Color:  Color,
// 				Angle:  0,
// 				Scale:  Scale,
// 			},
// 			SourceRect:  rl.NewRectangle(0, 0, 64, 64),
// 			IsMoving:    false,
// 			SpriteFrame: 0,
// 		},
// 		Position:     Position,
// 		Health:       100,
// 		IsDeflecting: false,
// 		DeflectTime:  0,
// 		XOffset:      -50 - int32(rand.Int32N(250)),
// 		YOffset:      int32(screenSize.Y) - (int32(rand.Int32N(int32(200*screenScale.Y)) + 100)),
// 	}

// 	return enemy
// }

// // Initializes Zombie entity
// func NewZombie(Sprite rl.Texture2D, Color rl.Color, Position rl.Vector2, Velocity rl.Vector2, Radius float32, Scale float32) Zombie {
// 	Zombie := Zombie{
// 		Enemy: Enemy{
// 			Body: physics.CircleBody{
// 				Position: Position,
// 				Velocity: Velocity,
// 				Radius:   Radius,
// 				Angle:    0,
// 			},
// 			Sprite: renderer.SpriteRenderer{
// 				Sprite: Sprite,
// 				Color:  Color,
// 				Scale:  Scale,
// 			},
// 			Damage:       10,
// 			Health:       100,
// 			IsDeflecting: false,
// 			DeflectTime:  0,
// 			DropChance:   rand.Float32(),
// 			XOffset:      int32(screenSize.X) + 50 + int32(rand.Int32N(250)),
// 			YOffset:      int32(screenSize.Y) - (int32(rand.Int32N(int32(200*screenScale.Y)) + 100)),
// 		},
// 	}
// 	return Zombie
// }

// Initializes Bat entity
func NewShootingEnemy(Sprite rl.Texture2D, Color rl.Color, Position rl.Vector2, Velocity rl.Vector2, Radius float32, Scale float32, Lane Lane) Shooting {
	newShooting := Shooting{
		Enemy: Enemy{
			Body: physics.CircleBody{
				Velocity: Velocity,
				Radius:   Radius,
				Angle:    0,
			},
			Sprite: renderer.CharacterSprite{
				Render: renderer.SpriteRenderer{
					Sprite: Sprite,
					Color:  Color,
					Angle:  0,
					Scale:  Scale,
				},
				SourceRect:  rl.NewRectangle(0, 0, 64, 64),
				IsMoving:    false,
				SpriteFrame: 0,
			},
			// Damage:       10,
			Position:     Position,
			Health:       100,
			IsDeflecting: false,
			DeflectTime:  0,
			// DropChance:   rand.Float32(),
			XOffset: -50 - int32(rand.Int32N(250)),
			YOffset: int32(screenSize.Y) - (int32(rand.Int32N(int32(200*screenScale.Y)) + 100)),
			Lane:    Lane,
		},
		Projectile:    []Projectile{},
		shootTimer:    0,
		shootInterval: 2,
	}
	return newShooting
}

func NewNormalEnemy(Sprite rl.Texture2D, Color rl.Color, Position rl.Vector2, Velocity rl.Vector2, Radius float32, Scale float32) Normal {
	newNormal := Normal{
		Enemy: Enemy{
			Body: physics.CircleBody{
				Velocity: Velocity,
				Radius:   Radius,
				Angle:    0,
			},
			Sprite: renderer.CharacterSprite{
				Render: renderer.SpriteRenderer{
					Sprite: Sprite,
					Color:  Color,
					Angle:  0,
					Scale:  Scale,
				},
				SourceRect:  rl.NewRectangle(0, 0, 100, 100),
				IsMoving:    false,
				SpriteFrame: 0,
			},
			// Damage:       10,
			Position:     Position,
			Health:       100,
			IsDeflecting: false,
			DeflectTime:  0,
			// DropChance:   rand.Float32(),
			XOffset: int32(screenSize.X) + 50 + int32(rand.Int32N(150)),
			YOffset: rand.Int32N(int32(screenSize.Y) - 100),
		},
	}
	return newNormal
}

func (enemy *Enemy) DrawSprite() {
	destRect := rl.NewRectangle(enemy.Position.X, enemy.Position.Y, 100*enemy.Sprite.Render.Scale, 100*enemy.Sprite.Render.Scale)
	origin := rl.Vector2Scale(rl.NewVector2(float32(enemy.Sprite.Render.Sprite.Width)/2, float32(enemy.Sprite.Render.Sprite.Height)/2), enemy.Sprite.Render.Scale)
	rl.DrawTexturePro(enemy.Sprite.Render.Sprite, enemy.Sprite.SourceRect,
		destRect,
		origin, enemy.Sprite.Render.Angle, enemy.Sprite.Render.Color)
	// rl.DrawCircle(int32(enemy.Position.X-float32(enemy.Sprite.Render.Sprite.Width/5*4)), int32(enemy.Position.Y+float32(enemy.Sprite.Render.Sprite.Height)), enemy.Body.Radius, rl.White)
}

func (enemy *Enemy) VelocityTick() {
	adjustedVel := rl.Vector2Scale(enemy.Body.Velocity, rl.GetFrameTime())
	enemy.Position = rl.Vector2Add(enemy.Position, adjustedVel)
}

func (enemy *Enemy) PhysicsUpdate() {
	enemy.VelocityTick()

}

func (enemies *Enemies) updateEnemyFrame() {
	for i := len(enemies.Shooting) - 1; i >= 0; i-- {
		libyan := &enemies.Shooting[i]

		libyan.Sprite.SourceRect.X = libyan.Sprite.SourceRect.Width * float32(libyan.Sprite.SpriteFrame)
		if libyan.Sprite.SpriteFrame > 3 {
			libyan.Sprite.SpriteFrame = 0
			libyan.Sprite.SourceRect = rl.NewRectangle(0, 0, 64, 64)

		}
		if frameCount%10 == 1 {
			libyan.Sprite.SpriteFrame++

		}
		// if frameCount%18 == 1 {
		// 	bullet.Sprite.SpriteFrame++

		// }
	}

}

// func (enemy *Enemy) getRectHitbox() rl.Rectangle {
// 	Scale := enemy.Sprite.Render.Scale
// 	X := enemy.Position.X
// 	Y := enemy.Position.Y
// 	Width := enemy.Body.Width
// 	Height := enemy.Body.Height

// 	scaledX := (X - (Width*Scale)/2) + 5*Scale
// 	scaledY := Y - (Height*Scale)/2 + 50*Scale
// 	scaledWidth := (Width * Scale) - 100*Scale
// 	scaledHeight := Height*Scale - 20*Scale

// 	hitBoxRect := rl.NewRectangle(scaledX, scaledY, scaledWidth, scaledHeight)
// 	// rl.DrawRectangle(int32(hitBoxRect.X), int32(hitBoxRect.Y), int32(hitBoxRect.Width), int32(hitBoxRect.Height), rl.White)
// 	return hitBoxRect
// }

// func (enemy *Enemy) Draw() {
// 	renderer.DrawCircleEntity(enemy.Position, &enemy.Sprite, enemy.Body.Angle)
// }
