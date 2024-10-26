package main

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Uses maps to index by string (name of sound)
type Audio struct {
	Sounds  map[string]rl.Sound
	Music   map[string]rl.Music
	isMuted bool
}

func NewAudio() Audio {
	return Audio{
		Sounds: make(map[string]rl.Sound),
		Music:  make(map[string]rl.Music),
	}
}

func (audio *Audio) loadAudio() {
	// Load and map music

	// Load and map sounds

	// Set Audio Volumes and Pitches
	audio.setAudioVolumes()
	audio.setAudioPitches()

	// Play background music

}

func (audio *Audio) setAudioVolumes() {

}

func (audio *Audio) setAudioPitches() {

}

func (audio *Audio) checkMute() {
	if rl.IsKeyPressed(rl.KeyM) && !audio.isMuted {
		rl.SetMasterVolume(0)
		audio.isMuted = !audio.isMuted

	} else if rl.IsKeyPressed(rl.KeyM) && audio.isMuted {
		rl.SetMasterVolume(1)
		audio.isMuted = !audio.isMuted
	}

	if audio.isMuted {
		rl.DrawText("Muted", int32(rl.GetScreenWidth())-140, int32(rl.GetScreenHeight())-50, 40, rl.White)
	}
}

func (audio *Audio) playWithRandPitch(sound rl.Sound) {
	// Get random pitch value from 0.7 to 1.4
	randPitch := rl.Clamp(rand.Float32()+0.7, 0.7, 1.4)

	rl.SetSoundPitch(sound, randPitch)
	rl.PlaySound(sound)
}
