package goeinstein

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type Command interface {
	DoAction()
}

type Widget struct {
	area *Area
}

var _ AreaWidgeter = (*Widget)(nil)

func (w *Widget) Close()                                          {}
func (w *Widget) OnMouseButtonDown(button uint8, x, y int32) bool { return false }
func (w *Widget) OnMouseButtonUp(button uint8, x, y int32) bool   { return false }
func (w *Widget) OnMouseMove(x, y int32) bool                     { return false }
func (w *Widget) Draw()                                           {}
func (w *Widget) SetParent(a *Area)                               { w.area = a }
func (w *Widget) OnKeyDown(sdl.Keycode, sdl.Scancode) bool        { return false }
func (w *Widget) DestroyByArea() bool                             { return true }

type Button struct {
	Widget

	left, top, width, height int32
	image, highlighted       *sdl.Surface
	mouseInside              bool
	command                  Command
}

func (b *Button) GetLeft() int32   { return b.left }
func (b *Button) GetTop() int32    { return b.top }
func (b *Button) GetWidth() int32  { return b.width }
func (b *Button) GetHeight() int32 { return b.height }

func (b *Button) MoveTo(x, y int32) {
	b.left = x
	b.top = y
}

func NewButton(x, y int32, name string, cmd Command) *Button {
	return NewButtonTransparent(x, y, name, cmd, true)
}

func NewButtonTransparent(x, y int32, name string, cmd Command, transparent bool) *Button {
	b := &Button{}
	b.image = LoadImageTransparent(name, transparent)
	b.highlighted = AdjustBrightnessTransparent(b.image, 1.5, transparent)

	b.left = x
	b.top = y
	b.width = b.image.W
	b.height = b.image.H

	b.mouseInside = false
	b.command = cmd

	return b
}

func NewButtonColor(x, y, w, h int32, font *Font, fR, fG, fB uint8, hR, hG, hB uint8, text string, cmd Command) *Button {
	b := &Button{}
	b.left = x
	b.top = y
	b.width = w
	b.height = h

	s, err := sdl.CreateRGBSurface(sdl.SWSURFACE, w, h, 24, 0x00FF0000, 0x0000FF00, 0x000000FF, 0)
	if err != nil {
		panic(fmt.Errorf("create RGB surface: %w", err))
	}
	src := &sdl.Rect{x, y, b.width, b.height}
	dst := &sdl.Rect{0, 0, b.width, b.height}
	SDL_BlitSurface(screen.GetSurface(), src, s, dst)

	tW, tH := font.GetSize(text)
	font.DrawSurface(s, (b.width-tW)/2, (b.height-tH)/2, fR, fG, fB, true, text)
	b.image = SDL_DisplayFormat(s)
	SDL_BlitSurface(screen.GetSurface(), src, s, dst)
	font.DrawSurface(s, (b.width-tW)/2, (b.height-tH)/2, hR, hG, hB, true, text)
	b.highlighted = SDL_DisplayFormat(s)
	s.Free()

	b.mouseInside = false
	b.command = cmd
	return b
}

func NewButtonTextBevel(x, y, w, h int32, font *Font, r, g, b uint8, bg string, text string, bevel bool, cmd Command) *Button {
	btn := &Button{}
	btn.left = x
	btn.top = y
	btn.width = w
	btn.height = h

	s := screen.GetSurface()
	btn.image = SDL_CreateRGBSurface(sdl.SWSURFACE, btn.width, btn.height, int32(s.Format.BitsPerPixel), s.Format.Rmask, s.Format.Gmask, s.Format.Bmask, s.Format.Amask)

	tile := LoadImageTransparent(bg, true)
	src := &sdl.Rect{0, 0, tile.W, tile.H}
	dst := &sdl.Rect{0, 0, tile.W, tile.H}
	for j := int32(0); j < btn.height; j += tile.H {
		for i := int32(0); i < btn.width; i += tile.W {
			dst.X = i
			dst.Y = j
			SDL_BlitSurface(tile, src, btn.image, dst)
		}
	}
	tile.Free()

	if bevel {
		SDL_LockSurface(btn.image)
		DrawBevel(btn.image, 0, 0, btn.width, btn.height, false, 1)
		DrawBevel(btn.image, 1, 1, btn.width-2, btn.height-2, true, 1)
		SDL_UnlockSurface(btn.image)
	}

	tW, tH := font.GetSize(text)
	font.DrawSurface(btn.image, (btn.width-tW)/2, (btn.height-tH)/2, r, g, b, true, text)

	btn.highlighted = AdjustBrightnessTransparent(btn.image, 1.5, false)
	SDL_SetColorKey(btn.image, true, GetCornerPixel(btn.image))
	SDL_SetColorKey(btn.highlighted, true, GetCornerPixel(btn.highlighted))

	btn.mouseInside = false
	btn.command = cmd
	return btn
}

