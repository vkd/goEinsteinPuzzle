package goeinstein

import (
	"fmt"
	"io"
	"math/rand"
)

func GetThingName(row int, thing Card) string {
	sym := func(s string) string {
		return string([]rune(s)[thing-1])
	}
	m := map[int]string{
		0: "123456",
		1: "ABCDEF",
		2: "ⅠⅡⅢⅣⅤⅥ",
		3: "⚀⚁⚂⚃⚄⚅",
		4: "△▽□◇⬠⭔",
		5: "+−÷×=√",
	}
	if s, ok := m[row]; ok {
		return sym(s)
	}
	var s string
	s += string('A' + rune(row))
	s += ToString(thing)
	return s
}

// NearRule
//
// A <> 5
type NearRule struct {
	Rule

	row1   int
	thing1 Card
	row2   int
	thing2 Card
}

var _ Ruler = (*NearRule)(nil)
var _ HintApplier = (*NearRule)(nil)

func (r *NearRule) GetShowOpts() ShowOptions { return SHOW_HORIZ }

func NewNearRule(puzzle SolvedPuzzle, rand *rand.Rand) *NearRule {
	r := &NearRule{}
	col1 := rand.Intn(PUZZLE_SIZE)
	r.row1 = rand.Intn(PUZZLE_SIZE)
	r.thing1 = puzzle[r.row1][col1]

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

	r.row2 = rand.Intn(PUZZLE_SIZE)
	r.thing2 = puzzle[r.row2][col2]
	return r
}

func NewNearRuleStream(stream io.Reader) *NearRule {
	r := &NearRule{}
	r.row1 = ReadInt(stream)
	r.thing1.ReadFrom(stream)
	r.row2 = ReadInt(stream)
	r.thing2.ReadFrom(stream)
	return r
}

func (r *NearRule) ApplyToCol(pos *Possibilities, col int, nearRow int, nearNum Card, thisRow int, thisNum Card) bool {
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
		if r.ApplyToCol(pos, i, r.row1, r.thing1, r.row2, r.thing2) {
			changed = true
		}
		if r.ApplyToCol(pos, i, r.row2, r.thing2, r.row1, r.thing1) {
			changed = true
		}
	}

	if changed {
		r.Apply(pos)
	}

	return changed
}

func (r *NearRule) ApplyHint(pos *Possibilities, re RuleExcluder) bool {
	var out bool
	if ci, ok := pos.GetCol(r.row1, r.thing1); ok {
		for i := 0; i < PUZZLE_SIZE; i++ {
			if i != ci-1 && i != ci+1 {
				if pos.IsPossible(i, r.row2, r.thing2) {
					pos.Exclude(i, r.row2, r.thing2)
					out = true
				}
			}
		}
		re.ExcludeRule(r)
		return out
	}
	if ci, ok := pos.GetCol(r.row2, r.thing2); ok {
		for i := 0; i < PUZZLE_SIZE; i++ {
			if i != ci-1 && i != ci+1 {
				if pos.IsPossible(i, r.row1, r.thing1) {
					pos.Exclude(i, r.row1, r.thing1)
					out = true
				}
			}
		}
		re.ExcludeRule(r)
		return out
	}
	return out
}

func (r *NearRule) GetAsText() string {
	return GetThingName(r.row1, r.thing1) + " is near to " + GetThingName(r.row2, r.thing2)
}

func (r *NearRule) Draw(x, y int32, iconSet *IconSet, h bool) {
	icon := iconSet.GetLargeIcon(r.row1, r.thing1, h)
	screen.Draw(x, y, icon)
	screen.Draw(x+icon.H, y, iconSet.GetNearHintIcon(h))
	screen.Draw(x+icon.H*2, y, iconSet.GetLargeIcon(r.row2, r.thing2, h))
}

func (r *NearRule) Save(stream io.Writer) {
	WriteString(stream, "near")
	WriteInt(stream, r.row1)
	r.thing1.WriteTo(stream)
	WriteInt(stream, r.row2)
	r.thing2.WriteTo(stream)
}

// DirectionRule
//
// A ... 5
type DirectionRule struct {
	Rule

	row1   int
	thing1 Card
	row2   int
	thing2 Card
}

var _ Ruler = (*DirectionRule)(nil)
var _ HintApplier = (*DirectionRule)(nil)

func (r *DirectionRule) GetShowOpts() ShowOptions { return SHOW_HORIZ }

func NewDirectionRule(puzzle SolvedPuzzle, rand *rand.Rand) *DirectionRule {
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
	r.thing1.ReadFrom(stream)
	r.row2 = ReadInt(stream)
	r.thing2.ReadFrom(stream)
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
	r.thing1.WriteTo(stream)
	WriteInt(stream, r.row2)
	r.thing2.WriteTo(stream)
}

func (r *DirectionRule) ApplyHint(pos *Possibilities, re RuleExcluder) bool {
	out := r.Apply(pos)
	if _, ok := pos.GetCol(r.row1, r.thing1); ok {
		re.ExcludeRule(r)
	}
	if _, ok := pos.GetCol(r.row2, r.thing2); ok {
		re.ExcludeRule(r)
	}
	return out
}

type OpenRule struct {
	Rule

	col, row int
	thing    Card
}

