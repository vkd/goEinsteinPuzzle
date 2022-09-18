package goeinstein

import (
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screen *Screen = NewScreen()

// var rndGen Random

var atexit = []func(){
	GetStorage().Close,
}

func initScreen() {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		panic(fmt.Errorf("Error initializing SDL: %w", err))
	}
	atexit = append(atexit, sdl.Quit)
	err = ttf.Init()
	if err != nil {
		panic(fmt.Errorf("Error initializing font engine: %w", err))
	}
	atexit = append(atexit, ttf.Quit)
	screen.SetMode(NewVideoMode(800, 600, 24, GetStorage().GetInt("fullscreen", 0) > 0))
	screen.InitCursors()

	mouse := LoadImage("cursor.bmp")
	SDL_SetColorKey(mouse, true, sdl.MapRGB(mouse.Format, 0, 0, 0))
	screen.SetMouseImage(mouse)
	SDL_FreeSurface(mouse)
	// SDL_WM_SetCaption("Einstein", NULL)

	screen.SetCursor(GetStorage().GetInt("niceCursor", 1) > 0)
}

func initAudio() {
	sound = NewSound()
	sound.SetVolume(float32(GetStorage().GetInt("volume", 20)) / 100.0)
}

func Main() error {
	defer func() {
		for i := len(atexit) - 1; i >= 0; i-- {
			atexit[i]()
		}
	}()

	runtime.LockOSThread()

	EnsureDirExists("./einstein")

	// LoadResources()
	initScreen()
	initAudio()
	Menu()
	GetStorage().Flush()

	screen.DoneCursors()
	return nil
}
