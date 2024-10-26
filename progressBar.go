package main

import rl "github.com/gen2brain/raylib-go/raylib"

type ProgressBar struct {
	X               int32
	Y               int32
	Width           int32
	Height          int32
	progress        float32
	BorderThickness int32
	colorTheme      *ColorTheme
}

func (pb *ProgressBar) SetProgress(newProgress float32) {
	pb.progress = newProgress
	if pb.progress < 0 {
		pb.progress = 0
	}
	if pb.progress > 1 {
		pb.progress = 1
	}
}

func (pb ProgressBar) DrawBar() {

	bothBorderEdges := pb.BorderThickness * 2
	rl.DrawRectangle(pb.X, pb.Y, pb.Width, pb.Height, rl.Black)

	// Draws inner progress bar with an X and Y offset, and decreased Width and Height
	// to make the outer progress bar rectangle look like the bar outline
	rl.DrawRectangle(pb.X+pb.BorderThickness,
		pb.Y+pb.BorderThickness,
		int32(pb.progress*(float32(pb.Width)-float32(bothBorderEdges))),
		pb.Height-bothBorderEdges,
		pb.colorTheme.blueBaseColor)
}

func NewProgressBar(newX, newY, newBorderThickness, newWidth, newHeight int32, newTheme *ColorTheme) ProgressBar {
	pb := ProgressBar{X: newX, Y: newY, BorderThickness: newBorderThickness, Width: newWidth, Height: newHeight}
	pb.colorTheme = newTheme
	pb.progress = 0
	return pb
}
