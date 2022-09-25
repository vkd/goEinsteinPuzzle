package goeinstein

import (
	"io"

	"github.com/veandco/go-sdl2/sdl"
)

//nolint:golint,nosnakecase,stylecheck
const (
	HORHINTS_HINTS_COLS  = 3
	HORHINTS_HINTS_ROWS  = 8
	HORHINTS_TILE_GAP_X  = 4
	HORHINTS_TILE_GAP_Y  = 4
	HORHINTS_TILE_X      = 348
	HORHINTS_TILE_Y      = 68
	HORHINTS_TILE_WIDTH  = 48
	HORHINTS_TILE_HEIGHT = 48
)

type HorHints struct {
	Widget

	iconSet       *IconSet
	rules         []Ruler
	excludedRules []Ruler
	numbersArr    []int
	showExcluded  bool
	highlighted   int
}

func NewHorHints(is *IconSet, r *Rules) *HorHints {
	h := &HorHints{}
	h.iconSet = is
	h.Reset(r)
	return h
}

func NewHorHintsStream(is *IconSet, rl *Rules, stream io.Reader) *HorHints {
	h := &HorHints{}
	h.iconSet = is

	qty := ReadInt(stream)

	for i := 0; i < qty; i++ {
		no := ReadInt(stream)
		h.numbersArr = append(h.numbersArr, no)
		r := GetRule(rl, no)
		excluded := ReadInt(stream) > 0
		if excluded {
			h.excludedRules = append(h.excludedRules, r)
			h.rules = append(h.rules, nil)
		} else {
			h.excludedRules = append(h.excludedRules, nil)
			h.rules = append(h.rules, r)
		}
	}

	h.showExcluded = ReadInt(stream) > 0

	x, y, _ := sdl.GetMouseState()
	h.highlighted = h.GetRuleNo(x, y)
	return h
}

func (h *HorHints) Reset(r *Rules) {
	h.rules = nil
	h.excludedRules = nil
	h.numbersArr = nil

	var no int
	for i, rule := range *r {
		if rule.GetShowOpts() == SHOW_HORIZ {
			h.rules = append(h.rules, (*r)[i])
			h.excludedRules = append(h.excludedRules, nil)
			h.numbersArr = append(h.numbersArr, no)
		}
		no++
	}

	h.showExcluded = false

	x, y, _ := sdl.GetMouseState()
	h.highlighted = h.GetRuleNo(x, y)
}

func (h *HorHints) Draw() {
	for i := 0; i < HORHINTS_HINTS_ROWS; i++ {
		for j := 0; j < HORHINTS_HINTS_COLS; j++ {
			h.DrawCellUpdate(j, i, true)
		}
	}
}

func (h *HorHints) DrawCell(col, row int) {
	h.DrawCellUpdate(col, row, true)
}

func (h *HorHints) DrawCellUpdate(col, row int, addToUpdate bool) {
	x := int32(HORHINTS_TILE_X + col*(HORHINTS_TILE_WIDTH*3+HORHINTS_TILE_GAP_X))
	y := int32(HORHINTS_TILE_Y + row*(HORHINTS_TILE_HEIGHT+HORHINTS_TILE_GAP_Y))

	var r Ruler
	no := row*HORHINTS_HINTS_COLS + col
	if no < len(h.rules) {
		if h.showExcluded {
			r = h.excludedRules[no]
		} else {
			r = h.rules[no]
		}
	}

	if options.AutoHints.value {
		if h.showExcluded {
			if r == nil && no < len(h.rules) {
				r = h.rules[no]
				if r != nil {
					switch r.(type) {
					case *DirectionRule:
					default:
						r = nil
					}
				}
			}
		} else if r != nil {
			switch r.(type) {
			case *DirectionRule:
				r = nil
			}
		}
	}

	if r != nil {
		r.Draw(x, y, h.iconSet, no == h.highlighted)
	} else {
		for i := int32(0); i < 3; i++ {
			screen.Draw(x+HORHINTS_TILE_HEIGHT*i, y, h.iconSet.GetEmptyHintIcon())
		}
	}

	if addToUpdate {
		screen.AddRegionToUpdate(x, y, HORHINTS_TILE_WIDTH*3, HORHINTS_TILE_HEIGHT) //nolint:gomnd
	}
}

