package goeinstein

import (
	"github.com/veandco/go-sdl2/sdl"
)

//nolint:golint
const (
	WIDTH         = 600
	HEIGHT        = 500
	CLIENT_WIDTH  = 570
	CLIENT_HEIGHT = 390
	START_X       = 115
	START_Y       = 100
)

func ShowDescription(parentArea *Area) {
	d := NewDescription(parentArea)
	d.Run()
}

type Description struct {
	widgets []AreaWidgeter

	prevCmd *CursorCommand
	nextCmd *CursorCommand

	area        *Area
	currentPage int

	titleFont  *Font
	buttonFont *Font
	textFont   *Font

	textHeight int
	text       *TextParser
}

func (d *Description) GetPage(no int) *TextPage { return d.text.GetPage(no) }

func NewDescription(parentArea *Area) *Description {
	d := &Description{
		area: NewArea(),
	}
	d.currentPage = 0
	// d.area.AddManaged(parentArea, false)
	d.titleFont = NewFont("nova.ttf", 26)
	d.buttonFont = NewFont("laudcn2.ttf", 14)
	d.textFont = NewFont("laudcn2.ttf", 16)
	d.textHeight = int(d.textFont.GetHeight("A") * 1.0)
	d.text = NewTextParser(string(resources.GetRef("rules.txt")), d.textFont, START_X, START_Y, CLIENT_WIDTH, CLIENT_HEIGHT)
	d.prevCmd = NewCursorCommand(-1, d, &d.currentPage)
	d.nextCmd = NewCursorCommand(1, d, &d.currentPage)
	return d
}

func (d *Description) Close() {
	d.DeleteWidgets()
	d.text.Close()
	d.titleFont.Close()
	d.buttonFont.Close()
	d.textFont.Close()
}

func (d *Description) DeleteWidgets() {
	for _, w := range d.widgets {
		d.area.Remove(w)
	}
	d.widgets = nil
}

func (d *Description) UpdateInfo() {
	d.DeleteWidgets()
	d.PrintPage()
	d.area.Draw()
}

func (d *Description) Run() {
	d.area.Add(NewWindow(100, 50, WIDTH, HEIGHT, "blue.bmp"))
	d.area.Add(NewLabelAligh(d.titleFont, 250, 60, 300, 40, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, 0, msg("rules")))
	d.area.Add(NewButtonText(110, 515, 80, 25, d.buttonFont, 255, 255, 0, "blue.bmp", msg("prev"), d.prevCmd))
	d.area.Add(NewButtonText(200, 515, 80, 25, d.buttonFont, 255, 255, 0, "blue.bmp", msg("next"), d.nextCmd))
	exitCmd := NewExitCommand(d.area)
	d.area.Add(NewButtonText(610, 515, 80, 25, d.buttonFont, 255, 255, 0, "blue.bmp", msg("close"), exitCmd))
	d.area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))
	d.PrintPage()
	d.area.Run()
}

func (d *Description) PrintPage() {
	page := d.text.GetPage(d.currentPage)
	if page == nil {
		return
	}
	ln := page.GetWidgetsCount()
	for i := 0; i < ln; i++ {
		w := page.GetWidget(i)
		if w != nil {
			d.widgets = append(d.widgets, w)
			d.area.Add(w)
		}
	}
}

type CursorCommand struct {
	step        int
	description *Description
	value       *int
}

var _ Command = (*CursorCommand)(nil)

func NewCursorCommand(s int, d *Description, v *int) *CursorCommand {
	c := &CursorCommand{
		description: d,
	}
	c.step = s
	c.value = v
	return c
}

func (c *CursorCommand) DoAction() {
	if c.value == nil && 0 > c.step {
		return
	}
	newPageNo := *c.value + c.step
	page := c.description.GetPage(newPageNo)
	if page != nil {
		*c.value = newPageNo
		c.description.UpdateInfo()
	}
}

type TextPage struct {
	widgets []AreaWidgeter
}

func NewTextPage() *TextPage {
	t := &TextPage{}
	return t
}

func (t *TextPage) GetWidget(no int) AreaWidgeter { return t.widgets[no] }
func (t *TextPage) GetWidgetsCount() int          { return len(t.widgets) }
func (t *TextPage) Add(widget AreaWidgeter)       { t.widgets = append(t.widgets, widget) }
func (t *TextPage) IsEmpty() bool                 { return len(t.widgets) == 0 }

