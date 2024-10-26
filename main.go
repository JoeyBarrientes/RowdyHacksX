package main

import (
	"fmt"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Global Variables
// Global Variables
var initScreenWidth int32 = 1920
var initScreenHeight int32 = 1080
var screenSize rl.Vector2 = rl.NewVector2(float32(initScreenWidth), float32(initScreenHeight))
var screenScale = rl.NewVector2((float32(screenSize.X) / float32(initScreenWidth)),
	(float32(screenSize.Y) / float32(initScreenHeight)))

var backgroundImage rl.Texture2D
var backgroundScale float32

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

func displayTitleScreen() {

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
