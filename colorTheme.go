package main

import rl "github.com/gen2brain/raylib-go/raylib"

type ColorTheme struct {
	backgroundColor   rl.Color
	baseColor         rl.Color
	brownBaseColor    rl.Color
	brownAccentColor  rl.Color
	orangeBaseColor   rl.Color
	orangeAccentColor rl.Color
	redBaseColor      rl.Color
	blueBaseColor     rl.Color
	accentColor       rl.Color
	textColor         rl.Color
	colorGradient     []rl.Color
}

// type colorGradient struct

func NewColorTheme(base, accent, text rl.Color) ColorTheme {
	ct := ColorTheme{
		backgroundColor:   rl.NewColor(30, 17, 16, 255),
		brownBaseColor:    rl.NewColor(114, 61, 70, 255),
		brownAccentColor:  rl.NewColor(141, 100, 94, 255),
		orangeBaseColor:   rl.NewColor(166, 47, 3, 255),
		orangeAccentColor: rl.NewColor(242, 92, 5, 255),
		redBaseColor:      rl.NewColor(128, 15, 47, 255),
		blueBaseColor:     rl.NewColor(144, 224, 239, 255),
		baseColor:         base,
		accentColor:       accent,
		textColor:         text,
		colorGradient: []rl.Color{
			rl.Brown,
			rl.Gray,
			rl.RayWhite,
		},
	}
	return ct
}
