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

const (
	sampleRate = 44100
)

type SoundPlayer struct {
	audioContext *audio.Context
	sounds       map[string]*audio.Player
}

func NewSoundPlayer() (*SoundPlayer, error) {
	audioContext := audio.NewContext(sampleRate)

	// Initialize sound map
	sounds := make(map[string]*audio.Player)

	// Load sound effects
	soundData := map[string][]byte{
		"teleport":    teleportData,
		"screwdriver": screwdriverData,
		"crash":       crashData,
		"gamestart":   gamestartData,
		"gameover":    gameoverData,
	}

	for name, data := range soundData {
		d, err := wav.Decode(audioContext, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		player, err := audio.NewPlayer(audioContext, d)
		if err != nil {
			return nil, err
		}

		sounds[name] = player
	}

	return &SoundPlayer{
		audioContext: audioContext,
		sounds:       sounds,
	}, nil
}

func (s *SoundPlayer) Play(name string) {
	if player, exists := s.sounds[name]; exists {
		player.Rewind()
		player.Play()
	}
}

func (s *SoundPlayer) Close() {
	for _, player := range s.sounds {
		player.Close()
	}
}
