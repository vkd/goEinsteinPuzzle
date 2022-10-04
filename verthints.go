package goeinstein

import (
	"io"

	"github.com/veandco/go-sdl2/sdl"
)

//nolint:golint,nosnakecase,stylecheck
const (
	VERTHINTS_TILE_NUM    = 15
	VERTHINTS_TILE_GAP    = 4
	VERTHINTS_TILE_X      = 12
	VERTHINTS_TILE_Y      = 495
	VERTHINTS_TILE_WIDTH  = 48
	VERTHINTS_TILE_HEIGHT = 48
)

type VertHints struct {
	Widget

	iconSet       *IconSet
	rules         []Ruler
	excludedRules []Ruler
	numbersArr    []int
	showExcluded  bool
	highlighted   int
}

func NewVertHints(is *IconSet, r *Rules) *VertHints {
	h := &VertHints{}
	h.iconSet = is
	h.Reset(r)
	return h
}

func NewVertHintsStream(is *IconSet, rl *Rules, stream io.Reader) *VertHints {
	v := &VertHints{}
	v.iconSet = is

	qty := ReadInt(stream)

	for i := 0; i < qty; i++ {
		no := ReadInt(stream)
		v.numbersArr = append(v.numbersArr, no)
		r := GetRule(rl, no)
		excluded := ReadInt(stream)
		if excluded > 0 {
			v.excludedRules = append(v.excludedRules, r)
			v.rules = append(v.rules, nil)
		} else {
			v.excludedRules = append(v.excludedRules, nil)
			v.rules = append(v.rules, r)
		}
	}

	v.showExcluded = ReadInt(stream) > 0

	x, y, _ := sdl.GetMouseState()
	v.highlighted = v.GetRuleNo(x, y)
	return v
}

func (v *VertHints) Reset(r *Rules) {
	v.rules = nil
	v.excludedRules = nil
	v.numbersArr = nil

	var no int
	for i, rule := range *r {
		if rule.GetShowOpts() == SHOW_VERT {
			v.rules = append(v.rules, (*r)[i])
			v.excludedRules = append(v.excludedRules, nil)
			v.numbersArr = append(v.numbersArr, no)
		}
		no++
	}

	v.showExcluded = false

	x, y, _ := sdl.GetMouseState()
	v.highlighted = v.GetRuleNo(x, y)
}

func (v *VertHints) Draw() {
	for i := 0; i < VERTHINTS_TILE_NUM; i++ {
		v.DrawCellUpdate(i, true)
	}
}

func (v *VertHints) DrawCell(col int) {
	v.DrawCellUpdate(col, true)
}

func (v *VertHints) DrawCellUpdate(col int, addToUpdate bool) {
	x := int32(VERTHINTS_TILE_X + col*(VERTHINTS_TILE_WIDTH+VERTHINTS_TILE_GAP))
	y := int32(VERTHINTS_TILE_Y)

	var r Ruler
	if col < len(v.rules) {
		if v.showExcluded {
			r = v.excludedRules[col]
		} else {
			r = v.rules[col]
		}
	}

	if options.AutoHints.value {
		if v.showExcluded {
			if r == nil && col < len(v.rules) {
				r = v.rules[col]
			}
		} else {
			r = nil
		}
	}

	if r != nil {
		r.Draw(x, y, v.iconSet, v.highlighted == col)
	} else {
		screen.Draw(x, y, v.iconSet.GetEmptyHintIcon())
		screen.Draw(x, y+VERTHINTS_TILE_HEIGHT, v.iconSet.GetEmptyHintIcon())
	}

	if addToUpdate {
		screen.AddRegionToUpdate(x, y, VERTHINTS_TILE_WIDTH, VERTHINTS_TILE_HEIGHT*2) //nolint:gomnd
	}
}

func (v *VertHints) OnMouseButtonDown(button uint8, x, y int32) bool {
	if button != 3 {
		return false
	}

	no := v.GetRuleNo(x, y)
	if no < 0 {
		return false
	}

	if no < len(v.rules) {
		if v.showExcluded {
			r := v.excludedRules[no]
			if r != nil {
				sound.Play("whizz.wav")
				v.rules[no] = r
				v.excludedRules[no] = nil
				v.DrawCell(no)
			}
		} else {
			v.Exclude(no)
		}
	}

	return true
}

func (v *VertHints) Exclude(no int) {
	r := v.rules[no]
	if r != nil {
		sound.Play("whizz.wav")
		v.rules[no] = nil
		v.excludedRules[no] = r
		v.DrawCell(no)
	}
}

func (v *VertHints) ToggleExcluded() {
	v.showExcluded = !v.showExcluded
	v.Draw()
}

func (v *VertHints) ExcludeRule(r Ruler) {
	rText := r.GetAsText()
	for vi, vr := range v.rules {
		if vr == nil {
			continue
		}
		if vr.GetAsText() == rText {
			v.Exclude(vi)
		}
	}
}

func (v *VertHints) OnMouseMove(x, y int32) bool {
	no := v.GetRuleNo(x, y)

	if no != v.highlighted {
		if no >= 0 && no < len(v.rules) {
			r := v.rules[no]
			if r != nil {
				r.OnMouseMove()
			}
		} else {
			Selected.Clear()
		}
		old := v.highlighted
		v.highlighted = no
		if v.IsActive(old) {
			v.DrawCell(old)
		}
		if v.IsActive(no) {
			v.DrawCell(no)
		}
	}

	return false
}

func (v *VertHints) GetRuleNo(x, y int32) int {
	if !IsInRect(x, y, VERTHINTS_TILE_X, VERTHINTS_TILE_Y, (VERTHINTS_TILE_WIDTH+VERTHINTS_TILE_GAP)*VERTHINTS_TILE_NUM, VERTHINTS_TILE_HEIGHT*2) { //nolint:gomnd
		return -1
	}

	x = x - VERTHINTS_TILE_X //nolint:gocritic
	y = y - VERTHINTS_TILE_Y //nolint:gocritic,ineffassign,staticcheck

	no := x / (VERTHINTS_TILE_WIDTH + VERTHINTS_TILE_GAP)
	if no*(VERTHINTS_TILE_WIDTH+VERTHINTS_TILE_GAP)+VERTHINTS_TILE_WIDTH < x {
		return -1
	}

	return int(no)
}

func (v *VertHints) IsActive(ruleNo int) bool {
	if ruleNo < 0 || ruleNo >= len(v.rules) {
		return false
	}
	var r Ruler
	if v.showExcluded {
		r = v.excludedRules[ruleNo]
	} else {
		r = v.rules[ruleNo]
	}
	return r != nil
}

func (v *VertHints) Save(stream io.Writer) {
	cnt := len(v.numbersArr)
	WriteInt(stream, cnt)
	for i := 0; i < cnt; i++ {
		WriteInt(stream, v.numbersArr[i])
		if v.rules[i] != nil {
			WriteInt(stream, 0)
		} else {
			WriteInt(stream, 1)
		}
	}
	if v.showExcluded {
		WriteInt(stream, 1)
	} else {
		WriteInt(stream, 0)
	}
}
