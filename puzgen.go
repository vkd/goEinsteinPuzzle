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

type SolvedPuzzle [PUZZLE_SIZE][PUZZLE_SIZE]int

type Possibilities struct {
	pos [PUZZLE_SIZE][PUZZLE_SIZE][PUZZLE_SIZE]int
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
			for element := 0; element < PUZZLE_SIZE; element++ {
				p.pos[col][row][element] = ReadInt(stream)
			}
		}
	}
	return p
}

func (p *Possibilities) Reset() {
	for i := 0; i < PUZZLE_SIZE; i++ {
		for j := 0; j < PUZZLE_SIZE; j++ {
			for k := 0; k < PUZZLE_SIZE; k++ {
				p.pos[i][j][k] = k + 1
			}
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

func (p *Possibilities) Exclude(col, row int, element int) {
	if p.pos[col][row][element-1] == 0 {
		return
	}

	p.pos[col][row][element-1] = 0
	p.CheckSingles(row)
}

func (p *Possibilities) Set(col, row int, element int) {
	for i := 0; i < PUZZLE_SIZE; i++ {
		if i != (element - 1) {
			p.pos[col][row][i] = 0
		} else {
			p.pos[col][row][i] = element
		}
	}

	for j := 0; j < PUZZLE_SIZE; j++ {
		if j != col {
			p.pos[j][row][element-1] = 0
		}
	}

	p.CheckSingles(row)
}

func (p *Possibilities) IsPossible(col, row int, element int) bool {
	return p.pos[col][row][element-1] == element
}

func (p *Possibilities) IsDefined(col, row int) bool {
	var solvedCnt, unsolvedCnt int
	for i := 0; i < PUZZLE_SIZE; i++ {
		if p.pos[col][row][i] == 0 {
			unsolvedCnt++
		} else {
			solvedCnt++
		}
	}
	return (unsolvedCnt == PUZZLE_SIZE-1) && (solvedCnt == 1)
}

func (p *Possibilities) GetDefined(col, row int) int {
	for i := 0; i < PUZZLE_SIZE; i++ {
		if p.pos[col][row][i] > 0 {
			return i + 1
		}
	}
	return 0
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

func (p *Possibilities) GetPosition(row int, element int) int {
	var cnt int
	lastPos := -1

	for i := 0; i < PUZZLE_SIZE; i++ {
		if p.pos[i][row][element-1] == element {
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

func (p *Possibilities) MakePossible(col, row int, element int) {
	p.pos[col][row][element-1] = element
}

func (p *Possibilities) Save(stream io.Writer) {
	for row := 0; row < PUZZLE_SIZE; row++ {
		for col := 0; col < PUZZLE_SIZE; col++ {
			for element := 0; element < PUZZLE_SIZE; element++ {
				WriteInt(stream, p.pos[col][row][element])
			}
		}
	}
}

func Shuffle(arr *[PUZZLE_SIZE]int) {
	var a, b, c int

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
			puzzle[i][j] = j + 1
		}
		Shuffle(&(*puzzle)[i])
	}

	GenRules(puzzle, rules)
	RemoveRules(puzzle, rules)
}

func OpenInitial(possib *Possibilities, rules *Rules) {
	for _, r := range *rules {
		if r.ApplyOnStart() {
			r.Apply(possib)
		}
	}
	if options.OpenInitials.value {
		var updated bool = true
		for updated {
			updated = false
			for _, r := range *rules {
				if oi, ok := r.(interface{ OpenInitials(*Possibilities) bool }); ok {
					if oi.OpenInitials(possib) {
						updated = true
					}
				}
			}
		}
	}
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
			WriteInt(stream, (*puzzle)[row][col])
		}
	}
}

func LoadPuzzle(puzzle *SolvedPuzzle, stream io.Reader) {
	for row := 0; row < PUZZLE_SIZE; row++ {
		for col := 0; col < PUZZLE_SIZE; col++ {
			puzzle[row][col] = ReadInt(stream)
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