func (t *TextPage) Close() {
	for _, w := range t.widgets {
		w.Close()
	}
}

type TextParser struct {
	tokenizer  *Tokenizer
	pages      []*TextPage
	font       *Font
	spaceWidth int32
	charHeight int32
	images     map[string]*sdl.Surface
	offsetX    int32
	offsetY    int32
	pageWidth  int32
	pageHeight int32
}

func NewTextParser(text string, font *Font, x, y int32, width, height int32) *TextParser {
	t := &TextParser{
		tokenizer: NewTokenizer(text),
		font:      font,
		images:    make(map[string]*sdl.Surface),
	}
	t.spaceWidth = font.GetWidth(" ")
	t.charHeight = font.GetWidth("A")
	t.offsetX = x
	t.offsetY = y
	t.pageWidth = width
	t.pageHeight = height
	return t
}

func (t *TextParser) Close() {
	for _, p := range t.pages {
		p.Close()
	}
	for _, v := range t.images {
		v.Free()
	}
}

func (t *TextParser) AddLine(page *TextPage, line *string, curPosY *int32, lineWidth *int32) {
	if 0 < len(*line) {
		page.Add(NewLabelShadow(t.font, t.offsetX, t.offsetY+(*curPosY), uint8(GetStorage().GetInt("text_red", 0)), uint8(GetStorage().GetInt("text_green", 0)), uint8(GetStorage().GetInt("text_blue", 100)), *line, false))
		*line = ""
		*curPosY += 10 + t.charHeight
		*lineWidth = 0
	}
}

func (t *TextParser) IsImage(name string) bool {
	ln := len(name)
	return (3 < ln) && ('$' == name[0]) && ('$' == name[ln-1])
}

func (t *TextParser) KeywordToImage(name string) string {
	return name[1 : len(name)-1]
}

func (t *TextParser) GetImage(name string) *sdl.Surface {
	img, ok := t.images[name]
	if !ok {
		img = LoadImage(name)
		t.images[name] = img
	}
	return img
}

func (t *TextParser) ParseNextPage() {
	if t.tokenizer.IsFinished() {
		return
	}

	var curPosY int32
	var lineWidth int32
	page := NewTextPage()
	var line string

	for {
		tn := t.tokenizer.GetNextToken()
		if Eof == tn.GetType() {
			break
		}
		if Para == tn.GetType() {
			if 0 < len(line) {
				t.AddLine(page, &line, &curPosY, &lineWidth)
			}
			if !page.IsEmpty() {
				curPosY += 10
			}
		} else if Word == tn.GetType() {
			word := tn.GetContent()
			if t.IsImage(word) {
				t.AddLine(page, &line, &curPosY, &lineWidth)
				image := t.GetImage(t.KeywordToImage(word))
				if ((image.H + curPosY) < t.pageHeight) || page.IsEmpty() {
					x := t.offsetX + (t.pageWidth-image.W)/2
					page.Add(NewPictureSurface(x, t.offsetY+curPosY, image))
					curPosY += image.H
				} else {
					t.tokenizer.Unget(tn)
					break
				}
			} else {
				width := t.font.GetWidth(word)
				if lineWidth+width > t.pageWidth {
					if lineWidth == 0 {
						line = word
						t.AddLine(page, &line, &curPosY, &lineWidth)
					} else {
						t.AddLine(page, &line, &curPosY, &lineWidth)
						if curPosY >= t.pageHeight {
							t.tokenizer.Unget(tn)
							break
						}
						line = word
						lineWidth = width
					}
				} else {
					lineWidth += width
					if len(line) > 0 {
						line += " "
						lineWidth += t.spaceWidth
					}
					line += word
				}
			}
		}
		if curPosY >= t.pageHeight {
			break
		}
	}
	t.AddLine(page, &line, &curPosY, &lineWidth)
	if !page.IsEmpty() {
		t.pages = append(t.pages, page)
	} else {
		page.Close()
	}
}

func (t *TextParser) GetPage(no int) *TextPage {
	for !t.tokenizer.IsFinished() && len(t.pages) <= no {
		t.ParseNextPage()
	}
	if len(t.pages) <= no {
		return nil
	}
	return t.pages[no]
}