func (h *HorHints) OnMouseButtonDown(button uint8, x, y int32) bool {
	if button != 3 {
		return false
	}

	no := h.GetRuleNo(x, y)
	if no < 0 {
		return false
	}
	row := no / HORHINTS_HINTS_COLS
	col := no - row*HORHINTS_HINTS_COLS

	if h.showExcluded {
		r := h.excludedRules[no]
		if r != nil {
			sound.Play("whizz.wav")
			h.rules[no] = r
			h.excludedRules[no] = nil
			h.DrawCell(col, row)
		}
	} else {
		h.Exclude(no)
	}

	return true
}

func (h *HorHints) ExcludeRule(r Ruler) {
	rText := r.GetAsText()
	for ri, r := range h.rules {
		if r == nil {
			continue
		}
		if r.GetAsText() == rText {
			h.Exclude(ri)
		}
	}
}

func (h *HorHints) Exclude(no int) {
	row := no / HORHINTS_HINTS_COLS
	col := no - row*HORHINTS_HINTS_COLS

	r := h.rules[no]
	if r != nil {
		sound.Play("whizz.wav")
		h.rules[no] = nil
		h.excludedRules[no] = r
		h.DrawCell(col, row)
	}
}

func (h *HorHints) ToggleExcluded() {
	h.showExcluded = !h.showExcluded
	h.Draw()
}

func (h *HorHints) OnMouseMove(x, y int32) bool {
	no := h.GetRuleNo(x, y)

	if no != h.highlighted {
		old := h.highlighted
		h.highlighted = no
		if h.IsActive(old) {
			row := old / HORHINTS_HINTS_COLS
			col := old - row*HORHINTS_HINTS_COLS
			h.DrawCell(col, row)
		}
		if h.IsActive(no) {
			row := no / HORHINTS_HINTS_COLS
			col := no - row*HORHINTS_HINTS_COLS
			h.DrawCell(col, row)
		}
	}

	return false
}

func (h *HorHints) GetRuleNo(x, y int32) int {
	if !IsInRect(x, y, HORHINTS_TILE_X, HORHINTS_TILE_Y, (HORHINTS_TILE_WIDTH*3+HORHINTS_TILE_GAP_X)*HORHINTS_HINTS_COLS, (HORHINTS_TILE_HEIGHT+HORHINTS_TILE_GAP_Y)*HORHINTS_HINTS_ROWS) {
		return -1
	}

	x = x - HORHINTS_TILE_X //nolint:gocritic
	y = y - HORHINTS_TILE_Y //nolint:gocritic

	col := x / (HORHINTS_TILE_WIDTH*3 + HORHINTS_TILE_GAP_X)
	if col*(HORHINTS_TILE_WIDTH*3+HORHINTS_TILE_GAP_X)+HORHINTS_TILE_WIDTH*3 < x {
		return -1
	}
	row := y / (HORHINTS_TILE_HEIGHT + HORHINTS_TILE_GAP_Y)
	if row*(HORHINTS_TILE_HEIGHT+HORHINTS_TILE_GAP_Y)+HORHINTS_TILE_HEIGHT < y {
		return -1
	}

	no := int(row*HORHINTS_HINTS_COLS + col)
	if no >= len(h.rules) {
		return -1
	}

	return no
}

func (h *HorHints) IsActive(ruleNo int) bool {
	if ruleNo < 0 || ruleNo >= len(h.rules) {
		return false
	}
	var r Ruler
	if h.showExcluded {
		r = h.excludedRules[ruleNo]
	} else {
		r = h.rules[ruleNo]
	}
	return r != nil
}

func (h *HorHints) Save(stream io.Writer) {
	cnt := len(h.numbersArr)
	WriteInt(stream, cnt)
	for i := 0; i < cnt; i++ {
		WriteInt(stream, h.numbersArr[i])
		if h.rules[i] != nil {
			WriteInt(stream, 0)
		} else {
			WriteInt(stream, 1)
		}
	}
	if h.showExcluded {
		WriteInt(stream, 1)
	} else {
		WriteInt(stream, 0)
	}
}