func NewButtonText(x, y, w, h int32, font *Font, r, g, b uint8, bg string, text string, cmd Command) *Button {
	btn := &Button{}
	btn.left = x
	btn.top = y
	btn.width = w
	btn.height = h

	s := screen.GetSurface()
	btn.image = SDL_CreateRGBSurface(sdl.SWSURFACE, btn.width, btn.height, int32(s.Format.BitsPerPixel), s.Format.Rmask, s.Format.Gmask, s.Format.Bmask, s.Format.Amask)

	tile := LoadImage(bg)
	src := &sdl.Rect{0, 0, tile.W, tile.H}
	dst := &sdl.Rect{0, 0, tile.W, tile.H}
	for j := int32(0); j < btn.height; j += tile.H {
		for i := int32(0); i < btn.width; i += tile.W {
			dst.X = i
			dst.Y = j
			SDL_BlitSurface(tile, src, btn.image, dst)
		}
	}
	SDL_FreeSurface(tile)

	SDL_LockSurface(btn.image)
	DrawBevel(btn.image, 0, 0, btn.width, btn.height, false, 1)
	DrawBevel(btn.image, 1, 1, btn.width-2, btn.height-2, true, 1)
	SDL_UnlockSurface(btn.image)

	tW, tH := font.GetSize(text)
	font.DrawSurface(btn.image, (btn.width-tW)/2, (btn.height-tH)/2, r, g, b, true, text)

	btn.highlighted = AdjustBrightnessTransparent(btn.image, 1.5, false)

	btn.mouseInside = false
	btn.command = cmd
	return btn
}

func (b *Button) Close() {
	SDL_FreeSurface(b.image)
	SDL_FreeSurface(b.highlighted)
}

func (b *Button) Draw() {
	if b.mouseInside {
		screen.Draw(b.left, b.top, b.highlighted)
	} else {
		screen.Draw(b.left, b.top, b.image)
	}
	screen.AddRegionToUpdate(b.left, b.top, b.width, b.height)
}

func (b *Button) GetBounds() (l, t, w, h int32) {
	return b.left, b.top, b.width, b.height
}

func (b *Button) OnMouseButtonDown(button uint8, x, y int32) bool {
	if IsInRect(x, y, b.left, b.top, b.width, b.height) {
		sound.Play("click.wav")
		if b.command != nil {
			b.command.DoAction()
		}
		return true
	}
	return false
}

func (b *Button) OnMouseMove(x, y int32) bool {
	in := IsInRect(x, y, b.left, b.top, b.width, b.height)
	if in != b.mouseInside {
		b.mouseInside = in
		b.Draw()
	}
	return false
}

type KeyAccel struct {
	Widget

	key     SDLKey
	command Command
}

func NewKeyAccel(sym sdl.Keycode, cmd Command) *KeyAccel {
	return &KeyAccel{
		key:     sym,
		command: cmd,
	}
}

func (ka *KeyAccel) OnKeyDown(k SDLKey, ch sdl.Scancode) bool {
	if ka.key == k {
		if ka.command != nil {
			ka.command.DoAction()
		}
		return true
	}
	return false
}

type TimerHandler interface {
	OnTimer()
}

