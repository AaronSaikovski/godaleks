// MIT License

// Copyright (c) 2025 - Aaron Saikovski <asaikovski@outlook.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package cmd

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

//go:embed assets/teleport.wav
var teleportData []byte

//go:embed assets/screwdriver.wav
var screwdriverData []byte

//go:embed assets/crash.wav
var crashData []byte

//go:embed assets/gameover.wav
var gameoverData []byte

//go:embed assets/gamestart.wav
var gamestartData []byte

//go:embed assets/dalek_emperor.wav
var dalekEmperorData []byte

const (
	sampleRate = 44100
)

type SoundPlayer struct {
	audioContext *audio.Context
	soundData    map[string][]byte // Raw WAV data for creating new players on demand
}

func NewSoundPlayer() (*SoundPlayer, error) {
	audioContext := audio.NewContext(sampleRate)

	// Store raw sound data for on-demand player creation
	sounds := map[string][]byte{
		"teleport":      teleportData,
		"screwdriver":   screwdriverData,
		"crash":         crashData,
		"gamestart":     gamestartData,
		"gameover":      gameoverData,
		"dalek_emperor": dalekEmperorData,
	}

	// Validate all WAV data can be decoded
	for name, data := range sounds {
		_, err := wav.Decode(audioContext, bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to decode sound %s: %w", name, err)
		}
	}

	return &SoundPlayer{
		audioContext: audioContext,
		soundData:    sounds,
	}, nil
}

func (s *SoundPlayer) Play(name string) {
	if s == nil || s.soundData == nil {
		return // Gracefully handle nil soundPlayer
	}
	data, exists := s.soundData[name]
	if !exists {
		return
	}
	// Create a fresh player each time to allow overlapping playback
	d, err := wav.Decode(s.audioContext, bytes.NewReader(data))
	if err != nil {
		return
	}
	player, err := audio.NewPlayer(s.audioContext, d)
	if err != nil {
		return
	}
	player.Play()
}
