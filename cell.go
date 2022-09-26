package goeinstein

import (
	"io"
)

type Cell [PUZZLE_SIZE]Card

func (c *Cell) Reset() {
	for i := range *c {
		(*c)[i] = Card(i + 1)
	}
}

func (c *Cell) Set(el Card) {
	for i := range *c {
		if i == int(el-1) {
			(*c)[i] = el
		} else {
			(*c)[i] = 0
		}
	}
}

func (c *Cell) Exclude(el Card) {
	(*c)[el-1] = 0
}

func (c *Cell) IsPossible(el Card) bool {
	return (*c)[el-1] == el
}

func (c *Cell) GetDefined() (Card, bool) {
	var card Card
	var found bool
	for i := range *c {
		if (*c)[i] != 0 {
			if found {
				return 0, false
			}
			card = (*c)[i]
			found = true
		}
	}
	return card, found
}

func (c *Cell) WriteTo(w io.Writer) {
	for i := range *c {
		WriteInt(w, int((*c)[i]))
	}
}

func (c *Cell) ReadFrom(r io.Reader) {
	for i := range *c {
		(*c)[i] = Card(ReadInt(r))
	}
}
