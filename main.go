package main

import (
	"fmt"
	"main/physics"
	"main/renderer"
	"math"
	"math/rand/v2"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Global Variables
var initScreenWidth int32 = 1920
var initScreenHeight int32 = 1080
var screenSize rl.Vector2 = rl.NewVector2(float32(initScreenWidth), float32(initScreenHeight))
var screenScale = rl.NewVector2((float32(screenSize.X) / float32(initScreenWidth)),
	(float32(screenSize.Y) / float32(initScreenHeight)))

// type Background struct {
// 	b
// }

var backgroundImage rl.Texture2D
var backgroundScale float32
var background rl.Texture2D
var midground rl.Texture2D
var foreground rl.Texture2D

var vehicleTexture rl.Texture2D
var frameCount int
var martyTexture rl.Texture2D
var Marty renderer.StillSprite

var audio Audio = NewAudio()

var theme ColorTheme

type GameScreen int

const (
	TITLE GameScreen = iota
	HOWTO
	GAMEPLAY
	HYPERJUMP
	GAMEOVER
	EXIT
)

type Background int

const (
	CITY Background = iota
	PREHISTORIC
)

type Lane int

const (
	TOP Lane = iota
	BOTTOM
)

var currentScreen GameScreen = TITLE
var currentLane Lane = TOP
var currentBackground Background = CITY

var score float32 = 0
var highScore float32 = 0
var hasSpawned bool = false
var vehicleRect rl.Rectangle
var hyperjump bool
var hudColor rl.Color = rl.NewColor(0, 0, 0, 128)

func main() {
	// Game variable initializations
	rl.InitWindow(initScreenWidth, initScreenHeight, "Game")
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.InitAudioDevice()
	audio.loadAudio()
	theme = NewColorTheme(rl.NewColor(176, 166, 170, 255), rl.NewColor(145, 63, 86, 255), rl.White)
	checkResize(&screenSize)

	vehicleTexture = rl.LoadTexture("textures/delorean.png")
	projectileTextures := rl.LoadTexture("textures/projectileBullet.png")
	vehiclePosition := rl.NewVector2(screenSize.X/8*7, screenSize.Y/9*6-float32(vehicleTexture.Height/2))
	DeLorean := NewVehicle(vehicleTexture, projectileTextures, rl.White, vehiclePosition, rl.NewVector2(0, 0), float32(vehicleTexture.Width), float32(vehicleTexture.Height), 4)
	martyTexture = rl.LoadTexture("textures/marty.png")
	Marty = renderer.NewStillSprite(rl.NewVector2(screenSize.X/2-float32(martyTexture.Width)/2, screenSize.Y/6*4-float32(martyTexture.Height)/2), martyTexture, rl.White, 0, 7)

	backgroundImage = rl.LoadTexture("textures/Bright/City3.png")
	background = rl.LoadTexture("textures/background.png")
	midground = rl.LoadTexture("textures/middleground.png")
	foreground = rl.LoadTexture("textures/foreground.png")

	libyanTexture := rl.LoadTexture("textures/enemylibyan.png")

	enemyTextures := []rl.Texture2D{}
	enemyTextures = append(enemyTextures, libyanTexture)

	enemies := Enemies{}

	var scrollingBack float32 = 0.0
	var scrollingMid float32 = 0.0
	var scrollingFore float32 = 0.0

	vehicleRect = DeLorean.getRectHitbox()

	checkResize(&screenSize)
	rl.SetTargetFPS(60)
	// Main game loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(theme.backgroundColor)

		scrollingBack -= 1.0
		scrollingMid -= 8.0
		scrollingFore -= 8.0

		if scrollingBack <= -float32(background.Width)*backgroundScale {
			scrollingBack = 0
		}
		if scrollingMid <= -float32(midground.Width)*backgroundScale {
			scrollingMid = 0
		}
		if scrollingFore <= -float32(foreground.Width)*backgroundScale {
			scrollingFore = 0
		}

		screenSize.X = float32(rl.GetScreenWidth())
		screenSize.Y = float32(rl.GetScreenHeight())

		checkResize(&screenSize)

		switch currentScreen {
		case TITLE:
			rl.DrawTextureEx(background, rl.NewVector2(scrollingBack, -20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(background, rl.NewVector2(float32(background.Width)*backgroundScale+scrollingBack, 0), 0.0, backgroundScale, rl.White)

			rl.DrawTextureEx(midground, rl.NewVector2(scrollingMid, 20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(midground, rl.NewVector2(float32(midground.Width)*backgroundScale+scrollingMid, 20), 0.0, backgroundScale, rl.White)

			rl.DrawTextureEx(foreground, rl.NewVector2(scrollingFore, 0), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(foreground, rl.NewVector2(float32(foreground.Width)*backgroundScale+scrollingFore, 0), 0.0, backgroundScale, rl.White)

			rl.DrawRectangle(0, 0, int32(screenSize.X), int32(screenSize.Y), hudColor)
			displayTitleScreen()
			audio.checkMute()
		case HOWTO:

			rl.DrawTextureEx(background, rl.NewVector2(scrollingBack, -20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(background, rl.NewVector2(float32(background.Width)*backgroundScale+scrollingBack, 0), 0.0, backgroundScale, rl.White)

			// Draw midzground image twice
			rl.DrawTextureEx(midground, rl.NewVector2(scrollingMid, 20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(midground, rl.NewVector2(float32(midground.Width)*backgroundScale+scrollingMid, 20), 0.0, backgroundScale, rl.White)

			// Draw foreground image twice
			rl.DrawTextureEx(foreground, rl.NewVector2(scrollingFore, 0), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(foreground, rl.NewVector2(float32(foreground.Width)*backgroundScale+scrollingFore, 0), 0.0, backgroundScale, rl.White)

			rl.DrawRectangle(0, 0, int32(screenSize.X), int32(screenSize.Y), hudColor)

			displayHowToScreen()
		case GAMEPLAY: // Main Game Loop State
			// rl.ClearBackground(rl.Black)

			rl.DrawTextureEx(background, rl.NewVector2(scrollingBack, -20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(background, rl.NewVector2(float32(background.Width)*backgroundScale+scrollingBack, 0), 0.0, backgroundScale, rl.White)

			// Draw midzground image twice
			rl.DrawTextureEx(midground, rl.NewVector2(scrollingMid, 20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(midground, rl.NewVector2(float32(midground.Width)*backgroundScale+scrollingMid, 20), 0.0, backgroundScale, rl.White)

			// Draw foreground image twice
			rl.DrawTextureEx(foreground, rl.NewVector2(scrollingFore, 0), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(foreground, rl.NewVector2(float32(foreground.Width)*backgroundScale+scrollingFore, 0), 0.0, backgroundScale, rl.White)

			// DeLorean.Position = vehiclePosition

			vehicleRect = DeLorean.getRectHitbox()
			DeLorean.updateFrame()
			// DeLorean.updateProjectileFrame()
			DeLorean.move()

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				DeLorean.shoot()
			}
			// rl.DrawRectangleLines(int32(vehicleRect.X), int32(vehicleRect.Y), int32(vehicleRect.Width), int32(vehicleRect.Height), rl.White)

			spawnEnemies(&enemies, &hasSpawned, enemyTextures, &DeLorean)
			drawEnemies(enemies)
			moveEnemies(&enemies)
			checkEnemyPlayerCollision(&enemies, &DeLorean, vehicleRect)
			checkEnemyProjectileCollision(&enemies, &DeLorean)
			shootProjectile(&enemies, &DeLorean)

			// if rl.IsKeyPressed(rl.KeySpace) {
			DeLorean.decreaseSpeed()

			enemies.updateEnemyFrame()
			drawProjectiles(&enemies, &DeLorean)
			updateProjectiles(&enemies, &DeLorean, vehicleRect)
			nextHyperjump(&DeLorean)
			// for i := range.
			// }
			// fmt.Println(DeLorean.Position.Y)

			DeLorean.DrawCharacter()
			DeLorean.drawBullets()
			DeLorean.updateBullets()
			DeLorean.despawnBullets()
			DeLorean.updateSpeed()
			audio.checkMute()
		case GAMEOVER: // Game Lose State
			displayGameOver()
			audio.checkMute()
		case EXIT: // Game Exit State
			rl.CloseWindow()
			return
		default: // Unknown Default State
			fmt.Println("Unknown screen")
		}
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func nextHyperjump(DeLorean *Vehicle) {
	if DeLorean.Hyperjump {
		fmt.Println("Ready to Jump")
		if rl.IsKeyPressed(rl.KeySpace) {
			DeLorean.Speed = 0
			DeLorean.SpeedTracker = 0
			DeLorean.Acceleration = 1
			if currentBackground == CITY {
				currentBackground = PREHISTORIC
				background = rl.LoadTexture("textures/background-prehistoric.png")
				midground = rl.LoadTexture("textures/middleground-prehistoric.png")
				foreground = rl.LoadTexture("textures/foreground-prehistoric.png")
				backgroundImage = rl.LoadTexture("textures/backgroundImage2.png")
				DeLorean.Hyperjump = false
			} else if currentBackground == PREHISTORIC {
				currentBackground = CITY
				backgroundImage = rl.LoadTexture("textures/Bright/City3.png")
				background = rl.LoadTexture("textures/background.png")
				midground = rl.LoadTexture("textures/middleground.png")
				foreground = rl.LoadTexture("textures/foreground.png")
				DeLorean.Hyperjump = false
			}
		}

	}
}

// Spawns bat and zombie entities off screen to the right
func spawnEnemies(enemies *Enemies, hasSpawned *bool, enemyTextures []rl.Texture2D, vehicle *Vehicle) {
	// Gets random number to determine
	// enemy type and number per spawn occurrence
	randEnemy := rand.IntN(1)
	randLane := rand.IntN(2)
	// laneOffset := 0
	// if randLane == 1 {
	// 	laneOffset = 100
	// }

	// Check the current game time to control enemy spawn frequency
	if int(rl.GetTime())%2 == 1 {
		if !*hasSpawned {
			var sprite rl.Texture2D
			spawnPosition := rl.NewVector2(-50, float32(750+130*randLane))

			if randEnemy == 0 { // Create bat enemy
				sprite = enemyTextures[0]
				libyan := NewShootingEnemy(sprite, rl.White, spawnPosition, rl.NewVector2(0, 0), 120*screenScale.X, 2.5*screenScale.X, Lane(randLane))
				// libyan.Position = rl.NewVector2(float32(libyan.XOffset), vehicle.Position.Y)
				fmt.Println(libyan.Lane)
				enemies.Shooting = append(enemies.Shooting, libyan)
				// }
			}
			// } else { // Create zombie enemy
			// 	sprite = enemyTextures[1]
			// 	for i := 0; i <= randAmount; i++ {
			// 		zombie := NewZombie(sprite, rl.White, initSpawnPosition, rl.NewVector2(0, 0), 60*screenScale.X, 6*screenScale.X)
			// 		zombie.Body.Position = rl.NewVector2(float32(zombie.XOffset), float32(zombie.YOffset))
			// 		enemies.Zombies = append(enemies.Zombies, zombie)
			// 	}
			// }

			// Ensures spawning only happens one frame per 2 seconds,
			// not every frame in one second per 2 seconds
			*hasSpawned = true
		}
	} else {
		*hasSpawned = false
	}
}

// Displays enemies to screen
func drawEnemies(enemies Enemies) {
	for _, libyan := range enemies.Shooting {
		libyan.DrawSprite()
	}
	// for _, zombies := range enemies.Zombies {
	// 	zombies.Draw()
	// }
}

func moveEnemies(enemies *Enemies) {
	// move each bat in slice
	for i := range enemies.Shooting {
		libyan := &enemies.Shooting[i]
		libyan.Body.Velocity.X = 150
		libyan.PhysicsUpdate()
	}

	// // move each zombie in slice
	// for i := range enemies.Zombies {
	// 	zombie := &enemies.Zombies[i]
	// 	if !zombie.IsDeflecting {
	// 		direction := rl.Vector2Subtract(knight.Body.Position, zombie.Body.Position)
	// 		direction = normalizeVector(direction)
	// 		(*zombie).Body.Velocity.X = direction.X * 300
	// 		(*zombie).Body.Velocity.Y = direction.Y * 150
	// 	} else {
	// 		zombie.DeflectTime += rl.GetFrameTime()
	// 		decaySpeed := rl.NewVector2(float32(math.Abs(float64(zombie.Body.Velocity.X*25))*float64(rl.GetFrameTime())),
	// 			float32(math.Abs(float64(zombie.Body.Velocity.Y*5))*float64(rl.GetFrameTime())))
	// 		zombie.Body.Velocity.X = applyVelocityDecay(zombie.Body.Velocity.X, decaySpeed.X)
	// 		// Y velocity is not changed

	// 		if zombie.DeflectTime >= deflectDuration {
	// 			zombie.IsDeflecting = false
	// 			zombie.DeflectTime = 0
	// 		}
	// 	}

	// 	zombie.Body.PhysicsUpdate()
	// }
}

// rl.NewVector2(libyan.Position.X-float32(libyan.Sprite.Render.Sprite.Width/5*4), libyan.Position.Y+float32(libyan.Sprite.Render.Sprite.Height))
// Creates and shoots a projectile from each bat toward the player at a random interval
func shootProjectile(enemies *Enemies, DeLorean *Vehicle) {
	for i := len(enemies.Shooting) - 1; i >= 0; i-- {
		libyan := &enemies.Shooting[i]
		// direction := rl.Vector2Subtract(DeLorean.Position, rl.NewVector2(libyan.Position.X-float32(libyan.Sprite.Render.Sprite.Width/5*4), libyan.Position.Y+float32(libyan.Sprite.Render.Sprite.Height)))
		// direction = normalizeVector(direction)

		libyan.shootTimer += rl.GetFrameTime()

		// Check if the bat can shoot
		if libyan.shootTimer >= libyan.shootInterval {
			// audio.playWithRandPitch(audio.Sounds["shot"])

			enemyBullet := Projectile{
				Body:     physics.NewCirclePhysicsBody(rl.NewVector2(100, 0), 15, 0),
				Position: rl.NewVector2(libyan.Position.X-float32(libyan.Sprite.Render.Sprite.Width/5*4), libyan.Position.Y+float32(libyan.Sprite.Render.Sprite.Height)),
				Speed:    600,
				Color:    rl.Red,
				Lane:     libyan.Lane,
			}
			if enemyBullet.Lane == DeLorean.Lane {
				enemyBullet.Color = rl.Red
			} else {
				enemyBullet.Color = rl.DarkGray
			}
			// pr := libyan.Projectile.NewProjectile(bat.Body.Position, bat.Body.Velocity, bat.Projectile.Speed, 10, 0, rl.White)
			enemyBullet.Body.Velocity.X = enemyBullet.Speed
			// pr.Body.Velocity.Y = direction.Y * 350
			libyan.Projectile = append(libyan.Projectile, enemyBullet)
			// enemyBullet.Position.X -= 35
			// libyan.Projectile = append(libyan.Projectile, enemyBullet)

			libyan.shootTimer = 0

			// Assign a new random shoot interval between 3 and 5 seconds
			libyan.shootInterval = 2.0
		}
	}
}

// Draws each projectile slice element to screen
func drawProjectiles(enemies *Enemies, DeLorean *Vehicle) {
	for i := len(enemies.Shooting) - 1; i >= 0; i-- {
		libyan := &enemies.Shooting[i]
		for _, projectile := range libyan.Projectile {
			projectile.Draw()
		}
	}
}

func updateProjectiles(enemies *Enemies, DeLorean *Vehicle, vehicleRect rl.Rectangle) {
	for i := len(enemies.Shooting) - 1; i >= 0; i-- {
		libyan := &enemies.Shooting[i]
		for j := range libyan.Projectile {
			// Use a pointer to each projectile to ensure updates apply to the actual object
			projectile := &libyan.Projectile[j]
			projectile.PhysicsUpdate()
			if projectile.Lane == DeLorean.Lane {
				projectile.Color = rl.Red
			} else {
				projectile.Color = rl.DarkGray
			}

			if rl.CheckCollisionCircleRec(projectile.Position, projectile.Body.Radius, vehicleRect) && DeLorean.Lane == projectile.Lane {
				DeLorean.SlowingDown = true
				libyan.Projectile = append(libyan.Projectile[:j], libyan.Projectile[j+1:]...)

				break
			}
		}
	}
}

// Handles player and enemy collision, and applies damage and despawns enemy if so
func checkEnemyPlayerCollision(enemies *Enemies, DeLorean *Vehicle, vehicleRect rl.Rectangle) {
	for i := len(enemies.Shooting) - 1; i >= 0; i-- {
		libyan := &enemies.Shooting[i]
		if libyan.Lane == DeLorean.Lane {
			libyan.Sprite.Render.Color = rl.White
		} else {
			libyan.Sprite.Render.Color = rl.DarkGray
		}
		if rl.CheckCollisionCircleRec(rl.NewVector2(libyan.Position.X-float32(libyan.Sprite.Render.Sprite.Width/5*4), libyan.Position.Y+float32(libyan.Sprite.Render.Sprite.Height)), libyan.Body.Radius, vehicleRect) && (libyan.Lane == DeLorean.Lane) {
			// rl.PlaySound(audio.Sounds["damaged"])
			// knight.Health -= bat.Damage
			// knight.Health = rl.Clamp(knight.Health, 0, 100)
			// *healthTracker = knight.Health / 100
			DeLorean.SlowingDown = true

			enemies.Shooting = append(enemies.Shooting[:i], enemies.Shooting[i+1:]...) // Remove Enemy after collision
		}
	}

	// for i := len(enemies.Zombies) - 1; i >= 0; i-- {
	// 	zombie := &enemies.Zombies[i]
	// 	if rl.CheckCollisionCircleRec(zombie.Body.Position, zombie.Body.Radius, knightRect) {
	// 		audio.playWithRandPitch(audio.Sounds["damaged"])
	// 		knight.Health -= zombie.Damage
	// 		knight.Health = rl.Clamp(knight.Health, 0, 100)
	// 		*healthTracker = knight.Health / 100

	// 		enemies.Zombies = append(enemies.Zombies[:i], enemies.Zombies[i+1:]...) // Remove Enemy after collision
	// 	}
	// }
}

// Handles reflected projectile and bat collision, despawns them, and chance of food spawn
func checkEnemyProjectileCollision(enemies *Enemies, DeLorean *Vehicle) {

	for i := len(enemies.Shooting) - 1; i >= 0; i-- {
		libyan := &enemies.Shooting[i]
		for j := len(DeLorean.Bullets) - 1; j >= 0; j-- {
			bullet := &(DeLorean.Bullets)[j]
			if rl.CheckCollisionCircles(rl.NewVector2(libyan.Position.X-float32(libyan.Sprite.Render.Sprite.Width/5*4), libyan.Position.Y+float32(libyan.Sprite.Render.Sprite.Height)),
				libyan.Body.Radius, bullet.Position, bullet.Body.Radius) && libyan.Lane == bullet.Lane {
				DeLorean.Bullets = append(DeLorean.Bullets[:j], DeLorean.Bullets[j+1:]...)
				enemies.Shooting = append(enemies.Shooting[:i], enemies.Shooting[i+1:]...)
				break
			}
		}
	}
}

// func checkBulletCollision(enemies *Enemies, DeLorean *Vehicle){
// 		for i := len(enemies.Shooting) - 1; i >= 0; i-- {
// 		libyan := &enemies.Shooting[i]
// 		for j := len(DeLorean.Bullets) - 1; j >= 0; j-- {
// 			bullet := &(DeLorean.Bullets)[j]
// }

// Implements title screen and UI elements
func displayTitleScreen() {

	title := "Route 88"
	titleWidth := rl.MeasureText(title, 50)
	rl.DrawText(title,
		int32(screenSize.X)/2-(titleWidth/2),
		int32(screenSize.Y)/4-(50/2),
		50, rl.White)

	playButton := NewButton(0, 0, 400, 100, 0.1, 0, 0, TITLE, HOWTO)
	playButton.X = int32(screenSize.X)/2 - playButton.Width/2
	playButton.Y = int32(screenSize.Y/2.5) - playButton.Height/2
	playButton.SetText("Play", 50)
	playTextWidth := rl.MeasureText(playButton.text, playButton.textSize)
	rl.DrawText(playButton.text,
		playButton.X+(playButton.Width/2)-(playTextWidth/2),
		playButton.Y+(playButton.Height/2)-(playButton.textSize/2),
		playButton.textSize,
		rl.White,
	)

	// howToButton := NewButton(playButton.X, playButton.Y+150, 400, 100, 0.1, 0, 0, TITLE, HOWTO)
	// howToButton.SetText("How To Play", 50)
	// howToTextWidth := rl.MeasureText(howToButton.text, howToButton.textSize)
	// rl.DrawText(howToButton.text,
	// 	howToButton.X+(howToButton.Width/2)-(howToTextWidth/2),
	// 	howToButton.Y+(howToButton.Height/2)-(howToButton.textSize/2),
	// 	howToButton.textSize,
	// 	rl.White,
	// )

	exitButton := NewButton(playButton.X, playButton.Y+150, 400, 100, 0.1, 0, 0, TITLE, EXIT)
	exitButton.SetText("Exit", 50)
	exitTextWidth := rl.MeasureText(exitButton.text, exitButton.textSize)
	rl.DrawText(exitButton.text,
		exitButton.X+(exitButton.Width/2)-(exitTextWidth/2),
		exitButton.Y+(exitButton.Height/2)-(exitButton.textSize/2),
		exitButton.textSize,
		rl.White,
	)

	addSwitchScreenAction(&playButton)
	// addSwitchScreenAction(&howToButton)
	addSwitchScreenAction(&exitButton)
	playButton.UpdateButton()
	// howToButton.UpdateButton()
	exitButton.UpdateButton()
}

type DialogBox struct {
	rl.Rectangle
	text     string
	textSize float32
}

var dialogCount int = 0

func displayHowToScreen() {
	Marty.Draw(0)
	Dialog := DialogBox{
		Rectangle: rl.NewRectangle(Marty.Position.X-400, Marty.Position.Y-600, 800, 400),
		text:      "",
		textSize:  50,
	}

	if dialogCount == 0 {
		Dialog.text = "Oh no! The Libyan\n" +
			"Terrorists have found Doc!\n" +
			"We need to go help him now!"
	} else if dialogCount == 1 {
		Dialog.Rectangle.X -= 150
		Dialog.Rectangle.Y -= 50
		Dialog.Width += 300
		Dialog.Height += 100
		Dialog.text = "Aim with MOUSE\n" +
			"Press LEFT MOUSE BUTTON to shoot\n" +
			"Press W to move one lane up\n" +
			"Press S to move one lane down\n\n" +
			"Once 88 MPH have been reached,\n" +
			"press SPACEBAR to begin the\n" +
			"hyperjump sequence!"
	} else if dialogCount == 2 {
		currentScreen = GAMEPLAY
	}

	// DialogRect := rl.NewRectangle(Marty.Position.X-400, Marty.Position.Y-600, 800, 400)

	if rl.IsKeyPressed(rl.KeySpace) || rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		dialogCount++
	}
	rl.DrawRectangleRounded(Dialog.Rectangle, 0.5, 1, rl.White)
	rl.DrawText(Dialog.text, Dialog.Rectangle.ToInt32().X+40, Dialog.Rectangle.ToInt32().Y+40, 50, rl.Black)
	// rl.DrawText(DialogText, DialogRect.ToInt32().X+40, DialogRect.ToInt32().Y+40, 50, rl.Black)
}

func displayGameOver() {

}

// Helper that decreases velocity
func applyVelocityDecay(velocity, decaySpeed float32) float32 {
	if velocity > 0 {
		velocity -= decaySpeed
		if velocity < 0 {
			velocity = 0
		}
	} else if velocity < 0 {
		velocity += decaySpeed
		if velocity > 0 {
			velocity = 0
		}
	}
	return velocity
}

var growthSpeed float32 = 0.5

// Helper that decreases velocity
func applyVelocityGrowth(velocity, growthSpeed float32) float32 {
	if velocity < 2 {
		velocity += growthSpeed
		if velocity > 2 {
			velocity = 2
		}
	}

	return velocity
}

// Normalizes a vector to -1 to 1
func normalizeVector(v rl.Vector2) rl.Vector2 {
	length := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
	if length == 0 {
		return rl.Vector2{X: 0, Y: 0}
	}
	return rl.Vector2{X: v.X / length, Y: v.Y / length}
}

// Added semi-functionality for resizing and full screen of game screen
func checkResize(screenSize *rl.Vector2) {
	screenSize.X = float32(rl.GetScreenWidth())
	screenSize.Y = float32(rl.GetScreenHeight())

	screenScale = rl.NewVector2((float32(screenSize.X) / float32(initScreenWidth)),
		(float32(screenSize.Y) / float32(initScreenHeight)))

	backgroundScale = float32(screenSize.Y) / float32(backgroundImage.Height)

	if rl.IsWindowResized() {
		rl.SetWindowSize(int(screenSize.X), int(screenSize.Y))
		if screenSize.X <= 1920 {
		} else if screenSize.X > 1920 {
		}
	}
}

// Draws game score to screen
func DisplayScore() {
	score += 10 * rl.GetFrameTime()
	ScoreStr := "Score:\t\t\t\t\t " + strconv.FormatFloat(float64(score), 'f', 0, 32)
	rl.DrawText(ScoreStr, 20, 45, 40, rl.White)
}

// Draws high score to score
func DisplayHighScore() {
	highScoreStr := "High Score:\t\t" + strconv.FormatFloat(float64(highScore), 'f', 0, 32)
	rl.DrawText(highScoreStr, 20, 95, 40, rl.White)
}
