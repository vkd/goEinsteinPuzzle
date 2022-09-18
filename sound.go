package goeinstein

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var sound *Sound = NewSound()

type Sound struct {
	disabled   bool
	chunkCache map[string]*mix.Chunk

	enableFx bool
	volume   float32
}

func NewSound() *Sound {
	s := &Sound{
		chunkCache: make(map[string]*mix.Chunk),
	}
	audioRate := 22050
	var audioFormat uint16 = sdl.AUDIO_S16
	audioChannels := 2
	audioBuffers := 1024

	err := mix.OpenAudio(audioRate, audioFormat, audioChannels, audioBuffers)
	if err != nil {
		log.Printf("Error on open audio: %v", err)
		s.disabled = true
	}
	return s
}

func (s *Sound) Close() {
	if !s.disabled {
		mix.CloseAudio()
	}
	for _, c := range s.chunkCache {
		c.Free()
	}
	mix.CloseAudio()
}

func (s *Sound) Play(name string) {
	if s.disabled || !s.enableFx {
		return
	}

	var chunk *mix.Chunk

	ch, ok := s.chunkCache[name]
	if ok {
		chunk = ch
	} else {
		bs := resources.GetRef(name)

		rw, err := sdl.RWFromMem(bs)
		if err != nil {
			panic(fmt.Errorf("rw from mem: %w", err))
		}

		chunk, err = mix.LoadWAVRW(rw, false)
		if err != nil {
			panic(fmt.Errorf("load WAV: %w", err))
		}
		s.chunkCache[name] = chunk
	}

	if chunk != nil {
		chunk.Volume(int(s.volume * 128))
		Mix_PlayChannel(chunk, -1, 0)
	}
	sdl.PumpEvents()
}

func (s *Sound) SetVolume(v float32) {
	s.volume = v
	s.enableFx = 0.01 < s.volume
}
