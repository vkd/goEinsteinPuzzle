package goeinstein

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Font struct {
	font *ttf.Font
	data []byte
}

func NewFont(name string, ptsize int) *Font {
	f := &Font{} //nolint:exhaustivestruct,exhaustruct
	f.data = resources.GetRef(name)
	if f.data == nil {
		panic(fmt.Errorf("%q not found", name))
	}
	op, err := sdl.RWFromMem(f.data)
	if err != nil {
		panic(fmt.Errorf("%q: RW from mem: %w", name, err))
	}
	f.font, err = ttf.OpenFontRW(op, 1, ptsize)
	if err != nil {
		panic(fmt.Errorf("%q: open font RW: %w", name, err))
	}

	if f.font == nil {
		panic(fmt.Sprintf("Error loading font %q", name))
	}
	return f
}

func (f *Font) Close() {
	f.font.Close()
}

func (f *Font) DrawSurface(s *sdl.Surface, x, y int32, r, g, b uint8, shadow bool, text string) {
	if text == "" {
		return
	}

	str := text

	if shadow {
		color := sdl.Color{1, 1, 1, 1}
		surface, err := f.font.RenderUTF8Blended(str, color)
		if err != nil {
			panic(fmt.Errorf("render shadow: %w", err))
		}
		src := &sdl.Rect{0, 0, surface.W, surface.H}
		dst := &sdl.Rect{x + 1, y + 1, surface.W, surface.H}
		err = surface.Blit(src, s, dst)
		if err != nil {
			panic(fmt.Errorf("blit surface shadow: %w", err))
		}
		surface.Free()
	}
	color := sdl.Color{r, g, b, 0}
	surface, err := f.font.RenderUTF8Blended(str, color)
	if err != nil {
		panic(fmt.Errorf("render text: %w", err))
	}
	src := &sdl.Rect{0, 0, surface.W, surface.H}
	dst := &sdl.Rect{x, y, surface.W, surface.H}
	err = surface.Blit(src, s, dst)
	if err != nil {
		panic(fmt.Errorf("blit surface: %w", err))
	}
	surface.Free()
}

func (f *Font) Draw(x, y int32, r, g, b uint8, shadow bool, text string) {
	f.DrawSurface(screen.GetSurface(), x, y, r, g, b, shadow, text)
}

func (f *Font) GetWidth(text string) int32 {
	w, _ := f.GetSize(text)
	return w
}

func (f *Font) GetHeight(text string) int32 {
	_, h := f.GetSize(text)
	return h
}

func (f *Font) GetSize(text string) (w, h int32) {
	wi, hi, err := f.font.SizeUTF8(text)
	if err != nil {
		panic(fmt.Errorf("size utf8: %w", err))
	}
	return int32(wi), int32(hi)
}
