package goeinstein

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

//nolint:golint,nosnakecase,stylecheck
const PUZZLE_SIZE = 6

type SolvedPuzzle [PUZZLE_SIZE][PUZZLE_SIZE]Card

type Possibilities struct {
	pos [PUZZLE_SIZE][PUZZLE_SIZE]Cell
}

func NewPossibilities() *Possibilities {
	p := &Possibilities{}
	p.Reset()
	return p
}

func NewPossibilitiesStream(stream io.Reader) *Possibilities {
	p := &Possibilities{}
	for row := 0; row < PUZZLE_SIZE; row++ {
		for col := 0; col < PUZZLE_SIZE; col++ {
			p.pos[col][row].ReadFrom(stream)
		}
	}
	return p
}

func (p *Possibilities) Reset() {
	for i := 0; i < PUZZLE_SIZE; i++ {
		for j := 0; j < PUZZLE_SIZE; j++ {
			p.pos[i][j].Reset()
		}
	}
}

func (p *Possibilities) CheckSingles(row int) {
	var cellsCnt [PUZZLE_SIZE]int // count of elements in cells
	var elsCnt [PUZZLE_SIZE]int   // total count of elements in row
	var elements [PUZZLE_SIZE]int // one element of each cell
	var elCells [PUZZLE_SIZE]int  // one cell of each element

	// check if there is only one element left in cell(col, row)
	for col := 0; col < PUZZLE_SIZE; col++ {
		for i := 0; i < PUZZLE_SIZE; i++ {
			if p.pos[col][row][i] > 0 {
				elsCnt[i]++
				elCells[i] = col
				cellsCnt[col]++
				elements[col] = i + 1
			}
		}
	}

	var changed bool

	// check for cells with single element
	for col := 0; col < PUZZLE_SIZE; col++ {
		if cellsCnt[col] == 1 && elsCnt[elements[col]-1] != 1 {
			// there is only one element in cell but it used somewhere else
			e := elements[col] - 1
			for i := 0; i < PUZZLE_SIZE; i++ {
				if i != col {
					p.pos[i][row][e] = 0
				}
			}
			changed = true
		}
	}

	// check for single element without exclusive cell
	for el := 0; el < PUZZLE_SIZE; el++ {
		if elsCnt[el] == 1 && cellsCnt[elCells[el]] != 1 {
			col := elCells[el]
			for i := 0; i < PUZZLE_SIZE; i++ {
				if i != el {
					p.pos[col][row][i] = 0
				}
			}
			changed = true
		}
	}

	if changed {
		p.CheckSingles(row)
	}
}

func (p *Possibilities) Exclude(col, row int, element Card) bool {
	if p.pos[col][row][element-1] == 0 {
		return false
	}

	p.pos[col][row][element-1] = 0
	p.CheckSingles(row)
	return true
}

func (p *Possibilities) Set(col, row int, element Card) {
	p.pos[col][row].Set(element)

	for j := 0; j < PUZZLE_SIZE; j++ {
		if j != col {
			p.pos[j][row].Exclude(element)
		}
	}

	p.CheckSingles(row)
}

func (p *Possibilities) IsPossible(col, row int, element Card) bool {
	return p.pos[col][row].IsPossible(element)
}

func (p *Possibilities) IsDefined(col, row int) bool {
	_, ok := p.pos[col][row].GetDefined()
	return ok
}

func (p *Possibilities) GetDefined(col, row int) (Card, bool) {
	return p.pos[col][row].GetDefined()
}

func (p *Possibilities) IsSolved() bool {
	for i := 0; i < PUZZLE_SIZE; i++ {
		for j := 0; j < PUZZLE_SIZE; j++ {
			if !p.IsDefined(i, j) {
				return false
			}
		}
	}
	return true
}

func (p *Possibilities) IsValid(puzzle *SolvedPuzzle) bool {
	for row := 0; row < PUZZLE_SIZE; row++ {
		for col := 0; col < PUZZLE_SIZE; col++ {
			if !p.IsPossible(col, row, puzzle[row][col]) {
				return false
			}
		}
	}
	return true
}

func (p *Possibilities) GetPosition(row int, element Card) int {
	var cnt int
	lastPos := -1

	for i := 0; i < PUZZLE_SIZE; i++ {
		if p.pos[i][row].IsPossible(element) {
			cnt++
			lastPos = i
		}
	}
	if cnt == 1 {
		return lastPos
	}
	return -1
}

func (p *Possibilities) Print() {
	for row := 0; row < PUZZLE_SIZE; row++ {
		fmt.Fprintf(os.Stdout, "%s ", string('A'+rune(row)))
		for col := 0; col < PUZZLE_SIZE; col++ {
			for i := 0; i < PUZZLE_SIZE; i++ {
				if p.pos[col][row][i] > 0 {
					fmt.Fprintf(os.Stdout, "%d", p.pos[col][row][i])
				} else {
					fmt.Fprint(os.Stdout, " ")
				}
			}
			fmt.Fprint(os.Stdout, "   ")
		}
		fmt.Fprint(os.Stdout, "\n")
	}
}

func (p *Possibilities) Save(stream io.Writer) {
	for row := 0; row < PUZZLE_SIZE; row++ {
		for col := 0; col < PUZZLE_SIZE; col++ {
			p.pos[col][row].WriteTo(stream)
		}
	}
}