type AreaWidgeter interface {
	// Area
	DestroyByArea() bool

	// Widget
	Draw()
	SetParent(*Area)
	OnMouseButtonDown(button uint8, x, y int32) bool
	OnMouseButtonUp(button uint8, x, y int32) bool
	OnMouseMove(x, y int32) bool
	OnKeyDown(sdl.Keycode, sdl.Scancode) bool
	Close()
}

type Area struct {
	Widget

	widgets           []AreaWidgeter
	notManagedWidgets map[AreaWidgeter]int
	terminate         bool
	time              uint64
	timer             TimerHandler
}

func (a *Area) DestroyByArea() bool { return false }

func NewArea() *Area {
	a := &Area{
		notManagedWidgets: make(map[AreaWidgeter]int),
	}
	a.timer = nil
	return a
}

func (a *Area) Close() {
	for _, w := range a.widgets {
		if w != nil && w.DestroyByArea() && (a.notManagedWidgets[w] == 0) {
			w.Close()
		}
	}
}

func (a *Area) Add(w AreaWidgeter) {
	a.AddManaged(w, true)
}

func (a *Area) AddManaged(w AreaWidgeter, managed bool) {
	a.widgets = append(a.widgets, w)
	if !managed {
		a.notManagedWidgets[w]++
	}
	w.SetParent(a)
}

func (a *Area) Remove(w AreaWidgeter) {
	for i := range a.widgets {
		if a.widgets[i] == w {
			a.widgets = append(a.widgets[:i], a.widgets[i+1:]...)
			break
		}
	}
	a.notManagedWidgets[w]++
}

func (a *Area) HandleEvent(event sdl.Event) {
	switch event := event.(type) {
	case *sdl.MouseButtonEvent:
		switch event.Type {
		case sdl.MOUSEBUTTONDOWN:
			for _, w := range a.widgets {
				if w.OnMouseButtonDown(event.Button, event.X, event.Y) {
					return
				}
			}
		case sdl.MOUSEBUTTONUP:
			for _, w := range a.widgets {
				if w.OnMouseButtonUp(event.Button, event.X, event.Y) {
					return
				}
			}
		}
	case *sdl.MouseMotionEvent:
		for _, w := range a.widgets {
			if w.OnMouseMove(event.X, event.Y) {
				return
			}
		}
	case *sdl.WindowEvent:
		switch event.Type {
		case sdl.WINDOWEVENT_EXPOSED:
			for _, w := range a.widgets {
				w.Draw()
			}
		}
	case *sdl.KeyboardEvent:
		switch event.Type {
		case sdl.KEYDOWN:
			for _, w := range a.widgets {
				if w.OnKeyDown(event.Keysym.Sym, event.Keysym.Scancode) {
					return
				}
			}
		}
	case *sdl.QuitEvent:
		os.Exit(0)
	}
}

func (a *Area) Run() {
	a.terminate = false
	var event sdl.Event

	var lastTimer uint64
	a.Draw()
	screen.ShowMouse()

	runTimer := a.timer != nil
	var dispetchEvent bool
	for !a.terminate {
		dispetchEvent = true
		if a.timer == nil {
			event = sdl.WaitEvent()
		} else {
			now := sdl.GetTicks64()
			if (now - lastTimer) > a.time {
				lastTimer = now
				runTimer = true
			}
			event = sdl.PollEvent()
			if event == nil {
				if !runTimer {
					sdl.Delay(20)
					continue
				} else {
					dispetchEvent = false
				}
			}
		}
		screen.HideMouse()
		if runTimer {
			if a.timer != nil {
				a.timer.OnTimer()
			}
			runTimer = false
		}
		if dispetchEvent {
			a.HandleEvent(event)
		}
		if !a.terminate {
			screen.ShowMouse()
			screen.Flush()
		}
	}
}

func (a *Area) FinishEventLoop() {
	a.terminate = true
}

func (a *Area) Draw() {
	for _, w := range a.widgets {
		w.Draw()
	}
}

func (a *Area) SetTimer(interval uint64, t TimerHandler) {
	a.time = interval
	a.timer = t
}

func (a *Area) UpdateMouse() {
	x, y, _ := sdl.GetMouseState()
	for _, w := range a.widgets {
		if w.OnMouseMove(x, y) {
			return
		}
	}
}

