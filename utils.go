package goeinstein

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func GetCornerPixel(surface *sdl.Surface) uint32 {
	SDL_LockSurface(surface)
	bpp := surface.Format.BytesPerPixel
	p := surface.Pixels()

	var binarier interface {
		Uint16([]byte) uint16
		Uint32([]byte) uint32
	} = binary.LittleEndian
	if sdl.BYTEORDER == sdl.BIG_ENDIAN {
		binarier = binary.BigEndian
	}

	var pixel uint32
	switch bpp {
	case 1:
		pixel = uint32(p[0])
	case 2:
		pixel = uint32(binarier.Uint16(p))
	case 3:
		switch sdl.BYTEORDER {
		case sdl.BIG_ENDIAN:
			pixel = uint32(p[2]) | uint32(p[1])<<8 | uint32(p[0])<<16
		case sdl.LIL_ENDIAN:
			pixel = uint32(p[0]) | uint32(p[1])<<8 | uint32(p[2])<<16
		}
	case 4:
		pixel = binarier.Uint32(p)
	}
	SDL_UnlockSurface(surface)
	return pixel
}

func LoadImage(name string) *sdl.Surface {
	return LoadImageTransparent(name, false)
}

func LoadImageTransparent(name string, transparent bool) *sdl.Surface {
	bmp := resources.GetRef(name)
	if bmp == nil {
		panic(fmt.Errorf("%q is not found", name))
	}

	op, err := sdl.RWFromMem(bmp)
	if err != nil {
		panic(fmt.Errorf("rw from mem: %w", err))
	}

	s, err := sdl.LoadBMPRW(op, false)
	if err != nil {
		panic(fmt.Errorf("load BMP: %w", err))
	}

	err = op.Free()
	if err != nil {
		panic(fmt.Errorf("free rw: %w", err))
	}

	resources.DelRef(bmp)

	screenS := s

	if transparent {
		err = screenS.SetColorKey(true, GetCornerPixel(screenS))
		if err != nil {
			panic(fmt.Errorf("set color key: %w", err))
		}
	}
	return screenS
}

func GetTimeOfDat() {
	panic("not implemented")
}

func DrawWallpaper(name string) {
	tile := LoadImage(name)
	src := &sdl.Rect{0, 0, tile.W, tile.H}
	dst := &sdl.Rect{0, 0, tile.W, tile.H}
	for y := int32(0); y < screen.GetHeight(); y += tile.H {
		for x := int32(0); x < screen.GetWidth(); x += tile.W {
			dst.X = x
			dst.Y = y
			err := tile.Blit(src, screen.GetSurface(), dst)
			if err != nil {
				panic(fmt.Errorf("blit (x=%d,y=%d,width=%d,height=%d): %w", x, y, screen.GetWidth(), screen.GetHeight(), err))
			}
		}
	}

	tile.Free()
}

func SetPixel(s *sdl.Surface, x, y int32, r, g, b uint8) error {
	bpp := s.Format.BytesPerPixel
	pixel := sdl.MapRGB(s.Format, r, g, b)

	p := s.Pixels()[y*s.Pitch+x*int32(bpp):]

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

	return nil
}

func GetPixel(surface *sdl.Surface, x, y int32) (r, g, b uint8) {
	bpp := surface.Format.BytesPerPixel
	p := surface.Pixels()[y*surface.Pitch+x*int32(bpp):]

	var binarier interface {
		Uint16([]byte) uint16
		Uint32([]byte) uint32
	} = binary.LittleEndian
	if sdl.BYTEORDER == sdl.BIG_ENDIAN {
		binarier = binary.BigEndian
	}

	var pixel uint32
	switch bpp {
	case 1:
		pixel = uint32(p[0])
	case 2:
		pixel = uint32(binarier.Uint16(p))
	case 3:
		switch sdl.BYTEORDER {
		case sdl.BIG_ENDIAN:
			pixel = uint32(p[2]) | uint32(p[1])<<8 | uint32(p[0])<<16
		case sdl.LIL_ENDIAN:
			pixel = uint32(p[0]) | uint32(p[1])<<8 | uint32(p[2])<<16
		}
	case 4:
		if len(p) < 4 {
			log.Printf("len: %d, p = %v, pixels: %d", len(p), p, len(surface.Pixels()))
			log.Printf("x, y = (%v, %v)", x, y)
			log.Printf("pitch: %v", surface.Pitch)
			log.Printf("bpp: %v", int32(bpp))
			log.Printf("pixels: %v", surface.Pixels())
		}
		pixel = binarier.Uint32(p)
	}
	return sdl.GetRGB(pixel, surface.Format)
}

var (
	gammaTable [256]uint8
	lastGamma  float64 = -1.0
)

func AdjustBrightness(image *sdl.Surface, x, y int32, k float64) {
	if lastGamma != k {
		for i := 0; i <= 255; i++ {
			gammaTable[i] = uint8(255*math.Pow(float64(i)/255.0, 1.0/k) + 0.5)
			if gammaTable[i] > 255 {
				gammaTable[i] = 255
			}
		}
		lastGamma = k
	}

	r, g, b := GetPixel(image, x, y)
	err := SetPixel(image, x, y, gammaTable[r], gammaTable[g], gammaTable[b])
	if err != nil {
		panic(fmt.Errorf("set pixel: %w", err))
	}
}