func (p *Possibilities) GetCol(row int, element Card) (int, bool) {
	for i := 0; i < PUZZLE_SIZE; i++ {
		if c, ok := p.GetDefined(i, row); ok && c == element {
			return i, true
		}
	}
	return 0, false
}

func Shuffle(arr *[PUZZLE_SIZE]Card) {
	var a, b int
	var c Card

	for i := 0; i < 30; i++ {
		a = rand.Intn(PUZZLE_SIZE)
		b = rand.Intn(PUZZLE_SIZE)
		c = arr[a]
		arr[a] = arr[b]
		arr[b] = c
	}
}

func CanSolve(puzzle *SolvedPuzzle, rules *Rules) bool {
	pos := NewPossibilities()
	var changed bool

	for {
		changed = false
		for _, rule := range *rules {
			if rule.Apply(pos) {
				changed = true
				if !pos.IsValid(puzzle) {
					fmt.Fprint(os.Stdout, "after error:\n")
					pos.Print()
					panic(fmt.Sprintf("Invalid possibilities after rule %s", rule.GetAsText()))
				}
			}
		}
		if !changed {
			break
		}
	}

	res := pos.IsSolved()
	return res
}

func RemoveRules(puzzle *SolvedPuzzle, rules *Rules) {
	var possible bool

	for {
		possible = false
		for ri := range *rules {
			excludedRules := append(append(Rules{}, (*rules)[:ri]...), (*rules)[ri+1:]...)
			if CanSolve(puzzle, &excludedRules) {
				possible = true
				*rules = excludedRules
				break
			}
		}
		if !possible {
			break
		}
	}
}

func GenRules(puzzle *SolvedPuzzle, rules *Rules) {
	var rulesDone bool

	for {
		rule := GenRule(puzzle)
		if rule != nil {
			s := rule.GetAsText()
			for _, r := range *rules {
				if r.GetAsText() == s {
					rule = nil
					break
				}
			}
			if rule != nil {
				*rules = append(*rules, rule)
				rulesDone = CanSolve(puzzle, rules)
			}
		}
		if rulesDone {
			break
		}
	}
}

func GenPuzzle(puzzle *SolvedPuzzle, rules *Rules) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < PUZZLE_SIZE; i++ {
		for j := 0; j < PUZZLE_SIZE; j++ {
			puzzle[i][j] = Card(j + 1)
		}
		Shuffle(&(*puzzle)[i])
	}

	GenRules(puzzle, rules)
	RemoveRules(puzzle, rules)
}

func OpenInitial(possib *Possibilities, rules *Rules, re RuleExcluder) {
	for _, r := range *rules {
		if r.ApplyOnStart() {
			r.Apply(possib)
		}
	}
	if options.AutoHints.value {
		rules.ApplyHints(possib, re)
	}
}

type HintApplier interface {
	ApplyHint(*Possibilities, RuleExcluder) bool
}

type RuleExcluder interface {
	ExcludeRule(r Ruler)
}

func GetHintsQty(rules *Rules, vert, horiz *int) {
	*vert = 0
	*horiz = 0

	for _, r := range *rules {
		so := r.GetShowOpts()
		switch so {
		case SHOW_VERT:
			*vert++
		case SHOW_HORIZ:
			*horiz++
		case SHOW_NOTHING:
		}
	}
}

func SavePuzzle(puzzle *SolvedPuzzle, stream io.Writer) {
	for row := 0; row < PUZZLE_SIZE; row++ {
		for col := 0; col < PUZZLE_SIZE; col++ {
			(*puzzle)[row][col].WriteTo(stream)
		}
	}
}

func LoadPuzzle(puzzle *SolvedPuzzle, stream io.Reader) {
	for row := 0; row < PUZZLE_SIZE; row++ {
		for col := 0; col < PUZZLE_SIZE; col++ {
			puzzle[row][col].ReadFrom(stream)
		}
	}
}

func GetRule(rules *Rules, no int) Ruler {
	var j int
	for i := range *rules {
		if j == no {
			return (*rules)[i]
		}
		j++
	}
	panic("Rule is not found")
}

type ShowOptions int8

//nolint:golint,nosnakecase,stylecheck
const (
	SHOW_VERT ShowOptions = iota
	SHOW_HORIZ
	SHOW_NOTHING
)

type Ruler interface {
	Close()
	GetAsText() string
	Apply(*Possibilities) bool
	ApplyOnStart() bool
	GetShowOpts() ShowOptions
	Draw(x, y int32, iconSet *IconSet, h bool)
	Save(io.Writer)
}

type Rule struct{}

func (r Rule) Close() {}

func (r Rule) ApplyOnStart() bool {
	return false
}

func (r Rule) GetShowOpts() ShowOptions {
	return SHOW_NOTHING
}

type Rules []Ruler

func (rs *Rules) ApplyHints(pos *Possibilities, re RuleExcluder) {
	updated := true
	for updated {
		updated = false
		for _, r := range *rs {
			if ra, ok := r.(HintApplier); ok {
				if ra.ApplyHint(pos, re) {
					updated = true
				}
			}
		}
	}
}
