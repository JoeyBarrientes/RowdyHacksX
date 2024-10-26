package main

import rl "github.com/gen2brain/raylib-go/raylib"

type PipCounter struct {
	X           int32
	Y           int32
	pipWidth    int32
	pipHeight   int32
	pipSpacing  int32
	maxInRow    int
	currentPips int
	maxPips     int
	colorTheme  *ColorTheme
}

func NewPipCounter(newX, newY, newPWidth, newPHeight, newSpacing int32, newRowMax, newPipMax int, newTheme *ColorTheme) PipCounter {
	pc := PipCounter{
		X:          newX,
		Y:          newY,
		pipWidth:   newPWidth,
		pipHeight:  newPHeight,
		pipSpacing: newSpacing,
		maxInRow:   newRowMax,
		maxPips:    newPipMax,
		colorTheme: newTheme,
	}
	pc.currentPips = pc.maxPips
	return pc
}

func (pc *PipCounter) SetPips(newPips int) {
	pc.currentPips = newPips
}

func (pc *PipCounter) AddPip() {
	pc.currentPips += 1
}

func (pc PipCounter) DrawPipCounter() {
	y := int32(0)
	x := int32(0)
	for i := 1; i <= pc.maxPips; i++ {
		color := pc.colorTheme.accentColor

		if i > pc.currentPips {
			color = pc.colorTheme.baseColor
		}

		rl.DrawRectangle(x+pc.X, y+pc.Y, pc.pipWidth, pc.pipHeight, color)
		x += (pc.pipSpacing + pc.pipWidth)
		if i%pc.maxInRow == 0 {
			y += pc.pipSpacing + pc.pipHeight
			x = 0
		}
	}
}
