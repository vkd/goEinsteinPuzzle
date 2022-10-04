package goeinstein

import (
	"github.com/veandco/go-sdl2/sdl"
)

type MenuBackground struct {
	Widget
}

func NewMenuBackground() *MenuBackground {
	return &MenuBackground{}
}

func (m *MenuBackground) Draw() {
	title := LoadImage("nova.bmp")
	screen.Draw(0, 0, title)
	title.Free()
	font := NewFont("nova.ttf", 28)
	s := msg("einsteinFlowix")
	width := font.GetWidth(s)
	font.Draw((screen.GetWidth()-width)/2, 30, 255, 255, 255, true, s)
	urlFont := NewFont("luximb.ttf", 16)
	s = "http://games.flowix.com"
	width = urlFont.GetWidth(s)
	urlFont.Draw((screen.GetWidth()-width)/2, 60, 255, 255, 0, true, s)
	screen.AddRegionToUpdate(0, 0, screen.GetWidth(), screen.GetHeight())
}

type NewGameCommand struct {
	area *Area
}

func NewNewGameCommand(a *Area) *NewGameCommand {
	n := &NewGameCommand{}
	n.area = a
	return n
}

func (n *NewGameCommand) DoAction() {
	game = NewGame()
	game.Run()
	n.area.UpdateMouse()
	n.area.Draw()
}

type LoadGameCommand struct {
	area *Area
}

var _ Command = (*LoadGameCommand)(nil)

func NewLoadGameCommand(a *Area) *LoadGameCommand {
	l := &LoadGameCommand{}
	l.area = a
	return l
}

func (l *LoadGameCommand) DoAction() {
	game = LoadGame(l.area)
	if game != nil {
		game.Run()
		game.Close()
	}
	l.area.UpdateMouse()
	l.area.Draw()
}

type TopScoresCommand struct {
	area *Area
}

var _ Command = (*TopScoresCommand)(nil)

func NewTopScoresCommand(a *Area) *TopScoresCommand {
	l := &TopScoresCommand{}
	l.area = a
	return l
}

func (l *TopScoresCommand) DoAction() {
	scores := NewTopScores()
	ShowScoresWindow(l.area, scores)
	l.area.UpdateMouse()
	l.area.Draw()
}

type RulesCommand struct {
	area *Area
}

var _ Command = (*RulesCommand)(nil)

func NewRulesCommand(a *Area) *RulesCommand {
	l := &RulesCommand{}
	l.area = a
	return l
}

func (l *RulesCommand) DoAction() {
	ShowDescription(l.area)
	l.area.UpdateMouse()
	l.area.Draw()
}

type OptionsCommand struct {
	area *Area
}

var _ Command = (*OptionsCommand)(nil)

func NewOptionsCommand(a *Area) *OptionsCommand {
	l := &OptionsCommand{}
	l.area = a
	return l
}

func (l *OptionsCommand) DoAction() {
	ShowOptionsWindow(l.area)
	l.area.UpdateMouse()
	l.area.Draw()
}

type AboutCommand struct {
	parentArea *Area
}

var _ Command = (*AboutCommand)(nil)

func NewAboutCommand(a *Area) *AboutCommand {
	l := &AboutCommand{}
	l.parentArea = a
	return l
}

func (l *AboutCommand) DoAction() {
	area := NewArea()
	titleFont := NewFont("nova.ttf", 26)
	font := NewFont("laudcn2.ttf", 14)
	urlFont := NewFont("luximb.ttf", 16)

	LABEL := func(pos int32, c uint8, f *Font, text string) {
		area.Add(NewLabelAligh(f, 220, pos, 360, 20, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, c, text))
	}
	area.Add(l.parentArea)
	area.Add(NewWindow(220, 160, 360, 280, "blue.bmp"))
	area.Add(NewLabelAligh(titleFont, 250, 165, 300, 40, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, 0, msg("about")))
	LABEL(240, 255, font, msg("einsteinPuzzle"))
	LABEL(260, 255, font, msg("version"))
	LABEL(280, 255, font, msg("copyright"))
	LABEL(330, 0, urlFont, "http://games.flowix.com")

	exitCmd := NewExitCommand(area)
	area.Add(NewButtonText(360, 400, 80, 25, font, 255, 255, 0, "blue.bmp", msg("ok"), exitCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))
	area.Add(NewKeyAccel(sdl.K_RETURN, exitCmd))
	area.Run()

	l.parentArea.UpdateMouse()
	l.parentArea.Draw()
}

func NewMenuButton(y int32, font *Font, text string, cmd Command) *Button {
	b := NewButtonColor(550, y, 220, 30, font, 0, 240, 240, 30, 255, 255, text, cmd)
	return b
}

func Menu() {
	area := NewArea()
	font := NewFont("laudcn2.ttf", 20)

	area.Add(NewMenuBackground())
	area.Draw()

	newGameCmd := NewNewGameCommand(area)
	area.Add(NewMenuButton(340, font, msg("newGame"), newGameCmd))
	loadGameCmd := NewLoadGameCommand(area)
	area.Add(NewMenuButton(370, font, msg("loadGame"), loadGameCmd))
	topScoresCmd := NewTopScoresCommand(area)
	area.Add(NewMenuButton(400, font, msg("topScores"), topScoresCmd))
	rulesCmd := NewRulesCommand(area)
	area.Add(NewMenuButton(430, font, msg("rules"), rulesCmd))
	optionsCmd := NewOptionsCommand(area)
	area.Add(NewMenuButton(460, font, msg("options"), optionsCmd))
	aboutCmd := NewAboutCommand(area)
	area.Add(NewMenuButton(490, font, msg("about"), aboutCmd))
	exitMenuCmd := NewExitCommand(area)
	area.Add(NewMenuButton(520, font, msg("exit"), exitMenuCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitMenuCmd))

	area.Draw()
	screen.AddRegionToUpdate(0, 0, screen.GetWidth(), screen.GetHeight())
	screen.Flush()

	area.Run()
}