type ExitCommand struct {
	area *Area
}

var _ Command = (*ExitCommand)(nil)

func NewExitCommand(a *Area) *ExitCommand {
	return &ExitCommand{
		area: a,
	}
}

func (e *ExitCommand) DoAction() {
	e.area.FinishEventLoop()
}

type AnyKeyAccel struct {
	Widget
	command Command
}

func NewAnyKeyAccelDefault() *AnyKeyAccel {
	return &AnyKeyAccel{
		command: nil,
	}
}

func NewAnyKeyAccel(cmd Command) *AnyKeyAccel {
	return &AnyKeyAccel{
		command: cmd,
	}
}

func (a *AnyKeyAccel) Close() {}

func (a *AnyKeyAccel) OnKeyDown(key SDLKey, ch sdl.Scancode) bool {
	if ((key >= sdl.K_NUMLOCKCLEAR) && (key <= sdl.K_APPLICATION)) || (key == sdl.K_TAB) || (key == sdl.K_UNKNOWN) {
		return false
	}

	if a.command != nil {
		a.command.DoAction()
	} else {
		a.area.FinishEventLoop()
	}
	return true
}

func (a *AnyKeyAccel) OnMouseButtonDown(button uint8, x, y int32) bool {
	if a.command != nil {
		a.command.DoAction()
	} else {
		a.area.FinishEventLoop()
	}
	return true
}

type Window struct {
	Widget
	left, top, width, height int32
	background               *sdl.Surface
}

func NewWindow(x, y, w, h int32, bg string) *Window {
	return NewWindowFrameRaised(x, y, w, h, bg, 4, true)
}

func NewWindowFrame(x, y, w, h int32, bg string, frameWidth int32) *Window {
	return NewWindowFrameRaised(x, y, w, h, bg, frameWidth, true)
}

func NewWindowRaised(x, y, w, h int32, bg string, raised bool) *Window {
	return NewWindowFrameRaised(x, y, w, h, bg, 4, raised)
}

func NewWindowFrameRaised(x, y, w, h int32, bg string, frameWidth int32, raised bool) *Window {
	wn := &Window{}
	wn.left = x
	wn.top = y
	wn.width = w
	wn.height = h

	s := screen.GetSurface()
	win := SDL_CreateRGBSurface(sdl.SWSURFACE, wn.width, wn.height, int32(s.Format.BitsPerPixel), s.Format.Rmask, s.Format.Gmask, s.Format.Bmask, s.Format.Amask)

	tile := LoadImage(bg)
	src := &sdl.Rect{0, 0, tile.W, tile.H}
	dst := &sdl.Rect{0, 0, tile.W, tile.H}
	for j := int32(0); j < wn.height; j += tile.H {
		for i := int32(0); i < wn.width; i += tile.W {
			dst.X = i
			dst.Y = j
			SDL_BlitSurface(tile, src, win, dst)
		}
	}
	SDL_FreeSurface(tile)

	SDL_LockSurface(win)
	k := 2.6
	f := 0.1
	for i := int32(0); i < frameWidth; i++ {
		var ltK, rbK float64
		if raised {
			ltK = k
			rbK = f
		} else {
			ltK = f
			rbK = k
		}
		for j := i; j < wn.height-i-1; j++ {
			AdjustBrightness(win, i, j, ltK)
		}
		for j := i; j < wn.width-i; j++ {
			AdjustBrightness(win, j, i, ltK)
		}
		for j := i + 1; j < wn.height-i; j++ {
			AdjustBrightness(win, wn.width-i-1, j, rbK)
		}
		for j := i; j < wn.width-i-1; j++ {
			AdjustBrightness(win, j, wn.height-i-1, rbK)
		}
		k -= 0.2
		f += 0.1
	}
	SDL_UnlockSurface(win)

	wn.background = SDL_DisplayFormat(win)
	SDL_FreeSurface(win)
	return wn
}

func (w *Window) Close() {
	SDL_FreeSurface(w.background)
}

func (w *Window) Draw() {
	screen.Draw(w.left, w.top, w.background)
	screen.AddRegionToUpdate(w.left, w.top, w.width, w.height)
}

