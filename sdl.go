//nolint:golint,stylecheck
package goeinstein

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func SDL_BlitSurface(src *sdl.Surface, srcRect *sdl.Rect, dst *sdl.Surface, dstRect *sdl.Rect) {
	err := src.Blit(srcRect, dst, dstRect)
	if err != nil {
		panic(fmt.Errorf("SDL_BlitSurface: %w", err))
	}
}

func SDL_DisplayFormat(s *sdl.Surface) *sdl.Surface {
	out, err := s.Convert(s.Format, 0)
	if err != nil {
		panic(fmt.Errorf("SDL_DisplayFormat: %w", err))
	}
	return out
}

//nolint:gocritic
func SDL_CreateRGBSurface(flags uint32, width, height, depth int32, Rmask, Gmask, Bmask, Amask uint32) *sdl.Surface {
	out, err := sdl.CreateRGBSurface(flags, width, height, depth, Rmask, Gmask, Bmask, Amask)
	if err != nil {
		panic(fmt.Errorf("SDL_CreateRGBSurface: %w", err))
	}
	return out
}

func SDL_FreeSurface(s *sdl.Surface) {
	s.Free()
}

func SDL_LockSurface(s *sdl.Surface) {
	err := s.Lock()
	if err != nil {
		panic(fmt.Errorf("SDL_LockSurface: %w", err))
	}
}

func SDL_UnlockSurface(s *sdl.Surface) {
	s.Unlock()
}

func SDL_SetColorKey(s *sdl.Surface, flag bool, key uint32) {
	err := s.SetColorKey(flag, key)
	if err != nil {
		panic(fmt.Errorf("SDL_SetColorKey: %w", err))
	}
}

type SDLKey = sdl.Keycode // *sdl.Keysym // *sdl.Keycode

func SDL_SetClipRect(s *sdl.Surface, r *sdl.Rect) {
	s.SetClipRect(r)
}

func SDL_FillRect(s *sdl.Surface, rect *sdl.Rect, color uint32) {
	err := s.FillRect(rect, color)
	if err != nil {
		panic(fmt.Errorf("SDL_FillRect: %w", err))
	}
}

func Mix_PlayChannel(c *mix.Chunk, channel, loops int) {
	_, err := c.Play(channel, loops)
	if err != nil {
		panic(fmt.Errorf("play chunk: %w", err))
	}
}

func SDL_CreateWindow(title string, w, h int32, flags uint32) *sdl.Window {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, flags)
	if err != nil {
		panic(fmt.Errorf("SDL_CreateWindow: %w", err))
	}
	return window
}

func SDL_GetWindowSurface(window *sdl.Window) *sdl.Surface {
	s, err := window.GetSurface()
	if err != nil {
		panic(fmt.Errorf("SDL_GetWindowSurface: %w", err))
	}
	return s
}

func SDL_WarpMouse(x, y int32) {
	err := sdl.WarpMouseGlobal(x, y)
	if err != nil {
		panic(fmt.Errorf("SDL_WarpMouse: %w", err))
	}
}

func SDL_SetCursor(cursor *sdl.Cursor) {
	sdl.SetCursor(cursor)
}

func SDL_FreeCursor(cursor *sdl.Cursor) {
	sdl.FreeCursor(cursor)
}

func SDL_UpdateRects(window *sdl.Window, rects []sdl.Rect) {
	err := window.UpdateSurfaceRects(rects)
	if err != nil {
		panic(fmt.Errorf("SDL_UpdateRects: %w", err))
	}
}

func SDL_DestroyWindow(window *sdl.Window) {
	err := window.Destroy()
	if err != nil {
		panic(fmt.Errorf("SDL_DestroyWindow: %w", err))
	}
}
