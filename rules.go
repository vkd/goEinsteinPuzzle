package goeinstein

import (
	"fmt"
	"io"
	"math/rand"
)

func GetThingName(row int, thing int) string {
	var s string
	s += string('A' + rune(row))
	s += ToString(thing)
	return s
}

type NearRule struct {
	Rule

	thing1 [2]int
	thing2 [2]int
}

var _ Ruler = (*NearRule)(nil)

func (r *NearRule) GetShowOpts() ShowOptions { return SHOW_HORIZ }

func NewNearRule(puzzle SolvedPuzzle) *NearRule {
	r := &NearRule{}
	col1 := rand.Intn(PUZZLE_SIZE)
	r.thing1[0] = rand.Intn(PUZZLE_SIZE)
	r.thing1[1] = puzzle[r.thing1[0]][col1]

	var col2 int
	if col1 == 0 {
		col2 = 1
	} else {
		if col1 == PUZZLE_SIZE-1 {
			col2 = PUZZLE_SIZE - 2
		} else {
			if rand.Intn(2) > 0 {
				col2 = col1 + 1
			} else {
				col2 = col1 - 1
			}
		}
	}

	r.thing2[0] = rand.Intn(PUZZLE_SIZE)
	r.thing2[1] = puzzle[r.thing2[0]][col2]
	return r
}

func NewNearRuleStream(stream io.Reader) *NearRule {
	r := &NearRule{}
	r.thing1[0] = ReadInt(stream)
	r.thing1[1] = ReadInt(stream)
	r.thing2[0] = ReadInt(stream)
	r.thing2[1] = ReadInt(stream)
	return r
}

func (r *NearRule) ApplyToCol(pos *Possibilities, col int, nearRow int, nearNum int, thisRow int, thisNum int) bool {
	var hasLeft, hasRight bool

	if col == 0 {
		hasLeft = false
	} else {
		hasLeft = pos.IsPossible(col-1, nearRow, nearNum)
	}
	if col == PUZZLE_SIZE-1 {
		hasRight = false
	} else {
		hasRight = pos.IsPossible(col+1, nearRow, nearNum)
	}

	if !hasRight && !hasLeft && pos.IsPossible(col, thisRow, thisNum) {
		pos.Exclude(col, thisRow, thisNum)
		return true
	}
	return false
}

func (r *NearRule) Apply(pos *Possibilities) bool {
	var changed bool

	for i := 0; i < PUZZLE_SIZE; i++ {
		if r.ApplyToCol(pos, i, r.thing1[0], r.thing1[1], r.thing2[0], r.thing2[1]) {
			changed = true
		}
		if r.ApplyToCol(pos, i, r.thing2[0], r.thing2[1], r.thing1[0], r.thing1[1]) {
			changed = true
		}
	}

	if changed {
		r.Apply(pos)
	}

	return changed
}

func (r *NearRule) GetAsText() string {
	return GetThingName(r.thing1[0], r.thing1[1]) + " is near to " + GetThingName(r.thing2[0], r.thing2[1])
}

func (r *NearRule) Draw(x, y int32, iconSet *IconSet, h bool) {
	icon := iconSet.GetLargeIcon(r.thing1[0], r.thing1[1], h)
	screen.Draw(x, y, icon)
	screen.Draw(x+icon.H, y, iconSet.GetNearHintIcon(h))
	screen.Draw(x+icon.H*2, y, iconSet.GetLargeIcon(r.thing2[0], r.thing2[1], h))
}

func (r *NearRule) Save(stream io.Writer) {
	WriteString(stream, "near")
	WriteInt(stream, r.thing1[0])
	WriteInt(stream, r.thing1[1])
	WriteInt(stream, r.thing2[0])
	WriteInt(stream, r.thing2[1])
}

type DirectionRule struct {
	Rule

	row1, thing1 int
	row2, thing2 int
}

var _ Ruler = (*DirectionRule)(nil)

func (r *DirectionRule) GetShowOpts() ShowOptions { return SHOW_HORIZ }

func NewDirectionRule(puzzle SolvedPuzzle) *DirectionRule {
	r := &DirectionRule{}
	r.row1 = rand.Intn(PUZZLE_SIZE)
	r.row2 = rand.Intn(PUZZLE_SIZE)
	col1 := rand.Intn(PUZZLE_SIZE - 1)
	col2 := rand.Intn(PUZZLE_SIZE-col1-1) + col1 + 1
	r.thing1 = puzzle[r.row1][col1]
	r.thing2 = puzzle[r.row2][col2]
	return r
}