func AdjustBrightnessTransparent(image *sdl.Surface, k float64, transparent bool) *sdl.Surface {
	if lastGamma != k {
		for i := 0; i <= 255; i++ {
			gammaTable[i] = uint8(255*math.Pow(float64(i)/255.0, 1.0/k) + 0.5)
			if gammaTable[i] > 255 {
				gammaTable[i] = 255
			}
		}
		lastGamma = k
	}

	s, err := image.Convert(image.Format, 0)
	if err != nil {
		panic(fmt.Errorf("convert image surface: %w", err))
	}

	SDL_LockSurface(s)
	for j := int32(0); j < s.H; j++ {
		for i := int32(0); i < s.W; i++ {
			r, g, b := GetPixel(s, i, j)
			err := SetPixel(s, i, j, gammaTable[r], gammaTable[g], gammaTable[b])
			if err != nil {
				panic(fmt.Errorf("set pixel (i=%d,j=%d): %w", i, j, err))
			}
		}
	}
	SDL_UnlockSurface(s)

	if transparent {
		SDL_SetColorKey(s, true, GetCornerPixel(s))
	}
	return s
}

type CenteredBitmap struct {
	Widget

	tile *sdl.Surface
	x, y int32
}

func NewCenteredBitmap(fileName string) *CenteredBitmap {
	c := &CenteredBitmap{}
	c.tile = LoadImage(fileName)
	c.x = (screen.GetWidth() - c.tile.W) / 2
	c.y = (screen.GetHeight() - c.tile.H) / 2
	return c
}

func (c *CenteredBitmap) Close() {
	c.tile.Free()
}

func (c *CenteredBitmap) Draw() {
	screen.Draw(c.x, c.y, c.tile)
	screen.AddRegionToUpdate(c.x, c.y, c.tile.W, c.tile.H)
}

func ShowWindow(parentArea *Area, fileName string) {
	area := NewArea()

	area.Add(parentArea)
	area.Add(NewCenteredBitmap(fileName))
	area.Add(NewAnyKeyAccelDefault())
	area.Run()
	sound.Play("click.wav")
}

func IsInRect(evX, evY, x, y, w, h int32) bool {
	return (evX >= x) && (evX < x+w) && (evY >= y) && (evY < y+h)
}

func SecToStr(time uint64) string {
	hours := time / 3600
	v := time - hours*3600
	minutes := v / 60
	seconds := v - minutes*60

	return fmt.Sprintf("%2.2d:%2.2d:%2.2d", hours, minutes, seconds)
}

func ShowMessageWindow(parentArea *Area, pattern string, width, height int32, font *Font, r, g, b uint8, msg string) {
	area := NewArea()

	x := (screen.GetWidth() - width) / 2
	y := (screen.GetHeight() - height) / 2

	area.Add(parentArea)
	area.Add(NewWindowFrame(x, y, width, height, pattern, 6))
	area.Add(NewLabelAligh(font, x, y, width, height, ALIGN_CENTER, ALIGN_MIDDLE, r, g, b, msg))
	area.Add(NewAnyKeyAccelDefault())
	area.Run()
	sound.Play("click.wav")
}

func DrawBevel(s *sdl.Surface, left, top, width, height int32, raised bool, size int32) {
	var k, f, kAdv, fAdv float64
	if raised {
		k = 2.6
		f = 0.1
		kAdv = -0.2
		fAdv = 0.1
	} else {
		f = 2.6
		k = 0.1
		fAdv = -0.2
		kAdv = 0.1
	}
	for i := int32(0); i < size; i++ {
		for j := i; j < height-i-1; j++ {
			AdjustBrightness(s, left+i, top+j, k)
		}
		for j := i; j < width-i; j++ {
			AdjustBrightness(s, left+j, top+i, k)
		}
		for j := i + 1; j < height-i; j++ {
			AdjustBrightness(s, left+width-i-1, top+j, f)
		}
		for j := i; j < width-i-1; j++ {
			AdjustBrightness(s, left+j, top+height-i-1, f)
		}
		k += kAdv
		f += fAdv
	}
}

func EnsureDirExists(fileName string) {
	_, err := os.Stat(fileName)
	if err == nil {
		return
	}
	if !os.IsNotExist(err) {
		panic(fmt.Errorf("unknown os.stat error: %w", err))
	}
	err = os.MkdirAll(fileName, fs.ModePerm)
	if err != nil {
		panic(fmt.Errorf("cannot create dir (%q): %w", fileName, err))
	}
}

func ReadInt(r io.Reader) int {
	buf := make([]byte, 4)
	n, err := r.Read(buf)
	if err != nil {
		panic(fmt.Errorf("readInt: %w", err))
	}
	if n != 4 {
		panic(fmt.Errorf("wrong len of read bytes: %d", n))
	}

	return int(buf[0]) + int(buf[1])*256 + int(buf[2])*256*256 + int(buf[3])*256*256*256
}

func ReadString(stream io.Reader) string {
	no := ReadInt(stream)
	if no <= 0 {
		panic(fmt.Errorf("wrong read string len (n=%d)", no))
	}
	bs := make([]byte, no)

	n, err := stream.Read(bs)
	if err != nil {
		panic(fmt.Errorf("read stream: %w", err))
	}
	if n != no {
		panic(fmt.Errorf("read less than expected: %d != %d", n, no))
	}
	return string(bs)
}

func WriteInt(w io.Writer, v int) {
	b := make([]byte, 4)
	var ib int

	for i := 0; i < 4; i++ {
		ib = v & 0xFF
		v >>= 8
		b[i] = byte(ib)
	}

	n, err := w.Write(b)
	if err != nil {
		panic(fmt.Errorf("write int: %w", err))
	}
	if n != 4 {
		panic(fmt.Errorf("amount of written bytes != 4"))
	}
}

func WriteString(stream io.Writer, value string) {
	WriteInt(stream, len(value))
	n, err := stream.Write([]byte(value))
	if err != nil {
		panic(fmt.Errorf("write string: %w", err))
	}
	if n != len(value) {
		panic(fmt.Errorf("write full string: %w", err))
	}
}