type HorAlign int8

//nolint:golint,stylecheck
const (
	ALIGN_LEFT HorAlign = iota
	ALIGN_CENTER
	ALIGN_RIGHT
)

type VerAlign int8

//nolint:golint,stylecheck
const (
	ALIGN_TOP VerAlign = iota
	ALIGN_MIDDLE
	ALIGN_BOTTOM
)

type Label struct {
	Widget

	font                     *Font
	text                     string
	left, top, width, height int32
	red, green, blue         uint8
	hAlign                   HorAlign
	vAlign                   VerAlign
	shadow                   bool
}

func NewLabel(f *Font, x, y int32, r, g, b uint8, s string) *Label {
	return NewLabelShadow(f, x, y, r, g, b, s, true)
}

func NewLabelShadow(f *Font, x, y int32, r, g, b uint8, s string, sh bool) *Label {
	l := &Label{
		text: s,
	}
	l.font = f
	l.left = x
	l.top = y
	l.red = r
	l.green = g
	l.blue = b
	l.hAlign = ALIGN_LEFT
	l.vAlign = ALIGN_TOP
	l.shadow = sh
	return l
}

func NewLabelAligh(f *Font, x, y, w, h int32, hA HorAlign, vA VerAlign, r, g, b uint8, s string) *Label {
	l := &Label{}
	l.text = s
	l.font = f
	l.left = x
	l.top = y
	l.red = r
	l.green = g
	l.blue = b
	l.hAlign = hA
	l.vAlign = vA
	l.width = w
	l.height = h
	l.shadow = true
	return l
}

func (l *Label) Draw() {
	var x, y int32
	w, h := l.font.GetSize(l.text)

	switch l.hAlign {
	case ALIGN_RIGHT:
		x = l.left + l.width - w
	case ALIGN_CENTER:
		x = l.left + (l.width-w)/2
	default:
		x = l.left
	}

	switch l.vAlign {
	case ALIGN_BOTTOM:
		y = l.top + l.height - h
	case ALIGN_MIDDLE:
		y = l.top + (l.height-h)/2
	default:
		y = l.top
	}

	l.font.Draw(x, y, l.red, l.green, l.blue, l.shadow, l.text)
	screen.AddRegionToUpdate(x, y, w, h)
}

//nolint:unused
type InputField struct {
	Window

	text             *string
	maxLength        int
	cursorPos        int
	red, green, blue uint8
	font             *Font
	lastCursor       uint64
	cursorVisible    bool
	lastChar         sdl.Keysym
	lastKeyUpdate    uint32
}

var _ TimerHandler = (*InputField)(nil)

func NewInputField(x, y, w, h int32, background string, s *string, maxLen int, r, g, b uint8, f *Font) *InputField {
	i := &InputField{
		Window: *NewWindowFrameRaised(x, y, w, h, background, 1, false),
		text:   s,
	}
	i.maxLength = maxLen
	i.red = r
	i.green = g
	i.blue = b
	i.font = f
	i.MoveCursor(len(*i.text))
	return i
}

func (i *InputField) Close() {}

func (i *InputField) Draw() {
	i.Window.Draw()

	rect := sdl.Rect{i.left + 1, i.top + 1, i.width - 2, i.height - 2}
	SDL_SetClipRect(screen.GetSurface(), &rect)

	i.font.Draw(i.left+1, i.top+1, i.red, i.green, i.blue, true, *i.text)

	if i.cursorVisible {
		var pos int32
		if i.cursorPos > 0 {
			pos += i.font.GetWidth((*i.text)[0:i.cursorPos])
		}
		for ii := int32(2); ii < i.height-2; ii++ {
			screen.SetPixel(i.left+pos, i.top+ii, i.red, i.green, i.blue)
			screen.SetPixel(i.left+pos+1, i.top+ii, i.red, i.green, i.blue)
		}
	}

	SDL_SetClipRect(screen.GetSurface(), nil)
}

func (i *InputField) SetParent(a *Area) {
	i.Window.SetParent(a)
	i.area.SetTimer(100, i)
}

