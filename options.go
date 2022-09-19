package goeinstein

import (
	"github.com/veandco/go-sdl2/sdl"
)

type OptionsChangedCommand struct {
	fullscreen *bool
	niceCursor *bool
	volume     *float32
	oldVolume  float32
	area       *Area
}

var _ Command = (*OptionsChangedCommand)(nil)

func NewOptionsChangedCommand(a *Area, fs *bool, ns *bool, v *float32) *OptionsChangedCommand {
	o := &OptionsChangedCommand{
		fullscreen: fs,
		niceCursor: ns,
		volume:     v,
		oldVolume:  *v,
	}
	o.area = a
	return o
}

func (o *OptionsChangedCommand) DoAction() {
	bool2Int := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}
	oldFullscreen := screen.fullScreen
	oldCursor := screen.niceCursor
	if *o.fullscreen != oldFullscreen {
		GetStorage().SetInt("fullscreen", bool2Int(*o.fullscreen))
		screen.SetMode(NewVideoMode(800, 600, 24, *o.fullscreen))
	}
	if *o.niceCursor != oldCursor {
		GetStorage().SetInt("niceCursor", bool2Int(*o.niceCursor))
		screen.SetCursor(*o.niceCursor)
	}
	if *o.volume != o.oldVolume {
		GetStorage().SetInt("volume", int(*o.volume*100.0))
		sound.SetVolume(*o.volume)
	}
	GetStorage().Flush()
	o.area.FinishEventLoop()
}

func ShowOptionsWindow(parentArea *Area) {
	titleFont := NewFont("nova.ttf", 26)
	font := NewFont("laudcn2.ttf", 14)

	fullscreen := screen.fullScreen
	niceCursor := screen.niceCursor
	volume := sound.volume

	area := NewArea()

	LABEL := func(y int32, s string) {
		area.Add(NewLabelAligh(font, 300, y, 300, 20, ALIGN_LEFT, ALIGN_MIDDLE, 255, 255, 255, msg(s)))
	}
	CHECKBOX := func(y int32, v *bool) {
		area.Add(NewCheckbox(265, y, 20, 20, font, 255, 255, 255, "blue.bmp", v))
	}
	OPTION := func(y int32, s string, v *bool) {
		LABEL(y, s)
		CHECKBOX(y, v)
	}

	area.Add(parentArea)
	area.Add(NewWindow(250, 170, 300, 260, "blue.bmp"))
	area.Add(NewLabelAligh(titleFont, 250, 175, 300, 40, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, 0, msg("options")))
	OPTION(260, "fullscreen", &fullscreen)
	OPTION(280, "niceCursor", &niceCursor)

	area.Add(NewLabelAligh(font, 265, 330, 300, 20, ALIGN_LEFT, ALIGN_MIDDLE, 255, 255, 255, msg("volume")))
	area.Add(NewSlider(360, 332, 160, 16, &volume))

	exitCmd := NewExitCommand(area)
	okCmd := NewOptionsChangedCommand(area, &fullscreen, &niceCursor, &volume)
	area.Add(NewButtonText(315, 390, 85, 25, font, 255, 255, 0, "blue.bmp", msg("ok"), okCmd))
	area.Add(NewButtonText(405, 390, 85, 25, font, 255, 255, 0, "blue.bmp", msg("cancel"), exitCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))
	area.Add(NewKeyAccel(sdl.K_RETURN, okCmd))
	area.Run()
}
