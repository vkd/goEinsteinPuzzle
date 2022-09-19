package goeinstein

import (
	"encoding/binary"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type VideoMode struct {
	width      int32
	height     int32
	bpp        uint8
	fullScreen bool
}

func NewVideoMode(w, h int32, bpp uint8, fullscreen bool) VideoMode {
	v := VideoMode{}
	v.width = w
	v.height = h
	v.bpp = bpp
	v.fullScreen = fullscreen
	return v
}

func (v VideoMode) GetWidth() int32    { return v.width }
func (v VideoMode) GetHeight() int32   { return v.height }
func (v VideoMode) GetBpp() uint8      { return v.bpp }
func (v VideoMode) IsFullScreen() bool { return v.fullScreen }

type Screen struct {
	screen              *sdl.Surface
	fullScreen          bool
	mouseImage          *sdl.Surface
	mouseSave           *sdl.Surface
	regions             []sdl.Rect
	mouseVisible        bool
	regionsList         []*sdl.Rect
	maxRegionsList      int
	saveX, saveY        int32
	niceCursor          bool
	cursor, emptyCursor *sdl.Cursor

	window *sdl.Window
}

func (s *Screen) GetSurface() *sdl.Surface { return s.screen }

func NewScreen() *Screen {
	s := &Screen{}
	s.screen = nil
	s.mouseImage = nil
	s.mouseSave = nil
	s.mouseVisible = false
	s.regionsList = nil
	s.maxRegionsList = 0
	return s
}

func (s *Screen) Close() {
	sdl.SetCursor(s.cursor)
	if s.mouseImage != nil {
		SDL_FreeSurface(s.mouseImage)
	}
	if s.mouseSave != nil {
		SDL_FreeSurface(s.mouseSave)
	}
	if s.regionsList != nil {
		s.regionsList = nil
	}
}

func (s *Screen) GetVideoMode() VideoMode {
	return NewVideoMode(s.screen.W, s.screen.H, s.screen.Format.BitsPerPixel, s.fullScreen)
}

func (s *Screen) SetFullscreen(isFullscreen bool) {
	m := s.GetVideoMode()
	m.fullScreen = isFullscreen
	s.SetMode(m)
}

func (s *Screen) SetMode(mode VideoMode) {
	s.fullScreen = mode.IsFullScreen()

	var flags uint32 = sdl.SWSURFACE
	if s.fullScreen {
		flags |= sdl.WINDOW_FULLSCREEN
	}
	if s.window != nil {
		SDL_DestroyWindow(s.window)
	}
	s.window = SDL_CreateWindow("Einstein", mode.GetWidth(), mode.GetHeight(), flags)
	s.screen = SDL_GetWindowSurface(s.window)
}

func (s *Screen) GetFullScreenModes() []VideoMode { return nil }

func (s *Screen) GetWidth() int32 {
	if s.screen != nil {
		return s.screen.W
	}
	panic("No video mode selected")
}

func (s *Screen) GetHeight() int32 {
	if s.screen != nil {
		return s.screen.H
	}
	panic("No video mode selected")
}

func (s *Screen) CenterMouse() {
	if s.screen != nil {
		SDL_WarpMouse(s.screen.W/2, s.screen.H/2)
	} else {
		panic("No video mode selected")
	}
}

func (s *Screen) SetMouseImage(image *sdl.Surface) {
	if s.mouseImage != nil {
		SDL_FreeSurface(s.mouseImage)
		s.mouseImage = nil
	}
	if s.mouseSave != nil {
		SDL_FreeSurface(s.mouseSave)
		s.mouseSave = nil
	}

	if image == nil {
		return
	}

	s.mouseImage = SDL_DisplayFormat(image)
	if s.mouseImage == nil {
		panic("Error creating surface")
	}
	// s.mouseSave = SDL_DisplayFormat(image)
	s.mouseSave = SDL_CreateRGBSurface(sdl.SWSURFACE, image.W, image.H, int32(s.screen.Format.BitsPerPixel), s.screen.Format.Rmask, s.screen.Format.Gmask, s.screen.Format.Bmask, screen.screen.Format.Amask)
	if s.mouseSave == nil {
		SDL_FreeSurface(s.mouseImage)
		panic("Error creating buffer surface")
	}
	SDL_SetColorKey(s.mouseImage, true, sdl.MapRGB(s.mouseImage.Format, 0, 0, 0))
}

func (s *Screen) HideMouse() {
	if !s.mouseVisible {
		return
	}
	if !s.niceCursor {
		s.mouseVisible = false
		return
	}

	if s.mouseSave != nil {
		src := sdl.Rect{0, 0, s.mouseSave.W, s.mouseSave.H}
		dst := sdl.Rect{s.saveX, s.saveY, s.mouseSave.W, s.mouseSave.H}
		if src.W > 0 {
			SDL_BlitSurface(s.mouseSave, &src, s.screen, &dst)
			s.AddRegionToUpdate(dst.X, dst.Y, dst.W, dst.H)
		}
	}
	s.mouseVisible = false
}

func (s *Screen) ShowMouse() {
	if s.mouseVisible {
		return
	}
	if !s.niceCursor {
		s.mouseVisible = true
		return
	}

	if s.mouseImage != nil && s.mouseSave != nil {
		x, y, _ := sdl.GetMouseState()
		s.saveX = x
		s.saveY = y
		src := &sdl.Rect{0, 0, s.mouseSave.W, s.mouseSave.H}
		dst := &sdl.Rect{x, y, s.mouseImage.W, s.mouseImage.H}
		if src.W > 0 {
			SDL_BlitSurface(s.screen, dst, s.mouseSave, src)
			SDL_BlitSurface(s.mouseImage, src, s.screen, dst)
			s.AddRegionToUpdate(dst.X, dst.Y, dst.W, dst.H)
		}
	}
	s.mouseVisible = true
}

func (s *Screen) UpdateMouse() {
	s.HideMouse()
	s.ShowMouse()
}

func (s *Screen) Flush() {
	if len(s.regions) == 0 {
		return
	}

	if s.regionsList == nil {
		s.regionsList = make([]*sdl.Rect, len(s.regions))
		s.maxRegionsList = len(s.regions)
	} else { //nolint:gocritic
		if s.maxRegionsList < len(s.regions) {
			r := make([]*sdl.Rect, len(s.regions))
			s.regionsList = r
			s.maxRegionsList = len(s.regions)
		}
	}

	for i := range s.regions {
		s.regionsList[i] = &s.regions[i]
	}

	SDL_UpdateRects(s.window, s.regions)
	s.regions = nil
}

func (s *Screen) AddRegionToUpdate(x, y, w, h int32) {
	if ((x >= s.GetWidth()) || (y >= s.GetHeight())) || (0 >= w) || (0 >= h) {
		return
	}
	if (x+w < 0) || (y+h < 0) {
		return
	}
	if x+w > s.GetWidth() {
		w = s.GetWidth() - x
	}
	if y+h > s.GetHeight() {
		h = s.GetHeight() - y
	}
	if 0 > x {
		w = w + x //nolint:gocritic
		x = 0
	}
	if 0 > y {
		h = h + y //nolint:gocritic
		y = 0
	}
	r := sdl.Rect{x, y, w, h}
	s.regions = append(s.regions, r)
}

func (s *Screen) SetPixel(x, y int32, r, g, b uint8) {
	SDL_LockSurface(s.screen)
	bpp := s.screen.Format.BytesPerPixel
	pixel := sdl.MapRGB(s.screen.Format, r, g, b)
	/* Here p is the address to the pixel we want to set */
	p := s.screen.Pixels()[y*s.screen.Pitch+x*int32(bpp):]

	var binarier interface {
		PutUint16([]byte, uint16)
		PutUint32([]byte, uint32)
	} = binary.LittleEndian
	if sdl.BYTEORDER == sdl.BIG_ENDIAN {
		binarier = binary.BigEndian
	}

	switch bpp {
	case 1:
		p[0] = byte(pixel)
	case 2:
		binarier.PutUint16(p, uint16(pixel))
	case 3:
		switch sdl.BYTEORDER {
		case sdl.BIG_ENDIAN:
			p[0] = byte(pixel >> 16)
			p[1] = byte(pixel >> 8)
			p[2] = byte(pixel)
		case sdl.LIL_ENDIAN:
			p[0] = byte(pixel)
			p[1] = byte(pixel >> 8)
			p[2] = byte(pixel >> 16)
		}
	case 4:
		binarier.PutUint32(p, pixel)
	}
	SDL_UnlockSurface(s.screen)
}

func (s *Screen) Draw(x, y int32, tile *sdl.Surface) {
	src := &sdl.Rect{0, 0, tile.W, tile.H}
	dst := &sdl.Rect{x, y, tile.W, tile.H}
	SDL_BlitSurface(tile, src, s.screen, dst)
}

func (s *Screen) SetCursor(nice bool) {
	if nice == s.niceCursor {
		return
	}

	oldVisible := s.mouseVisible
	if s.mouseVisible {
		s.HideMouse()
	}
	s.niceCursor = nice

	if s.niceCursor {
		SDL_SetCursor(s.emptyCursor)
	} else {
		SDL_SetCursor(s.cursor)
	}

	if oldVisible {
		s.ShowMouse()
	}
}

func (s *Screen) InitCursors() {
	s.cursor = sdl.GetCursor()
	var data, mask uint8
	s.emptyCursor = sdl.CreateCursor(&data, &mask, 8, 1, 0, 0)
}

func (s *Screen) DoneCursors() {
	if s.niceCursor {
		SDL_SetCursor(s.cursor)
	}
	SDL_FreeCursor(s.emptyCursor)
}

func (s *Screen) CreateSubimage(x, y, width, height int32) *sdl.Surface {
	srf, err := sdl.CreateRGBSurface(sdl.SWSURFACE, width, height, int32(s.screen.Format.BitsPerPixel), s.screen.Format.Rmask, s.screen.Format.Gmask, s.screen.Format.Bmask, s.screen.Format.Amask)
	if err != nil {
		panic(fmt.Errorf("Error creating buffer surface: %w", err))
	}
	src := &sdl.Rect{x, y, width, height}
	dst := &sdl.Rect{0, 0, width, height}
	SDL_BlitSurface(s.screen, src, srf, dst)
	return srf
}