func (i *InputField) OnTimer() {
	now := sdl.GetTicks64()
	if (now - i.lastCursor) > 1000 {
		i.cursorVisible = !i.cursorVisible
		i.lastCursor = now
		i.Draw()
	}
}

func (i *InputField) OnKeyDown(key SDLKey, translatedChar sdl.Scancode) bool {
	switch key {
	case sdl.K_BACKSPACE:
		if i.cursorPos > 0 {
			*i.text = (*i.text)[:i.cursorPos-1] + (*i.text)[i.cursorPos:]
			i.MoveCursor(i.cursorPos - 1)
		} else {
			i.MoveCursor(i.cursorPos)
		}
		i.Draw()
		return true

	case sdl.K_LEFT:
		if i.cursorPos > 0 {
			i.MoveCursor(i.cursorPos - 1)
		} else {
			i.MoveCursor(i.cursorPos)
		}
		i.Draw()
		return true

	case sdl.K_RIGHT:
		if i.cursorPos < len(*i.text) {
			i.MoveCursor(i.cursorPos + 1)
		} else {
			i.MoveCursor(i.cursorPos)
		}
		i.Draw()
		return true

	case sdl.K_HOME:
		i.MoveCursor(0)
		i.Draw()
		return true

	case sdl.K_END:
		i.MoveCursor(len(*i.text))
		i.Draw()
		return true

	case sdl.K_DELETE:
		if i.cursorPos < len(*i.text) {
			*i.text = (*i.text)[:i.cursorPos] + (*i.text)[i.cursorPos+1:]
		}
		i.MoveCursor(i.cursorPos)
		i.Draw()
		return true

	default:
	}

	if key >= sdl.K_a && key <= sdl.K_z {
		i.OnCharTyped(translatedChar)
	}
	return false
}

func (i *InputField) OnKeyUp(key sdl.Keycode) bool {
	return false
}

func (i *InputField) OnCharTyped(ch sdl.Scancode) {
	if len(*i.text) < i.maxLength {
		*i.text = (*i.text)[:i.cursorPos] + sdl.GetScancodeName(ch) + (*i.text)[i.cursorPos:]
		i.MoveCursor(i.cursorPos + 1)
	} else {
		i.MoveCursor(i.cursorPos)
	}
	i.Draw()
}

func (i *InputField) MoveCursor(pos int) {
	i.lastCursor = sdl.GetTicks64()
	i.cursorVisible = true
	i.cursorPos = pos
}

type Checkbox struct {
	Widget

	left, top, width, height         int32
	image, highlighted               *sdl.Surface
	checkedImage, checkedHighlighted *sdl.Surface
	checked                          *bool
	mouseInside                      bool
}

func (c *Checkbox) GetLeft() int32   { return c.left }
func (c *Checkbox) GetTop() int32    { return c.top }
func (c *Checkbox) GetWidth() int32  { return c.width }
func (c *Checkbox) GetHeight() int32 { return c.height }
func (c *Checkbox) MoveTo(x, y int32) {
	c.left = x
	c.top = y
}

func NewCheckbox(x, y, w, h int32, font *Font, r, g, b uint8, bg string, chk *bool) *Checkbox {
	c := &Checkbox{
		checked: chk,
	}
	c.left = x
	c.top = y
	c.width = w
	c.height = h

	s := screen.GetSurface()
	c.image = SDL_CreateRGBSurface(sdl.SWSURFACE, c.width, c.height, int32(s.Format.BitsPerPixel), s.Format.Rmask, s.Format.Gmask, s.Format.Bmask, s.Format.Amask)

	tile := LoadImage(bg)
	src := &sdl.Rect{0, 0, tile.W, tile.H}
	dst := &sdl.Rect{0, 0, tile.W, tile.H}
	for j := int32(0); j < c.height; j += tile.H {
		for i := int32(0); i < c.width; i += tile.W {
			dst.X = i
			dst.Y = j
			SDL_BlitSurface(tile, src, c.image, dst)
		}
	}
	SDL_FreeSurface(tile)

	SDL_LockSurface(c.image)
	DrawBevel(c.image, 0, 0, c.width, c.height, false, 1)
	DrawBevel(c.image, 1, 1, c.width-2, c.height-2, true, 1)
	SDL_UnlockSurface(c.image)

	c.highlighted = AdjustBrightnessTransparent(c.image, 1.5, false)

	c.checkedImage = SDL_DisplayFormat(c.image)
	tW, tH := font.GetSize("X")
	tH += 2
	tW += 2
	font.DrawSurface(c.checkedImage, (c.width-tW)/2, (c.height-tH)/2, r, g, b, true, "X")
	c.checkedHighlighted = AdjustBrightnessTransparent(c.checkedImage, 1.5, false)

	c.mouseInside = false
	return c
}