func NewDirectionRuleStream(stream io.Reader) *DirectionRule {
	r := &DirectionRule{}
	r.row1 = ReadInt(stream)
	r.thing1 = ReadInt(stream)
	r.row2 = ReadInt(stream)
	r.thing2 = ReadInt(stream)
	return r
}

func (r *DirectionRule) Apply(pos *Possibilities) bool {
	var changed bool

	for i := 0; i < PUZZLE_SIZE; i++ {
		if pos.IsPossible(i, r.row2, r.thing2) {
			pos.Exclude(i, r.row2, r.thing2)
			changed = true
		}
		if pos.IsPossible(i, r.row1, r.thing1) {
			break
		}
	}

	for i := PUZZLE_SIZE - 1; i >= 0; i-- {
		if pos.IsPossible(i, r.row1, r.thing1) {
			pos.Exclude(i, r.row1, r.thing1)
			changed = true
		}
		if pos.IsPossible(i, r.row2, r.thing2) {
			break
		}
	}

	return changed
}

func (r *DirectionRule) GetAsText() string {
	return GetThingName(r.row1, r.thing1) + " is from the left of " + GetThingName(r.row2, r.thing2)
}

func (r *DirectionRule) Draw(x, y int32, iconSet *IconSet, h bool) {
	icon := iconSet.GetLargeIcon(r.row1, r.thing1, h)
	screen.Draw(x, y, icon)
	screen.Draw(x+icon.H, y, iconSet.GetSideHintIcon(h))
	screen.Draw(x+icon.H*2, y, iconSet.GetLargeIcon(r.row2, r.thing2, h))
}

func (r *DirectionRule) Save(stream io.Writer) {
	WriteString(stream, "direction")
	WriteInt(stream, r.row1)
	WriteInt(stream, r.thing1)
	WriteInt(stream, r.row2)
	WriteInt(stream, r.thing2)
}

func (r *DirectionRule) OpenInitials(pos *Possibilities) {
	pos.Exclude(PUZZLE_SIZE-1, r.row1, r.thing1)
	pos.Exclude(0, r.row2, r.thing2)
}

type OpenRule struct {
	Rule

	col, row, thing int
}

var _ Ruler = (*OpenRule)(nil)

func (r *OpenRule) ApplyOnStart() bool                                  { return true }
func (r *OpenRule) Draw(x, y int32, iconSet *IconSet, highlighted bool) {}
func (r *OpenRule) GetShowOpts() ShowOptions                            { return SHOW_NOTHING }

func NewOpenRule(puzzle SolvedPuzzle) *OpenRule {
	r := &OpenRule{}
	r.col = rand.Intn(PUZZLE_SIZE)
	r.row = rand.Intn(PUZZLE_SIZE)
	r.thing = puzzle[r.row][r.col]
	return r
}

func NewOpenRuleStream(stream io.Reader) *OpenRule {
	r := &OpenRule{}
	r.col = ReadInt(stream)
	r.row = ReadInt(stream)
	r.thing = ReadInt(stream)
	return r
}

func (r *OpenRule) Apply(pos *Possibilities) bool {
	if !pos.IsDefined(r.col, r.row) {
		pos.Set(r.col, r.row, r.thing)
		return true
	}
	return false
}

func (r *OpenRule) GetAsText() string {
	return GetThingName(r.row, r.thing) + " is at column " + ToString(r.col+1)
}

func (r *OpenRule) Save(stream io.Writer) {
	WriteString(stream, "open")
	WriteInt(stream, r.col)
	WriteInt(stream, r.row)
	WriteInt(stream, r.thing)
}

type UnderRule struct {
	Rule

	row1, thing1, row2, thing2 int
}

var _ Ruler = (*UnderRule)(nil)

func (*UnderRule) GetShowOpts() ShowOptions { return SHOW_VERT }

func NewUnderRule(puzzle SolvedPuzzle) *UnderRule {
	r := &UnderRule{}
	col := rand.Intn(PUZZLE_SIZE)
	r.row1 = rand.Intn(PUZZLE_SIZE)
	r.thing1 = puzzle[r.row1][col]
	for {
		r.row2 = rand.Intn(PUZZLE_SIZE)
		if r.row2 != r.row1 {
			break
		}
	}
	r.thing2 = puzzle[r.row2][col]
	return r
}

