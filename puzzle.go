package goeinstein

import "github.com/veandco/go-sdl2/sdl"

//nolint:golint,nosnakecase,stylecheck
const (
	FIELD_OFFSET_X    = 12
	FIELD_OFFSET_Y    = 68
	FIELD_GAP_X       = 4
	FIELD_GAP_Y       = 4
	FIELD_TILE_WIDTH  = 48
	FIELD_TILE_HEIGHT = 48
)

type Puzzle struct {
	Widget

	possib                  *Possibilities
	iconSet                 *IconSet
	valid                   bool
	win                     bool
	solved                  *SolvedPuzzle
	hCol, hRow              int
	subHNo                  int
	winCommand, failCommand Command

	hinter Hinter
}

type Hinter interface {
	AutoHint(*Possibilities)
}

func (p *Puzzle) GetPossibilities() *Possibilities { return p.possib }
func (p *Puzzle) IsValid() bool                    { return p.valid }
func (p *Puzzle) Victory() bool                    { return p.win }

func NewPuzzle(is *IconSet, s *SolvedPuzzle, p *Possibilities, h Hinter) *Puzzle {
	pz := &Puzzle{}
	pz.iconSet = is
	pz.solved = s
	pz.possib = p
	pz.hinter = h

	pz.Reset()
	return pz
}

func (p *Puzzle) Close() {}

func (p *Puzzle) Reset() {
	p.valid = true
	p.win = false

	x, y, _ := sdl.GetMouseState()
	p.GetCellNo(x, y, &p.hCol, &p.hRow, &p.subHNo)
}

func (p *Puzzle) Draw() {
	for i := 0; i < PUZZLE_SIZE; i++ {
		for j := 0; j < PUZZLE_SIZE; j++ {
			p.DrawCellUpdate(i, j, true)
		}
	}
}

func (p *Puzzle) DrawCell(col, row int) {
	p.DrawCellUpdate(col, row, true)
}

func (p *Puzzle) DrawCellUpdate(col, row int, addToUpdate bool) {
	posX := int32(FIELD_OFFSET_X + col*(FIELD_TILE_WIDTH+FIELD_GAP_X))
	posY := int32(FIELD_OFFSET_Y + row*(FIELD_TILE_HEIGHT+FIELD_GAP_Y))

	if p.possib.IsDefined(col, row) {
		element := p.possib.GetDefined(col, row)
		if element > 0 {
			screen.Draw(posX, posY, p.iconSet.GetLargeIcon(row, element, (p.hCol == col) && (p.hRow == row)))
		}
	} else {
		screen.Draw(posX, posY, p.iconSet.GetEmptyFieldIcon())
		x := posX
		y := posY + (FIELD_TILE_HEIGHT / 6)
		for i := 0; i < 6; i++ {
			if p.possib.IsPossible(col, row, i+1) {
				screen.Draw(x, y, p.iconSet.GetSmallIcon(row, i+1, (p.hCol == col) && (p.hRow == row) && (i+1 == p.subHNo)))
			}
			if i == 2 {
				x = posX
				y += (FIELD_TILE_HEIGHT / 3)
			} else {
				x += (FIELD_TILE_WIDTH / 3)
			}
		}
	}
	if addToUpdate {
		screen.AddRegionToUpdate(posX, posY, FIELD_TILE_WIDTH, FIELD_TILE_HEIGHT)
	}
}

func (p *Puzzle) DrawRow(row int) {
	p.DrawRowUpdate(row, true)
}

func (p *Puzzle) DrawRowUpdate(row int, addToUpdate bool) {
	for i := 0; i < PUZZLE_SIZE; i++ {
		p.DrawCellUpdate(i, row, addToUpdate)
	}
}