func (c *Checkbox) Close() {
	SDL_FreeSurface(c.image)
	SDL_FreeSurface(c.highlighted)
	SDL_FreeSurface(c.checkedImage)
	SDL_FreeSurface(c.checkedHighlighted)
}

func (c *Checkbox) Draw() {
	if *c.checked {
		if c.mouseInside {
			screen.Draw(c.left, c.top, c.checkedHighlighted)
		} else {
			screen.Draw(c.left, c.top, c.checkedImage)
		}
	} else {
		if c.mouseInside {
			screen.Draw(c.left, c.top, c.highlighted)
		} else {
			screen.Draw(c.left, c.top, c.image)
		}
	}
	screen.AddRegionToUpdate(c.left, c.top, c.width, c.height)
}

func (c *Checkbox) GetBounds() (l, t, w, h int32) {
	l = c.left
	t = c.top
	w = c.width
	h = c.height
	return
}

func (c *Checkbox) OnMouseButtonDown(button uint8, x, y int32) bool {
	if IsInRect(x, y, c.left, c.top, c.width, c.height) {
		sound.Play("click.wav")
		*c.checked = !(*c.checked)
		c.Draw()
		return true
	}
	return false
}

func (c *Checkbox) OnMouseMove(x, y int32) bool {
	in := IsInRect(x, y, c.left, c.top, c.width, c.height)
	if in != c.mouseInside {
		c.mouseInside = in
		c.Draw()
	}
	return false
}

type Picture struct {
	Widget

	left, top, width, height int32
	image                    *sdl.Surface
	managed                  bool
}

func (p *Picture) GetLeft() int32   { return p.left }
func (p *Picture) GetTop() int32    { return p.top }
func (p *Picture) GetWidth() int32  { return p.width }
func (p *Picture) GetHeight() int32 { return p.height }

func NewPicture(x, y int32, name string, transparent bool) *Picture {
	p := &Picture{}
	p.image = LoadImageTransparent(name, transparent)
	p.left = x
	p.top = y
	p.width = p.image.W
	p.height = p.image.H
	p.managed = true
	return p
}

func NewPictureSurface(x, y int32, img *sdl.Surface) *Picture {
	p := &Picture{}
	p.image = img
	p.left = x
	p.top = y
	p.width = p.image.W
	p.height = p.image.H
	p.managed = false
	return p
}

func (p *Picture) Close() {
	if p.managed {
		SDL_FreeSurface(p.image)
	}
}

func (p *Picture) Draw() {
	screen.Draw(p.left, p.top, p.image)
	screen.AddRegionToUpdate(p.left, p.top, p.width, p.height)
}

func (p *Picture) MoveX(newX int32) {
	p.left = newX
}

func (p *Picture) GetBounds() (l, t, w, h int32) {
	l = p.left
	t = p.top
	w = p.width
	h = p.height
	return
}

type Slider struct {
	Widget

	left, top, width, height int32
	value                    *float32
	background               *sdl.Surface
	slider                   *sdl.Surface
	activeSlider             *sdl.Surface
	highlight                bool
	dragging                 bool
	dragOffsetX              int32
}

func NewSlider(x, y, w, h int32, v *float32) *Slider {
	s := &Slider{
		value: v,
	}
	s.left = x
	s.top = y
	s.width = w
	s.height = h
	s.background = nil
	s.CreateSlider(s.height)
	s.highlight = false
	s.dragging = false
	return s
}

