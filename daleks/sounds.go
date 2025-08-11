package daleks

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