var _ Ruler = (*OpenRule)(nil)

func (r *OpenRule) ApplyOnStart() bool                                  { return true }
func (r *OpenRule) Draw(x, y int32, iconSet *IconSet, highlighted bool) {}
func (r *OpenRule) GetShowOpts() ShowOptions                            { return SHOW_NOTHING }

func NewOpenRule(puzzle SolvedPuzzle, rand *rand.Rand) *OpenRule {
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
	r.thing.ReadFrom(stream)
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
	r.thing.WriteTo(stream)
}

// UnderRule
//
// 5
// A
type UnderRule struct {
	Rule

	row1, row2     int
	thing1, thing2 Card
}

var _ Ruler = (*UnderRule)(nil)
var _ HintApplier = (*UnderRule)(nil)

func (*UnderRule) GetShowOpts() ShowOptions { return SHOW_VERT }

func NewUnderRule(puzzle SolvedPuzzle, rand *rand.Rand) *UnderRule {
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
	r.thing1.ReadFrom(stream)
	r.row2 = ReadInt(stream)
	r.thing2.ReadFrom(stream)
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

func (r *UnderRule) ApplyHint(pos *Possibilities, re RuleExcluder) bool {
	out := r.Apply(pos)
	if _, ok := pos.GetCol(r.row1, r.thing1); ok {
		re.ExcludeRule(r)
	}
	if _, ok := pos.GetCol(r.row2, r.thing2); ok {
		re.ExcludeRule(r)
	}
	return out
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
	r.thing1.WriteTo(stream)
	WriteInt(stream, r.row2)
	r.thing2.WriteTo(stream)
}

type BetweenRule struct {
	Rule

	row1        int
	thing1      Card
	row2        int
	thing2      Card
	centerRow   int
	centerThing Card
}

var _ Ruler = (*BetweenRule)(nil)
var _ HintApplier = (*BetweenRule)(nil)

func (r *BetweenRule) GetShowOpts() ShowOptions { return SHOW_HORIZ }

func NewBetweenRule(puzzle SolvedPuzzle, rand *rand.Rand) *BetweenRule {
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
	r.thing1.ReadFrom(stream)
	r.row2 = ReadInt(stream)
	r.thing2.ReadFrom(stream)
	r.centerRow = ReadInt(stream)
	r.centerThing.ReadFrom(stream)
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
	r.thing1.WriteTo(stream)
	WriteInt(stream, r.row2)
	r.thing2.WriteTo(stream)
	WriteInt(stream, r.centerRow)
	r.centerThing.WriteTo(stream)
}

//nolint:gocyclo
func (r *BetweenRule) ApplyHint(pos *Possibilities, re RuleExcluder) bool {
	var out bool
	out = out || pos.Exclude(0, r.centerRow, r.centerThing)
	out = out || pos.Exclude(PUZZLE_SIZE-1, r.centerRow, r.centerThing)

	if ci, ok := pos.GetCol(r.row1, r.thing1); ok {
		for i := 0; i < PUZZLE_SIZE; i++ {
			if i != ci-1 && i != ci+1 {
				out = out || pos.Exclude(i, r.centerRow, r.centerThing)
			}
			if i != ci-2 && i != ci+2 {
				out = out || pos.Exclude(i, r.row2, r.thing2)
			}
		}
		if ci != 2 && ci != 3 {
			re.ExcludeRule(r)
			return out
		}
	}
	if ci, ok := pos.GetCol(r.centerRow, r.centerThing); ok {
		for i := 0; i < PUZZLE_SIZE; i++ {
			if i != ci-1 && i != ci+1 {
				out = out || pos.Exclude(i, r.row1, r.thing1)
				out = out || pos.Exclude(i, r.row2, r.thing2)
			}
		}
		if r.row1 == r.row2 {
			for i := Card(1); i <= PUZZLE_SIZE; i++ {
				if i != r.thing1 && i != r.thing2 {
					out = out || pos.Exclude(ci-1, r.row1, i)
					out = out || pos.Exclude(ci+1, r.row1, i)
				}
			}
			re.ExcludeRule(r)
			return out
		}
	}
	if ci, ok := pos.GetCol(r.row2, r.thing2); ok {
		for i := 0; i < PUZZLE_SIZE; i++ {
			if i != ci-1 && i != ci+1 {
				out = out || pos.Exclude(i, r.centerRow, r.centerThing)
			}
			if i != ci-2 && i != ci+2 {
				out = out || pos.Exclude(i, r.row1, r.thing1)
			}
		}
		if ci != 2 && ci != 3 {
			re.ExcludeRule(r)
			return out
		}
	}
	return out
}

func GenRule(puzzle *SolvedPuzzle, rand *rand.Rand) Ruler {
	a := rand.Intn(14)
	switch a {
	case 0, 1, 2, 3:
		return NewNearRule(*puzzle, rand)
	case 4:
		return NewOpenRule(*puzzle, rand)
	case 5, 6:
		return NewUnderRule(*puzzle, rand)
	case 7, 8, 9, 10:
		return NewDirectionRule(*puzzle, rand)
	case 11, 12, 13:
		return NewBetweenRule(*puzzle, rand)
	default:
		return GenRule(puzzle, rand)
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
