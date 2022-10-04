package goeinstein

import "github.com/veandco/go-sdl2/sdl"

type IconSet struct {
	smallIcons                    [6][6][2]*sdl.Surface
	largeIcons                    [6][6][2]*sdl.Surface
	emptyFieldIcon, emptyHintIcon *sdl.Surface
	nearHintIcon                  [2]*sdl.Surface
	sideHintIcon, betweenArrow    [2]*sdl.Surface

	BorderLarge *sdl.Surface
	BorderSmall *sdl.Surface
}

func NewIconSet() *IconSet {
	s := &IconSet{}
	buf := []rune("xy.bmp")

	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			buf[1] = rune('1' + j)
			buf[0] = rune('a' + i)
			s.smallIcons[i][j][0] = LoadImage("small-" + string(buf))
			s.smallIcons[i][j][1] = AdjustBrightnessTransparent(s.smallIcons[i][j][0], 1.5, false)
			buf[0] = rune('A' + i)
			s.largeIcons[i][j][0] = LoadImage("large-" + string(buf))
			s.largeIcons[i][j][1] = AdjustBrightnessTransparent(s.largeIcons[i][j][0], 1.5, false)
		}
	}
	s.emptyFieldIcon = LoadImage("tile.bmp")
	s.emptyHintIcon = LoadImage("hint-tile.bmp")
	s.nearHintIcon[0] = LoadImage("hint-near.bmp")
	s.nearHintIcon[1] = AdjustBrightnessTransparent(s.nearHintIcon[0], 1.5, false)
	s.sideHintIcon[0] = LoadImage("hint-side.bmp")
	s.sideHintIcon[1] = AdjustBrightnessTransparent(s.sideHintIcon[0], 1.5, false)
	s.betweenArrow[0] = LoadImageTransparent("betwarr.bmp", true)
	s.betweenArrow[1] = AdjustBrightnessTransparent(s.betweenArrow[0], 1.5, false)

	s.BorderLarge = LoadImageTransparent("border-large.bmp", true)
	s.BorderSmall = LoadImageTransparent("border-small.bmp", true)
	return s
}

func (s *IconSet) Close() {
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			for k := 0; k < 2; k++ {
				s.smallIcons[i][j][k].Free()
				s.largeIcons[i][j][k].Free()
			}
		}
	}
	s.emptyFieldIcon.Free()
	s.emptyHintIcon.Free()
	s.nearHintIcon[0].Free()
	s.nearHintIcon[1].Free()
	s.sideHintIcon[0].Free()
	s.sideHintIcon[1].Free()
	s.betweenArrow[0].Free()
	s.betweenArrow[1].Free()

	s.BorderLarge.Free()
	s.BorderSmall.Free()
}

func (s *IconSet) GetLargeIcon(row int, num Card, h bool) *sdl.Surface {
	var br int
	if h {
		br = 1
	}
	return s.largeIcons[row][num-1][br]
}

func (s *IconSet) GetSmallIcon(row int, num Card, h bool) *sdl.Surface {
	var br int
	if h {
		br = 1
	}
	return s.smallIcons[row][num-1][br]
}

func (s *IconSet) GetEmptyFieldIcon() *sdl.Surface { return s.emptyFieldIcon }
func (s *IconSet) GetEmptyHintIcon() *sdl.Surface  { return s.emptyHintIcon }

func (s *IconSet) GetNearHintIcon(h bool) *sdl.Surface {
	var br int
	if h {
		br = 1
	}
	return s.nearHintIcon[br]
}

func (s *IconSet) GetSideHintIcon(h bool) *sdl.Surface {
	var br int
	if h {
		br = 1
	}
	return s.sideHintIcon[br]
}

func (s *IconSet) GetBetweenArrow(h bool) *sdl.Surface {
	var br int
	if h {
		br = 1
	}
	return s.betweenArrow[br]
}