func (p *Puzzle) OnMouseButtonDown(button uint8, x, y int32) bool {
	var col, row, element int

	if !p.GetCellNo(x, y, &col, &row, &element) {
		return false
	}

	if !p.possib.IsDefined(col, row) {
		// 	if button == 3 {
		// 		for i := 1; i <= PUZZLE_SIZE; i++ {
		// 			p.possib.MakePossible(col, row, i)
		// 			p.DrawCell(col, row)
		// 		}
		// 	}
		// } else {
		if element == -1 {
			return false
		}
		if button == 1 {
			if p.possib.IsPossible(col, row, element) {
				p.possib.Set(col, row, element)
				sound.Play("laser.wav")
			}
		} else if button == 3 {
			if p.possib.IsPossible(col, row, element) {
				p.possib.Exclude(col, row, element)
				sound.Play("whizz.wav")
			}
			// else {
			// 	p.possib.MakePossible(col, row, element)
			// }
		}

		if options.AutoHints.value {
			p.hinter.AutoHint(p.possib)
		}

		p.Draw()
	}

	valid := p.possib.IsValid(p.solved)
	if !valid {
		p.OnFail()
	} else { //nolint:gocritic
		if p.possib.IsSolved() && p.valid {
			p.OnVictory()
		}
	}

	return true
}

func (p *Puzzle) OnFail() {
	if p.failCommand != nil {
		p.failCommand.DoAction()
	}
}

func (p *Puzzle) OnVictory() {
	if p.winCommand != nil {
		p.winCommand.DoAction()
	}
}

func (p *Puzzle) GetCellNo(x, y int32, col, row *int, subNo *int) bool {
	*col = -1
	*row = -1
	*subNo = -1

	if !IsInRect(x, y, FIELD_OFFSET_X, FIELD_OFFSET_Y, (FIELD_TILE_WIDTH+FIELD_GAP_X)*PUZZLE_SIZE, (FIELD_TILE_HEIGHT+FIELD_GAP_Y)*PUZZLE_SIZE) {
		return false
	}

	x = x - FIELD_OFFSET_X //nolint:gocritic
	y = y - FIELD_OFFSET_Y //nolint:gocritic

	*col = int(x) / (FIELD_TILE_WIDTH + FIELD_GAP_X)
	if (*col)*(FIELD_TILE_WIDTH+FIELD_GAP_X)+FIELD_TILE_WIDTH < int(x) {
		return false
	}
	*row = int(y) / (FIELD_TILE_HEIGHT + FIELD_GAP_Y)
	if *row*(FIELD_TILE_HEIGHT+FIELD_GAP_Y)+FIELD_TILE_HEIGHT < int(y) {
		return false
	}

	x = x - int32(*col)*(FIELD_TILE_WIDTH+FIELD_GAP_X) //nolint:gocritic
	y = y - int32(*row)*(FIELD_TILE_HEIGHT+FIELD_GAP_Y) - FIELD_TILE_HEIGHT/6
	if (y < 0) || (y >= (FIELD_TILE_HEIGHT/3)*2) {
		return true
	}
	cCol := int(x) / (FIELD_TILE_WIDTH / 3)
	if cCol >= 3 {
		*col = -1
		*row = -1
		return false
	}
	cRow := int(y) / (FIELD_TILE_HEIGHT / 3)
	*subNo = cRow*3 + cCol + 1

	return true
}

func (p *Puzzle) OnMouseMove(x, y int32) bool {
	oldCol := p.hCol
	oldRow := p.hRow
	oldElement := p.subHNo

	p.GetCellNo(x, y, &p.hCol, &p.hRow, &p.subHNo)
	if (p.hCol != oldCol) || (p.hRow != oldRow) || (p.subHNo != oldElement) {
		if (oldCol != -1) && (oldRow != -1) {
			p.DrawCell(oldCol, oldRow)
		}
		if (p.hCol != -1) && (p.hRow != -1) {
			p.DrawCell(p.hCol, p.hRow)
		}
	}

	return false
}

func (p *Puzzle) SetCommand(win, fail Command) {
	p.winCommand = win
	p.failCommand = fail
}