func NewUnderRuleStream(stream io.Reader) *UnderRule {
	r := &UnderRule{}
	r.row1 = ReadInt(stream)
	r.thing1 = ReadInt(stream)
	r.row2 = ReadInt(stream)
	r.thing2 = ReadInt(stream)
	return r
}

func (r *UnderRule) Apply(pos *Possibilities) bool {
	var changed bool

	for i := 0; i < PUZZLE_SIZE; i++ {
		if !pos.IsPossible(i, r.row1, r.thing1) && pos.IsPossible(i, r.row2, r.thing2) {
			pos.Exclude(i, r.row2, r.thing2)
			changed = true
		}
		if !pos.IsPossible(i, r.row2, r.thing2) && pos.IsPossible(i, r.row1, r.thing1) {
			pos.Exclude(i, r.row1, r.thing1)
			changed = true
		}
	}

	return changed
}

func (r *UnderRule) GetAsText() string {
	return GetThingName(r.row1, r.thing1) + " is the same column as " + GetThingName(r.row2, r.thing2)
}

func (r *UnderRule) Draw(x, y int32, iconSet *IconSet, h bool) {
	icon := iconSet.GetLargeIcon(r.row1, r.thing1, h)
	screen.Draw(x, y, icon)
	screen.Draw(x, y+icon.H, iconSet.GetLargeIcon(r.row2, r.thing2, h))
}

func (r *UnderRule) Save(stream io.Writer) {
	WriteString(stream, "under")
	WriteInt(stream, r.row1)
	WriteInt(stream, r.thing1)
	WriteInt(stream, r.row2)
	WriteInt(stream, r.thing2)
}

type BetweenRule struct {
	Rule

	row1, thing1           int
	row2, thing2           int
	centerRow, centerThing int
}

var _ Ruler = (*BetweenRule)(nil)

func (r *BetweenRule) GetShowOpts() ShowOptions { return SHOW_HORIZ }

func NewBetweenRule(puzzle SolvedPuzzle) *BetweenRule {
	r := &BetweenRule{}
	r.centerRow = rand.Intn(PUZZLE_SIZE)
	r.row1 = rand.Intn(PUZZLE_SIZE)
	r.row2 = rand.Intn(PUZZLE_SIZE)

	centerCol := rand.Intn(PUZZLE_SIZE-2) + 1
	r.centerThing = puzzle[r.centerRow][centerCol]
	if rand.Intn(2) > 0 {
		r.thing1 = puzzle[r.row1][centerCol-1]
		r.thing2 = puzzle[r.row2][centerCol+1]
	} else {
		r.thing1 = puzzle[r.row1][centerCol+1]
		r.thing2 = puzzle[r.row2][centerCol-1]
	}
	return r
}

func NewBetweenRuleStream(stream io.Reader) *BetweenRule {
	r := &BetweenRule{}
	r.row1 = ReadInt(stream)
	r.thing1 = ReadInt(stream)
	r.row2 = ReadInt(stream)
	r.thing2 = ReadInt(stream)
	r.centerRow = ReadInt(stream)
	r.centerThing = ReadInt(stream)
	return r
}

