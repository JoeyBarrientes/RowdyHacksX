package main

import (
	"fmt"
	"math"
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

var currentScreen GameScreen = TITLE

var score float32 = 0
var highScore float32 = 0

func main() {
	// Game variable initializations
	rl.InitWindow(initScreenWidth, initScreenHeight, "Game")
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.InitAudioDevice()
	audio.loadAudio()
	theme = NewColorTheme(rl.NewColor(176, 166, 170, 255), rl.NewColor(145, 63, 86, 255), rl.White)
	checkResize(&screenSize)

	vehicleTexture = rl.LoadTexture("textures/delorean.png")
	DeLorean := NewVehicle(vehicleTexture, rl.White, rl.NewVector2(screenSize.X/8*7, screenSize.Y/2-float32(vehicleTexture.Height/2)), rl.NewVector2(0, 0), float32(vehicleTexture.Width), float32(vehicleTexture.Height), 4)

	checkResize(&screenSize)
	rl.SetTargetFPS(60)
	// Main game loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(theme.backgroundColor)

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
			DeLorean.Position = rl.NewVector2(screenSize.X/8*7, screenSize.Y/2-float32(vehicleTexture.Height/2))
			DeLorean.DrawCharacter()
			DeLorean.Sprite.SourceRect.X = DeLorean.Sprite.SourceRect.Width * float32(DeLorean.Sprite.SpriteFrame)
			frameCount++
			if DeLorean.Sprite.SpriteFrame > 2 {
				DeLorean.Sprite.SpriteFrame = 0
				DeLorean.Sprite.SourceRect = rl.NewRectangle(0, 0, 64, 64)
			}
			if frameCount%10 == 1 {
				DeLorean.Sprite.SpriteFrame++

			}
			if frameCount == 30 {
				frameCount = 0
			}
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				DeLorean.shoot()
			}

			// if rl.IsKeyPressed(rl.KeySpace) {
			DeLorean.decreaseSpeed()
			// }

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