func (s *Slider) Close() {
	if s.background != nil {
		SDL_FreeSurface(s.background)
	}
	if s.slider != nil {
		SDL_FreeSurface(s.slider)
	}
	if s.activeSlider != nil {
		SDL_FreeSurface(s.activeSlider)
	}
}

func (s *Slider) Draw() {
	if s.background == nil {
		s.CreateBackground()
	}
	screen.Draw(s.left, s.top, s.background)
	screen.AddRegionToUpdate(s.left, s.top, s.width, s.height)
	posX := s.ValueToX(*s.value)
	var srf *sdl.Surface
	if s.highlight {
		srf = s.activeSlider
	} else {
		srf = s.slider
	}
	screen.Draw(s.left+posX, s.top, srf)
}

func (s *Slider) CreateBackground() {
	s.background = screen.CreateSubimage(s.left, s.top, s.width, s.height)
	y := s.height / 2
	SDL_LockSurface(s.background)
	DrawBevel(s.background, 0, y-2, s.width, 4, false, 1)
	SDL_UnlockSurface(s.background)
}

func (s *Slider) CreateSlider(size int32) {
	srf := screen.GetSurface()
	image := SDL_CreateRGBSurface(sdl.SWSURFACE, size, size, int32(srf.Format.BitsPerPixel), srf.Format.Rmask, srf.Format.Gmask, srf.Format.Bmask, srf.Format.Amask)

	tile := LoadImage("blue.bmp")
	src := &sdl.Rect{0, 0, tile.W, tile.H}
	dst := &sdl.Rect{0, 0, tile.W, tile.H}
	for j := int32(0); j < size; j += tile.H {
		for i := int32(0); i < size; i += tile.W {
			dst.X = i
			dst.Y = j
			SDL_BlitSurface(tile, src, image, dst)
		}
	}
	SDL_FreeSurface(tile)

	SDL_LockSurface(image)
	DrawBevel(image, 0, 0, size, size, false, 1)
	DrawBevel(image, 1, 1, size-2, size-2, true, 1)
	SDL_UnlockSurface(image)

	s.activeSlider = AdjustBrightnessTransparent(image, 1.5, false)
	s.slider = SDL_DisplayFormat(image)

	SDL_FreeSurface(image)
}

func (s *Slider) OnMouseButtonDown(button uint8, x, y int32) bool {
	in := IsInRect(x, y, s.left, s.top, s.width, s.height)
	if in {
		sliderX := s.ValueToX(*s.value)
		hl := IsInRect(x, y, s.left+sliderX, s.top, s.height, s.height)
		if hl {
			s.dragging = true
			s.dragOffsetX = x - s.left - sliderX
		}
	}
	return in
}

func (s *Slider) OnMouseButtonUp(button uint8, x, y int32) bool {
	if s.dragging {
		s.dragging = false
		return true
	}
	return false
}

func (s *Slider) ValueToX(value float32) int32 {
	if value < 0 {
		v := float32(0.0)
		s.value = &v
	}
	if value > 1 {
		v := float32(1.0)
		s.value = &v
	}
	return int32(float32(s.width-s.height) * value)
}

func (s *Slider) XToValue(pos int32) float32 {
	if 0 > pos {
		pos = 0
	}
	if (s.width - s.height) < pos {
		pos = s.width - s.height
	}
	return float32(pos) / float32(s.width-s.height)
}

func (s *Slider) OnMouseMove(x, y int32) bool {
	if s.dragging {
		val := s.XToValue(x - s.left - s.dragOffsetX)
		if val != *s.value {
			*s.value = val
			s.Draw()
			sound.SetVolume(*s.value)
		}
		return true
	}

	in := IsInRect(x, y, s.left, s.top, s.width, s.height)
	if in {
		sliderX := s.ValueToX(*s.value)
		hl := IsInRect(x, y, s.left+sliderX, s.top, s.height, s.height)
		if hl != s.highlight {
			s.highlight = hl
			s.Draw()
		}
	} else { //nolint:gocritic
		if s.highlight {
			s.highlight = false
			s.Draw()
		}
	}
	return in
}