func (r *BetweenRule) Apply(pos *Possibilities) bool {
	var changed bool

	if pos.IsPossible(0, r.centerRow, r.centerThing) {
		changed = true
		pos.Exclude(0, r.centerRow, r.centerThing)
	}

	if pos.IsPossible(PUZZLE_SIZE-1, r.centerRow, r.centerThing) {
		changed = true
		pos.Exclude(PUZZLE_SIZE-1, r.centerRow, r.centerThing)
	}

	var goodLoop bool
	for {
		goodLoop = false

		for i := 1; i < PUZZLE_SIZE-1; i++ {
			if pos.IsPossible(i, r.centerRow, r.centerThing) {
				if !((pos.IsPossible(i-1, r.row1, r.thing1) &&
					pos.IsPossible(i+1, r.row2, r.thing2)) || (pos.IsPossible(i-1, r.row2, r.thing2) &&
					pos.IsPossible(i+1, r.row1, r.thing1))) {
					pos.Exclude(i, r.centerRow, r.centerThing)
					goodLoop = true
				}
			}
		}

		for i := 0; i < PUZZLE_SIZE; i++ {
			var leftPossible, rightPossible bool

			if pos.IsPossible(i, r.row2, r.thing2) {
				if i < 2 {
					leftPossible = false
				} else {
					leftPossible = pos.IsPossible(i-1, r.centerRow, r.centerThing) && pos.IsPossible(i-2, r.row1, r.thing1)
				}
				if i >= PUZZLE_SIZE-2 {
					rightPossible = false
				} else {
					rightPossible = pos.IsPossible(i+1, r.centerRow, r.centerThing) && pos.IsPossible(i+2, r.row1, r.thing1)
				}
				if !leftPossible && !rightPossible {
					pos.Exclude(i, r.row2, r.thing2)
					goodLoop = true
				}
			}

			if pos.IsPossible(i, r.row1, r.thing1) {
				if i < 2 {
					leftPossible = false
				} else {
					leftPossible = pos.IsPossible(i-1, r.centerRow, r.centerThing) && pos.IsPossible(i-2, r.row2, r.thing2)
				}
				if i >= PUZZLE_SIZE-2 {
					rightPossible = false
				} else {
					rightPossible = pos.IsPossible(i+1, r.centerRow, r.centerThing) && pos.IsPossible(i+2, r.row2, r.thing2)
				}
				if !leftPossible && !rightPossible {
					pos.Exclude(i, r.row1, r.thing1)
					goodLoop = true
				}
			}
		}
		if goodLoop {
			changed = true
		}
		if !goodLoop {
			break
		}
	}
	return changed
}

func (r *BetweenRule) GetAsText() string {
	return GetThingName(r.centerRow, r.centerThing) + " is between " + GetThingName(r.row1, r.thing1) + " and " + GetThingName(r.row2, r.thing2)
}

func (r *BetweenRule) Draw(x, y int32, iconSet *IconSet, h bool) {
	icon := iconSet.GetLargeIcon(r.row1, r.thing1, h)
	screen.Draw(x, y, icon)
	screen.Draw(x+icon.W, y, iconSet.GetLargeIcon(r.centerRow, r.centerThing, h))
	screen.Draw(x+icon.W*2, y, iconSet.GetLargeIcon(r.row2, r.thing2, h))
	arrow := iconSet.GetBetweenArrow(h)
	screen.Draw(x+icon.W-(arrow.W-icon.W)/2, y+0, arrow)
}

func (r *BetweenRule) Save(stream io.Writer) {
	WriteString(stream, "between")
	WriteInt(stream, r.row1)
	WriteInt(stream, r.thing1)
	WriteInt(stream, r.row2)
	WriteInt(stream, r.thing2)
	WriteInt(stream, r.centerRow)
	WriteInt(stream, r.centerThing)
}

func (r *BetweenRule) OpenInitials(pos *Possibilities) {
	pos.Exclude(0, r.centerRow, r.centerThing)
	pos.Exclude(PUZZLE_SIZE-1, r.centerRow, r.centerThing)
}

func GenRule(puzzle *SolvedPuzzle) Ruler {
	a := rand.Intn(14)
	switch a {
	case 0, 1, 2, 3:
		return NewNearRule(*puzzle)
	case 4:
		return NewOpenRule(*puzzle)
	case 5, 6:
		return NewUnderRule(*puzzle)
	case 7, 8, 9, 10:
		return NewDirectionRule(*puzzle)
	case 11, 12, 13:
		return NewBetweenRule(*puzzle)
	default:
		return GenRule(puzzle)
	}
}

func SaveRules(rules *Rules, stream io.Writer) {
	WriteInt(stream, len(*rules))
	for _, rule := range *rules {
		rule.Save(stream)
	}
}

func LoadRules(rules *Rules, stream io.Reader) {
	no := ReadInt(stream)

	for i := 0; i < no; i++ {
		ruleType := ReadString(stream)
		var r Ruler
		if ruleType == "near" { //nolint:gocritic
			r = NewNearRuleStream(stream)
		} else if ruleType == "open" {
			r = NewOpenRuleStream(stream)
		} else if ruleType == "under" {
			r = NewUnderRuleStream(stream)
		} else if ruleType == "direction" {
			r = NewDirectionRuleStream(stream)
		} else if ruleType == "between" {
			r = NewBetweenRuleStream(stream)
		} else {
			panic(fmt.Errorf("invalid rule type: %q", ruleType))
		}
		*rules = append(*rules, r)
	}
}
