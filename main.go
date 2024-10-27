package main

import (
	"fmt"
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

var backgroundImage rl.Texture2D
var backgroundScale float32
var background rl.Texture2D
var midground rl.Texture2D
var foreground rl.Texture2D

var vehicleTexture rl.Texture2D
var frameCount int

var audio Audio = NewAudio()

var theme ColorTheme

type GameScreen int

const (
	TITLE GameScreen = iota
	HOWTO
	GAMEPLAY
	GAMEOVER
	EXIT
)

type Lane int

const (
	TOP Lane = iota
	BOTTOM
)

var currentScreen GameScreen = TITLE
var currentLane Lane = TOP

var score float32 = 0
var highScore float32 = 0
var hasSpawned bool = false

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
		case TITLE: // Title Screen State
			displayTitleScreen()
			audio.checkMute()
		case HOWTO:
			displayHowToScreen()
		case GAMEPLAY: // Main Game Loop State
			// rl.ClearBackground(rl.Black)

			rl.DrawTextureEx(background, rl.NewVector2(scrollingBack, -20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(background, rl.NewVector2(float32(background.Width)*backgroundScale+scrollingBack, 0), 0.0, backgroundScale, rl.White)

			// Draw midground image twice
			rl.DrawTextureEx(midground, rl.NewVector2(scrollingMid, 20), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(midground, rl.NewVector2(float32(midground.Width)*backgroundScale+scrollingMid, 20), 0.0, backgroundScale, rl.White)

			// Draw foreground image twice
			rl.DrawTextureEx(foreground, rl.NewVector2(scrollingFore, 0), 0.0, backgroundScale, rl.White)
			rl.DrawTextureEx(foreground, rl.NewVector2(float32(foreground.Width)*backgroundScale+scrollingFore, 0), 0.0, backgroundScale, rl.White)

			// DeLorean.Position = vehiclePosition
			DeLorean.DrawCharacter()
			DeLorean.updateFrame()
			// DeLorean.updateProjectileFrame()
			DeLorean.move()

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				DeLorean.shoot()
			}

			spawnEnemies(&enemies, &hasSpawned, enemyTextures, &DeLorean)
			drawEnemies(enemies)
			moveEnemies(&enemies)

			// if rl.IsKeyPressed(rl.KeySpace) {
			DeLorean.decreaseSpeed()
			// }
			fmt.Println(DeLorean.Position.Y)

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
	if int(rl.GetTime())%8 == 1 {
		if !*hasSpawned {
			var sprite rl.Texture2D
			spawnPosition := rl.NewVector2(-50, float32(730+130*randLane))

			if randEnemy == 0 { // Create bat enemy
				sprite = enemyTextures[0]
				// for i := 0; i <= randAmount; i++ {
				libyan := NewShootingEnemy(sprite, rl.White, spawnPosition, rl.NewVector2(0, 0), 45*screenScale.X, 3*screenScale.X)
				// libyan.Position = rl.NewVector2(float32(libyan.XOffset), vehicle.Position.Y)
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
		// if !libyan.IsDeflecting {
		// 	direction := rl.Vector2Subtract(knight.Body.Position, bat.Body.Position)
		// 	direction = normalizeVector(direction)
		// 	bat.Body.Velocity.X = direction.X * 150
		// 	bat.Body.Velocity.Y = direction.Y * 150
		// } else {
		// 	bat.DeflectTime += rl.GetFrameTime()
		// 	decaySpeed := rl.NewVector2(float32(math.Abs(float64(bat.Body.Velocity.X*25))*float64(rl.GetFrameTime())),
		// 		float32(math.Abs(float64(bat.Body.Velocity.Y*5))*float64(rl.GetFrameTime())))
		// 	bat.Body.Velocity.X = applyVelocityDecay(bat.Body.Velocity.X, decaySpeed.X)
		// 	// Y velocity is not changed

		// 	if bat.DeflectTime >= deflectDuration {
		// 		bat.IsDeflecting = false
		// 		bat.DeflectTime = 0
		// 	}
		// }

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

// Handles player and enemy collision, and applies damage and despawns enemy if so
func checkEnemyPlayerCollision(enemies *Enemies, DeLorean *Vehicle, knightRect rl.Rectangle) {
	for i := len(enemies.Shooting) - 1; i >= 0; i-- {
		libyan := &enemies.Shooting[i]
		if rl.CheckCollisionCircleRec(libyan.Position, libyan.Body.Radius, knightRect) {
			// rl.PlaySound(audio.Sounds["damaged"])
			// knight.Health -= bat.Damage
			// knight.Health = rl.Clamp(knight.Health, 0, 100)
			// *healthTracker = knight.Health / 100

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

// Implements title screen and UI elements
func displayTitleScreen() {

	title := "Route 88"
	titleWidth := rl.MeasureText(title, 50)
	rl.DrawText(title,
		int32(screenSize.X)/2-(titleWidth/2),
		int32(screenSize.Y)/4-(50/2),
		50, rl.White)

	playButton := NewButton(0, 0, 400, 100, 0.1, 0, 0, TITLE, GAMEPLAY)
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

func displayHowToScreen() {

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
