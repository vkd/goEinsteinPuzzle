package goeinstein

import (
	"github.com/veandco/go-sdl2/sdl"
)

var options = struct {
	Fullscreen *BoolConfig
	NiceCursor *BoolConfig
	AutoHints  *BoolConfig
}{
	Fullscreen: NewBoolConfigCmd("fullscreen", false, screen.SetFullscreen),
	NiceCursor: NewBoolConfigCmd("niceCursor", true, screen.SetCursor),
	AutoHints:  NewBoolConfig("autoHints", false),
}

type OptionsChangedCommand struct {
	volume    *float32
	oldVolume float32
}

var _ Command = (*OptionsChangedCommand)(nil)

func NewOptionsChangedCommand(v *float32) *OptionsChangedCommand {
	o := &OptionsChangedCommand{
		volume:    v,
		oldVolume: *v,
	}
	return o
}

func (o *OptionsChangedCommand) DoAction() {
	if *o.volume != o.oldVolume {
		GetStorage().SetInt("volume", int(*o.volume*100.0))
		sound.SetVolume(*o.volume)
	}
	GetStorage().Flush()
}

func ShowOptionsWindow(parentArea *Area) {
	titleFont := NewFont("nova.ttf", 26)
	font := NewFont("laudcn2.ttf", 14)

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

	x := int32(240)
	var checkboxCommands Commands
	for _, ch := range []*BoolConfig{
		options.Fullscreen,
		options.NiceCursor,
		options.AutoHints,
	} {
		OPTION(x, ch.name, &ch.value)
		checkboxCommands = append(checkboxCommands, ch)
		x += 20
	}

	area.Add(NewLabelAligh(font, 265, 330, 300, 20, ALIGN_LEFT, ALIGN_MIDDLE, 255, 255, 255, msg("volume")))
	area.Add(NewSlider(360, 332, 160, 16, &volume))

	exitCmd := area.FinishCommand()
	okCmd := Combine(
		checkboxCommands,
		NewOptionsChangedCommand(&volume),
		exitCmd,
	)
	area.Add(NewButtonText(315, 390, 85, 25, font, 255, 255, 0, "blue.bmp", msg("ok"), okCmd))
	area.Add(NewButtonText(405, 390, 85, 25, font, 255, 255, 0, "blue.bmp", msg("cancel"), exitCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))
	area.Add(NewKeyAccel(sdl.K_RETURN, okCmd))
	area.Run()
}

type BoolConfig struct {
	name     string
	value    bool
	oldValue bool
	onSaveFn BoolConfigOnSaveFunc
}

type BoolConfigOnSaveFunc func(bool)

var boolToInt = map[bool]int{
	false: 0,
	true:  1,
}

func NewBoolConfigCmd(name string, v bool, onSave BoolConfigOnSaveFunc) *BoolConfig {
	v = GetStorage().GetInt(name, boolToInt[v]) > 0
	return &BoolConfig{
		name:     name,
		value:    v,
		oldValue: v,
		onSaveFn: onSave,
	}
}

func NewBoolConfig(name string, v bool) *BoolConfig {
	return NewBoolConfigCmd(name, v, nil)
}

func (b *BoolConfig) Value() bool { return b.value }

func (b *BoolConfig) Save() {
	if b.value == b.oldValue {
		return
	}

	b.oldValue = b.value
	GetStorage().SetInt(b.name, boolToInt[b.value])
	if b.onSaveFn != nil {
		b.onSaveFn(b.value)
	}
}

func (b *BoolConfig) SaveCommand() Command {
	return FnCommand(b.Save)
}

func (b *BoolConfig) DoAction() {
	b.Save()
}
